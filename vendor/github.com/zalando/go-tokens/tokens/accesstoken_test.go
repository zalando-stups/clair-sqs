package tokens

import (
	"testing"
	"time"
)

func TestAccessToken(t *testing.T) {
	for _, test := range []struct {
		validUntil time.Time
		want       bool
	}{
		{time.Now().Add(8 * time.Hour), false},
		{time.Now().Add(-10 * time.Second), true},
	} {
		at := &AccessToken{Token: "foo", validUntil: test.validUntil}
		got := at.Expired()
		if got != test.want {
			t.Errorf("Unexpected expiration status. Wanted %v, got %v\n", test.want, got)
		}
	}
}

func TestStringer(t *testing.T) {
	for _, test := range []struct {
		id         string
		expiration int
		want       string
	}{
		{"foo", 0, "foo expires in 0 second(s)"},
		{"bar", 0, "bar expires in 0 second(s)"},
		{"baz", 42, "baz expires in 42 second(s)"},
	} {
		at := &AccessToken{Token: test.id, ExpiresIn: test.expiration}
		got := at.String()
		if got != test.want {
			t.Errorf("Unexpected result. Wanted %q, got %q\n", test.want, got)
		}
	}

}
