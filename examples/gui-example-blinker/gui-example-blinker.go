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
	go blinker(mux.MakeEnv(), image.Rect(100, 100, 350, 250))
	go blinker(mux.MakeEnv(), image.Rect(450, 100, 700, 250))
	go blinker(mux.MakeEnv(), image.Rect(100, 350, 350, 500))
	go blinker(mux.MakeEnv(), image.Rect(450, 350, 700, 500))

	// we use the master env now, win is used by the mux
	for event := range env.Events() {
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
	env.Draw() <- redraw(true)

	for event := range env.Events() {
		switch event.Name() {
		case gui.EventMouseLeftDown:
			if event.Data().(image.Point).In(r) {
				// user clicked on the rectangle we blink 3 times
				for i := 0; i < 3; i++ {
					env.Draw() <- redraw(false)
					time.Sleep(time.Second / 3)

					env.Draw() <- redraw(true)
					time.Sleep(time.Second / 3)
				}
			}
		}
	}

	close(env.Draw())
}
