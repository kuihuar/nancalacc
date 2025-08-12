package auth

import (
	"time"

	gocache "github.com/patrickmn/go-cache"
)

// Cache 缓存接口
type Cache interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{}, ttl time.Duration)
	Delete(key string)
}

// LocalCache 本地缓存实现
type LocalCache struct {
	cache *gocache.Cache
}

// NewLocalCache 创建新的本地缓存实例
func NewLocalCache() Cache {
	// 默认清理间隔为10分钟，默认过期时间为1小时
	return &LocalCache{
		cache: gocache.New(1*time.Hour, 10*time.Minute),
	}
}

func (lc *LocalCache) Get(key string) (interface{}, bool) {
	return lc.cache.Get(key)
}

func (lc *LocalCache) Set(key string, value interface{}, ttl time.Duration) {
	lc.cache.Set(key, value, ttl)
}

func (lc *LocalCache) Delete(key string) {
	lc.cache.Delete(key)
}
