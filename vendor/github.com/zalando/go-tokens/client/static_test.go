package client

import "testing"

func TestStaticClientCredentialsProvider(t *testing.T) {
	cp := NewStaticClientCredentialsProvider("foo", "bar")
	c, err := cp.Get()

	if err != nil {
		t.Fatal("Unexpected error creating the static client credentials provider")
	}

	if c.Id() != "foo" {
		t.Errorf("Unpexpected client id. Wanted 'foo' but got %q", c.Id())
	}

	if c.Secret() != "bar" {
		t.Errorf("Unpexpected client secret. Wanted 'bar' but got %q", c.Secret())
	}
}
