package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/patrickmn/go-cache"
)

type InvalidateCache interface {
	InvalidateCache()
}

type CachingAuthenticator interface {
	Authenticator
	InvalidateCache
}

type localCacheAuthenticator struct {
	delegate        Authenticator
	cache           *cache.Cache
	cacheKey        string
	ttl             time.Duration
	cleanupInterval time.Duration
}

type CacheOption func(*localCacheAuthenticator)

var (
	DefaultExpiration      = 30 * time.Minute
	DefaultCleanupInterval = 1 * time.Hour
	DefaultCacheKey        = "default_access_token"
)

func WithTTL(d time.Duration) CacheOption {
	return func(c *localCacheAuthenticator) {
		c.ttl = d
	}
}

func WithCleanupInterval(d time.Duration) CacheOption {
	return func(c *localCacheAuthenticator) {
		c.cleanupInterval = d
	}
}

func WithKey(key string) CacheOption {
	return func(c *localCacheAuthenticator) {
		c.cacheKey = key
	}
}
func NewLocalCachedAuthenticator(delegate Authenticator, opts ...CacheOption) CachingAuthenticator {

	c := &localCacheAuthenticator{
		delegate:        delegate,
		ttl:             DefaultExpiration,      // 默认TTL
		cleanupInterval: DefaultCleanupInterval, // 默认清理间隔
		cacheKey:        DefaultCacheKey,
	}

	for _, opt := range opts {
		opt(c)
	}

	c.cache = cache.New(c.ttl, c.cleanupInterval)
	return c
}

// 6. Impleted
func (c *localCacheAuthenticator) GetAccessToken(ctx context.Context) (*AccessTokenResp, error) {
	if cached, found := c.cache.Get(c.cacheKey); found {
		if token, ok := cached.(*AccessTokenResp); ok {
			fmt.Println("GetAccessToken from cache")
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
	return token, nil
}

func (c *localCacheAuthenticator) InvalidateCache() {
	c.cache.Delete(c.cacheKey)
}
