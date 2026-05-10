package middleware

import (
	"maps"
	"net/http"
	"slices"
	"strconv"
	"time"

	"golang.org/x/sync/singleflight"

	"github.com/AnatoleLucet/traefik-test/pkg/cache"
	"github.com/AnatoleLucet/traefik-test/pkg/header"
)

// Cache is a simple TTL cache middleware.
func Cache(next http.Handler, ttl time.Duration, ch *cache.Cache) http.Handler {
	var group singleflight.Group

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !isRequestCacheable(r) {
			next.ServeHTTP(w, r)
			return
		}

		// always assume its a miss until we know otherwise
		w.Header().Set("X-Cache", "MISS")

		if entry, ok := ch.Get(r); ok {
			writeEntry(w, entry)
			return
		}

		served := false // check if we're the one serving the request or if another goroutine already did
		v, _, _ := group.Do(r.Method+" "+r.URL.String(), func() (any, error) {
			served = true

			rr := newResponseRecorder(w)
			next.ServeHTTP(rr, r)

			entry := cache.Entry{
				Vary:      rr.Header().Get("Vary"),
				Body:      rr.Body(),
				Status:    rr.Status(),
				Headers:   header.StripHopHeaders(rr.Header().Clone()),
				CreatedAt: time.Now(),
				ExpiresAt: time.Now().Add(ttl),
			}

			if isResponseCacheable(rr.Status(), rr.Header()) {
				ch.Store(r, entry)
			}

			return entry, nil
		})

		if !served {
			flushEntry(w, v.(cache.Entry))
		}
	})
}

func isRequestCacheable(r *http.Request) bool {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		return false
	}

	cc := header.Split(r.Header.Get("Cache-Control"))
	if slices.Contains(cc, "no-cache") || slices.Contains(cc, "no-store") {
		return false
	}

	return true
}

func isResponseCacheable(status int, headers http.Header) bool {
	if headers.Get("Vary") == "*" {
		return false
	}

	cc := header.Split(headers.Get("Cache-Control"))
	if slices.Contains(cc, "no-cache") ||
		slices.Contains(cc, "no-store") ||
		slices.Contains(cc, "private") {
		return false
	}

	if status < 200 || status >= 300 {
		return false
	}

	return true
}

func writeEntry(w http.ResponseWriter, entry cache.Entry) {
	age := int(time.Since(entry.CreatedAt).Seconds())

	maps.Copy(w.Header(), entry.Headers)
	w.Header().Set("X-Cache", "HIT")
	w.Header().Set("Age", strconv.Itoa(age))

	w.WriteHeader(entry.Status)
	w.Write(entry.Body)
}

func flushEntry(w http.ResponseWriter, entry cache.Entry) {
	maps.Copy(w.Header(), entry.Headers)
	w.WriteHeader(entry.Status)
	w.Write(entry.Body)
}
