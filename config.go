package staticgen

import (
	"github.com/tj/go-config"
	"time"
)

// Config is the static website generator configuration.
type Config struct {
	// URL is the target website to crawl. Defaults to "http://127.0.0.1:3000".
	URL string `json:"url"`

	// Dir is the static website output directory. Defaults to "build".
	Dir string `json:"dir"`

	// Command is the optional server command executed before crawling.
	Command string `json:"command"`

	// Pages is a list of paths added to crawl, typically
	// including unlinked pages such as error pages,
	// landing pages and so on.
	Pages []string `json:"pages"`

	// Concurrency is the number of concurrent pages to crawl. Defaults to 30.
	Concurrency int `json:"concurrency"`

	// Allow404 can be enabled to opt-in to pages resulting in a 404,
	// which otherwise lead to an error.
	Allow404 bool `json:"allow_404"`

	// ResourceTimeout specifies a time limit in seconds for requests for any resources.
	// The timeout includes connection time, any redirects, and reading the response body.
	// Set to 0 to disable. Defaults to 10.
	ResourceTimeout time.Duration `json:"resource_timeout"`
}

// Load configuration from the given path.
func (c *Config) Load(path string) error {
	if c.URL == "" {
		c.URL = "http://127.0.0.1:3000"
	}

	if c.Dir == "" {
		c.Dir = "build"
	}

	if c.Concurrency == 0 {
		c.Concurrency = 30
	}

	if c.ResourceTimeout == 0 {
		c.ResourceTimeout = 10
	}

	err := config.Load(path, c)
	if err != nil {
		return err
	}

	return nil
}
