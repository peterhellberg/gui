package gui

import "testing"

func TestEvent(t *testing.T) {
	t.Run("Name", func(t *testing.T) {
		name := "test-event"

		e := event{name: name}

		if got, want := e.Name(), name; got != want {
			t.Fatalf("e.Name() = %q, want %q", got, want)
		}
	})

	t.Run("Data", func(t *testing.T) {
		data := 123

		e := event{data: data}

		if got, want := e.Data().(int), data; got != want {
			t.Fatalf("e.Data().(int) = %d, want %d", got, want)
		}
	})
}

func TestMakeEventsChan(t *testing.T) {
	out, in := makeEventsChan()

	name := "test-event"

	in <- event{name: name}

	e := <-out

	if got, want := e.Name(), name; got != want {
		t.Fatalf("e.Name() = %q, want %q", got, want)
	}
}
