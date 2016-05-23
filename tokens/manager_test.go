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

func TestOptions(t *testing.T) {

	if _, err := Manage("", []ManagementRequest{}); err != ErrMissingURL {
		t.Error("Expected an ErrMissingURL error for empty URL")
	}

	if _, err := Manage("dontcare", []ManagementRequest{}); err != ErrNoManagementRequests {
		t.Error("Expected an ErrNoManagementRequests error for empty management requests")
	}

	reqs := []ManagementRequest{NewPasswordRequest("test", "uid", "team")}
	if _, err := Manage("dontcare", reqs, RefreshPercentageThreshold(0.8), WarningPercentageThreshold(0.5)); err == nil {
		t.Error("Expected an error with initial invalid options")
	}

	server1, url := testServer(http.StatusInternalServerError)
	defer server1.Close()

	if _, err := Manage(url, reqs); err == nil {
		t.Error("Expected an error with initial invalid response from token endpoint")
	}

	server2, url := testServer(http.StatusOK)
	defer server2.Close()

	manager, err := Manage(
		url,
		reqs,
		RefreshPercentageThreshold(0.3),
		WarningPercentageThreshold(0.5),
	)

	if err != nil {
		t.Fatal(err)
	}

	if err = RefreshPercentageThreshold(0.4)(manager); err != nil {
		t.Error("Unable to set refresh threshold")
	}

	if err = RefreshPercentageThreshold(0.0)(manager); err == nil {
		t.Error("Expected an error setting a 0 percent refresh threshold")
	}

	if err = RefreshPercentageThreshold(1.0)(manager); err == nil {
		t.Error("Expected an error setting a refresh threshold higher or equal to 100 percent")
	}

	if err = RefreshPercentageThreshold(0.8)(manager); err == nil {
		t.Error("Expected an error setting a refresh threshold higher than the warning threshold")
	}

	if manager.tokenRefresher.refreshPercentageThreshold != 0.4 {
		t.Errorf("Unexpected refresh threshold. Wanted 0.4, got %f\n", manager.tokenRefresher.refreshPercentageThreshold)
	}

	if err = WarningPercentageThreshold(0.99)(manager); err != nil {
		t.Error("Unable to set warning threshold")
	}

	if err = WarningPercentageThreshold(0.10)(manager); err == nil {
		t.Error("Expected an error setting a warning threshold lower than the refresh threshold")
	}

	if manager.tokenRefresher.warningPercentageThreshold != 0.99 {
		t.Errorf("Unexpected warning threshold. Wanted 0.99, got %f\n", manager.tokenRefresher.warningPercentageThreshold)
	}
}

func testServer(status int) (*httptest.Server, string) {
	handler := func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(status)
		//if status == http.StatusOK {
		fmt.Fprint(w, testToken)
		//}
	}

	server := httptest.NewServer(http.HandlerFunc(handler))

	return server, fmt.Sprintf("http://%s", server.Listener.Addr())
}
