package data

import (
	"context"
	"encoding/json"
	"nancalacc/internal/biz"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
)

type redisCacheService struct {
	client *redis.Client
	log    *log.Helper
}

// NewAccounterRepo .
func NewRedisCacheService(client *redis.Client, logger log.Logger) biz.CacheService {
	return &redisCacheService{
		client: client,
		log:    log.NewHelper(logger),
	}
}

func (r *redisCacheService) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	p, err := json.Marshal(value)
	if err != nil {
		return err
	}
	r.log.Infof("Set key: %s, value: %s", key, string(p))
	return r.client.Set(ctx, key, p, expiration).Err()
}

func (r *redisCacheService) Get(ctx context.Context, key string) (interface{}, bool, error) {
	p, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, false, nil
		}
		return nil, false, err
	}
	var v interface{}
	err = json.Unmarshal(p, &v)
	if err != nil {
		return nil, false, err
	}
	return v, true, nil
}

func (r *redisCacheService) Del(ctx context.Context, key string) error {
	r.client.Del(ctx, key)
	return nil
}

// func (r *redisRepo) GetWithCachePenetrationProtection(
// 	ctx context.Context,
// 	key string,
// 	dest interface{},
// 	fallback func() (interface{}, error),
// 	ttl time.Duration,
// ) error {
// 	// 1. 先查缓存
// 	err := r.Get(ctx, key, dest)
// 	if err == nil {
// 		return nil
// 	}

// 	// 2. 获取分布式锁（防击穿）
// 	lockKey := "lock:" + key
// 	if !r.Lock(ctx, lockKey, 10*time.Second) {
// 		return errors.New("操作过于频繁")
// 	}
// 	defer r.Unlock(ctx, lockKey)

// 	// 3. 回源查询
// 	data, err := fallback()
// 	if err != nil {
// 		return err
// 	}

// 	// 4. 写入缓存
// 	return r.Set(ctx, key, data, ttl)
// }

// func (r *redisRepo) Lock(ctx context.Context, key string, expiration time.Duration) bool {
// 	return r.data.redis.SetNX(ctx, key, 1, expiration).Val()
// }

// func (r *redisRepo) Unlock(ctx context.Context, key string) {
// 	r.data.redis.Del(ctx, key)
// }

// other method：HSet/HGet、LPush/LRange、Incr 等...
