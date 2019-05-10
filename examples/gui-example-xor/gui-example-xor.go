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
