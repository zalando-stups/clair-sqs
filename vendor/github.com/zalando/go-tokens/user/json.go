package user

import (
	"encoding/json"
	"io/ioutil"
)

type jsonFileUserCredentialsProvider struct {
	fileName string
}

func NewJSONFileUserCredentialsProvider(fileName string) CredentialsProvider {
	return &jsonFileUserCredentialsProvider{fileName}
}

func (cp *jsonFileUserCredentialsProvider) Get() (Credentials, error) {
	buf, err := ioutil.ReadFile(cp.fileName)
	if err != nil {
		return nil, err
	}

	var credentials jsonFileUserCredentials
	err = json.Unmarshal(buf, &credentials)
	if err != nil {
		return nil, err
	}
	return credentials, nil
}

type jsonFileUserCredentials struct {
	User string `json:"application_username"`
	Pass string `json:"application_password"`
}

func (uc jsonFileUserCredentials) Username() string {
	return uc.User
}

func (uc jsonFileUserCredentials) Password() string {
	return uc.Pass
}
