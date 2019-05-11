package stream

import (
	"sync"
	"time"
)

// Event to be windowed.
type Event struct {
	Date      time.Time `json:"date"`
	Type      string    `json:"type"`
	SessionID string    `json:"sessionID"`
}

// EventWindow on a stream.
type EventWindow struct {
	m      sync.Mutex
	maxAge time.Duration
	events []Event
	exit   func(head Event, tail []Event)
}

// NewEventWindow creates a window against events.
func NewEventWindow(maxAge time.Duration, exit func(head Event, tail []Event)) *EventWindow {
	return &EventWindow{
		maxAge: maxAge,
		events: []Event{},
		exit:   exit,
	}
}

// Push an event into the Windower.
func (ew *EventWindow) Push(e Event) {
	defer ew.m.Unlock()
	ew.m.Lock()
	removeBefore := e.Date.Add(-ew.maxAge)
	var i int
	var previous Event
	for i, previous = range ew.events {
		if previous.Date.After(removeBefore) {
			break
		}
	}
	for j := 0; j < i; j++ {
		ew.exit(ew.events[j], ew.events[j:])
	}
	ew.events = append(ew.events[i:], e)
}

// Window returns the window of current events.
func (ew *EventWindow) Window() (events []Event) {
	return ew.events
}
