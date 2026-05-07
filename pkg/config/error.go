package config

import "errors"

var (
	ErrReadFile  = errors.New("failed to read config file")
	ErrUnmarshal = errors.New("failed to unmarshal config file")
)
