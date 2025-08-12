package auth

import (
	"context"
	"fmt"
	"nancalacc/internal/conf"
	"time"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dingtalkoauth2_1_0 "github.com/alibabacloud-go/dingtalk/oauth2_1_0"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
)

const (
	DingtalkAuthType = "dingtalk"
)

type DingTalkAuthenticator interface {
	Authenticator
}

type DingTalkAuth struct {
	AppKey    string
	AppSecret string
	Endpoint  string
	//Timeout     string
	dingtalkCli *dingtalkoauth2_1_0.Client
	cache       Cache
}

func NewDingTalkAuthenticator() DingTalkAuthenticator {

	cfg := conf.Get().GetAuth().GetDingtalk()
	config := &openapi.Config{
		Protocol: tea.String("https"),
		RegionId: tea.String("central"),
	}
	client, err := dingtalkoauth2_1_0.NewClient(config)
	if err != nil {
		fmt.Printf("NewClient err: %v", err)
	}
	return &DingTalkAuth{
		Endpoint:    cfg.Endpoint,
		AppKey:      cfg.AppKey,
		AppSecret:   cfg.AppSecret,
		dingtalkCli: client,
		cache:       NewLocalCache(),
	}
}

func (r *DingTalkAuth) GetAccessToken(ctx context.Context) (*AccessTokenResp, error) {
	// 尝试从缓存获取
	cacheKey := fmt.Sprintf("dingtalk_token_%s", r.AppKey)
	if cached, found := r.cache.Get(cacheKey); found {
		if token, ok := cached.(*AccessTokenResp); ok {
			return token, nil
		}
	}

	// 缓存中没有，从API获取
	token, err := r.getAccessTokenFromAPI(ctx)
	if err != nil {
		return nil, err
	}

	// 缓存token，提前5分钟过期
	cacheTTL := time.Duration(token.ExpiresIn-300) * time.Second
	if cacheTTL > 0 {
		r.cache.Set(cacheKey, token, cacheTTL)
	}

	return token, nil
}

func (r *DingTalkAuth) getAccessTokenFromAPI(ctx context.Context) (*AccessTokenResp, error) {
	request := &dingtalkoauth2_1_0.GetAccessTokenRequest{
		AppKey:    tea.String(r.AppKey),
		AppSecret: tea.String(r.AppSecret),
	}

	var res *AccessTokenResp
	var accessToken dingtalkoauth2_1_0.GetAccessTokenResponseBody

	tryErr := func() error {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				err := r
				fmt.Printf("恢复的错误: %v\n", err)
			}
		}()

		response, err := r.dingtalkCli.GetAccessToken(request)
		if err != nil {
			return err
		}

		accessToken = *response.Body
		return nil
	}()

	if tryErr != nil {
		// 处理错误
		var sdkErr = &tea.SDKError{}
		if _t, ok := tryErr.(*tea.SDKError); ok {
			sdkErr = _t
		} else {
			sdkErr.Message = tea.String(tryErr.Error())
		}

		if !tea.BoolValue(util.Empty(sdkErr.Code)) && !tea.BoolValue(util.Empty(sdkErr.Message)) {
			return res, fmt.Errorf("获取access_token失败: [%s] %s", *sdkErr.Code, *sdkErr.Message)
		}
		return res, fmt.Errorf("获取access_token失败: %s", *sdkErr.Message)
	}

	return &AccessTokenResp{
		AccessToken: *accessToken.AccessToken,
		ExpiresIn:   int(*accessToken.ExpireIn),
	}, nil
}

// InvalidateCache 清除缓存
func (r *DingTalkAuth) InvalidateCache() {
	cacheKey := fmt.Sprintf("dingtalk_token_%s", r.AppKey)
	r.cache.Delete(cacheKey)
}
