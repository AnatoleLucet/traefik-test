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

type HTTPEntrypoint struct {
	server *http.Server
}

func NewHTTPEntrypoint(cfg config.Server, handler http.Handler) *HTTPEntrypoint {
	return &HTTPEntrypoint{
		server: &http.Server{
			Addr:    net.JoinHostPort(cfg.Host, string(cfg.Ports.HTTP)),
			Handler: handler,
		},
	}
}

func (h *HTTPEntrypoint) Serve() error {
	logger.Infof("Starting HTTP server on \"%s\".", h.server.Addr)

	err := h.server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("%w: %w", ErrServeHTTP, err)
	}

	return nil
}

func (h *HTTPEntrypoint) Shutdown(ctx context.Context) error {
	logger.Infof("Shutting down HTTP server.")

	err := h.server.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrShutdownHTTP, err)
	}

	return nil
}
