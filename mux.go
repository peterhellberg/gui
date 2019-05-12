package gui

import (
	"image"
	"image/draw"
	"sync"
)

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
			env.Draw(d)
		}

		env.Close()
	}()

	go func() {
		for e := range env.Events() {
			if e.Name() == "resize" {
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

// Env creates a new virtual Env that interacts with the root Env of the Mux.
func (mux *Mux) Env() Env {
	return mux.makeEnv(false)
}

type muxEnv struct {
	events <-chan Event
	draw   chan<- func(draw.Image) image.Rectangle
}

func (m *muxEnv) Events() <-chan Event {
	return m.events
}

func (m *muxEnv) Draw(fn func(draw.Image) image.Rectangle) {
	m.draw <- fn
}

func (m *muxEnv) Close() {
	close(m.draw)
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
