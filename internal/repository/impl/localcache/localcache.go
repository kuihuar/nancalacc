package localcache

import (
	"context"
	"nancalacc/internal/repository/contracts"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/patrickmn/go-cache"
)

type localCacheRepository struct {
	client *cache.Cache
	logger log.Logger
}

var (
	DefaultExpiration = 5 * time.Minute
	CleanupInterval   = 10 * time.Minute
)

func NewLocalCacheRepository(logger log.Logger) contracts.CacheRepository {
	return &localCacheRepository{
		client: cache.New(DefaultExpiration, CleanupInterval),
		logger: logger,
	}
}

func (s *localCacheRepository) Set(ctx context.Context, key string, value interface{}, d time.Duration) error {
	s.client.Set(key, value, d)
	return nil
}

func (s *localCacheRepository) Get(ctx context.Context, key string) (interface{}, bool, error) {
	res, ok := s.client.Get(key)
	if !ok {
		return nil, false, nil
	}
	return res, true, nil
}

func (s *localCacheRepository) Del(ctx context.Context, key string) error {
	s.client.Delete(key)
	return nil
}
