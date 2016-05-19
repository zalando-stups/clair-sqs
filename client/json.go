package client

import (
	"encoding/json"
	"io/ioutil"
)

type jsonFileClientCredentialsProvider struct {
	fileName string
}

// NewJSONFileClientCredentialsProvider returns a client.CredentialsProvider that reads both client id and
// secret from a JSON file stored in the filesystem
func NewJSONFileClientCredentialsProvider(fileName string) CredentialsProvider {
	return &jsonFileClientCredentialsProvider{fileName}
}

func (cp *jsonFileClientCredentialsProvider) Get() (Credentials, error) {
	buf, err := ioutil.ReadFile(cp.fileName)
	if err != nil {
		return nil, err
	}

	var credentials jsonFileClientCredentials
	err = json.Unmarshal(buf, &credentials)
	if err != nil {
		return nil, err
	}
	return credentials, nil
}

type jsonFileClientCredentials struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

func (uc jsonFileClientCredentials) Id() string {
	return uc.ClientId
}

func (uc jsonFileClientCredentials) Secret() string {
	return uc.ClientSecret
}
