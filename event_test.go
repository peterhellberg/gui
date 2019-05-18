package gui

import (
	"image"
	"testing"
)

func TestMakeEventsChan(t *testing.T) {
	name := "close"

	out, in := makeEventsChan()

	in <- EventClose{}

	e := <-out

	if got, want := e.Name(), name; got != want {
		t.Fatalf("e.Name() = %q, want %q", got, want)
	}
}

func TestEventNames(t *testing.T) {
	for _, tt := range []struct {
		event Event
		name  string
	}{
		{EventResize{}, "resize"},
		{EventClose{}, "close"},
		{EventMouseMove{}, "mouse/move"},
		{EventMouseScroll{}, "mouse/scroll"},
		{EventMouseLeftDown{}, "mouse/left/down"},
		{EventMouseLeftUp{}, "mouse/left/up"},
		{EventMouseMiddleDown{}, "mouse/middle/down"},
		{EventMouseMiddleUp{}, "mouse/middle/up"},
		{EventMouseRightDown{}, "mouse/right/down"},
		{EventMouseRightUp{}, "mouse/right/up"},
		{EventKeyboardChar{}, "keyboard/char"},
		{EventKeyboardDown{}, "keyboard/down"},
		{EventKeyboardUp{}, "keyboard/up"},
		{EventKeyboardRepeat{}, "keyboard/repeat"},
	} {
		if got, want := tt.event.Name(), tt.name; got != want {
			t.Fatalf("tt.event.Name() = %q, want %q", got, want)
		}
	}
}

func TestEventData(t *testing.T) {
	t.Run("EventResize", func(t *testing.T) {
		e := EventResize{image.Rect(0, 1, 2, 3)}

		if got, want := e.Data().(image.Rectangle), e.Rectangle; got != want {
			t.Fatalf("e.Data().(image.Rectangle) = %v, want %v", got, want)
		}
	})

	t.Run("EventClose", func(t *testing.T) {
		e := EventClose{}

		if got := e.Data(); got != nil {
			t.Fatalf("e.Data() = %v", got)
		}
	})

	t.Run("EventMouseMove", func(t *testing.T) {
		e := EventMouseMove{image.Pt(1, 2)}

		if got, want := e.Data().(image.Point), e.Point; got != want {
			t.Fatalf("e.Data().(image.Point) = %v, want %v", got, want)
		}
	})

	t.Run("EventMouseScroll", func(t *testing.T) {
		e := EventMouseScroll{image.Pt(1, 2)}

		if got, want := e.Data().(image.Point), e.Point; got != want {
			t.Fatalf("e.Data().(image.Point) = %v, want %v", got, want)
		}
	})

	t.Run("EventMouseLeftDown", func(t *testing.T) {
		e := EventMouseLeftDown{image.Pt(1, 2)}

		if got, want := e.Data().(image.Point), e.Point; got != want {
			t.Fatalf("e.Data().(image.Point) = %v, want %v", got, want)
		}
	})

	t.Run("EventMouseLeftUp", func(t *testing.T) {
		e := EventMouseLeftUp{image.Pt(1, 2)}

		if got, want := e.Data().(image.Point), e.Point; got != want {
			t.Fatalf("e.Data().(image.Point) = %v, want %v", got, want)
		}
	})

	t.Run("EventMouseMiddleDown", func(t *testing.T) {
		e := EventMouseMiddleDown{image.Pt(1, 2)}

		if got, want := e.Data().(image.Point), e.Point; got != want {
			t.Fatalf("e.Data().(image.Point) = %v, want %v", got, want)
		}
	})

	t.Run("EventMouseMiddleUp", func(t *testing.T) {
		e := EventMouseMiddleUp{image.Pt(1, 2)}

		if got, want := e.Data().(image.Point), e.Point; got != want {
			t.Fatalf("e.Data().(image.Point) = %v, want %v", got, want)
		}
	})

	t.Run("EventMouseRightDown", func(t *testing.T) {
		e := EventMouseRightDown{image.Pt(1, 2)}

		if got, want := e.Data().(image.Point), e.Point; got != want {
			t.Fatalf("e.Data().(image.Point) = %v, want %v", got, want)
		}
	})

	t.Run("EventMouseRightUp", func(t *testing.T) {
		e := EventMouseRightUp{image.Pt(1, 2)}

		if got, want := e.Data().(image.Point), e.Point; got != want {
			t.Fatalf("e.Data().(image.Point) = %v, want %v", got, want)
		}
	})

	t.Run("EventKeyboardChar", func(t *testing.T) {
		e := EventKeyboardChar{'x'}

		if got, want := e.Data().(rune), 'x'; got != want {
			t.Fatalf("e.Data().(rune) = %v, want %v", got, want)
		}
	})

	t.Run("EventKeyboardDown", func(t *testing.T) {
		e := EventKeyboardDown{"escape"}

		if got, want := e.Data().(string), "escape"; got != want {
			t.Fatalf("e.Data().(string) = %v, want %v", got, want)
		}
	})

	t.Run("EventKeyboardUp", func(t *testing.T) {
		e := EventKeyboardUp{"escape"}

		if got, want := e.Data().(string), "escape"; got != want {
			t.Fatalf("e.Data().(string) = %v, want %v", got, want)
		}
	})

	t.Run("EventKeyboardRepeat", func(t *testing.T) {
		e := EventKeyboardRepeat{"escape"}

		if got, want := e.Data().(string), "escape"; got != want {
			t.Fatalf("e.Data().(string) = %v, want %v", got, want)
		}
	})
}
