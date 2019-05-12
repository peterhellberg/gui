package gui

import (
	"image"
	"image/draw"
	"testing"
)

func TestWindowIsEnv(t *testing.T) {
	var _ Env = &Window{}
}

type mockEnv struct {
	EventsFn func() <-chan Event
	DrawFn   func(func(draw.Image) image.Rectangle)
}

func (env *mockEnv) Events() <-chan Event {
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
