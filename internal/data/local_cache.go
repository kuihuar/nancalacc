package data

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/patrickmn/go-cache"
)

type LocalCacheService struct {
	client *cache.Cache
	logger log.Logger
}

var (
	DefaultExpiration = 5 * time.Minute
	CleanupInterval   = 10 * time.Minute
	Prefix            = "task:"
)

func NewCLocalCacheService(defaultExpiration, cleanupInterval time.Duration, logger log.Logger) *LocalCacheService {
	return &LocalCacheService{
		client: cache.New(defaultExpiration, cleanupInterval),
		logger: logger,
	}
}

func (s *LocalCacheService) Set(ctx context.Context, key string, value interface{}, d time.Duration) error {
	s.client.Set(key, value, d)
	return nil
}

func (s *LocalCacheService) Get(ctx context.Context, key string) (interface{}, bool, error) {
	res, ok := s.client.Get(key)
	if !ok {
		return nil, false, nil
	}
	return res, true, nil
}

// other methods...
