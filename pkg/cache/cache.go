package cache

import (
	"net/http"
	"sync"
	"time"
)

const janitorInterval = time.Minute

type Entry struct {
	Status  int
	Body    []byte
	Headers http.Header

	CreatedAt time.Time
	ExpiresAt time.Time
}

type Cache struct {
	entries sync.Map
}

func New() *Cache {
	c := &Cache{}

	go c.janitor(janitorInterval)
	return c
}

func (c *Cache) Get(key string) (Entry, bool) {
	val, ok := c.entries.Load(key)
	if !ok {
		return Entry{}, false
	}

	entry := val.(Entry)
	if time.Now().After(entry.ExpiresAt) {
		c.entries.Delete(key)
		return Entry{}, false
	}

	return entry, true
}

func (c *Cache) Set(key string, entry Entry) {
	c.entries.Store(key, entry)
}

func (c *Cache) janitor(interval time.Duration) {
	for range time.Tick(interval) {
		c.entries.Range(func(k, v any) bool {
			if time.Now().After(v.(Entry).ExpiresAt) {
				c.entries.Delete(k)
			}

			return true
		})
	}
}
