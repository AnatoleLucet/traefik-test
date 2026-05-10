package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/AnatoleLucet/traefik-test/pkg/cache"
)

func TestCache(t *testing.T) {
	t.Run("miss and store", func(t *testing.T) {
		c := cache.New(t.Context())
		called := 0
		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called++
			w.Header().Set("X-From-Next", "true")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("hello"))
		})

		mw := Cache(next, time.Minute, c)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/foo", nil)
		mw.ServeHTTP(rec, req)

		assert.Equal(t, 1, called)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "MISS", rec.Header().Get("X-Cache"))
		assert.Equal(t, "true", rec.Header().Get("X-From-Next"))
		assert.Equal(t, "hello", rec.Body.String())
	})

	t.Run("hit", func(t *testing.T) {
		c := cache.New(t.Context())
		called := 0
		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called++
			w.Header().Set("X-From-Next", "true")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("hello"))
		})

		mw := Cache(next, time.Minute, c)

		// miss
		rec1 := httptest.NewRecorder()
		req1 := httptest.NewRequest(http.MethodGet, "/foo", nil)
		mw.ServeHTTP(rec1, req1)
		require.Equal(t, 1, called)

		// hit
		rec2 := httptest.NewRecorder()
		req2 := httptest.NewRequest(http.MethodGet, "/foo", nil)
		mw.ServeHTTP(rec2, req2)

		assert.Equal(t, 1, called, "next handler should not be called on cache hit")
		assert.Equal(t, http.StatusOK, rec2.Code)
		assert.Equal(t, "HIT", rec2.Header().Get("X-Cache"))
		assert.Equal(t, "hello", rec2.Body.String())
	})

	t.Run("non-cacheable request passes through", func(t *testing.T) {
		c := cache.New(t.Context())
		called := 0
		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called++
			w.WriteHeader(http.StatusCreated)
		})

		mw := Cache(next, time.Minute, c)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/foo", nil) // POST is not cacheable
		req.Header.Set("Cache-Control", "no-store")
		mw.ServeHTTP(rec, req)

		assert.Equal(t, 1, called)
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Empty(t, rec.Header().Get("X-Cache"))
	})

	t.Run("non-cacheable response is not stored", func(t *testing.T) {
		c := cache.New(t.Context())
		called := 0
		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called++
			w.Header().Set("Cache-Control", "no-store")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("secret"))
		})

		mw := Cache(next, time.Minute, c)

		// miss
		rec1 := httptest.NewRecorder()
		req1 := httptest.NewRequest(http.MethodGet, "/foo", nil)
		mw.ServeHTTP(rec1, req1)
		require.Equal(t, 1, called)

		// miss because upstream response is not cacheable
		rec2 := httptest.NewRecorder()
		req2 := httptest.NewRequest(http.MethodGet, "/foo", nil)
		mw.ServeHTTP(rec2, req2)

		assert.Equal(t, 2, called)
		assert.Equal(t, "MISS", rec2.Header().Get("X-Cache"))
	})
}
