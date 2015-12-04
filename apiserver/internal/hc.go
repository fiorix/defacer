package internal

import (
	"errors"
	"os"

	"github.com/lazywei/go-opencv/opencv"
)

func DefaultHaarCascade() (*opencv.HaarCascade, error) {
	f, err := restoreAsset("haarcascade_frontalface_alt.xml")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	defer os.Remove(f.Name())
	hc := opencv.LoadHaarClassifierCascade(f.Name())
	if hc == nil {
		return nil, errors.New("failed to load haar cascade")
	}
	return hc, nil
}
