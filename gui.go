package gui

import (
	"image"
	"image/draw"
	"runtime"
	"sync"
	"time"
	"unsafe"

	"github.com/faiface/mainthread"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

// Event constants
const (
	EventMouseDown      = "mouse/down"
	EventMouseMove      = "mouse/move"
	EventMouseScroll    = "mouse/scroll"
	EventMouseUp        = "mouse/up"
	EventKeyboardChar   = "keyboard/char"
	EventKeyboardDown   = "keyboard/down"
	EventKeyboardUp     = "keyboard/up"
	EventKeyboardRepeat = "keyboard/repeat"
	EventResize         = "resize"
	EventClose          = "close"
)

var buttonNames = map[glfw.MouseButton]string{
	glfw.MouseButtonLeft:   "left",
	glfw.MouseButtonRight:  "right",
	glfw.MouseButtonMiddle: "middle",
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

type options struct {
	title         string
	width, height int
	resizable     bool
}

// Option is a functional option to the window constructor New.
type Option func(*options)

// Title option sets the title (caption) of the window.
func Title(title string) Option {
	return func(o *options) {
		o.title = title
	}
}

// Size option sets the width and height of the window.
func Size(width, height int) Option {
	return func(o *options) {
		o.width = width
		o.height = height
	}
}

// Resizable option makes the window resizable by the user.
func Resizable(b bool) Option {
	return func(o *options) {
		o.resizable = b
	}
}

// New creates a new window with all the supplied options.
//
// The default title is empty and the default size is 640x480.
func New(opts ...Option) (*Window, error) {
	o := options{
		title:     "",
		width:     640,
		height:    480,
		resizable: false,
	}
	for _, opt := range opts {
		opt(&o)
	}

	eventsOut, eventsIn := makeEventsChan()

	w := &Window{
		eventsOut: eventsOut,
		eventsIn:  eventsIn,
		draw:      make(chan func(draw.Image) image.Rectangle),
		newSize:   make(chan image.Rectangle),
		finish:    make(chan struct{}),
	}

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

func makeGLFWWindow(o *options) (*glfw.Window, error) {
	if err := glfw.Init(); err != nil {
		return nil, err
	}

	glfw.WindowHint(glfw.DoubleBuffer, glfw.False)

	if o.resizable {
		glfw.WindowHint(glfw.Resizable, glfw.True)
	} else {
		glfw.WindowHint(glfw.Resizable, glfw.False)
	}

	return glfw.CreateWindow(o.width, o.height, o.title, nil, nil)
}

// Env is an interactive graphical environment, such as a window.
type Env interface {
	Events() <-chan Event
	Draw() chan<- func(draw.Image) image.Rectangle
}

// Window is an Env that handles an actual graphical window.
type Window struct {
	eventsOut <-chan Event
	eventsIn  chan<- Event
	draw      chan func(draw.Image) image.Rectangle

	newSize chan image.Rectangle
	finish  chan struct{}

	w     *glfw.Window
	img   *image.RGBA
	ratio int
}

// Events returns the events channel of the window.
func (w *Window) Events() <-chan Event { return w.eventsOut }

// Draw returns the draw channel of the window.
func (w *Window) Draw() chan<- func(draw.Image) image.Rectangle { return w.draw }

func (w *Window) eventThread() {
	var moX, moY int

	w.w.SetCursorPosCallback(func(_ *glfw.Window, x, y float64) {
		moX, moY = int(x), int(y)
		w.eventsIn <- event{EventMouseMove, []int{moX * w.ratio, moY * w.ratio}}
	})

	w.w.SetMouseButtonCallback(func(_ *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
		b, ok := buttonNames[button]
		if !ok {
			return
		}

		switch action {
		case glfw.Press:
			w.eventsIn <- event{EventMouseDown, []interface{}{moX * w.ratio, moY * w.ratio, b}}
		case glfw.Release:
			w.eventsIn <- event{EventMouseUp, []interface{}{moX * w.ratio, moY * w.ratio, b}}
		}
	})

	w.w.SetScrollCallback(func(_ *glfw.Window, xoff, yoff float64) {
		w.eventsIn <- event{EventMouseScroll, []int{int(xoff), int(yoff)}}
	})

	w.w.SetCharCallback(func(_ *glfw.Window, r rune) {
		w.eventsIn <- event{EventKeyboardChar, r}
	})

	w.w.SetKeyCallback(func(_ *glfw.Window, key glfw.Key, _ int, action glfw.Action, _ glfw.ModifierKey) {
		k, ok := keyNames[key]
		if !ok {
			return
		}

		switch action {
		case glfw.Press:
			w.eventsIn <- event{EventKeyboardDown, k}
		case glfw.Release:
			w.eventsIn <- event{EventKeyboardUp, k}
		case glfw.Repeat:
			w.eventsIn <- event{EventKeyboardRepeat, k}
		}
	})

	w.w.SetFramebufferSizeCallback(func(_ *glfw.Window, width, height int) {
		r := image.Rect(0, 0, width, height)
		w.newSize <- r
		w.eventsIn <- event{EventResize, []int{r.Min.X, r.Min.Y, r.Max.X, r.Max.Y}}
	})

	w.w.SetCloseCallback(func(_ *glfw.Window) {
		w.eventsIn <- event{EventClose, nil}
	})

	w.eventsIn <- event{EventResize, w.img.Bounds()}

	for {
		select {
		case <-w.finish:
			close(w.eventsIn)
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

// Event includes its name and args.
type Event interface {
	Name() string
	Args() interface{}
}

type event struct {
	name string
	args interface{}
}

func (e event) Name() string {
	return e.name
}

func (e event) Args() interface{} {
	return e.args
}

func makeEventsChan() (<-chan Event, chan<- Event) {
	out, in := make(chan Event), make(chan Event)

	go func() {
		var queue []Event

		for {
			x, ok := <-in
			if !ok {
				close(out)
				return
			}
			queue = append(queue, x)

			for len(queue) > 0 {
				select {
				case out <- queue[0]:
					queue = queue[1:]
				case x, ok := <-in:
					if !ok {
						for _, x := range queue {
							out <- x
						}
						close(out)
						return
					}
					queue = append(queue, x)
				}
			}
		}
	}()

	return out, in
}

// Mux can be used to multiplex an Env.
type Mux struct {
	mu         sync.Mutex
	lastResize Event
	eventsIns  []chan<- Event
	draw       chan<- func(draw.Image) image.Rectangle
}

// NewMux creates a new Mux that multiplexes the given Env.
func NewMux(env Env) (mux *Mux, master Env) {
	drawChan := make(chan func(draw.Image) image.Rectangle)
	mux = &Mux{draw: drawChan}
	master = mux.makeEnv(true)

	go func() {
		for d := range drawChan {
			env.Draw() <- d
		}
		close(env.Draw())
	}()

	go func() {
		for e := range env.Events() {
			if e.Name() == EventResize {
				mux.mu.Lock()
				mux.lastResize = e
				mux.mu.Unlock()
			}

			mux.mu.Lock()
			for _, eventsIn := range mux.eventsIns {
				eventsIn <- e
			}
			mux.mu.Unlock()
		}

		mux.mu.Lock()
		for _, eventsIn := range mux.eventsIns {
			close(eventsIn)
		}
		mux.mu.Unlock()
	}()

	return mux, master
}

// MakeEnv creates a new virtual Env that interacts with the root Env of the Mux.
func (mux *Mux) MakeEnv() Env {
	return mux.makeEnv(false)
}

type muxEnv struct {
	events <-chan Event
	draw   chan<- func(draw.Image) image.Rectangle
}

func (m *muxEnv) Events() <-chan Event {
	return m.events
}

func (m *muxEnv) Draw() chan<- func(draw.Image) image.Rectangle {
	return m.draw
}

func (mux *Mux) makeEnv(master bool) Env {
	eventsOut, eventsIn := makeEventsChan()
	drawChan := make(chan func(draw.Image) image.Rectangle)

	env := &muxEnv{eventsOut, drawChan}

	mux.mu.Lock()
	mux.eventsIns = append(mux.eventsIns, eventsIn)
	if mux.lastResize != nil {
		eventsIn <- mux.lastResize
	}
	mux.mu.Unlock()

	go func() {
		func() {
			defer func() {
				if recover() != nil {
					for range drawChan {
					}
				}
			}()

			for d := range drawChan {
				mux.draw <- d // !
			}
		}()

		if master {
			mux.mu.Lock()
			for _, eventsIn := range mux.eventsIns {
				close(eventsIn)
			}
			mux.eventsIns = nil
			close(mux.draw)
			mux.mu.Unlock()
		} else {
			mux.mu.Lock()
			i := -1
			for i = range mux.eventsIns {
				if mux.eventsIns[i] == eventsIn {
					break
				}
			}
			if i != -1 {
				mux.eventsIns = append(mux.eventsIns[:i], mux.eventsIns[i+1:]...)
			}
			mux.mu.Unlock()
		}
	}()

	return env
}

// Run calls mainthread.Run
func Run(run func()) {
	mainthread.Run(run)
}
