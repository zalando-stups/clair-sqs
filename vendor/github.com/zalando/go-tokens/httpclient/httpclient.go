package httpclient

import (
	"net"
	"net/http"
	"time"
)

const (
	defaultHTTPClientTimeout    = 10 * time.Second
	defaultHTTPClientTLSTimeout = 10 * time.Second
)

var (
	// UserAgent can be used to specify the User-Agent header sent on every request that used this package's
	// http.Client
	UserAgent = "go-tokens"
)

// Default returns a new http.Client with KeepAlive disabled. That means no connection pooling.
// Use it only for one time requests where performance is not a concern
// It use some settings from the options package: options.HttpClientTimeout and options.HttpClientTlsTimeout
func Default() *http.Client {
	return New(defaultHTTPClientTimeout, defaultHTTPClientTLSTimeout)
}

// New returns a new http.Client with specific timeouts from its arguments. KeepAlive is disabled.
// That means no connection pooling. Use it only for one time requests where performance is not a concern
func New(timeout time.Duration, tlsTimeout time.Duration) *http.Client {
	return &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			Proxy:               http.ProxyFromEnvironment,
			DisableKeepAlives:   true,
			Dial:                (&net.Dialer{Timeout: timeout}).Dial,
			TLSHandshakeTimeout: tlsTimeout}}
}
