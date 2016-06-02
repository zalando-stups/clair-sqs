package tokens

import (
	"encoding/json"
	"github.com/kr/pretty"
	"reflect"
	"strings"
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

// 2009-11-10 23:00:00 UTC
func TestStringer(t *testing.T) {
	for _, test := range []struct {
		id         string
		validUntil time.Time
		want       string
	}{
		{"foo", time.Date(2009, 11, 10, 23, 0, 0, 0, time.UTC), "f...o valid until 2009-11-10 23:00:00 +0000 UTC"},
		{"bar", time.Date(2009, 11, 10, 23, 42, 0, 0, time.UTC), "b...r valid until 2009-11-10 23:42:00 +0000 UTC"},
		{"baz", time.Date(2009, 11, 10, 23, 0, 42, 0, time.UTC), "b...z valid until 2009-11-10 23:00:42 +0000 UTC"},
	} {
		at := &AccessToken{Token: test.id, validUntil: test.validUntil}
		got := at.String()
		if got != test.want {
			t.Errorf("Unexpected result. Wanted %q, got %q\n", test.want, got)
		}
	}
}

func TestUnmarshalling(t *testing.T) {
	for _, test := range []struct {
		payload   string
		want      *AccessToken
		wantError error
	}{
		{"%%", nil, nil},
	} {
		var got *AccessToken
		err := json.NewDecoder(strings.NewReader(test.payload)).Decode(got)
		pretty.Println(err, got)
		if test.wantError != nil && test.wantError != err {
			t.Errorf("Unexpected error condition. Wanted %v but got %v\n", test.wantError, err)
		} else {
			if !reflect.DeepEqual(test.want, got) {
				t.Errorf("Unexpected access token. Wanted %v but got %v\n", test.want, got)
			}
		}
	}
}
