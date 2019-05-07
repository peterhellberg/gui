# peterhellberg/gui

Minimal GUI in Go based on <https://github.com/faiface/gui>

> NOTE: This is just my take on how to handle events (interface instead of string)
> and you should most likely use <https://github.com/faiface/gui>

## Dependencies

- <https://github.com/faiface/mainthread>
- <https://github.com/go-gl/gl>
- <https://github.com/go-gl/glfw>

## Example

```go
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
	win, err := gui.New(gui.Title("gui-xor"), gui.Size(512, 512))
	if err != nil {
		return
	}

	win.Draw() <- func(dst draw.Image) image.Rectangle {
		r := dst.Bounds()

		gfx.EachPixel(r, func(x, y int) {
			c := uint8(x ^ y)

			dst.Set(x, y, gfx.ColorNRGBA(c, c%192, c, 255))
		})

		return r
	}

	for event := range win.Events() {
		switch event.Name() {
		case gui.EventClose, gui.EventKeyboardUp:
			close(win.Draw())
		default:
			gfx.Log("Event: %+v", event)
		}
	}
}
```

![gui-xor](https://user-images.githubusercontent.com/565124/57329314-d007cc00-7113-11e9-892b-e4c75401004f.png)
