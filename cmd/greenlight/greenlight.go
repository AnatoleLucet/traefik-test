package main

import (
	"context"
	"net/http"
	"os"

	"github.com/AnatoleLucet/traefik-test/pkg/cache"
	"github.com/AnatoleLucet/traefik-test/pkg/config"
	"github.com/AnatoleLucet/traefik-test/pkg/logger"
	"github.com/AnatoleLucet/traefik-test/pkg/proxy"
	"github.com/AnatoleLucet/traefik-test/pkg/router"
	"github.com/AnatoleLucet/traefik-test/pkg/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		logger.Errorf("Failed to load config: %v", err)
		os.Exit(1)
	}

	app, err := NewApp(cfg)
	if err != nil {
		logger.Errorf("Failed to initialize app: %v", err)
		os.Exit(1)
	}

	err = app.Serve()
	if err != nil {
		logger.Errorf("Stopping server. Cause: %v", err)
		app.Shutdown()
		os.Exit(1)
	}
}

type App struct {
	ctx    context.Context
	cancel context.CancelFunc

	proxy *proxy.Proxy
	cache *cache.Cache

	router *router.Router
	server *server.Server
}

func NewApp(cfg config.Config) (*App, error) {
	ctx, cancel := context.WithCancel(context.Background())

	app := &App{
		ctx:    ctx,
		cancel: cancel,
		proxy:  proxy.New(),
		cache:  cache.New(ctx),
	}

	rules, err := router.Compiler{Cache: app.cache, Proxy: app.proxy}.CompileRules(cfg.Rules)
	if err != nil {
		return nil, err
	}
	app.router = router.New(rules)

	var entrypoints []server.Entrypoint
	if cfg.Server.Ports.HTTP.Enabled() {
		entrypoints = append(entrypoints, server.NewHTTPEntrypoint(cfg.Server, app))
	}
	if cfg.Server.Ports.HTTPS.Enabled() {
		entrypoints = append(entrypoints, server.NewHTTPSEntrypoint(cfg.Server, app))
	}
	app.server = server.New(entrypoints)

	return app, nil
}

func (app *App) Serve() error {
	return app.server.Serve()
}

func (app *App) Shutdown() error {
	app.cancel()
	return app.server.Shutdown()
}

func (app *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rl, ok := app.router.Match(router.Request{
		Host:   r.Host,
		Path:   r.URL.Path,
		Method: r.Method,
	})
	if !ok {
		http.Error(w, "Bad Gateway", http.StatusBadGateway)
		return
	}

	rl.Handler.ServeHTTP(w, r)
}
