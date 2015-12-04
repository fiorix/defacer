package apiserver

import (
	"bytes"
	"testing"

	"github.com/fiorix/doger/apiserver/internal"
)

func TestDeface(t *testing.T) {
	df, err := NewDefacer(nil)
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
