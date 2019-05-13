package main

import (
	"image"
	"image/color"
	"testing"
)

func TestUpdate(t *testing.T) {
	dst := image.NewRGBA(image.Rect(0, 0, 512, 512))

	update(dst)

	w, h, c := 110, 140, color.RGBA{226, 34, 226, 255}

	if got, want := dst.At(w, h), c; got != want {
		t.Fatalf("dst.At(%d, %d) = %v, want %v", w, h, got, want)
	}
}
