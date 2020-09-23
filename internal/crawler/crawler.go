// Package crawler provides a website crawler.
package crawler

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	dom "github.com/PuerkitoBio/goquery"
	"github.com/tj/staticgen/internal/deduplicator"
)

// atImportRe is the regexp used for parsing @import directives.
var atImportRe = regexp.MustCompile(`@import *"([^"]+)"`)

// A Target is a target URL to crawl, with optional Parent page URL.
type Target struct {
	Parent *url.URL
	URL    *url.URL
}

// A Resource is representation of the response to a Target request,
// for a particular page or asset.
type Resource struct {
	Target
	StatusCode int
	Duration   time.Duration
	Body       io.ReadCloser
	Error      error
}

// A Crawler is in charge of visiting or "crawling"
// all pages and assets of a particular URL.
type Crawler struct {
	URL         *url.URL
	Concurrency int
	Allow404    bool
	HTTPClient  *http.Client

	pending    sync.WaitGroup
	resources  chan Resource
	targets    chan Target
	duplicates deduplicator.Deduplicator
	done       chan struct{}
}

// Run starts the crawling process and waits for completion.
func (c *Crawler) Run(ctx context.Context) error {
	err := c.Start(ctx)
	if err != nil {
		return fmt.Errorf("starting: %w", err)
	}

	err = c.Wait()
	if err != nil {
		return fmt.Errorf("waiting: %w", err)
	}

	return nil
}

// Start crawling workers asynchronously. Use Wait() to block until completion.
func (c *Crawler) Start(ctx context.Context) error {
	// defaults
	if c.Concurrency == 0 {
		c.Concurrency = 1
	}

	if c.HTTPClient == nil {
		c.HTTPClient = http.DefaultClient
	}

	// setup
	c.resources = make(chan Resource)
	c.targets = make(chan Target)
	c.done = make(chan struct{})

	// initial page
	c.Queue(c.URL)

	// start workers
	ctx, cancel := context.WithCancel(ctx)
	for i := 0; i < c.Concurrency; i++ {
		go c.crawl(ctx)
	}

	// wait for completion
	go func() {
		c.pending.Wait()
		cancel()
		close(c.done)
	}()

	return nil
}

// Queue a given URL. This method is non-blocking.
func (c *Crawler) Queue(u *url.URL) {
	c.pending.Add(1)
	go func() {
		c.targets <- Target{URL: u}
	}()
}

// Wait for all pending targets to be crawled.
func (c *Crawler) Wait() error {
	<-c.done
	return nil
}

// Resources returns a channel of resources visited by the crawler.
func (c *Crawler) Resources() <-chan Resource {
	return c.resources
}

// crawl all targets, discover additional links,
// and publish resources visited to the Resources() channel.
func (c *Crawler) crawl(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case t := <-c.targets:
			urls, r, err := c.visit(ctx, t)

			// send resource error
			if err != nil {
				r.Error = err
				select {
				case c.resources <- r:
					c.pending.Done()
				case <-ctx.Done():
					return
				}
				return
			}

			// queue urls
			urls = c.duplicates.Filter(urls)
			c.pending.Add(len(urls))
			go c.queue(urls, t)

			// send resource
			select {
			case c.resources <- r:
				c.pending.Done()
			case <-ctx.Done():
				return
			}
		}
	}
}

// visit a target and return any additional targets to crawl.
func (c *Crawler) visit(ctx context.Context, t Target) ([]*url.URL, Resource, error) {
	start := time.Now()
	r := Resource{Target: t}

	// request
	req, err := http.NewRequest("GET", t.URL.String(), nil)
	if err != nil {
		return nil, r, err
	}

	req = req.WithContext(ctx)

	// response
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, r, err
	}

	r.StatusCode = res.StatusCode
	r.Duration = time.Since(start)
	r.Body = res.Body

	if c.Allow404 && res.StatusCode == 404 {
		return nil, r, nil
	}

	// http error
	if res.StatusCode >= 300 {
		return nil, r, fmt.Errorf("%s response", res.Status)
	}

	// file handling
	switch filepath.Ext(t.URL.Path) {
	case ".css":
		defer res.Body.Close()
		var buf bytes.Buffer
		body := io.TeeReader(res.Body, &buf)
		urls, err := visitCSS(body, c.URL, t.URL)
		r.Body = ioutil.NopCloser(&buf)
		return urls, r, err
	case ".html", ".htm", "":
		defer res.Body.Close()
		var buf bytes.Buffer
		body := io.TeeReader(res.Body, &buf)
		urls, err := visitHTML(body, c.URL, t.URL)
		r.Body = ioutil.NopCloser(&buf)
		return urls, r, err
	default:
		return nil, r, nil
	}
}

// queue the given urls with parent target.
func (c *Crawler) queue(urls []*url.URL, t Target) {
	for _, u := range urls {
		c.targets <- Target{
			URL:    u,
			Parent: t.URL,
		}
	}
}

// visitCSS returns targets found in a CSS file.
func visitCSS(r io.Reader, root, u *url.URL) ([]*url.URL, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return parseImports(b, root, u)
}

// parseImports returns CSS imports.
func parseImports(b []byte, root, u *url.URL) (urls []*url.URL, err error) {
	matches := atImportRe.FindAllSubmatch(b, -1)
	for _, m := range matches {
		s := string(m[1])

		target, err := url.Parse(s)
		if err != nil {
			return nil, fmt.Errorf("parsing css import %q: %w", s, err)
		}

		resolved := u.ResolveReference(target)

		if follow(root, resolved) {
			urls = append(urls, resolved)
		}
	}

	return
}

// visitHTML returns targets found in an HTML file.
func visitHTML(r io.Reader, root, u *url.URL) ([]*url.URL, error) {
	doc, err := dom.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}

	return parseLinks(doc, root, u)
}

// parseLinks returns resolved target urls in the document.
func parseLinks(doc *dom.Document, root, u *url.URL) (urls []*url.URL, err error) {
	doc.Find("a, link").Each(func(i int, s *dom.Selection) {
		href := s.AttrOr("href", s.AttrOr("src", ""))
		if href == "" {
			return
		}

		target, err := url.Parse(href)
		if err != nil {
			return
		}

		target.Fragment = ""
		target.RawQuery = ""

		resolved := u.ResolveReference(target)

		if follow(root, resolved) {
			urls = append(urls, resolved)
		}
	})

	return urls, nil
}

// follow returns true if URL u should be followed.
func follow(root, u *url.URL) bool {
	// invalid scheme
	if u.Scheme != "https" && u.Scheme != "http" {
		return false
	}

	// cross domain
	if u.Host != root.Host {
		return false
	}

	// path prefix
	if !strings.HasPrefix(u.Path, root.Path) {
		return false
	}

	return true
}
