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
