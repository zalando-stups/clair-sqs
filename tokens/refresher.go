package tokens

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/zalando/go-tokens/client"
	"github.com/zalando/go-tokens/httpclient"
	"github.com/zalando/go-tokens/user"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type refresher struct {
	url                       string
	httpClient                *http.Client
	userCredentialsProvider   user.CredentialsProvider
	clientCredentialsProvider client.CredentialsProvider

	tokenHolder      *holder
	refreshScheduler *scheduler

	refreshPercentageThreshold float64
	warningPercentageThreshold float64
}

const (
	defaultRefreshPercentageThreshold = 0.6
	defaultWarningPercentageThreshold = 0.8
	retryDelay                        = 10 * time.Second
)

func newRefresher(url string, ucp user.CredentialsProvider, ccp client.CredentialsProvider, h *holder) *refresher {
	r := &refresher{
		httpClient: httpclient.Default(),
		url:        url,
		userCredentialsProvider:   ucp,
		clientCredentialsProvider: ccp,

		refreshPercentageThreshold: defaultRefreshPercentageThreshold,
		warningPercentageThreshold: defaultWarningPercentageThreshold,

		tokenHolder: h,
	}
	r.refreshScheduler = newScheduler(r.refreshToken)
	return r
}

func (r *refresher) refreshTokens(requests []ManagementRequest) error {
	for _, tokenRequest := range requests {
		if err := r.doRefreshToken(tokenRequest); err != nil {
			return err
		}
	}
	return nil
}

func (r *refresher) doRefreshToken(tr ManagementRequest) error {
	uc, err := r.userCredentialsProvider.Get()
	if err != nil {
		return err
	}

	cc, err := r.clientCredentialsProvider.Get()
	if err != nil {
		return err
	}

	c := make(url.Values)
	c.Set("grant_type", tr.grantType)
	c.Set("scope", strings.Join(tr.scopes, " "))
	if tr.grantType == "password" {
		c.Set("username", uc.Username())
		c.Set("password", uc.Password())
	}

	req, err := http.NewRequest("POST", r.url, strings.NewReader(c.Encode()))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", httpclient.UserAgent)
	req.Header.Set("Authorization", "Basic "+basicAuth(cc.Id(), cc.Secret()))

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("Error getting token: %d - %v", resp.StatusCode, resp.Body)
	}

	buf, err := ioutil.ReadAll(resp.Body)

	at := new(AccessToken)
	if err = json.Unmarshal(buf, at); err != nil {
		return fmt.Errorf("Invalid token response: %v", err)
	}

	at.issuedAt = time.Now().Add(-1 * time.Second)
	at.validUntil = at.issuedAt.Add(time.Duration(at.ExpiresIn) * time.Second)
	delta := float64(at.ExpiresIn) * r.refreshPercentageThreshold
	r.refreshScheduler.scheduleTokenRefresh(tr, time.Duration(int64(delta))*time.Second)

	r.tokenHolder.set(tr.id, at)
	return nil
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

// This is the callback function the scheduler will run when the timer expires
func (r *refresher) refreshToken(tr ManagementRequest) {
	if err := r.doRefreshToken(tr); err != nil {
		r.refreshScheduler.scheduleTokenRefresh(tr, retryDelay)
	}
}
