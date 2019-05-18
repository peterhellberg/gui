# peterhellberg/gui

[![Build Status](https://travis-ci.org/peterhellberg/gui.svg?branch=master)](https://travis-ci.org/peterhellberg/gui)
[![Go Report Card](https://goreportcard.com/badge/github.com/peterhellberg/gui?style=flat)](https://goreportcard.com/report/github.com/peterhellberg/gui)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/peterhellberg/gui)

Minimal GUI in Go, it was initially based on <https://github.com/faiface/gui> but has since diverged from the original project.

## Dependencies

This package has a few third party dependencies:

- <https://github.com/faiface/mainthread> - Run stuff on the main thread in Go
- <https://github.com/go-gl/gl> - Go bindings for OpenGL (generated via glow)
- <https://github.com/go-gl/glfw> - Go bindings for GLFW 3

## Examples

### Minimal

![gui-minimal](https://user-images.githubusercontent.com/565124/57968481-83ba5880-796b-11e9-8339-934a8d7a542c.png)

[embedmd]:# (examples/gui-example-minimal/gui-example-minimal.go)
```go
package main

import (
	"image"
	"image/draw"

	"github.com/peterhellberg/gui"
)

func main() {
	gui.Run(func() {
		win, err := gui.Open(gui.Title("gui-minimal"))
		if err != nil {
			panic(err)
		}

		for event := range win.Events() {
			switch event.(type) {
			case gui.EventClose:
				win.Close()
			case gui.EventResize:
				win.Draw(func(dst draw.Image) image.Rectangle {
					return dst.Bounds()
				})
			}
		}
	})
}
```

### XOR

![gui-xor](https://user-images.githubusercontent.com/565124/57329314-d007cc00-7113-11e9-892b-e4c75401004f.png)

[embedmd]:# (examples/gui-example-xor/gui-example-xor.go)
```go
package main

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/peterhellberg/gui"
)

func main() {
	gui.Run(loop)
}

func loop() {
	win, err := gui.Open(
		gui.Title("gui-xor"),
		gui.Size(512, 512),
		gui.Decorated(true),
		gui.Resizable(true),
	)
	if err != nil {
		panic(err)
	}

	for event := range win.Events() {
		switch event := event.(type) {
		case gui.EventClose:
			win.Close()
		case gui.EventKeyboardDown:
			if event.Key == "escape" {
				win.Close()
			}
		case gui.EventKeyboardChar:
			if event.Char == 'q' {
				win.Close()
			}
		case gui.EventResize:
			win.Draw(update)
		}

		gui.Log("Event: %+v", event)
	}
}

func update(dst draw.Image) image.Rectangle {
	bounds := dst.Bounds()

	for x := 0; x < bounds.Max.X; x++ {
		for y := 0; y < bounds.Max.Y; y++ {
			c := uint8(x ^ y)

			dst.Set(x, y, color.NRGBA{c, c % 192, c, 255})
		}
	}

	return bounds
}
```
## Blinker

![gui-blinker](https://user-images.githubusercontent.com/565124/57541634-c10d5d80-734f-11e9-8774-14c71ea920f1.png)

[embedmd]:# (examples/gui-example-blinker/gui-example-blinker.go)
```go
package main

import (
	"image"
	"image/draw"
	"time"

	"github.com/peterhellberg/gui"
)

func main() {
	gui.Run(loop)
}

func loop() {
	win, err := gui.Open(
		gui.Title("gui-blinker"),
		gui.Size(800, 600),
	)
	if err != nil {
		panic(err)
	}

	mux, env := gui.NewMux(win)

	// we create four blinkers, each with its own Env from the mux
	go blinker(mux.Env(), image.Rect(100, 100, 350, 250))
	go blinker(mux.Env(), image.Rect(450, 100, 700, 250))
	go blinker(mux.Env(), image.Rect(100, 350, 350, 500))
	go blinker(mux.Env(), image.Rect(450, 350, 700, 500))

	// we use the master env now, win is used by the mux
	for event := range env.Events() {
		switch event := event.(type) {
		case gui.EventClose:
			win.Close()
		case gui.EventKeyboardDown:
			if event.Key == "escape" {
				win.Close()
			}
		case gui.EventKeyboardChar:
			if event.Char == 'q' {
				win.Close()
			}
		}
	}
}

func blinker(env gui.Env, r image.Rectangle) {
	// redraw takes a bool and produces a draw command
	redraw := func(visible bool) func(draw.Image) image.Rectangle {
		return func(dst draw.Image) image.Rectangle {
			if visible {
				draw.Draw(dst, r, image.White, image.ZP, draw.Src)
			} else {
				draw.Draw(dst, r, image.Black, image.ZP, draw.Src)
			}

			return r
		}
	}

	// first we draw a white rectangle
	env.Draw(redraw(true))

	for event := range env.Events() {
		switch event := event.(type) {
		case gui.EventMouseLeftDown:
			if event.In(r) {
				// user clicked on the rectangle we blink 3 times
				for i := 0; i < 3; i++ {
					env.Draw(redraw(false))
					time.Sleep(time.Second / 3)

					env.Draw(redraw(true))
					time.Sleep(time.Second / 3)
				}
			}
		}
	}

	env.Close()
}
```
