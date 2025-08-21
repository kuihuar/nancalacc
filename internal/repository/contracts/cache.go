package contracts

import (
	"context"
	"time"
)

// CacheRepository 缓存数据访问接口
type CacheRepository interface {
	Get(ctx context.Context, key string) (interface{}, bool, error)
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Del(ctx context.Context, key string) error
}
