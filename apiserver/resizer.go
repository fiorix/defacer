package apiserver

import (
	"image"
	"time"

	"github.com/nfnt/resize"
)

// ImageResizer is an object that can resize images to a given size.
type ImageResizer interface {
	Resize(size image.Point) image.Image
}

// NewImageResizer stores the given image and returns an ImageResizer
// that can return different sizes of the stored image.
func NewImageResizer(m image.Image) ImageResizer {
	ir := &imageResizer{
		Image: m,
		Inbox: make(chan *imageResizerReq, 1000),
	}
	go ir.coalesce()
	return ir
}

type imageResizerReq struct {
	Size image.Point
	Resp chan image.Image
}

type imageResizer struct {
	Image image.Image
	Inbox chan *imageResizerReq
}

func (ir *imageResizer) Resize(size image.Point) image.Image {
	req := &imageResizerReq{
		Size: size,
		Resp: make(chan image.Image),
	}
	defer close(req.Resp)
	ir.Inbox <- req
	return <-req.Resp
}

type imageResizerBatch map[image.Point][]chan image.Image

func (ir *imageResizer) coalesce() {
	backoff := 10 * time.Millisecond
	batch := make(imageResizerBatch)
	for {
		select {
		case req := <-ir.Inbox:
			batch[req.Size] = append(batch[req.Size], req.Resp)
			backoff = 10 * time.Millisecond
		case <-time.After(backoff):
			if len(batch) == 0 {
				backoff *= backoff
				break
			}
			ir.dispatch(batch)
			batch = make(imageResizerBatch)
		}
	}
}

func (ir *imageResizer) dispatch(batch imageResizerBatch) {
	for size, resp := range batch {
		go ir.resize(size, resp)
	}
	batch = nil
}

func (ir *imageResizer) resize(size image.Point, callers []chan image.Image) {
	img := defaultImageCache.Get(&size)
	if img == nil {
		img = resize.Resize(
			uint(size.X),
			uint(size.Y),
			ir.Image,
			resize.Bicubic,
		)
		defaultImageCache.Set(&size, img)
	}
	for n, resp := range callers {
		resp <- img
		if n > 0 {
			defacerImageResizeCoalesceSum.Inc()
		}
	}
}
