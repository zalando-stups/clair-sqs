package tokens

import (
	"fmt"
	"time"
)

// AccessToken type holds the minimum amount of data from an OAuth access token
type AccessToken struct {
	Token      string `json:"access_token"`
	ExpiresIn  int    `json:"expires_in"`
	issuedAt   time.Time
	validUntil time.Time
}

// Expired returns true if the access token is no longer valid. It depends on the local clock current time
func (at *AccessToken) Expired() bool {
	return at.validUntil.Before(time.Now())
}

// String implements the Stringer interface for pretty printing
func (at *AccessToken) String() string {
	return fmt.Sprintf("%s expires in %d second(s)", at.Token, at.ExpiresIn)
}
