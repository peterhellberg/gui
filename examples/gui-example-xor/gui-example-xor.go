package main

import (
	"image"
	"image/draw"

	"github.com/peterhellberg/gfx"
	"github.com/peterhellberg/gui"
)

func main() {
	gui.Run(loop)
}

func loop() {
	win, err := gui.New(
		gui.Title("qui-xor"),
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

		gfx.Log("Event: %+v", event)
	}
}

func update(dst draw.Image) image.Rectangle {
	gfx.EachPixel(dst.Bounds(), func(x, y int) {
		c := uint8(x ^ y)

		dst.Set(x, y, gfx.ColorNRGBA(c, c%192, c, 255))
	})

	return dst.Bounds()
}
