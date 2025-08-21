package repository

import (
	"nancalacc/internal/data"
	"nancalacc/internal/repository/contracts"
	"nancalacc/internal/repository/impl/localcache"
	"nancalacc/internal/repository/impl/mysql"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet 是 repository 的依赖注入集合
var ProviderSet = wire.NewSet(
	NewAccountRepository,
	NewTaskRepository,
	NewCacheRepository,
)

// NewAccountRepository 创建账户Repository
func NewAccountRepository(data *data.Data, logger log.Logger) contracts.AccountRepository {
	syncDB, err := data.GetSyncDB()
	if err != nil {
		logger.Log(log.LevelError, "msg", "failed to get sync database", "error", err)
		return nil
	}
	return mysql.NewAccountRepository(syncDB, logger)
}

// NewTaskRepository 创建任务Repository
func NewTaskRepository(data *data.Data, logger log.Logger) contracts.TaskRepository {
	mainDB, err := data.GetMainDB()
	if err != nil {
		logger.Log(log.LevelError, "msg", "failed to get sync database", "error", err)
		return nil
	}
	return mysql.NewTaskRepository(mainDB, logger)
}

// NewCacheRepos	itory 创建缓存Repository
func NewCacheRepository(logger log.Logger) contracts.CacheRepository {
	return localcache.NewLocalCacheRepository(logger)
}
