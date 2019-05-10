# peterhellberg/gui

[![Build Status](https://travis-ci.org/peterhellberg/gui.svg?branch=master)](https://travis-ci.org/peterhellberg/gui)
[![Go Report Card](https://goreportcard.com/badge/github.com/peterhellberg/gui?style=flat)](https://goreportcard.com/report/github.com/peterhellberg/gui)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/peterhellberg/gui)

Minimal GUI in Go initially based on <https://github.com/faiface/gui>

> NOTE: This is just my take on how to handle events (interface instead of string)
> and you should most likely use <https://github.com/faiface/gui>

## Dependencies

- <https://github.com/faiface/mainthread>
- <https://github.com/go-gl/gl>
- <https://github.com/go-gl/glfw>

## Example

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
		return
	}

	for event := range win.Events() {
		switch event.Name() {
		case gui.EventClose:
			win.Close()
		case gui.EventKeyboardDown:
			if event.Data().(string) == "escape" {
				win.Close()
			}
		case gui.EventKeyboardChar:
			if event.Data().(rune) == 'q' {
				win.Close()
			}
		case gui.EventResize:
			win.Draw() <- update
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

![gui-xor](https://user-images.githubusercontent.com/565124/57329314-d007cc00-7113-11e9-892b-e4c75401004f.png)
