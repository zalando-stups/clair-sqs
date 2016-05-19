package tokens

import (
	"log"
)

// ManagementRequest holds data for named token requests that can be passed to the token manager
type ManagementRequest struct {
	id        string
	grantType string
	scopes    []string
}

const passwordGrantType = "password"

// NewPasswordRequest returns a named token request which uses the specific password grant type.
// It accepts any amount of scopes
func NewPasswordRequest(id string, scopes ...string) ManagementRequest {
	return newRequest(id, passwordGrantType, scopes...)
}

func newRequest(id string, grantType string, scopes ...string) ManagementRequest {
	t := ManagementRequest{
		id:        id,
		grantType: grantType,
		scopes:    make([]string, 0, len(scopes)),
	}

	for _, scope := range scopes {
		if len(scope) < 1 {
			log.Printf("Empty scope in management request %q dropped\n", id)
		} else {
			t.scopes = append(t.scopes, scope)
		}
	}
	return t
}
