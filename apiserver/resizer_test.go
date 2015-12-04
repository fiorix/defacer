package apiserver

import (
	"image"
	"testing"

	"github.com/fiorix/defacer/apiserver/internal"
)

func TestImageResizer(t *testing.T) {
	overlay, err := internal.DefaultDefaceImage()
	if err != nil {
		t.Fatal(err)
	}
	ir := NewImageResizer(overlay)
	im := ir.Resize(image.Point{101, 102})
	size := im.Bounds().Max
	if size.X != 101 || size.Y != 102 {
		t.Fatal("unexpected size:", size)
	}
}
