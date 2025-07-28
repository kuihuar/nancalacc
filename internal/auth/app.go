package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"nancalacc/internal/conf"
	"nancalacc/pkg/httputil"
	stdurl "net/url"
	"strings"
	"time"

	"github.com/patrickmn/go-cache"
)

const (
	AppAuthType = "app"
)

type AppAuthenticator struct {
	clientId     string
	clientSecret string
	url          string
}

const (
	AppAuthPath = "/openapi/oauth2/token"
	grantType   = "client_credentials"
)

var (
	tokenCache = cache.New(3600*time.Minute, 7200*time.Minute)
)

// [POST] {配置域名}/openapi/oauth2/token

func NewAppAuthenticator(cfg *conf.Service) Authenticator {
	return &AppAuthenticator{
		clientId:     cfg.Auth.App.ClientId,
		clientSecret: cfg.Auth.App.ClientSecret,
		url:          cfg.Auth.App.AuthUrl,
	}
}

func (a *AppAuthenticator) GetAccessToken(ctx context.Context) (*AccessTokenResp, error) {
	clientId := a.clientId
	clientSecret := a.clientSecret

	url := a.url

	// if token, found := tokenCache.Get(clientId); found {
	// 	return token.(*AccessTokenResp), nil
	// }
	_, err := stdurl.Parse(url)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %v", err)
	}

	if !strings.Contains(url, "https") && !strings.Contains(url, "http") {
		return nil, fmt.Errorf("domain must be https or http")
	}
	if strings.Contains(url, "https") {
		clientSecret = MakeSECSecret(clientId, clientSecret, time.Now())
	}

	uri := fmt.Sprintf("%s%s", url, AppAuthPath)

	dataStr := fmt.Sprintf(`grant_type=%s&client_id=%s&client_secret=%s`, grantType, clientId, clientSecret)
	//fmt.Printf("Postpre uri: %s, dataStr: %s\n", uri, dataStr)
	data := []byte(dataStr)
	bs, err := httputil.Post(uri, data, 5*time.Second)
	//fmt.Printf("PostAfter err: %v,bs: %s\n", err, string(bs))
	if err != nil {
		return nil, err
	}
	var resp *AccessTokenResp
	err = json.Unmarshal(bs, &resp)
	if err != nil {
		return nil, err
	}
	tokenCache.Set(clientId, resp, time.Duration(resp.ExpiresIn)*time.Second)
	return resp, nil
}
