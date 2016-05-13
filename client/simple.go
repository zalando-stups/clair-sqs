package client

type simpleUserCredentials struct {
}

func (suc *simpleUserCredentials) Username() string {
	return ""
}

func (suc *simpleUserCredentials) Password() string {
	return ""
}
