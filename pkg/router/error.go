package router

import "errors"

var (
	ErrCompileHandler     = errors.New("failed to compile handler")
	ErrCompileMiddlewares = errors.New("failed to compile middlewares")

	ErrInvalidTTL = errors.New("invalid TTL value")
)
