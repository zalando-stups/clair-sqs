package user

import (
	"encoding/json"
	"io/ioutil"
)

type provider struct {
	path string
}

// NewJSONFileUserCredentialsProvider returns a user.CredentialsProvider that reads both username and
// password from a JSON file stored in the specified filesystem path.
// The contents of such file should follow the following specifications:
//		{"application_username":"foo","application_password":"bar"}
func NewJSONFileUserCredentialsProvider(path string) CredentialsProvider {
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
	User string `json:"application_username"`
	Pass string `json:"application_password"`
}

func (uc credentials) Username() string {
	return uc.User
}

func (uc credentials) Password() string {
	return uc.Pass
}
