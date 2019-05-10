package gui

import (
	"image"
	"image/draw"
)

// Env is an interactive graphical environment, such as a window.
type Env interface {
	Events() <-chan Event
	Draw() chan<- func(draw.Image) image.Rectangle
}
