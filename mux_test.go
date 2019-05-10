package gui

import (
	"image"
	"image/color"
	"image/draw"
	"testing"
)

func TestNewMux(t *testing.T) {
	NewMux(&mockEnv{
		EventsFn: func() <-chan Event {
			ch := make(chan Event, 1)

			ch <- event{name: "resize"}

			return ch
		},
	})
}

func TestEnv(t *testing.T) {
	mux := &Mux{}
	mux.Env()
}

func TestMuxEnvEvents(t *testing.T) {
	me := &muxEnv{}
	me.Events()
}

func TestMuxEnvDraw(t *testing.T) {
	me := &muxEnv{
		draw: make(chan func(draw.Image) image.Rectangle, 1),
	}

	me.Draw() <- func(dst draw.Image) image.Rectangle {
		dst.Set(0, 0, color.RGBA{255, 0, 0, 255})

		return dst.Bounds()
	}
}
