package main

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/AnatoleLucet/traefik-test/pkg/config"
	"github.com/AnatoleLucet/traefik-test/pkg/router"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error starting server: %v\n", err)
		os.Exit(1)
	}

	rtr := router.New(cfg.Rules)

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error accepting connection: %v\n", err)
			continue
		}

		go handler(conn, rtr)
	}
}

func handler(conn net.Conn, rtr *router.Router) {
	defer conn.Close()

	req, err := http.ReadRequest(bufio.NewReader(conn))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading request: %v\n", err)
		return
	}

	rule, ok := rtr.Match(router.Request{
		Host:   req.Host,
		Path:   req.URL.Path,
		Method: req.Method,
	})
	if ok {
		fmt.Fprintf(conn, "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n\r\nMatched rule: %+v", rule)
	} else {
		fmt.Fprint(conn, "HTTP/1.1 502 Bad Gateway\r\n\r\n")
	}
}
