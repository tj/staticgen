// Package staticgen provides static website generation from a live server.
package staticgen

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/apex/log"

	"github.com/tj/staticgen/internal/crawler"
)

// Target is a target URL.
type Target struct {
	Parent *url.URL
	URL    *url.URL
}

// Generator is a static website generator.
type Generator struct {
	// Config used for crawling and producing the static website.
	Config

	// HTTPClient ...
	HTTPClient *http.Client

	// crawler
	crawler crawler.Crawler
	wg      sync.WaitGroup

	// server command
	cmd *exec.Cmd
	out bytes.Buffer

	// events
	events chan<- Event
}

// Run starts the configured server command, starts to perform crawling,
// and waits for completion before shutting down the configured server.
func (g *Generator) Run(ctx context.Context) error {
	err := g.Start(ctx)
	if err != nil {
		return fmt.Errorf("starting: %w", err)
	}

	err = g.Wait()
	if err != nil {
		return fmt.Errorf("waiting: %w", err)
	}

	defer g.emit(EventStopCrawl{})
	if err := g.stopCommand(ctx); err != nil {
		return fmt.Errorf("stopping: %w", err)
	}

	return nil
}

// Start loads configuration from ./static.json, starts the
// configured server, and begins the crawling process.
func (g *Generator) Start(ctx context.Context) error {
	// load configuration
	err := g.Config.Load("static.json")
	if err != nil {
		return fmt.Errorf("loading configuration: %w", err)
	}

	// remove output dir
	err = os.RemoveAll(g.Dir)
	if err != nil {
		return fmt.Errorf("removing output directory: %w", err)
	}

	// create output dir
	err = os.MkdirAll(g.Dir, 0755)
	if err != nil {
		return fmt.Errorf("creating output directory: %w", err)
	}

	// start command
	err = g.startCommand(ctx)
	if err != nil {
		return fmt.Errorf("starting command: %w", err)
	}

	// parse url
	u, err := url.Parse(g.URL)
	if err != nil {
		return fmt.Errorf("parsing url: %w", err)
	}

	// setup crawler
	g.crawler = crawler.Crawler{
		URL:         u,
		Concurrency: g.Concurrency,
		HTTPClient:  g.HTTPClient,
	}

	// start workers
	ctx, cancel := context.WithCancel(ctx)
	for i := 0; i < g.Concurrency; i++ {
		g.wg.Add(1)
		go func() {
			g.saveLoop(ctx)
			g.wg.Done()
		}()
	}

	// start crawling
	g.emit(EventStartCrawl{})
	err = g.crawler.Start(ctx)
	if err != nil {
		cancel()
		return fmt.Errorf("starting crawler: %w", err)
	}

	// queue pages
	g.queuePages(u)

	// wait for crawling to complete,
	// then exit the save loops.
	go func() {
		g.crawler.Wait()
		cancel()
	}()

	return nil
}

// Wait for crawling to complete.
func (g *Generator) Wait() error {
	g.wg.Wait()
	return nil
}

// Report registers a channel for reporting on events.
func (g *Generator) Report(ch chan<- Event) {
	g.events = ch
}

// queuePages queues the configured Pages relative to the given url.
func (g *Generator) queuePages(u *url.URL) {
	for _, p := range g.Pages {
		t := *u
		t.Path = p
		g.crawler.Queue(&t)
	}
}

// saveLoop saves the crawler resources to disk.
func (g *Generator) saveLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case r := <-g.crawler.Resources():
			err := g.save(r)
			if err != nil {
				log.WithError(err).WithField("url", r.URL.String()).Error("error saving")
			}
		}
	}
}

// save a resource to disk.
func (g *Generator) save(r crawler.Resource) error {
	// determine local path for the file
	ext := filepath.Ext(r.URL.Path)
	dir, file := path.Split(r.URL.Path)
	dst := filepath.Join(g.Dir, dir, file)

	// save html into directories with index.html
	if file != "index.html" && (ext == ".html" || ext == "") {
		file = strings.Replace(file, ".html", "", 1)
		dst = filepath.Join(g.Dir, dir, file, "index.html")
	}

	g.emit(EventVisitedResource{
		Target:     Target(r.Target),
		Duration:   r.Duration,
		StatusCode: r.StatusCode,
		Error:      r.Error,
		Filename:   dst,
	})

	// request error, don't copy to disk
	if r.Error != nil {
		return nil
	}

	defer r.Body.Close()
	return writeFile(r.Body, dst)
}

// startCommand starts the configured server command.
func (g *Generator) startCommand(ctx context.Context) error {
	if g.Command == "" {
		return nil
	}

	g.emit(EventStartingServer{
		Command: g.Command,
		URL:     g.URL,
	})

	// start
	g.cmd = exec.Command("sh", "-c", g.Command)
	g.cmd.Env = append(os.Environ(), "STATICGEN=1")
	g.cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	g.cmd.Stdout = &g.out
	g.cmd.Stderr = &g.out
	err := g.cmd.Start()
	if err != nil {
		return err
	}

	// timeout
	ctx, cancel := context.WithTimeout(ctx, time.Second*15)
	defer cancel()

	// wait
	err = waitForListen(ctx, g.URL)
	if err != nil {
		return fmt.Errorf("waiting for app to start: %w", err)
	}

	g.emit(EventStartedServer{
		Command: g.Command,
		URL:     g.URL,
	})

	return nil
}

// stopCommand stops the configured command.
func (g *Generator) stopCommand(ctx context.Context) error {
	if g.Command == "" {
		return nil
	}

	g.emit(EventStoppingServer{})

	pgid, err := syscall.Getpgid(g.cmd.Process.Pid)
	if err != nil {
		return fmt.Errorf("getting process group id: %w", err)
	}

	err = syscall.Kill(-pgid, syscall.SIGTERM)
	if err != nil {
		return fmt.Errorf("kill: %w", err)
	}

	_ = g.cmd.Wait()
	return nil
}

// emit an event.
func (g *Generator) emit(e Event) {
	if g.events != nil {
		g.events <- e
	}
}

// writeFile writes to filename and ensures the directory exists.
func writeFile(r io.Reader, filename string) error {
	dir := filepath.Dir(filename)

	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	f, err := os.Create(filename)
	if err != nil {
		return err
	}

	_, err = io.Copy(f, r)
	if err != nil {
		return err
	}

	return f.Close()
}

// waitForListen blocks until `addr` is listening.
func waitForListen(ctx context.Context, url string) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Second):
			if isListening(url) {
				return nil
			}
		}
	}
}

// isListening returns true if there's an HTTP server listening on url.
func isListening(url string) bool {
	_, err := http.Head(url)
	return err == nil
}
