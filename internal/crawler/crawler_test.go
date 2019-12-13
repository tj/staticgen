package crawler_test

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/tj/assert"

	"github.com/tj/staticgen/internal/crawler"
)

// Test .
func TestCrawler(t *testing.T) {
	u, _ := url.Parse("https://apex.sh")

	c := crawler.Crawler{
		URL:         u,
		Concurrency: 10,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	go func() {
		for {
			select {
			case r := <-c.Resources():
				if err := r.Error; err == nil {
					fmt.Printf("GET %s -> %s (%s)\n", r.URL, http.StatusText(r.StatusCode), r.Duration.Round(time.Millisecond))
					io.Copy(ioutil.Discard, r.Body)
					r.Body.Close()
				} else {
					fmt.Printf("GET %s -> error %s\n", r.URL, r.Error)
				}
			case <-ctx.Done():
				fmt.Printf("reporter done\n")
				return
			}
		}
	}()

	err := c.Run(ctx)
	assert.NoError(t, err, "wait")

	fmt.Printf("done\n")
}
