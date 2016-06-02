package httpclient

import (
	"net/http"
	"testing"
)

func TestHttpClient(t *testing.T) {
	c := Default()

	if c.Timeout != defaultHTTPClientTimeout {
		t.Errorf("Unexpected timeout value. Wanted %v, got %v\n", defaultHTTPClientTimeout, c.Timeout)
	}

	transport, ok := c.Transport.(*http.Transport)
	if !ok {
		t.Error("Client does not have the standard Transport")
	}

	if !transport.DisableKeepAlives {
		t.Error("Client did not disable keep alive")
	}

	if transport.TLSHandshakeTimeout != defaultHTTPClientTLSTimeout {
		t.Errorf("Unexpected TLS timeout value. Wanted %v, got %v\n", defaultHTTPClientTLSTimeout, transport.TLSHandshakeTimeout)
	}
}
