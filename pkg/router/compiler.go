package router

import (
	"fmt"
	"net/http"
	"time"

	"github.com/AnatoleLucet/traefik-test/pkg/cache"
	"github.com/AnatoleLucet/traefik-test/pkg/config"
	"github.com/AnatoleLucet/traefik-test/pkg/handler"
	"github.com/AnatoleLucet/traefik-test/pkg/middleware"
	"github.com/AnatoleLucet/traefik-test/pkg/proxy"
)

type CompiledRule struct {
	Rule    config.Rule
	Handler http.Handler
}

// Compiler is responsible for "compiling" rules into static handlers.
type Compiler struct {
	Cache *cache.Cache
	Proxy *proxy.Proxy
}

func (c Compiler) CompileRule(rule config.Rule) (CompiledRule, error) {
	handler, err := c.compileHandler(rule)
	if err != nil {
		return CompiledRule{}, fmt.Errorf("%w: %w", ErrCompileHandler, err)
	}

	handler, err = c.compileMiddlewares(rule, handler)
	if err != nil {
		return CompiledRule{}, fmt.Errorf("%w: %w", ErrCompileMiddlewares, err)
	}

	return CompiledRule{
		Rule:    rule,
		Handler: handler,
	}, nil
}

func (c Compiler) CompileRules(rules []config.Rule) ([]CompiledRule, error) {
	var compiled []CompiledRule
	for _, rule := range rules {
		cmp, err := c.CompileRule(rule)
		if err != nil {
			return nil, err
		}

		compiled = append(compiled, cmp)
	}

	return compiled, nil
}

func (c Compiler) compileHandler(rule config.Rule) (http.Handler, error) {
	if rule.Then.Forward.Enabled() {
		return handler.Forward(string(rule.Then.Forward), c.Proxy), nil
	}

	if rule.Then.Redirect.Enabled() {
		return handler.Redirect(string(rule.Then.Redirect)), nil
	}

	if rule.Then.Respond.Enabled() {
		return handler.Respond(rule.Then.Respond.Body, rule.Then.Respond.Status, rule.Then.Respond.Headers), nil
	}

	return handler.Respond("Bad Gateway", http.StatusBadGateway, nil), nil
}

func (c Compiler) compileMiddlewares(rule config.Rule, handler http.Handler) (http.Handler, error) {
	if rule.Middleware.Cache.Enabled() {
		ttl, err := time.ParseDuration(rule.Middleware.Cache.TTL)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", ErrInvalidTTL, err)
		}

		handler = middleware.Cache(handler, ttl, c.Cache)
	}

	// if rule.Middleware.Logger.Enabled() {
	//   ...
	// }

	return handler, nil
}
