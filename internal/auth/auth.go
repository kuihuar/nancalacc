package auth

import "context"

type Authenticator interface {
	GetAccessToken(ctx context.Context) (*AccessTokenResp, error)
}
