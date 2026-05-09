package handler

import (
	"fmt"
	"net/http"

	"github.com/AnatoleLucet/traefik-test/pkg/proxy"
)

func Forward(upstream string, proxy *proxy.Proxy) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// proxy.ForwardHTTP(w, r, target)
		fmt.Fprintf(w, "Forwarding to %s\n", upstream)
	})
}
