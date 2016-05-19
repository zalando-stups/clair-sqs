package tokens

import (
	"errors"
	"github.com/zalando/go-tokens/client"
	"github.com/zalando/go-tokens/user"
	"os"
	"path"
)

type tokensManager struct {
	tokenRequests  []ManagementRequest
	tokenRefresher *refresher
	tokenHolder    *holder
}

var (
	// ErrMissingURL is returned whenever the manager receives an empty URL
	ErrMissingURL = errors.New("Missing OAuth2 Token URL")
	// ErrNoManagementRequests is returned whenever the manager is created without any management requests
	ErrNoManagementRequests = errors.New("No token requests")
	// ErrInvalidRefreshThreshold is returned if an invalid refresh threshold is used. It should be between 0.0 and 1.0
	ErrInvalidRefreshThreshold = errors.New("Invalid refresh threshold")
	// ErrInvalidWarningThreshold is returned if an invalid refresh threshold is used. It should be between 0.0 and 1.0
	ErrInvalidWarningThreshold = errors.New("Invalid warning threshold")

	// ErrTokenNotAvailable is returned when a named token is not available
	ErrTokenNotAvailable = errors.New("no token available")
	// ErrTokenExpired is returned when a named token is found but has expired
	ErrTokenExpired = errors.New("token expired")
)

// RefreshPercentageThreshold returns a function that can set the refresh threshold on a tokensManager
func RefreshPercentageThreshold(threshold float64) func(*tokensManager) error {
	return func(t *tokensManager) error {
		if threshold <= 0 || threshold > 1 {
			return ErrInvalidRefreshThreshold
		}
		t.tokenRefresher.refreshPercentageThreshold = threshold
		return nil
	}
}

// WarningPercentageThreshold returns a function that can set the warning threshold on a tokensManager
func WarningPercentageThreshold(threshold float64) func(*tokensManager) error {
	return func(t *tokensManager) error {
		if threshold <= 0 || threshold > 1 || threshold < t.tokenRefresher.refreshPercentageThreshold {
			return ErrInvalidWarningThreshold
		}
		t.tokenRefresher.warningPercentageThreshold = threshold
		return nil
	}
}

// Manage is the main function of the token manager. It accepts management requests that will be retrieved from
// the url parameter and, optionally, configured with a set of options.
// It loads the initial set of tokens synchronously and will fail if any of those requests also fail
func Manage(url string, requests []ManagementRequest, options ...func(*tokensManager) error) (*tokensManager, error) {
	if url == "" {
		return nil, ErrMissingURL
	}

	if len(requests) < 1 {
		return nil, ErrNoManagementRequests
	}

	userCredentialsFile := path.Join(os.Getenv("CREDENTIALS_DIR"), "user.json")
	clientCredentialsFile := path.Join(os.Getenv("CREDENTIALS_DIR"), "client.json")
	th := newHolder()
	t := tokensManager{
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
		if err := t.SetOption(options...); err != nil {
			return nil, err
		}
	}

	if err := t.tokenRefresher.refreshTokens(requests); err != nil {
		return nil, err
	}

	return &t, nil
}

func (t *tokensManager) SetOption(options ...func(*tokensManager) error) error {
	for _, opt := range options {
		if err := opt(t); err != nil {
			return err
		}
	}
	return nil
}

// Get allows you to get a named token from the manager. It checks if a token has expired
func (t *tokensManager) Get(tokenID string) (*AccessToken, error) {
	at := t.tokenHolder.get(tokenID)
	if at == nil {
		return nil, ErrTokenNotAvailable
	}

	if at.Expired() {
		return nil, ErrTokenExpired
	}

	return at, nil
}
