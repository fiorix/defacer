//go:generate go-bindata -pkg=internal ./assets

package internal

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/lazywei/go-opencv/opencv"
)

func tempDir() string {
	if v := os.Getenv("TMPDIR"); v != "" {
		return v
	}
	return "/tmp"
}

func restoreAsset(name string) (*os.File, error) {
	b, err := Asset(filepath.Join("assets", name))
	if err != nil {
		return nil, err
	}
	f, err := ioutil.TempFile(tempDir(), "cvdata")
	if err != nil {
		return nil, err
	}
	f.Write(b)
	f.Seek(0, 0)
	return f, nil
}

func imageAsset(name string) (*opencv.IplImage, error) {
	b, err := Asset(filepath.Join("assets", name))
	if err != nil {
		return nil, err
	}
	img := opencv.DecodeImageMem(b)
	if img == nil {
		return nil, errors.New("failed to load image asset")
	}
	return img, nil
}
