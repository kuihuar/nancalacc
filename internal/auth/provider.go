package auth

import (
	"github.com/google/wire"
)

var AuthProviderSet = wire.NewSet(NewAppAuthenticator, NewUserAuthenticator, NewLocalCachedAuthenticator)
