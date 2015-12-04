package apiserver

import (
	"errors"
	"image"
	"image/draw"
	"io"
	"io/ioutil"
	"runtime"

	"github.com/lazywei/go-opencv/opencv"
	"github.com/nfnt/resize"

	"github.com/fiorix/defacer/apiserver/internal"
)

// A Defacer can scan and deface people's faces in images.
type Defacer interface {
	// Deface reads binary image bytes from a given reader
	// and returns a defaced image.
	Deface(io.Reader) (image.Image, error)
}

// NewDefacer creates and initializes a new Defacer.
func NewDefacer(overlay image.Image) (Defacer, error) {
	var err error
	if overlay == nil {
		overlay, err = internal.DefaultDefaceImage()
		if err != nil {
			return nil, err
		}
	}
	hc, err := internal.DefaultHaarCascade()
	if err != nil {
		return nil, err
	}
	df := &defacer{
		Overlay:     overlay,
		HaarCascade: hc,
	}
	return df, nil
}

type defacer struct {
	Overlay     image.Image
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
	draw.Draw(dst, dst.Bounds(), img, b.Min, draw.Src)
	for _, rect := range faces {
		rw := uint(rect.Max.X - rect.Min.X)
		rh := uint(rect.Max.Y - rect.Min.Y)
		img = resize.Resize(rw, rh, df.Overlay, resize.Bicubic)
		b = img.Bounds()
		draw.DrawMask(dst, rect, img, b.Min, img, b.Min, draw.Over)
	}
	return dst, nil
}

// scan reads binary image data from the given reader and scans for
// faces, returning a slice of rectangles where faces were detected.
func (df *defacer) scan(src io.Reader) (m image.Image, r []image.Rectangle, err error) {
	b, err := ioutil.ReadAll(src)
	if err != nil {
		return nil, nil, err
	}
	img := opencv.DecodeImageMem(b)
	if img == nil {
		return nil, nil, errors.New("failed to load source image")
	}
	faces := df.HaarCascade.DetectObjects(img)
	if faces == nil {
		return img.ToImage(), []image.Rectangle{}, nil
	}
	fr := make([]image.Rectangle, len(faces))
	for i, rect := range faces {
		fr[i] = image.Rectangle{
			image.Point{rect.X(), rect.Y()},
			image.Point{rect.X() + rect.Width(), rect.Y() + rect.Height()},
		}
	}
	return img.ToImage(), fr, nil
}

type defacerPool struct {
	Inbox chan *defacerReq
}

type defacerReq struct {
	Reader io.Reader
	Resp   chan *defacerResp
}

type defacerResp struct {
	Image image.Image
	Error error
}

// NewDefacerPool creates a pool of Defacers, and
// implements the Defacer interface.
func NewDefacerPool(overlay image.Image, nworkers uint) (Defacer, error) {
	dp := &defacerPool{
		Inbox: make(chan *defacerReq, nworkers),
	}
	i := uint(0)
	for ; i < nworkers; i++ {
		df, err := NewDefacer(overlay)
		if err != nil {
			close(dp.Inbox)
			return nil, err
		}
		go dp.run(df)
	}
	return dp, nil
}

func (dp *defacerPool) Deface(r io.Reader) (image.Image, error) {
	req := &defacerReq{
		Reader: r,
		Resp:   make(chan *defacerResp),
	}
	dp.Inbox <- req
	resp := <-req.Resp
	return resp.Image, resp.Error
}

func (dp *defacerPool) run(df Defacer) {
	runtime.LockOSThread()
	for req := range dp.Inbox {
		img, err := df.Deface(req.Reader)
		req.Resp <- &defacerResp{
			Image: img,
			Error: err,
		}
		defaceCounter.Inc()
	}
}
