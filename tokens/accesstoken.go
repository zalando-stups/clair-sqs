package tokens

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"time"
)

// AccessToken type holds the minimum amount of data from an OAuth access token
type AccessToken struct {
	Token      string `json:"access_token"`
	ExpiresIn  int    `json:"expires_in"`
	issuedAt   time.Time
	validUntil time.Time
}

var (
	ErrMissingToken      = errors.New("missing token")
	ErrInvalidToken      = errors.New("invalid token string")
	ErrMissingExpiration = errors.New("missing expiration")
	ErrInvalidExpiration = errors.New("invalid expiration integer value")
)

// Expired returns true if the access token is no longer valid. It depends on the local clock current time
func (at *AccessToken) Expired() bool {
	return at.validUntil.Before(time.Now())
}

// String implements the Stringer interface for pretty printing
func (at *AccessToken) String() string {
	obfuscated := at.Token
	l := len(at.Token)
	if l > 2 {
		chunkSize := int(math.Max(float64(l)/10.0, 1.0))
		obfuscated = at.Token[:chunkSize] + "..." + at.Token[len(at.Token)-chunkSize:]
	}
	return fmt.Sprintf("%s valid until %v", obfuscated, at.validUntil)
}

func (at *AccessToken) RefreshIn(threshold float64) time.Duration {
	delta := float64(at.ExpiresIn) * threshold
	return time.Duration(int64(delta)) * time.Second
}

// UnmarshalJSON is used to unmarshal an AccessToken entry from the input bytes and also adding
// the extra validity fields
func (at *AccessToken) UnmarshalJSON(data []byte) error {
	var buf map[string]interface{}
	if err := json.Unmarshal(data, &buf); err != nil {
		return err
	}

	t, has := buf["access_token"]
	if !has {
		return ErrMissingToken
	}

	accessToken, ok := t.(string)
	if !ok {
		return ErrInvalidToken
	}

	e, has := buf["expires_in"]
	if !has {
		return ErrMissingExpiration
	}

	expiresIn, ok := e.(float64)
	if !ok {
		return ErrInvalidExpiration
	}

	issueTime := time.Now().Add(-1 * time.Second)

	*at = AccessToken{
		Token:      accessToken,
		ExpiresIn:  int(expiresIn),
		issuedAt:   issueTime,
		validUntil: issueTime.Add(time.Duration(expiresIn) * time.Second),
	}

	return nil
}
