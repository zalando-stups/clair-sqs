package tokens

import (
	"fmt"
	"time"
)

type AccessToken struct {
	Token      string `json:"access_token"`
	ExpiresIn  int    `json:"expires_in"`
	issuedAt   time.Time
	validUntil time.Time
}

func (at *AccessToken) Expired() bool {
	return at.validUntil.Before(time.Now())
}

func (at *AccessToken) String() string {
	return fmt.Sprintf("%s expires in %d second(s)", at.Token, at.ExpiresIn)
}
