package tokens

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

const testToken = `{"access_token":"header.claims.sig","token_type":"Bearer","expires_in":4,` +
	`"scope":"uid","realm":"/services"}`

func TestManager(t *testing.T) {
	server, url := testServer(http.StatusOK)
	defer server.Close()

	os.Setenv("CREDENTIALS_DIR", "testdata")
	tr := NewPasswordRequest("test", "uid", "team")
	manager, err := Manage(
		url,
		[]ManagementRequest{tr},
	)
	if err != nil {
		t.Fatal(err)
	}

	at, err := manager.Get("test")
	if err != nil {
		t.Error(err)
	}

	at.validUntil = time.Now().Add(-1 * time.Hour)

	if _, err = manager.Get("test"); err != ErrTokenExpired {
		t.Error("Expected an ErrTokenExpired error for expired named token")
	}

	if _, err = manager.Get("missing"); err != ErrTokenNotAvailable {
		t.Error("Expected an ErrTokenNotAvailable error for missing named token")
	}
}

func testServer(status int) (*httptest.Server, string) {
	handler := func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(status)
		fmt.Fprint(w, testToken)
	}

	server := httptest.NewServer(http.HandlerFunc(handler))

	return server, fmt.Sprintf("http://%s", server.Listener.Addr())
}
