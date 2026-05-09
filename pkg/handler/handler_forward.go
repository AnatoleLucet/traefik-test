package handler

import (
	"net/http"
	"strings"

	"github.com/AnatoleLucet/traefik-test/pkg/proxy"
)

func Forward(upstream string, p *proxy.Proxy) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		target := buildTarget(upstream, r)
		p.Forward(target, w, r)
	})
}

func buildTarget(upstream string, r *http.Request) string {
	scheme, host, hasScheme := strings.Cut(upstream, "://")
	if !hasScheme {
		host = scheme
		if r.TLS != nil {
			scheme = "https"
		} else {
			scheme = "http"
		}
	}

	return scheme + "://" + host
}
