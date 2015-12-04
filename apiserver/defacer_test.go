package apiserver

import (
	"bytes"
	"testing"

	"github.com/fiorix/defacer/apiserver/internal"
)

func TestDefacer(t *testing.T) {
	overlay, err := internal.DefaultDefaceImage()
	if err != nil {
		t.Fatal(err)
	}
	df, err := NewDefacer(NewImageResizer(overlay))
	if err != nil {
		t.Fatal(err)
	}
	src, err := internal.DefaultFaceBytes()
	if err != nil {
		t.Fatal(err)
	}
	_, err = df.Deface(bytes.NewBuffer(src))
	if err != nil {
		t.Fatal(err)
	}
}

func TestDefacerPool(t *testing.T) {
	overlay, err := internal.DefaultDefaceImage()
	if err != nil {
		t.Fatal(err)
	}
	df, err := NewDefacerPool(NewImageResizer(overlay), 1)
	if err != nil {
		t.Fatal(err)
	}
	src, err := internal.DefaultFaceBytes()
	if err != nil {
		t.Fatal(err)
	}
	_, err = df.Deface(bytes.NewBuffer(src))
	if err != nil {
		t.Fatal(err)
	}
}
