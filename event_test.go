package gui

import "testing"

func TestMakeEventsChan(t *testing.T) {
	name := "close"

	out, in := makeEventsChan()

	in <- EventClose{}

	e := <-out

	if got, want := e.Name(), name; got != want {
		t.Fatalf("e.Name() = %q, want %q", got, want)
	}
}
