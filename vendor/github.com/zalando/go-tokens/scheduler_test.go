package tokens

import (
	"testing"
	"time"
)

func TestScheduling(t *testing.T) {
	var at *AccessToken
	s := NewScheduler(func(mr ManagementRequest) {
		at = &AccessToken{Token: "foo", ExpiresIn: 42}
	})

	runner = func(d time.Duration, f func()) *time.Timer {
		f()
		return nil
	}

	mr := ManagementRequest{id: "bar"}
	if err := s.scheduleTokenRefresh(mr, 0); err != nil {
		t.Error(err)
	}

	if at == nil {
		t.Error("Failed to refresh the token")
	}

	if at.Token != "foo" || at.ExpiresIn != 42 {
		t.Error("Wrong token received from callback")
	}

	s.Stop()
}

func TestRescheduleFailure(t *testing.T) {
	var at *AccessToken
	s := NewScheduler(func(mr ManagementRequest) {
		at = &AccessToken{Token: "foo", ExpiresIn: 42}
	})

	runner = func(d time.Duration, f func()) *time.Timer {
		// Don't execute anything
		return nil
	}
	mr := ManagementRequest{id: "bar"}
	if err := s.scheduleTokenRefresh(mr, 0); err != nil {
		t.Error(err)
	}

	if err := s.scheduleTokenRefresh(mr, 0); err == nil {
		t.Error("Rescheduling should have failed")
	}
}
