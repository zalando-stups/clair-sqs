package tokens

import (
	"fmt"
	"github.com/zalando/go-tokens/client"
	"github.com/zalando/go-tokens/user"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRefresher(t *testing.T) {
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
	tr := NewPasswordRequest("test", "uid", "team")
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

func TestRefresherFailure(t *testing.T) {
	handler := func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	url := fmt.Sprintf("http://%s", server.Listener.Addr())
	th := NewHolder()
	for _, test := range []struct {
		u string
		ucp user.CredentialsProvider
		ccp client.CredentialsProvider
	}{
		{
			u: url,
			ucp: user.NewJSONFileUserCredentialsProvider("testdata/user.json"),
			ccp: client.NewJSONFileClientCredentialsProvider("testdata/client.json"),
		},
		{
			u: url,
			ucp: user.NewJSONFileUserCredentialsProvider("missing-file.json"),
			ccp: client.NewJSONFileClientCredentialsProvider("testdata/client.json"),
		},
		{
			u: url,
			ucp: user.NewJSONFileUserCredentialsProvider("testdata/user.json"),
			ccp: client.NewJSONFileClientCredentialsProvider("missing-file.json"),
		},
		{
			u: "http://192.168.0.%31/",
			ucp: user.NewJSONFileUserCredentialsProvider("testdata/user.json"),
			ccp: client.NewJSONFileClientCredentialsProvider("testdata/client.json"),
		},
	} {
		r := NewRefresher(test.u, test.ucp, test.ccp, th)

		err := r.refreshTokens([]ManagementRequest{NewPasswordRequest("test", "uid", "team")})
		if err == nil {
			t.Error("Refresh should have failed")
		}

	}
}
