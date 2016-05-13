package tokens

import (
	"errors"
	"github.com/zalando/go-tokens/client"
	"github.com/zalando/go-tokens/user"
	"os"
	"path"
)

type tokens struct {
	tokenRequests  []ManagementRequest
	tokenRefresher *refresher
	tokenHolder    *holder
}

var (
	ErrMissingUrl              = errors.New("Missing OAuth2 Token URL")
	ErrNoTokenRequests         = errors.New("No token requests")
	ErrInvalidRefreshThreshold = errors.New("Invalid refresh threshold")
	ErrInvalidWarningThreshold = errors.New("Invalid warning threshold")
)

func RefreshPercentageThreshold(threshold float64) func(*tokens) error {
	return func(t *tokens) error {
		if threshold <= 0 || threshold > 1 {
			return ErrInvalidRefreshThreshold
		}
		t.tokenRefresher.refreshPercentageThreshold = threshold
		return nil
	}
}

func WarningPercentageThreshold(threshold float64) func(*tokens) error {
	return func(t *tokens) error {
		if threshold <= 0 || threshold > 1 || threshold < t.tokenRefresher.refreshPercentageThreshold {
			return ErrInvalidWarningThreshold
		}
		t.tokenRefresher.warningPercentageThreshold = threshold
		return nil
	}
}

func Manage(url string, requests []ManagementRequest, options ...func(*tokens) error) (*tokens, error) {
	if url == "" {
		return nil, ErrMissingUrl
	}

	if len(requests) < 1 {
		return nil, ErrNoTokenRequests
	}

	userCredentialsFile := path.Join(os.Getenv("CREDENTIALS_DIR"), "user.json")
	clientCredentialsFile := path.Join(os.Getenv("CREDENTIALS_DIR"), "client.json")
	th := NewHolder()
	t := tokens{
		tokenRequests: requests,

		tokenRefresher: NewRefresher(
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

func (t *tokens) SetOption(options ...func(*tokens) error) error {
	for _, opt := range options {
		if err := opt(t); err != nil {
			return err
		}
	}
	return nil
}

var (
	ErrTokenNotAvailable = errors.New("no token available")
	ErrTokenExpired      = errors.New("token expired")
)

func (t *tokens) Get(tokenId string) (*AccessToken, error) {
	at := t.tokenHolder.get(tokenId)
	if at == nil {
		return nil, ErrTokenNotAvailable
	}

	if at.Expired() {
		return nil, ErrTokenExpired
	}

	return at, nil
}
