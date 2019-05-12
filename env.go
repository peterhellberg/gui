package gui

// Env is an interactive graphical environment, such as a window.
type Env interface {
	Events() <-chan Event
	Draw(DrawFunc)
	Close()
}
