package gui

import "image"

// Event includes its name and data.
type Event interface {
	Name() string
	Data() interface{}
}

// EventResize event
type EventResize struct {
	image.Rectangle
}

// Name of event
func (r EventResize) Name() string {
	return "resize"
}

// Data for event
func (r EventResize) Data() interface{} {
	return r.Rectangle
}

// Close event
type EventClose struct{}

// Name of event
func (c EventClose) Name() string {
	return "close"
}

// Data for event
func (c EventClose) Data() interface{} {
	return nil
}

// EventMouseMove event
type EventMouseMove struct {
	image.Point
}

// Name of event
func (mm EventMouseMove) Name() string {
	return "mouse"
}

// Data for event
func (mm EventMouseMove) Data() interface{} {
	return mm.Point
}

// EventMouseScroll event
type EventMouseScroll struct {
	image.Point
}

// Name of event
func (ms EventMouseScroll) Name() string {
	return "mouse/scroll"
}

// Data for event
func (ms EventMouseScroll) Data() interface{} {
	return ms.Point
}

// EventMouseLeftDown event
type EventMouseLeftDown struct {
	image.Point
}

// Name of event
func (mld EventMouseLeftDown) Name() string {
	return "mouse/left/down"
}

// Data for event
func (mld EventMouseLeftDown) Data() interface{} {
	return mld.Point
}

// EventMouseLeftUp event
type EventMouseLeftUp struct {
	image.Point
}

// Name of event
func (mlu EventMouseLeftUp) Name() string {
	return "mouse/left/up"
}

// Data for event
func (mlu EventMouseLeftUp) Data() interface{} {
	return mlu.Point
}

// EventMouseMiddleDown event
type EventMouseMiddleDown struct {
	image.Point
}

// Name of event
func (mmd EventMouseMiddleDown) Name() string {
	return "mouse/middle/down"
}

// Data for event
func (mmd EventMouseMiddleDown) Data() interface{} {
	return mmd.Point
}

// EventMouseMiddleUp event
type EventMouseMiddleUp struct {
	image.Point
}

// Name of event
func (mmu EventMouseMiddleUp) Name() string {
	return "mouse/middle/up"
}

// Data for event
func (mmu EventMouseMiddleUp) Data() interface{} {
	return mmu.Point
}

// EventMouseRightDown event
type EventMouseRightDown struct {
	image.Point
}

// Name of event
func (mrd EventMouseRightDown) Name() string {
	return "mouse/right/down"
}

// Data for event
func (mrd EventMouseRightDown) Data() interface{} {
	return mrd.Point
}

// EventMouseRightUp event
type EventMouseRightUp struct {
	image.Point
}

// Name of event
func (mru EventMouseRightUp) Name() string {
	return "mouse/right/up"
}

// Data for event
func (mru EventMouseRightUp) Data() interface{} {
	return mru.Point
}

// EventKeyboardChar event
type EventKeyboardChar struct {
	Char rune
}

// Name of event
func (kc EventKeyboardChar) Name() string {
	return "keyboard/char"
}

// Data for event
func (kc EventKeyboardChar) Data() interface{} {
	return kc.Char
}

// EventKeyboardDown event
type EventKeyboardDown struct {
	Key string
}

// Name of event
func (kd EventKeyboardDown) Name() string {
	return "keyboard/down"
}

// Data for event
func (kd EventKeyboardDown) Data() interface{} {
	return kd.Key
}

// EventKeyboardUp event
type EventKeyboardUp struct {
	Key string
}

// Name of event
func (ku EventKeyboardUp) Name() string {
	return "keyboard/up"
}

// Data for event
func (ku EventKeyboardUp) Data() interface{} {
	return ku.Key
}

// EventKeyboardRepeat event
type EventKeyboardRepeat struct {
	Key string
}

// Name of event
func (kr EventKeyboardRepeat) Name() string {
	return "keyboard/repeat"
}

// Data for event
func (kr EventKeyboardRepeat) Data() interface{} {
	return kr.Key
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
