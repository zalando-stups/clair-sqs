package tokens

import (
	"errors"
	"github.com/zalando/go-tokens/client"
	"github.com/zalando/go-tokens/user"
)

var (
	// ErrInvalidRefreshThreshold is returned if an invalid refresh threshold is used. It should be between 0.0 and 1.0
	ErrInvalidRefreshThreshold = errors.New("Invalid refresh threshold")
	// ErrInvalidWarningThreshold is returned if an invalid refresh threshold is used. It should be between 0.0 and 1.0
	ErrInvalidWarningThreshold = errors.New("Invalid warning threshold")

	// ErrInvalidUserCredentialsProvider is returned if there is an attempt to set a nil credentials provider
	ErrInvalidUserCredentialsProvider = errors.New("Invalid user credentials provider")
	// ErrInvalidClientCredentialsProvider is returned if there is an attempt to set a nil credentials provider
	ErrInvalidClientCredentialsProvider = errors.New("Invalid client credentials provider")
)

// RefreshPercentageThreshold returns a function that can set the refresh threshold on a tokensManager
func RefreshPercentageThreshold(threshold float64) func(*Manager) error {
	return func(t *Manager) error {
		if threshold <= 0 || threshold >= 1 || threshold > t.tokenRefresher.warningPercentageThreshold {
			return ErrInvalidRefreshThreshold
		}
		t.tokenRefresher.refreshPercentageThreshold = threshold
		return nil
	}
}

// WarningPercentageThreshold returns a function that can set the warning threshold on a tokensManager. It should be
// higher than the refresh threshold
func WarningPercentageThreshold(threshold float64) func(*Manager) error {
	return func(t *Manager) error {
		if threshold <= 0 || threshold > 1 || threshold < t.tokenRefresher.refreshPercentageThreshold {
			return ErrInvalidWarningThreshold
		}
		t.tokenRefresher.warningPercentageThreshold = threshold
		return nil
	}
}

// UserCredentialsProvider returns a function that changes the User credentials provider for a token manager
func UserCredentialsProvider(ucp user.CredentialsProvider) func(*Manager) error {
	return func(t *Manager) error {
		if ucp == nil {
			return ErrInvalidUserCredentialsProvider
		}
		t.tokenRefresher.userCredentialsProvider = ucp
		return nil
	}
}

// ClientCredentialsProvider returns a function that changes the Client credentials provider for a token manager
func ClientCredentialsProvider(ccp client.CredentialsProvider) func(*Manager) error {
	return func(t *Manager) error {
		if ccp == nil {
			return ErrInvalidUserCredentialsProvider
		}
		t.tokenRefresher.clientCredentialsProvider = ccp
		return nil
	}
}
