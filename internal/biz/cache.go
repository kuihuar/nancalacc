package biz

import (
	"context"
	"time"
)

type RedisCacheRepo interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string, dest interface{}) error
	GetWithCachePenetrationProtection(
		ctx context.Context,
		key string,
		dest interface{},
		fallback func() (interface{}, error),
		ttl time.Duration,
	) error
	Lock(ctx context.Context, key string, expiration time.Duration) bool
	Unlock(ctx context.Context, key string)
}
