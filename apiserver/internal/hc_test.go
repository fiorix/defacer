package internal

import "testing"

func TestDefaultHaarCascade(t *testing.T) {
	_, err := DefaultHaarCascade()
	if err != nil {
		t.Fatal(err)
	}
}
