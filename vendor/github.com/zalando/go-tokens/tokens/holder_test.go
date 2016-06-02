package tokens

import (
	"testing"
)

func TestCaching(t *testing.T) {
	holder := newHolder()
	foo := &AccessToken{Token: "foo", ExpiresIn: 42}
	bar := &AccessToken{Token: "bar", ExpiresIn: 0}

	old := holder.set("foo", foo)
	if old != nil {
		t.Error("Set returned a non nil value for a new key")
	}

	v := holder.get("foo")
	if v == nil {
		t.Error("Failed to retrieve value for `foo` from holder")
	}

	v = holder.get("bar")
	if v != nil {
		t.Error("Retrieved value for `bar` from holder")
	}

	old = holder.set("foo", bar)
	if old != foo {
		t.Errorf("Old value was incorrect. Expected %v, got %v\n", foo, old)
	}

	holder.shutdown()
}

func TestPanic(t *testing.T) {
	defer func() {
		r := recover()
		err, ok := r.(error)
		if !ok {
			t.Error("Panic didn't contain an error")
		}
		if err.Error() != "Unknown operation: -1" {
			t.Errorf("Recovered from a different error: %v\n", err)
		}
	}()
	doOp(nil, &operation{op: -1})
}
