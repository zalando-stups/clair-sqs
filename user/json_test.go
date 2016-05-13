package user

import (
	"testing"
)

func TestWorkingJSONCredentials(t *testing.T) {
	cp := NewJSONFileUserCredentialsProvider("testdata/user.json")
	c, err := cp.Get()

	if err != nil {
		t.Fatalf("Failed to get JSON file credentials: %v", err)
	}

	if c.Username() != "go-tokens-user" {
		t.Errorf("Unexpected username. Wanted \"go-tokens-user\", got %q", c.Username())
	}

	if c.Password() != "fake-password" {
		t.Errorf("Unexpected password. Wanted \"fake-password\", got %q", c.Password())
	}
}

func TestInvalidCredentialsFile(t *testing.T) {
	for _, test := range []struct {
		fileName string
	}{
		{"missing.json"},
		{"testdata/broken.json"},
	} {
		cp := NewJSONFileUserCredentialsProvider(test.fileName)
		_, err := cp.Get()
		if err == nil {
			t.Error("Expected an error for invalid credentials file")
		}
	}
}
