package data

import (
	"context"
	"nancalacc/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
)

func NewRedisClient(c *conf.Data, logger log.Logger) (*redis.Client, error) {

	if !c.Redis.Enable {
		logger.Log(log.LevelWarn, "redis Enable", c.Redis.Enable)
		return &redis.Client{}, nil
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     c.Redis.Addr,
		Password: c.Redis.Password,
		DB:       int(c.Redis.Db),
		// PoolSize: int(c.Redis.Pool_size),
	})

	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		logger.Log(log.LevelError, "ping redis failed", err)
		return nil, nil
	}
	return rdb, nil
}
