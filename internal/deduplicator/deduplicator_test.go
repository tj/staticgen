package deduplicator_test

import (
	"net/url"
	"testing"

	"github.com/tj/assert"

	"github.com/tj/staticgen/internal/deduplicator"
)

// Test filtering.
func TestDeduplicator_Filter(t *testing.T) {
	var d deduplicator.Deduplicator

	apex, _ := url.Parse("https://apex.sh")
	netflix, _ := url.Parse("https://netflix.com")
	youtube, _ := url.Parse("https://youtube.com")
	facebook, _ := url.Parse("https://facebook.com")

	urls := []*url.URL{apex, netflix, youtube}

	urls = d.Filter(urls)
	assert.Len(t, urls, 3)

	urls = d.Filter(urls)
	assert.Len(t, urls, 0)

	urls = append(urls, facebook)
	urls = d.Filter(urls)
	assert.Len(t, urls, 1)
}
