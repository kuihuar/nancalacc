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
)

const (
	AppAuthType = "app"
)

type WpsAppAuthenticator interface {
	Authenticator
}

type WpsAppAuth struct {
	clientId     string
	clientSecret string
	url          string
	cache        Cache
}

const (
	AppAuthPath = "/openapi/oauth2/token"
	grantType   = "client_credentials"
)

func NewWpsAppAuthenticator() WpsAppAuthenticator {
	cfg := conf.Get().GetAuth().GetWpsapp()
	return &WpsAppAuth{
		clientId:     cfg.ClientId,
		clientSecret: cfg.ClientSecret,
		url:          cfg.AuthUrl,
		cache:        NewLocalCache(),
	}
}

func (a *WpsAppAuth) GetAccessToken(ctx context.Context) (*AccessTokenResp, error) {
	// 尝试从缓存获取
	cacheKey := fmt.Sprintf("wps_token_%s", a.clientId)
	if cached, found := a.cache.Get(cacheKey); found {
		if token, ok := cached.(*AccessTokenResp); ok {
			return token, nil
		}
	}

	// 缓存中没有，从API获取
	token, err := a.getAccessTokenFromAPI(ctx)
	if err != nil {
		return nil, err
	}

	// 缓存token，提前5分钟过期
	cacheTTL := time.Duration(token.ExpiresIn-300) * time.Second
	if cacheTTL > 0 {
		a.cache.Set(cacheKey, token, cacheTTL)
	}

	return token, nil
}

func (a *WpsAppAuth) getAccessTokenFromAPI(ctx context.Context) (*AccessTokenResp, error) {
	clientId := a.clientId
	clientSecret := a.clientSecret
	url := a.url

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
	data := []byte(dataStr)
	bs, err := httputil.Post(uri, data, 5*time.Second)
	if err != nil {
		return nil, err
	}
	var resp *AccessTokenResp
	err = json.Unmarshal(bs, &resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// InvalidateCache 清除缓存
func (a *WpsAppAuth) InvalidateCache() {
	cacheKey := fmt.Sprintf("wps_token_%s", a.clientId)
	a.cache.Delete(cacheKey)
}
