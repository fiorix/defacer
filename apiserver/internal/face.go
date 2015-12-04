package internal

import (
	"bytes"
	"image"
	"image/png"
	"path/filepath"

	"github.com/lazywei/go-opencv/opencv"
)

func DefaultFace() (*opencv.IplImage, error) {
	return imageAsset("face.jpg")
}

func DefaultFaceBytes() ([]byte, error) {
	return Asset(filepath.Join("assets", "face.jpg"))
}

func DefaultDeface() (*opencv.IplImage, error) {
	return imageAsset("deface.png")
}

func DefaultDefaceBytes() ([]byte, error) {
	return Asset(filepath.Join("assets", "deface.png"))
}

func DefaultDefaceImage() (image.Image, error) {
	b, err := DefaultDefaceBytes()
	if err != nil {
		return nil, err
	}
	return png.Decode(bytes.NewBuffer(b))
}
