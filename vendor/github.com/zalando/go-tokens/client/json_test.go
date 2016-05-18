package client

import (
	"testing"
)

func TestWorkingJSONCredentials(t *testing.T) {
	cp := NewJSONFileClientCredentialsProvider("testdata/client.json")
	c, err := cp.Get()

	if err != nil {
		t.Fatalf("Failed to get JSON file credentials: %v", err)
	}

	if c.Id() != "go-tokens-client" {
		t.Errorf("Unexpected client-id. Wanted \"go-tokens-client\", got %q", c.Id())
	}

	if c.Secret() != "fake-secret" {
		t.Errorf("Unexpected secret. Wanted \"fake-secret\", got %q", c.Secret())
	}
}

func TestInvalidCredentialsFile(t *testing.T) {
	for _, test := range []struct {
		fileName string
	}{
		{"missing.json"},
		{"testdata/broken.json"},
	} {
		cp := NewJSONFileClientCredentialsProvider(test.fileName)
		_, err := cp.Get()
		if err == nil {
			t.Error("Expected an error for invalid credentials file")
		}
	}
}
