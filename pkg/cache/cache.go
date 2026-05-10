package cache

import (
	"context"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/AnatoleLucet/traefik-test/pkg/header"
)

const janitorInterval = time.Minute

type Entry struct {
	Vary string

	Status  int
	Body    []byte
	Headers http.Header

	CreatedAt time.Time
	ExpiresAt time.Time
}

// Cache is an in-memory cache for HTTP responses, keyed by method + URL + Vary headers.
type Cache struct {
	varies  sync.Map // map[basekey][]header
	entries sync.Map // map[fullkey]Entry
}

func New(ctx context.Context) *Cache {
	c := &Cache{}
	go c.janitor(ctx, janitorInterval)
	return c
}

func (c *Cache) Get(req *http.Request) (Entry, bool) {
	bk := baseKey(req)
	vary, ok := c.getVary(bk)
	if !ok {
		return Entry{}, false
	}

	fk := fullKey(req, vary)
	entry, ok := c.getEntry(fk)
	if !ok {
		return Entry{}, false
	}

	if time.Now().After(entry.ExpiresAt) {
		c.entries.Delete(fk)
		return Entry{}, false
	}

	return entry, true
}

func (c *Cache) Store(req *http.Request, entry Entry) {
	if entry.Vary == "*" {
		return // can't cache responses that vary on all headers
	}

	vary := header.ParseVary(entry.Vary)
	sort.Strings(vary) // ensure stable order for keys

	bk := baseKey(req)
	c.varies.Store(bk, vary)

	fk := fullKey(req, vary)
	c.entries.Store(fk, entry)
}

func (c *Cache) getVary(base string) ([]string, bool) {
	val, ok := c.varies.Load(base)
	if !ok {
		return nil, false
	}

	return val.([]string), true
}

func (c *Cache) getEntry(key string) (Entry, bool) {
	val, ok := c.entries.Load(key)
	if !ok {
		return Entry{}, false
	}

	return val.(Entry), true
}

func (c *Cache) janitor(ctx context.Context, interval time.Duration) {
	t := time.NewTicker(interval)
	defer t.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			c.cleanup()
		}
	}
}

func (c *Cache) cleanup() {
	c.entries.Range(func(k, v any) bool {
		if time.Now().After(v.(Entry).ExpiresAt) {
			c.entries.Delete(k)
		}

		return true
	})
}

// "<method>\x00<url>"
func baseKey(req *http.Request) string {
	return req.Method + "\x00" + req.URL.String()
}

// "<method>\x00<url>\x00<vary-header>:<req-header>\x00..."
func fullKey(req *http.Request, vary []string) string {
	if len(vary) == 0 {
		return baseKey(req)
	}

	var key strings.Builder
	key.WriteString(baseKey(req))
	for _, h := range vary {
		key.WriteByte(0x00)
		key.WriteString(h)
		key.WriteByte(':')
		key.WriteString(strings.Join(req.Header[h], ","))
	}

	return key.String()
}
