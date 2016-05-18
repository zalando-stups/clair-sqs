package user

import "testing"

func TestNoOpUserCredentials(t *testing.T) {
	var uc Credentials = new(noOpUserCredentials)

	if uc.Username() != "" {
		t.Error("Wrong username. Expected an empty username but got %q", uc.Username())
	}

	if uc.Password() != "" {
		t.Error("Wrong username. Expected an empty username but got %q", uc.Username())
	}
}
