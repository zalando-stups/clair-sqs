package tokens

import (
	"errors"
	"github.com/zalando/go-tokens/client"
	"github.com/zalando/go-tokens/user"
	"os"
	"path"
)

// The Manager takes care of background refreshing of access tokens which can be requested
// at any time, in a thread safe way, for usage on your applications
type Manager struct {
	tokenRequests  []ManagementRequest
	tokenRefresher *refresher
	tokenHolder    *holder
}

var (
	// ErrMissingURL is returned whenever the manager receives an empty URL
	ErrMissingURL = errors.New("Missing OAuth2 Token URL")
	// ErrNoManagementRequests is returned whenever the manager is created without any management requests
	ErrNoManagementRequests = errors.New("No token requests")

	// ErrTokenNotAvailable is returned when a named token is not available
	ErrTokenNotAvailable = errors.New("no token available")
	// ErrTokenExpired is returned when a named token is found but has expired
	ErrTokenExpired = errors.New("token expired")
)

// Manage is the main function of the token manager. It accepts management requests that will be retrieved from
// the url parameter and, optionally, configured with a set of options.
// It loads the initial set of tokens synchronously and will fail if any of those requests also fail
func Manage(url string, requests []ManagementRequest, options ...func(*Manager) error) (*Manager, error) {
	if url == "" {
		return nil, ErrMissingURL
	}

	if len(requests) < 1 {
		return nil, ErrNoManagementRequests
	}

	userCredentialsFile := path.Join(os.Getenv("CREDENTIALS_DIR"), "user.json")
	clientCredentialsFile := path.Join(os.Getenv("CREDENTIALS_DIR"), "client.json")
	th := newHolder()
	t := &Manager{
		tokenRequests: requests,

		tokenRefresher: newRefresher(
			url,
			user.NewJSONFileUserCredentialsProvider(userCredentialsFile),
			client.NewJSONFileClientCredentialsProvider(clientCredentialsFile),
			th,
		),

		tokenHolder: th,
	}

	if len(options) > 0 {
		if err := t.setOption(options...); err != nil {
			t.Close()
			return nil, err
		}
	}

	tks, err := t.tokenRefresher.refreshTokens(requests)
	if err != nil {
		t.Close()
		return nil, err
	}

	t.tokenRefresher.start(requests, tks)

	return t, nil
}

func (t *Manager) setOption(options ...func(*Manager) error) error {
	for _, opt := range options {
		if err := opt(t); err != nil {
			return err
		}
	}
	return nil
}

// Get allows you to get a named token from the manager. It checks if a token has expired
func (t *Manager) Get(tokenID string) (*AccessToken, error) {
	at := t.tokenHolder.get(tokenID)
	if at == nil {
		return nil, ErrTokenNotAvailable
	}

	if at.Expired() {
		return nil, ErrTokenExpired
	}

	return at, nil
}

func (t *Manager) Close() {
	t.tokenRefresher.stop()
	t.tokenHolder.shutdown()
}
