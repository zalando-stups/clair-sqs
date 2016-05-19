/*
Package user holds the interfaces for the User credentials
*/
package user

// Credentials is the interface for any type that is able to return a username and password
type Credentials interface {
	Username() string
	Password() string
}

// CredentialsProvider is the interface for any type that is able to return Credentials
type CredentialsProvider interface {
	Get() (Credentials, error)
}
