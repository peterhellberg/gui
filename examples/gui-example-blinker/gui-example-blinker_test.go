package main

import (
	"image"
	"image/draw"
	"testing"

	"github.com/peterhellberg/gui"
)

func TestBlinker(t *testing.T) {
	env := &mockEnv{
		EventsFn: func() <-chan gui.Event {
			ch := make(chan gui.Event)

			close(ch)

			return ch
		},
		DrawFn: func(func(draw.Image) image.Rectangle) {},
	}

	blinker(env, image.ZR)
}

type mockEnv struct {
	EventsFn func() <-chan gui.Event
	DrawFn   func(func(draw.Image) image.Rectangle)
}

func (env *mockEnv) Events() <-chan gui.Event {
	if env.EventsFn == nil {
		panic("*mockEnv.Events called, but it is not mocked")
	}

	return env.EventsFn()
}

func (env *mockEnv) Draw(fn func(draw.Image) image.Rectangle) {
	if env.DrawFn == nil {
		panic("*mockEnv.Draw called, but it is not mocked")
	}

	env.DrawFn(fn)
}

func (env *mockEnv) Close() {}
