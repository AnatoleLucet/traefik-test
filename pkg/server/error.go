package server

import "errors"

var (
	ErrServeHTTP     = errors.New("error while serving HTTP")
	ErrShutdownHTTP  = errors.New("failed to shutdown HTTP server")
	ErrServeHTTPS    = errors.New("error while serving HTTPS")
	ErrShutdownHTTPS = errors.New("failed to shutdown HTTPS server")
)
