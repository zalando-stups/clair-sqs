package client

type static struct {
	id     string
	secret string
}

// NewStaticClientCredentialsProvider returns a client.CredentialsProvider that returns the id and secret
// used in the arguments clientId and clientSecret, respectively
func NewStaticClientCredentialsProvider(clientID string, clientSecret string) CredentialsProvider {
	return &static{id: clientID, secret: clientSecret}
}

func (cp *static) Get() (Credentials, error) {
	return cp, nil
}

func (cp *static) Id() string {
	return cp.id
}

func (cp *static) Secret() string {
	return cp.secret
}
