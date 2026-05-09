package proxy

import (
	"io"
	"maps"
	"net"
	"net/http"
	"time"
)

var hopByHop = []string{
	"Connection",
	"Keep-Alive",
	"Transfer-Encoding",
	"Te",
	"Trailers",
	"Upgrade",
	"Proxy-Authorization",
	"Proxy-Authenticate",
}

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
	req, err := http.NewRequestWithContext(r.Context(), r.Method, upstream+r.URL.RequestURI(), r.Body)
	if err != nil {
		http.Error(w, "Bad Gateway", http.StatusBadGateway)
		return
	}

	maps.Copy(req.Header, r.Header)
	removeHopByHop(req.Header)
	setXForwarded(req, r)

	res, err := p.client.Do(req)
	if err != nil {
		http.Error(w, "Bad Gateway", http.StatusBadGateway)
		return
	}
	defer res.Body.Close()

	removeHopByHop(res.Header)
	maps.Copy(w.Header(), res.Header)
	w.WriteHeader(res.StatusCode)
	io.Copy(w, res.Body)
}

func removeHopByHop(h http.Header) {
	for _, key := range hopByHop {
		h.Del(key)
	}
}

func setXForwarded(out *http.Request, in *http.Request) {
	ip, _, _ := net.SplitHostPort(in.RemoteAddr)
	if prior := in.Header.Get("X-Forwarded-For"); prior != "" {
		ip = prior + ", " + ip
	}

	proto := "http"
	if in.TLS != nil {
		proto = "https"
	}

	out.Header.Set("X-Forwarded-For", ip)
	out.Header.Set("X-Forwarded-Host", in.Host)
	out.Header.Set("X-Forwarded-Proto", proto)
}
