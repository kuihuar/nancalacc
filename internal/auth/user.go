package auth

import "context"

type UserAuthenticator struct {
}

func NewUserAuthenticator() *UserAuthenticator {
	return &UserAuthenticator{}
}

func (au *UserAuthenticator) GetAccessToken(ctx context.Context) (string, error) {
	return "", nil
}
