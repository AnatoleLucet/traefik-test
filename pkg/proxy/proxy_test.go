package proxy

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProxy_Forward(t *testing.T) {
	t.Run("forwards request and writes response", func(t *testing.T) {
		upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/some/path?key=val", r.URL.RequestURI())
			w.Header().Set("X-Custom", "from-upstream")
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte("upstream-body"))
		}))
		defer upstream.Close()

		p := New()
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/some/path?key=val", strings.NewReader("client-body"))

		p.Forward(upstream.URL, rec, req)

		res := rec.Result()
		defer res.Body.Close()
		body, _ := io.ReadAll(res.Body)

		assert.Equal(t, http.StatusCreated, res.StatusCode)
		assert.Equal(t, "from-upstream", res.Header.Get("X-Custom"))
		assert.Equal(t, "upstream-body", string(body))
	})

	t.Run("returns 502 when upstream is unreachable", func(t *testing.T) {
		p := New()
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/", nil)

		p.Forward("http://127.0.0.1:1", rec, req)

		res := rec.Result()
		defer res.Body.Close()
		body, _ := io.ReadAll(res.Body)

		assert.Equal(t, http.StatusBadGateway, res.StatusCode)
		assert.Contains(t, string(body), "Bad Gateway")
	})

	t.Run("strips hop headers and sets x-forwarded headers", func(t *testing.T) {
		upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Empty(t, r.Header.Get("Connection"), "Connection is a hop header and should be stripped")
			assert.NotEmpty(t, r.Header.Get("X-Forwarded-For"))
			assert.NotEmpty(t, r.Header.Get("X-Forwarded-Host"))
			assert.NotEmpty(t, r.Header.Get("X-Forwarded-Proto"))
			w.WriteHeader(http.StatusOK)
		}))
		defer upstream.Close()

		p := New()
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Connection", "keep-alive")
		req.Host = "client.example.com"

		p.Forward(upstream.URL, rec, req)

		res := rec.Result()
		defer res.Body.Close()
		require.Equal(t, http.StatusOK, res.StatusCode)
	})
}
