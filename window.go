package gui

import (
	"image"
	"image/draw"
	"runtime"
	"time"
	"unsafe"

	"github.com/faiface/mainthread"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

// Run calls mainthread.Run
func Run(fn func()) {
	mainthread.Run(fn)
}

func newWindow() *Window {
	out, in := makeEventsChan()

	return &Window{
		out:     out,
		in:      in,
		draw:    make(chan func(draw.Image) image.Rectangle),
		newSize: make(chan image.Rectangle),
		finish:  make(chan struct{}),
	}
}

// Window is an Env that handles an actual graphical window.
type Window struct {
	out  <-chan Event
	in   chan<- Event
	draw chan func(draw.Image) image.Rectangle

	newSize chan image.Rectangle
	finish  chan struct{}

	w     *glfw.Window
	img   *image.RGBA
	ratio int
}

// Open a new window with all the supplied options.
//
// The default title is empty and the default size is 640x480.
func Open(opts ...Option) (*Window, error) {
	o := newOptions(opts...)
	w := newWindow()

	var err error

	mainthread.Call(func() {
		w.w, err = makeGLFWWindow(&o)
	})

	if err != nil {
		return nil, err
	}

	mainthread.Call(func() {
		// HiDPI hack
		width, _ := w.w.GetFramebufferSize()

		w.ratio = width / o.width
		if w.ratio != 1 {
			o.width /= w.ratio
			o.height /= w.ratio
		}

		w.w.Destroy()
		w.w, err = makeGLFWWindow(&o)
	})

	if err != nil {
		return nil, err
	}

	w.img = image.NewRGBA(image.Rect(0, 0, o.width*w.ratio, o.height*w.ratio))

	go func() {
		runtime.LockOSThread()
		w.openGLThread()
	}()

	mainthread.CallNonBlock(w.eventThread)

	return w, nil
}

// Events returns the events channel of the window.
func (w *Window) Events() <-chan Event { return w.out }

// Draw returns the draw channel of the window.
func (w *Window) Draw() chan<- func(draw.Image) image.Rectangle { return w.draw }

// Close closes the draw channel
func (w *Window) Close() {
	close(w.Draw())
}

func (w *Window) eventThread() {
	var moX, moY int

	w.w.SetCursorPosCallback(func(_ *glfw.Window, x, y float64) {
		moX, moY = int(x), int(y)
		w.in <- EventMouseMove{image.Point{moX * w.ratio, moY * w.ratio}}
	})

	w.w.SetMouseButtonCallback(func(_ *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
		pos := image.Point{moX * w.ratio, moY * w.ratio}

		switch {
		case button == glfw.MouseButtonLeft && action == glfw.Press:
			w.in <- EventMouseLeftDown{pos}
		case button == glfw.MouseButtonLeft && action == glfw.Release:
			w.in <- EventMouseLeftUp{pos}
		case button == glfw.MouseButtonMiddle && action == glfw.Press:
			w.in <- EventMouseMiddleDown{pos}
		case button == glfw.MouseButtonMiddle && action == glfw.Release:
			w.in <- EventMouseMiddleUp{pos}
		case button == glfw.MouseButtonRight && action == glfw.Press:
			w.in <- EventMouseRightDown{pos}
		case button == glfw.MouseButtonRight && action == glfw.Release:
			w.in <- EventMouseRightUp{pos}
		}
	})

	w.w.SetScrollCallback(func(_ *glfw.Window, xoff, yoff float64) {
		w.in <- EventMouseScroll{image.Point{int(xoff), int(yoff)}}
	})

	w.w.SetCharCallback(func(_ *glfw.Window, r rune) {
		w.in <- EventKeyboardChar{r}
	})

	w.w.SetKeyCallback(func(_ *glfw.Window, key glfw.Key, _ int, action glfw.Action, _ glfw.ModifierKey) {
		k, ok := keyNames[key]
		if !ok {
			return
		}

		switch action {
		case glfw.Press:
			w.in <- EventKeyboardDown{k}
		case glfw.Release:
			w.in <- EventKeyboardUp{k}
		case glfw.Repeat:
			w.in <- EventKeyboardRepeat{k}
		}
	})

	w.w.SetFramebufferSizeCallback(func(_ *glfw.Window, width, height int) {
		r := image.Rect(0, 0, width, height)
		w.newSize <- r
		w.in <- EventResize{r}
	})

	w.w.SetCloseCallback(func(_ *glfw.Window) {
		w.in <- EventClose{}
	})

	w.in <- EventResize{w.img.Bounds()}

	for {
		select {
		case <-w.finish:
			close(w.in)
			w.w.Destroy()
			return
		default:
			glfw.WaitEventsTimeout(1.0 / 30)
		}
	}
}

func (w *Window) openGLThread() {
	w.w.MakeContextCurrent()
	gl.Init()

	w.openGLFlush(w.img.Bounds())

loop:
	for {
		var totalR image.Rectangle

		select {
		case r := <-w.newSize:
			img := image.NewRGBA(r)
			draw.Draw(img, w.img.Bounds(), w.img, w.img.Bounds().Min, draw.Src)
			w.img = img
			totalR = totalR.Union(r)

		case d, ok := <-w.draw:
			if !ok {
				close(w.finish)
				return
			}
			r := d(w.img)
			totalR = totalR.Union(r)
		}

		for {
			select {
			case <-time.After(time.Second / 960):
				w.openGLFlush(totalR)
				totalR = image.ZR
				continue loop

			case r := <-w.newSize:
				img := image.NewRGBA(r)
				draw.Draw(img, w.img.Bounds(), w.img, w.img.Bounds().Min, draw.Src)
				w.img = img
				totalR = totalR.Union(r)

			case d, ok := <-w.draw:
				if !ok {
					close(w.finish)
					return
				}
				r := d(w.img)
				totalR = totalR.Union(r)
			}
		}
	}
}

func (w *Window) openGLFlush(r image.Rectangle) {
	bounds := w.img.Bounds()

	r = r.Intersect(bounds)

	if r.Empty() {
		return
	}

	tmp := image.NewRGBA(r)

	draw.Draw(tmp, r, w.img, r.Min, draw.Src)

	gl.DrawBuffer(gl.FRONT)
	gl.Viewport(
		int32(bounds.Min.X),
		int32(bounds.Min.Y),
		int32(bounds.Dx()),
		int32(bounds.Dy()),
	)
	gl.RasterPos2d(
		-1+2*float64(r.Min.X)/float64(bounds.Dx()),
		+1-2*float64(r.Min.Y)/float64(bounds.Dy()),
	)
	gl.PixelZoom(1, -1)
	gl.DrawPixels(
		int32(r.Dx()),
		int32(r.Dy()),
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		unsafe.Pointer(&tmp.Pix[0]),
	)
	gl.Flush()
}

func makeGLFWWindow(o *options) (*glfw.Window, error) {
	if err := glfw.Init(); err != nil {
		return nil, err
	}

	glfw.WindowHint(glfw.DoubleBuffer, glfw.True)

	if o.resizable {
		glfw.WindowHint(glfw.Resizable, glfw.True)
	} else {
		glfw.WindowHint(glfw.Resizable, glfw.False)
	}

	if o.decorated {
		glfw.WindowHint(glfw.Decorated, glfw.True)
	} else {
		glfw.WindowHint(glfw.Decorated, glfw.False)
	}

	return glfw.CreateWindow(o.width, o.height, o.title, nil, nil)
}

var keyNames = map[glfw.Key]string{
	glfw.KeyLeft:         "left",
	glfw.KeyRight:        "right",
	glfw.KeyUp:           "up",
	glfw.KeyDown:         "down",
	glfw.KeyEscape:       "escape",
	glfw.KeySpace:        "space",
	glfw.KeyBackspace:    "backspace",
	glfw.KeyDelete:       "delete",
	glfw.KeyEnter:        "enter",
	glfw.KeyTab:          "tab",
	glfw.KeyHome:         "home",
	glfw.KeyEnd:          "end",
	glfw.KeyPageUp:       "pageup",
	glfw.KeyPageDown:     "pagedown",
	glfw.KeyLeftShift:    "shift",
	glfw.KeyRightShift:   "shift",
	glfw.KeyLeftControl:  "ctrl",
	glfw.KeyRightControl: "ctrl",
	glfw.KeyLeftAlt:      "alt",
	glfw.KeyRightAlt:     "alt",
}
