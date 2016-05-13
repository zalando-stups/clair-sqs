package tokens

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestTokens(t *testing.T) {
	handler := func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"access_token":"header.claims.sig","token_type":"Bearer","expires_in":4,"scope":"uid","realm":"/services"}`)
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	url := fmt.Sprintf("http://%s", server.Listener.Addr())

	os.Setenv("CREDENTIALS_DIR", "testdata")
	tr := NewRequest("test", "password", "uid", "team")
	tks, err := Manage(
		url,
		[]ManagementRequest{tr},
	)
	if err != nil {
		t.Fatal(err)
	}

	_, err = tks.Get("test")
	if err != nil {
		t.Error(err)
	}
}
