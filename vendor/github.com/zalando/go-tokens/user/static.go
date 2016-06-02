package user

type static struct {
	username string
	password string
}

// NewStaticUserCredentialsProvider returns a user.CredentialsProvider that returns the username and password
// used in the arguments u and p, respectively
func NewStaticUserCredentialsProvider(u string, p string) CredentialsProvider {
	return &static{username: u, password: p}
}

func (cp *static) Get() (Credentials, error) {
	return cp, nil
}

func (cp *static) Username() string {
	return cp.username
}

func (cp *static) Password() string {
	return cp.password
}
