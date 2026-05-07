package server

import (
	"context"
	"errors"
	"time"
)

type Entrypoint interface {
	Serve() error
	Shutdown(ctx context.Context) error
}

// Server orchestrates the lifecycle of multiple entrypoints (protocols).
type Server struct {
	entrypoints []Entrypoint
}

func New(entrypoints []Entrypoint) *Server {
	return &Server{
		entrypoints: entrypoints,
	}
}

func (s *Server) Serve() error {
	errc := make(chan error, len(s.entrypoints))
	for _, ep := range s.entrypoints {
		go func() {
			errc <- ep.Serve()
		}()
	}

	var errs []error
	for range len(s.entrypoints) {
		errs = append(errs, <-errc)
	}

	return errors.Join(errs...)
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var errs []error
	for _, ep := range s.entrypoints {
		errs = append(errs, ep.Shutdown(ctx))
	}

	return errors.Join(errs...)
}
