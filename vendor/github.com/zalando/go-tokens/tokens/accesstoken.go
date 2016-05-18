package tokens

import (
	"time"
)

/*
{
	"access_token":"header.claims.sig",
	"token_type":"Bearer",
	"expires_in":28800,
	"scope":"uid",
	"realm":"/services"
}
*/

type AccessToken struct {
	Token      string `json:"access_token"`
	ExpiresIn  int    `json:"expires_in"`
	issuedAt   time.Time
	validUntil time.Time
}

func (at *AccessToken) Expired() bool {
	return at.validUntil.Before(time.Now())
}
