package tokens

import (
	"log"
)

type ManagementRequest struct {
	id        string
	grantType string
	scopes    []string
}

func NewRequest(id string, grantType string, scopes ...string) ManagementRequest {
	t := ManagementRequest{
		id:        id,
		grantType: grantType,
		scopes:    make([]string, 0, len(scopes)),
	}

	for _, scope := range scopes {
		if len(scope) < 1 {
			log.Println("Empty scope in token request dropped")
		} else {
			t.scopes = append(t.scopes, scope)
		}
	}
	return t
}
