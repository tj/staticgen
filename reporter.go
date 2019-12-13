package staticgen

import (
	"net/http"
	"time"

	"github.com/apex/log"
	"github.com/dustin/go-humanize"
)

// A Reporter outputs a human-friendly report of events.
type Reporter struct {
	count int64
	start time.Time
}

// Report on the given event channel. Returns a channel which is
// closed when the event channel is closed, and all reporting has
// been completed.
func (r *Reporter) Report(ch <-chan Event) <-chan struct{} {
	done := make(chan struct{})

	go func() {
		defer close(done)
		for e := range ch {
			switch e := e.(type) {
			case EventStartCrawl:
				r.start = time.Now()
			case EventStartingServer:
				log.Infof("Starting server with command %q", e.Command)
				log.Infof("Waiting for server to listen on %s", e.URL)
			case EventStartedServer:
				log.Infof("Server is listening for requests")
			case EventStoppingServer:
				log.Infof("Stopping server, sending SIGTERM")
			case EventVisitedResource:
				r.count++
				if e.Error == nil {
					log.Infof("GET %s —— %s —— %s (%s)", e.URL, e.Filename, http.StatusText(e.StatusCode), e.Duration.Round(time.Millisecond))
				} else {
					log.Errorf("GET %s —— %s (error: %s)", e.URL, http.StatusText(e.StatusCode), e.Error)
				}
			case EventStopCrawl:
				log.Infof("Completed %s resources in %s", humanize.Comma(r.count), time.Since(r.start).Round(time.Millisecond))
			}
		}
	}()

	return done
}
