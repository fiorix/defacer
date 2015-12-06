package apiserver

import (
	"errors"
	"image"
	"image/draw"
	"io"
	"runtime"
	"sync"

	"github.com/lazywei/go-opencv/opencv"

	"github.com/fiorix/defacer/apiserver/internal"
)

// A Defacer can scan and deface people's faces in images.
type Defacer interface {
	// Deface reads binary image bytes from a given reader
	// and returns a defaced version of the image.
	Deface(io.Reader) (image.Image, error)
}

// NewDefacer creates and initializes a new Defacer.
func NewDefacer(resizer ImageResizer) (Defacer, error) {
	hc, err := internal.DefaultHaarCascade()
	if err != nil {
		return nil, err
	}
	df := &defacer{
		Resizer:     resizer,
		HaarCascade: hc,
	}
	return df, nil
}

type defacer struct {
	sync.Mutex
	Resizer     ImageResizer
	HaarCascade *opencv.HaarCascade
}

// Deface implements the Defacer interface.
func (df *defacer) Deface(r io.Reader) (image.Image, error) {
	img, faces, err := df.scan(r)
	if err != nil {
		return nil, err
	}
	b := img.Bounds()
	dst := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(dst, dst.Bounds(), image.Transparent, image.ZP, draw.Src)
	draw.DrawMask(dst, dst.Bounds(), img, b.Min, img, b.Min, draw.Over)
	switch len(faces) {
	case 0: // nothing to do
	case 1:
		df.draw(nil, nil, dst, faces[0])
	default:
		mu, wg := &sync.Mutex{}, &sync.WaitGroup{}
		for _, rect := range faces {
			wg.Add(1)
			go df.draw(mu, wg, dst, rect)
		}
		wg.Wait()
	}
	return dst, nil
}

// scan reads binary image data from the given reader and scans for
// faces, returning a slice of rectangles where faces were detected.
func (df *defacer) scan(src io.Reader) (m image.Image, r []image.Rectangle, err error) {
	img, _, err := image.Decode(src)
	if err != nil {
		return nil, nil, err
	}
	df.Lock()
	defer df.Unlock()
	cvimg := opencv.FromImage(img)
	if cvimg == nil {
		return nil, nil, errors.New("failed to load source image")
	}
	faces := df.HaarCascade.DetectObjects(cvimg)
	if faces == nil {
		return img, []image.Rectangle{}, nil
	}
	fr := make([]image.Rectangle, len(faces))
	for i, rect := range faces {
		x, y, w, h := rect.X(), rect.Y(), rect.Width(), rect.Height()
		fr[i] = image.Rectangle{
			image.Point{
				roundDown(x),
				roundDown(y),
			},
			image.Point{
				roundUp(x + w),
				roundUp(y + h),
			},
		}
	}
	return img, fr, nil
}

// draw blends the deface image onto dst, of the size of the given rectangle.
func (df *defacer) draw(mu *sync.Mutex, wg *sync.WaitGroup, dst draw.Image, r image.Rectangle) {
	size := image.Point{r.Max.X - r.Min.X, r.Max.Y - r.Min.Y}
	img := df.Resizer.Resize(size)
	b := img.Bounds()
	if mu != nil {
		mu.Lock()
		defer mu.Unlock()
	}
	draw.DrawMask(dst, r, img, b.Min, img, b.Min, draw.Over)
	if wg != nil {
		wg.Done()
	}
	defacerImageDefaceSum.Inc()
}

type defacerPool struct {
	Inbox   chan *defacerReq
	Resizer ImageResizer
}

type defacerReq struct {
	Reader io.Reader
	Resp   chan *defacerResp
}

type defacerResp struct {
	Image image.Image
	Error error
}

// NewDefacerPool creates a pool of Defacers.
func NewDefacerPool(resizer ImageResizer, workers uint) (Defacer, error) {
	dp := &defacerPool{
		Inbox:   make(chan *defacerReq, workers),
		Resizer: resizer,
	}
	i := uint(0)
	wg := &sync.WaitGroup{}
	errc := make(chan error, 1)
	defer close(errc)
	for ; i < workers; i++ {
		wg.Add(1)
		go dp.run(wg, errc)
	}
	wg.Wait()
	select {
	case err := <-errc:
		close(dp.Inbox)
		return nil, err
	default:
		return dp, nil
	}
}

func (dp *defacerPool) Deface(r io.Reader) (image.Image, error) {
	req := &defacerReq{
		Reader: r,
		Resp:   make(chan *defacerResp),
	}
	defer close(req.Resp)
	dp.Inbox <- req
	resp := <-req.Resp
	return resp.Image, resp.Error
}

func (dp *defacerPool) run(wg *sync.WaitGroup, errc chan error) {
	runtime.LockOSThread()
	df, err := NewDefacer(dp.Resizer)
	if err != nil {
		select {
		case errc <- err:
		default:
		}
		wg.Done()
		return
	}
	wg.Done()
	for req := range dp.Inbox {
		img, err := df.Deface(req.Reader)
		req.Resp <- &defacerResp{
			Image: img,
			Error: err,
		}
	}
}
