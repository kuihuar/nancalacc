package repository

import (
	"nancalacc/internal/repository/contracts"
	"nancalacc/internal/repository/impl/localcache"
	"nancalacc/internal/repository/impl/mysql"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"gorm.io/gorm"
)

// ProviderSet 是 repository 的依赖注入集合
var ProviderSet = wire.NewSet(
	NewAccountRepository,
	NewTaskRepository,
	NewCacheRepository,
)

// NewAccountRepository 创建账户Repository
func NewAccountRepository(db *gorm.DB, logger log.Logger) contracts.AccountRepository {
	return mysql.NewAccountRepository(db, logger)
}

// NewTaskRepository 创建任务Repository
func NewTaskRepository(db *gorm.DB, logger log.Logger) contracts.TaskRepository {
	return mysql.NewTaskRepository(db, logger)
}

// NewCacheRepos	itory 创建缓存Repository
func NewCacheRepository(logger log.Logger) contracts.CacheRepository {
	return localcache.NewLocalCacheRepository(logger)
}
