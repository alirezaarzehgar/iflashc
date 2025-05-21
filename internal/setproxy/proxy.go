package setproxy

import (
	"fmt"
	"net"
	"net/http"

	"golang.org/x/net/proxy"
)

func NewSocks5Client(addr string, auth *proxy.Auth) (*http.Client, error) {
	s5Client, err := proxy.SOCKS5("tcp", addr, auth, &net.Dialer{})
	if err != nil {
		return nil, fmt.Errorf("failed to create dialer: %w", err)
	}

	contextDialer, ok := s5Client.(proxy.ContextDialer)
	if !ok {
		contextDialer = &net.Dialer{}
	}

	return &http.Client{
		Transport: &http.Transport{
			DialContext: contextDialer.DialContext,
		},
	}, nil
}
