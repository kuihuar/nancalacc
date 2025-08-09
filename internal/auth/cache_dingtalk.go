package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/patrickmn/go-cache"
)

type DingtalkCacheConfig struct {
	ttl             time.Duration
	cleanupInterval time.Duration
	cacheKey        string
}

// 实现通用配置接口
func (d *DingtalkCacheConfig) SetTTL(ttl time.Duration) {
	d.ttl = ttl
}

func (d *DingtalkCacheConfig) SetCleanupInterval(interval time.Duration) {
	d.cleanupInterval = interval
}

func (d *DingtalkCacheConfig) SetKey(key string) {
	d.cacheKey = key
}

type DingtalkCacheAuthenticator struct {
	DingtalkCacheConfig
	delegate Authenticator
	cache    *cache.Cache
}

func NewDingtalkCacheAuthenticator(
	delegate Authenticator,
	opts ...CacheOption[*DingtalkCacheConfig],
) *DingtalkCacheAuthenticator {
	config := &DingtalkCacheConfig{
		ttl:             30 * time.Minute, // 默认值
		cleanupInterval: 1 * time.Hour,
		cacheKey:        "dingtalk_token",
	}

	for _, opt := range opts {
		opt(config)
	}

	return &DingtalkCacheAuthenticator{
		DingtalkCacheConfig: *config,
		delegate:            delegate,
		cache:               cache.New(config.ttl, config.cleanupInterval),
	}
}

func (c *DingtalkCacheAuthenticator) GetAccessToken(ctx context.Context) (*AccessTokenResp, error) {
	if cached, found := c.cache.Get(c.cacheKey); found {
		if token, ok := cached.(*AccessTokenResp); ok {
			fmt.Println("GetAccessToken from cache")
			//fmt.Printf("cacheKey: %s, token: %v\n", c.cacheKey, token)
			return token, nil
		}
	}

	token, err := c.delegate.GetAccessToken(ctx)
	if err != nil {
		return nil, err
	}

	cacheTTL := time.Duration(token.ExpiresIn-10) * time.Second
	if cacheTTL <= 0 {
		cacheTTL = cache.DefaultExpiration
	}

	c.cache.Set(c.cacheKey, token, cacheTTL)
	fmt.Println("GetAccessToken from api")
	//fmt.Printf("cacheKey: %s, token: %v, cacheTTL: %d\n", c.cacheKey, token, cacheTTL)
	return token, nil
}

func (c *DingtalkCacheAuthenticator) InvalidateCache() {
	c.cache.Delete(c.cacheKey)
}

// 使用示例：
// auth := NewDingtalkAuthenticator(baseAuth,
//     WithTTL[*DingtalkConfig](10*time.Minute),
//     WithKey[*DingtalkConfig]("custom_key"),
// )
