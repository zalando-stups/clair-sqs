package user

import "testing"

func TestStaticUserCredentialsProvider(t *testing.T) {
	cp := NewStaticUserCredentialsProvider("foo", "bar")
	c, err := cp.Get()

	if err != nil {
		t.Fatal("Unexpected error creating the static user credentials provider")
	}

	if c.Username() != "foo" {
		t.Errorf("Unpexpected username. Wanted 'foo' but got %q", c.Username())
	}

	if c.Password() != "bar" {
		t.Errorf("Unpexpected password. Wanted 'bar' but got %q", c.Password())
	}
}
