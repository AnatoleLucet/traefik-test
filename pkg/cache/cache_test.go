package cache

import (
	"context"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCache_Get(t *testing.T) {
	t.Run("miss on empty cache", func(t *testing.T) {
		c := New(context.Background())
		req := &http.Request{
			Method: http.MethodGet,
			URL:    &url.URL{Scheme: "http", Host: "example.com", Path: "/foo"},
		}

		entry, ok := c.Get(req)

		assert.False(t, ok)
		assert.Equal(t, Entry{}, entry)
	})

	t.Run("hit after store", func(t *testing.T) {
		c := New(context.Background())
		req := &http.Request{
			Method: http.MethodGet,
			URL:    &url.URL{Scheme: "http", Host: "example.com", Path: "/foo"},
		}
		c.Store(req, Entry{Body: []byte("bar"), ExpiresAt: time.Now().Add(time.Minute)})

		entry, ok := c.Get(req)

		assert.True(t, ok)
		assert.Equal(t, []byte("bar"), entry.Body)
	})

	t.Run("miss on expired entry", func(t *testing.T) {
		c := New(context.Background())
		req := &http.Request{
			Method: http.MethodGet,
			URL:    &url.URL{Scheme: "http", Host: "example.com", Path: "/foo"},
		}
		c.Store(req, Entry{Body: []byte("bar"), ExpiresAt: time.Now().Add(-time.Minute)})

		entry, ok := c.Get(req)

		assert.False(t, ok)
		assert.Equal(t, Entry{}, entry)
	})
}

func TestCache_Vary(t *testing.T) {
	t.Run("hit with matching vary header", func(t *testing.T) {
		c := New(context.Background())
		req := &http.Request{
			Method: http.MethodGet,
			URL:    &url.URL{Scheme: "http", Host: "example.com", Path: "/foo"},
			Header: http.Header{"Accept-Encoding": []string{"gzip"}},
		}
		c.Store(req, Entry{Vary: "Accept-Encoding", Body: []byte("gzip-body"), ExpiresAt: time.Now().Add(time.Minute)})

		entry, ok := c.Get(req)

		assert.True(t, ok)
		assert.Equal(t, []byte("gzip-body"), entry.Body)
	})

	t.Run("miss with different vary header value", func(t *testing.T) {
		c := New(context.Background())
		storeReq := &http.Request{
			Method: http.MethodGet,
			URL:    &url.URL{Scheme: "http", Host: "example.com", Path: "/foo"},
			Header: http.Header{"Accept-Encoding": []string{"gzip"}},
		}
		getReq := &http.Request{
			Method: http.MethodGet,
			URL:    &url.URL{Scheme: "http", Host: "example.com", Path: "/foo"},
			Header: http.Header{"Accept-Encoding": []string{"deflate"}},
		}
		c.Store(storeReq, Entry{Vary: "Accept-Encoding", Body: []byte("gzip-body"), ExpiresAt: time.Now().Add(time.Minute)})

		_, ok := c.Get(getReq)

		assert.False(t, ok)
	})
}

func TestCache_janitor(t *testing.T) {
	t.Run("removes expired entries periodically", func(t *testing.T) {
		c := &Cache{}
		go c.janitor(context.Background(), 10*time.Millisecond)
		req := &http.Request{
			Method: http.MethodGet,
			URL:    &url.URL{Scheme: "http", Host: "example.com", Path: "/foo"},
		}
		c.Store(req, Entry{Body: []byte("bar"), ExpiresAt: time.Now().Add(5 * time.Millisecond)})

		entry, ok := c.Get(req)
		assert.True(t, ok)
		assert.Equal(t, []byte("bar"), entry.Body)

		time.Sleep(20 * time.Millisecond)

		// should be removed by janitor
		_, ok = c.Get(req)
		assert.False(t, ok)
	})
}
