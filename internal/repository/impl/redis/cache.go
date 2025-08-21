package redis

import (
	"context"
	"encoding/json"
	"nancalacc/internal/data"
	"nancalacc/internal/repository/contracts"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
)

type redisRepository struct {
	client *redis.Client
	log    log.Logger
}

// NewCacheRepository 创建缓存Repository
func NewCacheRepository(data *data.Data, logger log.Logger) contracts.CacheRepository {
	return &redisRepository{

		client: data.GetRedis(),
		log:    logger,
	}
}

func (r *redisRepository) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	p, err := json.Marshal(value)
	if err != nil {
		return err
	}
	r.log.Log(log.LevelInfo, "msg", "Set key: %s, value: %s", key, string(p))
	return r.client.Set(ctx, key, p, expiration).Err()
}

func (r *redisRepository) Get(ctx context.Context, key string) (interface{}, bool, error) {
	p, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			// 未找到对应的key，返回false
			return nil, false, nil
		}
		// 其他错误，返回错误信息
		return nil, false, err
	}
	var v interface{}
	err = json.Unmarshal(p, &v)
	if err != nil {
		return nil, false, err
	}
	return v, true, nil
}

func (r *redisRepository) Del(ctx context.Context, key string) error {
	r.client.Del(ctx, key)
	return nil
}
