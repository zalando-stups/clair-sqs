package client

type Credentials interface {
	Id() string
	Secret() string
}

type CredentialsProvider interface {
	Get() (Credentials, error)
}
