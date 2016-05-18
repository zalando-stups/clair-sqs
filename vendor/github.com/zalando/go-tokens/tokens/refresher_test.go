package tokens

import (
	"fmt"
	"github.com/zalando/go-tokens/client"
	"github.com/zalando/go-tokens/user"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewRefresher(t *testing.T) {
	handler := func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"access_token":"header.claims.sig","token_type":"Bearer","expires_in":4,"scope":"uid","realm":"/services"}`)
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	url := fmt.Sprintf("http://%s", server.Listener.Addr())
	th := NewHolder()
	r := NewRefresher(
		url,
		user.NewJSONFileUserCredentialsProvider("testdata/user.json"),
		client.NewJSONFileClientCredentialsProvider("testdata/client.json"),
		th,
	)
	tr := NewRequest("test", "password", "uid", "team")
	err := r.doRefreshToken(tr)
	if err != nil {
		t.Error(err)
	}

	at := th.get("test")
	if at == nil {
		t.Fatal("Failed to get token 'test' from the token holder")
	}

	if at.Token != "header.claims.sig" {
		t.Error(`Invalid token. Wanted "header.claims.sig", got %q`, at.Token)
	}

	if at.ExpiresIn != 4 {
		t.Error(`Invalid expiration time. Wanted 4, got %d`, at.ExpiresIn)
	}
}
