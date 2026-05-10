package proxy

import (
	"io"
	"maps"
	"net/http"
	"time"

	"github.com/AnatoleLucet/traefik-test/pkg/header"
	"github.com/AnatoleLucet/traefik-test/pkg/logger"
)

const (
	maxIdleConns        = 100
	maxIdleConnsPerHost = 10
	idleConnTimeout     = 90 * time.Second
)

type Proxy struct {
	client *http.Client
}

func New() *Proxy {
	return &Proxy{
		client: &http.Client{
			Transport: &http.Transport{
				MaxIdleConns:        maxIdleConns,
				MaxIdleConnsPerHost: maxIdleConnsPerHost,
				IdleConnTimeout:     idleConnTimeout,
			},
		},
	}
}

func (p *Proxy) Forward(upstream string, w http.ResponseWriter, r *http.Request) {
	target := upstream + r.URL.RequestURI()

	req, err := http.NewRequestWithContext(r.Context(), r.Method, target, r.Body)
	if err != nil {
		logger.Errorf("Failed to forward request: %v", err)
		http.Error(w, "Bad Gateway", http.StatusBadGateway)
		return
	}

	maps.Copy(req.Header, r.Header)
	req.Header = header.StripHopHeaders(req.Header)
	header.SetXForwarded(req, r)

	res, err := p.client.Do(req)
	if err != nil {
		logger.Errorf("Failed to forward request: %v", err)
		http.Error(w, "Bad Gateway", http.StatusBadGateway)
		return
	}
	defer res.Body.Close()

	res.Header = header.StripHopHeaders(res.Header)
	maps.Copy(w.Header(), res.Header)
	w.WriteHeader(res.StatusCode)
	io.Copy(w, res.Body)
}
