// cache_options.go
package auth

import "time"

// 通用缓存配置接口
type CacheConfig interface {
	SetTTL(time.Duration)
	SetCleanupInterval(time.Duration)
	SetKey(string)
}

// 泛型 Option 类型
type CacheOption[T CacheConfig] func(T)

// 通用 Option 构造函数
func WithTTL[T CacheConfig](ttl time.Duration) CacheOption[T] {
	return func(c T) {
		c.SetTTL(ttl)
	}
}

func WithCleanupInterval[T CacheConfig](interval time.Duration) CacheOption[T] {
	return func(c T) {
		c.SetCleanupInterval(interval)
	}
}

func WithKey[T CacheConfig](key string) CacheOption[T] {
	return func(c T) {
		c.SetKey(key)
	}
}
