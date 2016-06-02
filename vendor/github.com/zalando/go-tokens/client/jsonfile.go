package client

import (
	"encoding/json"
	"io/ioutil"
)

type provider struct {
	path string
}

// NewJSONFileClientCredentialsProvider returns a client.CredentialsProvider that reads both client id and
// secret from a JSON file stored in the specified filesystem path
// The contents of such file should follow the following specifications:
//		{"client_id":"foo","client_secret":"bar"}
func NewJSONFileClientCredentialsProvider(path string) CredentialsProvider {
	return &provider{path}
}

func (cp *provider) Get() (Credentials, error) {
	buf, err := ioutil.ReadFile(cp.path)
	if err != nil {
		return nil, err
	}

	var credentials credentials
	err = json.Unmarshal(buf, &credentials)
	if err != nil {
		return nil, err
	}
	return credentials, nil
}

type credentials struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

func (uc credentials) Id() string {
	return uc.ClientId
}

func (uc credentials) Secret() string {
	return uc.ClientSecret
}
