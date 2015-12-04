package internal

import "testing"

func TestDefaultFace(t *testing.T) {
	_, err := DefaultFace()
	if err != nil {
		t.Fatal(err)
	}
}

func TestDefaultFaceBytes(t *testing.T) {
	_, err := DefaultFaceBytes()
	if err != nil {
		t.Fatal(err)
	}
}

func TestDefaultDeface(t *testing.T) {
	_, err := DefaultDeface()
	if err != nil {
		t.Fatal(err)
	}
}

func TestDefaultDefaceBytes(t *testing.T) {
	_, err := DefaultDefaceBytes()
	if err != nil {
		t.Fatal(err)
	}
}

func TestDefaultDefaceImage(t *testing.T) {
	_, err := DefaultDefaceImage()
	if err != nil {
		t.Fatal(err)
	}
}
