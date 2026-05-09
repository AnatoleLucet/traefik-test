package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCache_Get(t *testing.T) {
	t.Run("miss on empty cache", func(t *testing.T) {
		c := New()

		entry, ok := c.Get("foo")

		assert.False(t, ok)
		assert.Equal(t, Entry{}, entry)
	})

	t.Run("hit after set", func(t *testing.T) {
		c := New()
		c.Set("foo", Entry{Body: []byte("bar"), ExpiresAt: time.Now().Add(time.Minute)})

		entry, ok := c.Get("foo")

		assert.True(t, ok)
		assert.Equal(t, []byte("bar"), entry.Body)
	})

	t.Run("miss on expired entry", func(t *testing.T) {
		c := New()
		c.Set("foo", Entry{Body: []byte("bar"), ExpiresAt: time.Now().Add(-time.Minute)})

		entry, ok := c.Get("foo")

		assert.False(t, ok)
		assert.Equal(t, Entry{}, entry)
	})
}

func TestCache_janitor(t *testing.T) {
	t.Run("removes expired entries periodically", func(t *testing.T) {
		c := &Cache{}
		go c.janitor(10 * time.Millisecond)
		c.Set("foo", Entry{Body: []byte("bar"), ExpiresAt: time.Now().Add(5 * time.Millisecond)})

		entry, ok := c.Get("foo")
		assert.True(t, ok)
		assert.Equal(t, []byte("bar"), entry.Body)

		time.Sleep(20 * time.Millisecond)

		// should be removed by janitor
		_, ok = c.Get("foo")
		assert.False(t, ok)
	})
}
