package server

import "context"

type Entrypoint interface {
	Serve() error
	Shutdown(ctx context.Context) error
}
