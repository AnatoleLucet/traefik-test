package cache

import (
	"net/http"
	"sync"
	"time"
)

type Entry struct {
	Status  int
	Body    []byte
	Headers http.Header

	CreatedAt time.Time
	ExpiresAt time.Time
}

type Cache struct {
	mu      sync.RWMutex
	entries map[string]Entry
}

func New() *Cache {
	return &Cache{
		entries: make(map[string]Entry),
	}
}

func (c *Cache) Get(key string) (Entry, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, ok := c.entries[key]
	if !ok || time.Now().After(entry.ExpiresAt) {
		return Entry{}, false
	}

	return entry, true
}

func (c *Cache) Set(key string, entry Entry) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[key] = entry
}
