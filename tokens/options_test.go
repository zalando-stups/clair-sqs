package tokens

import (
	"github.com/zalando/go-tokens/client"
	"github.com/zalando/go-tokens/user"
	"net/http"
	"testing"
)

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
		UserCredentialsProvider(user.NewStaticUserCredentialsProvider("user", "pw")),
		ClientCredentialsProvider(client.NewStaticClientCredentialsProvider("id", "secret")),
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

	if err = UserCredentialsProvider(nil)(manager); err == nil {
		t.Error("Manager should not have accepted a nil user credentials provider")
	}

	if err = ClientCredentialsProvider(nil)(manager); err == nil {
		t.Error("Manager should not have accepted a nil client credentials provider")
	}
}
