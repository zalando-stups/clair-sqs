package user

type Credentials interface {
	Username() string
	Password() string
}

type CredentialsProvider interface {
	Get() (Credentials, error)
}
