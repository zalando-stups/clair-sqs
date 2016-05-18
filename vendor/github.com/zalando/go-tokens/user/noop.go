package user

type noOpUserCredentials int

func (_ *noOpUserCredentials) Username() string {
	return ""
}

func (_ *noOpUserCredentials) Password() string {
	return ""
}
