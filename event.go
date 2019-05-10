package gui

// Event constants
const (
	EventClose           = "close"
	EventKeyboardChar    = "keyboard/char"
	EventKeyboardDown    = "keyboard/down"
	EventKeyboardRepeat  = "keyboard/repeat"
	EventKeyboardUp      = "keyboard/up"
	EventMouseLeftDown   = "mouse/left/down"
	EventMouseLeftUp     = "mouse/left/up"
	EventMouseMiddleDown = "mouse/middle/down"
	EventMouseMiddleUp   = "mouse/middle/up"
	EventMouseMove       = "mouse/move"
	EventMouseRightDown  = "mouse/right/down"
	EventMouseRightUp    = "mouse/right/up"
	EventMouseScroll     = "mouse/scroll"
	EventResize          = "resize"
)

// Event includes its name and data.
type Event interface {
	Name() string
	Data() interface{}
}

type event struct {
	name string
	data interface{}
}

func (e event) Name() string {
	return e.name
}

func (e event) Data() interface{} {
	return e.data
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
