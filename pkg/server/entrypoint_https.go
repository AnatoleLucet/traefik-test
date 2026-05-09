package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/AnatoleLucet/traefik-test/pkg/config"
)

type HTTPSEntrypoint struct {
	key  string
	cert string

	server *http.Server
}

func NewHTTPSEntrypoint(cfg config.Server, handler http.Handler) *HTTPSEntrypoint {
	return &HTTPSEntrypoint{
		key:  cfg.TLS.Key,
		cert: cfg.TLS.Cert,
		server: &http.Server{
			Addr:    net.JoinHostPort(cfg.Host, string(cfg.Ports.HTTPS)),
			Handler: handler,
		},
	}
}

func (h *HTTPSEntrypoint) Serve() error {
	fmt.Printf("Starting HTTPS server on \"%s\".\n", h.server.Addr)

	err := h.server.ListenAndServeTLS(h.cert, h.key)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("%w: %w", ErrServeHTTPS, err)
	}

	return nil
}

func (h *HTTPSEntrypoint) Shutdown(ctx context.Context) error {
	fmt.Printf("Shutting down HTTPS server.\n")

	err := h.server.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrShutdownHTTPS, err)
	}

	return nil
}
