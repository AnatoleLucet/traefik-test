package header

import (
	"net"
	"net/http"
	"strings"
)

var hopHeaders = []string{
	"Connection",
	"Keep-Alive",
	"Transfer-Encoding",
	"Te",
	"Trailers",
	"Upgrade",
	"Proxy-Authorization",
	"Proxy-Authenticate",
}

func StripHopHeaders(header http.Header) http.Header {
	clone := header.Clone()
	for _, h := range hopHeaders {
		clone.Del(h)
	}

	return clone
}

func SetXForwarded(out *http.Request, in *http.Request) {
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

func ParseVary(vary string) []string {
	if vary == "" || vary == "*" {
		return nil
	}

	parts := strings.Split(vary, ",")

	headers := make([]string, 0, len(parts))
	for _, p := range parts {
		h := strings.TrimSpace(p)
		if h != "" && h != "*" {
			headers = append(headers, http.CanonicalHeaderKey(h))
		}
	}

	return headers
}

func Split(header string) []string {
	parts := strings.Split(header, ",")

	headers := make([]string, 0, len(parts))
	for _, p := range parts {
		h := strings.TrimSpace(p)
		if h != "" {
			headers = append(headers, strings.ToLower(h))
		}
	}

	return headers
}
