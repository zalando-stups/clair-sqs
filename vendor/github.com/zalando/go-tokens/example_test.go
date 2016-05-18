package main

import (
	"fmt"
	"github.com/zalando/go-tokens/tokens"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
)

func ExampleTokens() {
	os.Setenv("CREDENTIALS_DIR", "tokens/testdata")
	handler := func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"access_token":"header.claims.sig","token_type":"Bearer","expires_in":4,"scope":"uid","realm":"/services"}`)
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	url := fmt.Sprintf("http://%s", server.Listener.Addr())
	tokens, err := tokens.Manage(url, []tokens.ManagementRequest{tokens.NewRequest("test", "password", "read")})
	if err != nil {
		log.Fatal(err)
	}

	at, err := tokens.Get("test")
	if err != nil {
		log.Println(err)
	}

	log.Println(at)
}
