// Package deduplicator provides URL deduplication.
package deduplicator

import (
	"net/url"
	"strings"
	"sync"
)

// A Deduplicator manages the deduplication of URLs, allowing you
// to filter out those which you've already visited. The zero value
// is valid.
type Deduplicator struct {
	mu      sync.Mutex
	visited map[string]struct{}
}

// Filter implementation.
func (d *Deduplicator) Filter(urls []*url.URL) (filtered []*url.URL) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.visited == nil {
		d.visited = make(map[string]struct{})
	}

	for _, u := range urls {
		u = normalize(u)

		_, ok := d.visited[u.String()]
		if ok {
			continue
		}

		d.visited[u.String()] = struct{}{}
		filtered = append(filtered, u)
	}

	return
}

// normalize returns a URL with its path normalized,
// stripping the tailing "/" if present, treating
// "/blog/" and "/blog" as the same.
func normalize(u *url.URL) *url.URL {
	if strings.HasSuffix(u.Path, "/") {
		c := *u
		c.Path = strings.TrimRight(c.Path, "/")
		return &c
	}
	return u
}
