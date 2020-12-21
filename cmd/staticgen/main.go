package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/apex/httplog"
	"github.com/apex/log"
	"github.com/apex/log/handlers/text"
	"github.com/pkg/browser"
	"github.com/tj/kingpin"

	"github.com/tj/staticgen"
)

// version of staticgen.
var version string

// main.
func main() {
	app := kingpin.New("staticgen", "Static website generator")
	dir := app.Flag("chdir", "Change working directory").Short('C').Default(".").String()

	log.SetHandler(text.Default)

	app.PreAction(func(_ *kingpin.ParseContext) error {
		return os.Chdir(*dir)
	})

	generateCmd(app)
	serveCmd(app)
	versionCmd(app)

	_, err := app.Parse(os.Args[1:])
	if err != nil {
		log.Fatalf("error: %s\n", err)
	}
}

// generateCmd command.
func generateCmd(app *kingpin.Application) {
	cmd := app.Command("generate", "Generate static website").Default()
	timeout := cmd.Flag("timeout", "Timeout of website generation").Short('t').Default("15m").String()
	cmd.Action(func(_ *kingpin.ParseContext) error {
		// generator
		g := staticgen.Generator{}

		// parse timeout
		d, err := time.ParseDuration(*timeout)
		if err != nil {
			return fmt.Errorf("parsing duration: %w", err)
		}

		// timeout
		ctx, cancel := context.WithTimeout(context.Background(), d)
		defer cancel()

		// trap interrupt and quit
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, os.Interrupt)
		go func() {
			log.Infof("Received signal %s — quitting\n", <-ch)
			cancel()
		}()

		// reporting
		events := make(chan staticgen.Event, 1000)

		var r staticgen.Reporter
		g.Report(events)
		done := r.Report(events)

		// start
		err = g.Run(ctx)
		if err != nil {
			return fmt.Errorf("crawling: %w", err)
		}

		close(events)
		<-done
		return nil
	})
}

// serveCmd command.
func serveCmd(app *kingpin.Application) {
	cmd := app.Command("serve", "Serve the generated website")
	addr := cmd.Flag("address", "Bind address").Default("localhost:3000").String()
	cmd.Action(func(_ *kingpin.ParseContext) error {
		var c staticgen.Config

		err := c.Load("static.json")
		if err != nil {
			return fmt.Errorf("loading configuration: %w", err)
		}

		_ = browser.OpenURL("http://" + *addr)

		server := http.FileServer(http.Dir(c.Dir))
		fmt.Printf("Starting static file server on %s\n", *addr)
		return http.ListenAndServe(*addr, httplog.New(server))
	})
}

// versionCmd command.
func versionCmd(app *kingpin.Application) {
	cmd := app.Command("version", "Output the version.").Hidden()
	cmd.Action(func(_ *kingpin.ParseContext) error {
		fmt.Printf("Staticgen %s\n", version)
		return nil
	})
}
