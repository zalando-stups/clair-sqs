package tokens

import "testing"

func TestRequest(t *testing.T) {
	r := NewRequest("example", "type", "read", "write")

	if r.id != "example" {
		t.Errorf("Invalid request Id. Wanted %q, got %q", "example", r.id)
	}

	if len(r.scopes) != 2 {
		t.Errorf("Invalid scopes count. Wanted 2, got %d", len(r.scopes))
	}

	s1 := r.scopes[0]
	if s1 != "read" {
		t.Errorf("Invalid scope. Wanted %q, got %q", "read", s1)
	}

	s2 := r.scopes[1]
	if s2 != "write" {
		t.Errorf("Invalid scope. Wanted %q, got %q", "write", s2)
	}
}

//func TestInvalidRequests(t *testing.T) {
//	for _, test := range []struct {
//		id string
//		grantType string
//		scopes []string
//		wanted error
//	}{
//		{"", "", nil, ErrInvalidRequestId},
//		{"foo", "", nil, ErrInvalidGrantType},
//		{"foo", "bar", nil, ErrNoScopes},
//		{"foo", "bar", []string{}, ErrNoScopes},
//		{"foo", "bar", []string{""}, ErrInvalidScope},
//	} {
//		_, err := New(test.id, test.scopes...)
//		if err != test.wanted {
//			t.Error("Invalid result. Wamted %v, got %v", test.wanted, err)
//		}
//	}
//
//}
