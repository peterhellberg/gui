package gui

// Option is a functional option to the window constructor New.
type Option func(*options)

type options struct {
	title         string
	width, height int
	resizable     bool
	decorated     bool
}

func newOptions(opts ...Option) options {
	o := options{
		title:     "",
		width:     640,
		height:    480,
		resizable: false,
		decorated: true,
	}

	for _, opt := range opts {
		opt(&o)
	}

	return o
}

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
func Resizable(resizable bool) Option {
	return func(o *options) {
		o.resizable = resizable
	}
}

// Decorated options controls if the window should have any chrome.
func Decorated(decorated bool) Option {
	return func(o *options) {
		o.decorated = decorated
	}
}
