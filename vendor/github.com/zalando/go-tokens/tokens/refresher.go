package tokens

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/zalando/go-tokens/client"
	"github.com/zalando/go-tokens/httpclient"
	"github.com/zalando/go-tokens/user"
	"log"
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

func (r *refresher) refreshTokens(requests []ManagementRequest) (map[string]*AccessToken, error) {
	tks := make(map[string]*AccessToken)
	for _, tokenRequest := range requests {
		if at, err := r.doRefreshToken(tokenRequest); err != nil {
			return nil, err
		} else {
			tks[tokenRequest.id] = at
		}
	}
	return tks, nil
}

func (r *refresher) doRefreshToken(tr ManagementRequest) (*AccessToken, error) {
	uc, err := r.userCredentialsProvider.Get()
	if err != nil {
		return nil, err
	}

	cc, err := r.clientCredentialsProvider.Get()
	if err != nil {
		return nil, err
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
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", httpclient.UserAgent)
	req.Header.Set("Authorization", "Basic "+basicAuth(cc.Id(), cc.Secret()))

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("Error getting token: %d - %v", resp.StatusCode, resp.Body)
	}

	at := new(AccessToken)
	if err = json.NewDecoder(resp.Body).Decode(at); err != nil {
		return nil, fmt.Errorf("Invalid token response: %v", err)
	}

	r.tokenHolder.set(tr.id, at)
	return at, nil
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

// This is the callback function the scheduler will run when the timer expires
func (r *refresher) refreshToken(tr ManagementRequest) {
	var d = retryDelay
	if at, err := r.doRefreshToken(tr); err == nil {
		d = at.RefreshIn(r.refreshPercentageThreshold)
	}
	if err := r.refreshScheduler.scheduleTokenRefresh(tr, d); err != nil {
		log.Println(err)
	}
}

func (r *refresher) start(requests []ManagementRequest, tks map[string]*AccessToken) {
	for _, req := range requests {
		at := tks[req.id]
		delta := at.RefreshIn(r.refreshPercentageThreshold)
		if err := r.refreshScheduler.scheduleTokenRefresh(req, delta); err != nil {
			log.Println(err)
		}
	}
}

func (r *refresher) stop() {
	r.refreshScheduler.stop()
}
