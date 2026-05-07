package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
)

const DEFAULT_CONFIG_PATH = "greenlight.yml"

func Load() (*Config, error) {
	pth := os.Getenv("CONFIG_PATH")
	if pth == "" {
		pth = DEFAULT_CONFIG_PATH
	}

	return LoadFromFile(filepath.Clean(pth))
}

func LoadFromFile(pth string) (*Config, error) {
	cfg, err := os.ReadFile(pth)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrReadFile, err)
	}

	var config Config
	err = yaml.Unmarshal(cfg, &config)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrUnmarshal, err)
	}

	// NOTE: ideally we'd fully validate the config after parsing it.
	// but it feels a bit much for this simple project.
	// i dont want to clutter the assignement.
	//
	// err = validateConfig(config)
	// if err != nil {
	// 	return nil, fmt.Errorf("%w: %w", ErrInvalidConfig, err)
	// }

	return &config, nil
}
