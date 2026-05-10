package middleware

import (
	"maps"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/AnatoleLucet/traefik-test/pkg/cache"
)

// Cache is a simple TTL cache middleware.
func Cache(next http.Handler, ttl time.Duration, ch *cache.Cache) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// always assume its a miss until we know otherwise
		w.Header().Set("X-Cache", "MISS")

		if !isRequestCacheable(r) {
			next.ServeHTTP(w, r)
			return
		}

		key := r.Method + " " + r.URL.String()
		if cached, ok := ch.Get(key); ok {
			age := int(time.Since(cached.CreatedAt).Seconds())

			maps.Copy(w.Header(), cached.Headers)
			w.Header().Set("X-Cache", "HIT")
			w.Header().Set("Age", strconv.Itoa(age))

			w.WriteHeader(cached.Status)
			w.Write(cached.Body)
			return
		}

		rr := newResponseRecorder(w)
		next.ServeHTTP(rr, r)

		if isResponseCacheable(rr.Status(), rr.Header()) {
			ch.Set(key, cache.Entry{
				Body:      rr.Body(),
				Status:    rr.Status(),
				Headers:   rr.Header().Clone(),
				CreatedAt: time.Now(),
				ExpiresAt: time.Now().Add(ttl),
			})
		}
	})
}

func isRequestCacheable(r *http.Request) bool {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		return false
	}

	cc := r.Header.Get("Cache-Control")
	if strings.Contains(cc, "no-cache") {
		return false
	}

	return true
}

func isResponseCacheable(status int, headers http.Header) bool {
	// vary unsupported for now
	if headers.Get("Vary") != "" {
		return false
	}

	cc := headers.Get("Cache-Control")
	if strings.Contains(cc, "no-cache") ||
		strings.Contains(cc, "no-store") ||
		strings.Contains(cc, "private") {
		return false
	}

	if status < 200 || status >= 300 {
		return false
	}

	return true
}
