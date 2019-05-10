package gui

import "testing"

func TestNewOptions(t *testing.T) {
	want := options{
		title:     "test-title",
		width:     100,
		height:    200,
		resizable: true,
		decorated: true,
	}

	got := newOptions(
		Title(want.title),
		Size(want.width, want.height),
		Resizable(want.resizable),
		Decorated(want.decorated),
	)

	if got != want {
		t.Fatalf("got = %v, want %v", got, want)
	}
}

func TestOptions(t *testing.T) {
	o := &options{}

	t.Run("Title", func(t *testing.T) {
		title := "test-title"

		Title(title)(o)

		if got, want := o.title, title; got != want {
			t.Fatalf("o.title = %q, want %q", got, want)
		}
	})

	t.Run("Size", func(t *testing.T) {
		width, height := 100, 200

		Size(width, height)(o)

		if got, want := o.width, width; got != want {
			t.Fatalf("o.width = %d, want %d", got, want)
		}

		if got, want := o.height, height; got != want {
			t.Fatalf("o.height = %d, want %d", got, want)
		}
	})

	t.Run("Resizable", func(t *testing.T) {
		resizable := true

		Resizable(resizable)(o)

		if got, want := o.resizable, resizable; got != want {
			t.Fatalf("o.resizable = %v, want %v", got, want)
		}
	})

	t.Run("Decorated", func(t *testing.T) {
		decorated := true

		Decorated(decorated)(o)

		if got, want := o.decorated, decorated; got != want {
			t.Fatalf("o.decorated = %v, want %v", got, want)
		}
	})
}
