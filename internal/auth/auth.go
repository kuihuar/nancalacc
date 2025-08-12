package auth

import "context"

// Authenticator 认证器接口
type Authenticator interface {
	GetAccessToken(ctx context.Context) (*AccessTokenResp, error)
	InvalidateCache()
}

type AccessTokenResp struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
	// CreatedAt   time.Time `json:"created_at,omitempty"`
}
