package redis

// import (
// 	"context"
// 	"encoding/json"
// 	"errors"
// 	"nancalacc/internal/data"
// 	"time"
// )

// type RedisRepo struct {
// 	data *data.Data
// }

// func (r *RedisRepo) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
// 	p, err := json.Marshal(value)
// 	if err != nil {
// 		return err
// 	}
// 	return r.data.redis.Set(ctx, key, p, expiration).Err()
// }

// func (r *RedisRepo) Get(ctx context.Context, key string, dest interface{}) error {
// 	p, err := r.data.redis.Get(ctx, key).Bytes()
// 	if err != nil {
// 		return err
// 	}
// 	return json.Unmarshal(p, dest)
// }

// func (r *RedisRepo) GetWithCachePenetrationProtection(
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

// 其他封装方法：HSet/HGet、LPush/LRange、Incr 等...
