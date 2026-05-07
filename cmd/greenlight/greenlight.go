package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/AnatoleLucet/traefik-test/pkg/config"
	"github.com/AnatoleLucet/traefik-test/pkg/router"
	"github.com/AnatoleLucet/traefik-test/pkg/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	app := NewApp(cfg)

	err = app.Serve()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Stopping server. Cause: %v\n", err)
		app.Shutdown()
		os.Exit(1)
	}
}

type App struct {
	router *router.Router
	// cache *cache.Cache

	server *server.Server
}

func NewApp(cfg config.Config) *App {
	app := &App{}

	var entrypoints []server.Entrypoint
	if cfg.Server.Ports.HTTP != "" {
		entrypoints = append(entrypoints, server.NewHTTPEntrypoint(cfg.Server, app.HandleHTTP))
	}
	if cfg.Server.Ports.HTTPS != "" {
		entrypoints = append(entrypoints, server.NewHTTPSEntrypoint(cfg.Server, app.HandleHTTP))
	}

	app.router = router.New(cfg.Rules)
	app.server = server.New(entrypoints)

	return app
}

func (app *App) Serve() error {
	return app.server.Serve()
}

func (app *App) Shutdown() error {
	return app.server.Shutdown()
}

func (app *App) HandleHTTP(w http.ResponseWriter, r *http.Request) {
	rule, ok := app.router.Match(router.Request{
		Host:   r.Host,
		Path:   r.URL.Path,
		Method: r.Method,
	})

	if ok {
		fmt.Fprintf(w, "Matched rule: %+v", rule)
	} else {
		http.Error(w, "Bad Gateway", http.StatusBadGateway)
	}
}
