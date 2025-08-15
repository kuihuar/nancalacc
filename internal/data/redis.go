package data

import (
	"context"
	"nancalacc/internal/conf"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
)

func NewRedisClient(c *conf.Data, logger log.Logger) (*redis.Client, error) {
	// 检查 Redis 是否启用
	if !c.Redis.Enable {
		logger.Log(log.LevelInfo, "msg", "redis is disabled, skipping initialization")
		return nil, nil
	}

	// logger := integration.CreateLogger()
	rdb := redis.NewClient(&redis.Options{
		Addr:         c.Redis.Addr,
		Password:     c.Redis.Password, // 添加密码字段
		ReadTimeout:  c.Redis.ReadTimeout.AsDuration(),
		WriteTimeout: c.Redis.WriteTimeout.AsDuration(),
	})
	timeout, cancelFunc := context.WithTimeout(context.Background(), time.Second*2)
	defer cancelFunc()
	err := rdb.Ping(timeout).Err()
	if err != nil {
		logger.Log(log.LevelError, "msg", "redis ping failed", "error", err)
		return nil, err
	}
	return rdb, nil
}
