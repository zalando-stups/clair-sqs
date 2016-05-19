package tokens

import "testing"

func TestRequest(t *testing.T) {
	r := NewPasswordRequest("example", "read", "write", "")

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
