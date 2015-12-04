package apiserver

import (
	"image"
	"testing"

	"github.com/fiorix/defacer/apiserver/internal"
)

func TestImageCache(t *testing.T) {
	overlay, err := internal.DefaultDefaceImage()
	if err != nil {
		t.Fatal(err)
	}
	size := image.Point{10, 20}
	m := defaultImageCache.Get(&size)
	if m != nil {
		t.Fatal("unexpected image from cache")
	}
	defaultImageCache.Set(&size, overlay)
	m = defaultImageCache.Get(&size)
	if m != overlay {
		t.Fatal("image missing from cache")
	}
}
