package tokens

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
	// Default global instance of a custom http.Client using the defaults from the options package
	Default = DefaultHTTPClient()
	// UserAgent can be used to specify the User-Agent header sent on every request that used this package's
	// http.Client
	UserAgent = "go-tokens"
)

// DefaultHTTPClient returns a new http.Client with KeepAlive disabled. That means no connection pooling.
// Use it only for one time requests where performance is not a concern
// It use some settings from the options package: options.HttpClientTimeout and options.HttpClientTlsTimeout
func DefaultHTTPClient() *http.Client {
	return NewHTTPClient(defaultHTTPClientTimeout, defaultHTTPClientTLSTimeout)
}

// NewHTTPClient returns a new http.Client with specific timeouts from its arguments. KeepAlive is disabled.
// That means no connection pooling. Use it only for one time requests where performance is not a concern
func NewHTTPClient(timeout time.Duration, tlsTimeout time.Duration) *http.Client {
	return &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			Proxy:               http.ProxyFromEnvironment,
			DisableKeepAlives:   true,
			Dial:                (&net.Dialer{Timeout: timeout}).Dial,
			TLSHandshakeTimeout: tlsTimeout}}
}
