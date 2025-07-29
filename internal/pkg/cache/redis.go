package cache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(client *redis.Client) *RedisCache {
	return &RedisCache{client: client}
}

func (r *RedisCache) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *RedisCache) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	return r.client.Set(ctx, key, value, ttl).Err()
}
