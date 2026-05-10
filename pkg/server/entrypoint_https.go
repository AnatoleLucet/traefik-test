package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/AnatoleLucet/traefik-test/pkg/config"
	"github.com/AnatoleLucet/traefik-test/pkg/logger"
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
	logger.Infof("Starting HTTPS server on \"%s\".", h.server.Addr)

	err := h.server.ListenAndServeTLS(h.cert, h.key)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("%w: %w", ErrServeHTTPS, err)
	}

	return nil
}

func (h *HTTPSEntrypoint) Shutdown(ctx context.Context) error {
	logger.Infof("Shutting down HTTPS server.")

	err := h.server.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrShutdownHTTPS, err)
	}

	return nil
}
