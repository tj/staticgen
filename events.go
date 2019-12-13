package staticgen

import (
	"time"
)

// Event is an event.
type Event interface {
	event()
}

// EventStartingServer .
type EventStartingServer struct {
	Command string
	URL     string
}

// EventStartedServer .
type EventStartedServer struct {
	Command string
	URL     string
}

// EventStoppingServer .
type EventStoppingServer struct{}

// EventStartCrawl .
type EventStartCrawl struct{}

// EventStopCrawl .
type EventStopCrawl struct{}

// EventVisitedResource .
type EventVisitedResource struct {
	Target
	Duration   time.Duration
	StatusCode int
	Filename   string
	Error      error
}

// event implementation.
func (e EventStartingServer) event()  {}
func (e EventStartedServer) event()   {}
func (e EventStoppingServer) event()  {}
func (e EventStartCrawl) event()      {}
func (e EventStopCrawl) event()       {}
func (e EventVisitedResource) event() {}
