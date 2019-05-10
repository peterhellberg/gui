package gui

import "testing"

func TestNewWindow(t *testing.T) {
	if got := newWindow(); got == nil {
		t.Fatalf("expected *Window, got nil")
	}
}
