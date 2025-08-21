package data

import (
	"context"
	"fmt"
	"sync"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	"gorm.io/gorm"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(
	NewRedisClient,
	NewDatabaseFactory,     // 数据库工厂
	NewDatabaseInitializer, // 数据库初始化器
	NewDataWithFactory,     // 使用工厂创建数据层
)

// DatabaseType 数据库类型
type DatabaseType string

const (
	MainDBType  DatabaseType = "main"  // 主数据库
	SyncDBType  DatabaseType = "sync"  // 同步数据库
	UserDBType  DatabaseType = "user"  // 用户数据库
	LogDBType   DatabaseType = "log"   // 日志数据库
	CacheDBType DatabaseType = "cache" // 缓存数据库
	SagaDBType  DatabaseType = "saga"  // Saga 分布式事务数据库
)

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Type     DatabaseType
	Name     string
	DB       *gorm.DB
	Config   interface{}
	IsActive bool
}

// DatabaseManager 数据库管理器
type DatabaseManager struct {
	databases map[DatabaseType]*DatabaseConfig
	mu        sync.RWMutex
}

// NewDatabaseManager 创建数据库管理器
func NewDatabaseManager() *DatabaseManager {
	return &DatabaseManager{
		databases: make(map[DatabaseType]*DatabaseConfig),
	}
}

// RegisterDatabase 注册数据库
func (dm *DatabaseManager) RegisterDatabase(dbType DatabaseType, name string, db *gorm.DB, config interface{}) {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	dm.databases[dbType] = &DatabaseConfig{
		Type:     dbType,
		Name:     name,
		DB:       db,
		Config:   config,
		IsActive: db != nil,
	}
}

// GetDatabase 获取数据库
func (dm *DatabaseManager) GetDatabase(dbType DatabaseType) (*gorm.DB, error) {
	dm.mu.RLock()
	defer dm.mu.RUnlock()

	if config, exists := dm.databases[dbType]; exists && config.IsActive {
		return config.DB, nil
	}
	return nil, fmt.Errorf("database %s not found or not active", dbType)
}

// GetDatabaseConfig 获取数据库配置
func (dm *DatabaseManager) GetDatabaseConfig(dbType DatabaseType) (*DatabaseConfig, error) {
	dm.mu.RLock()
	defer dm.mu.RUnlock()

	if config, exists := dm.databases[dbType]; exists {
		return config, nil
	}
	return nil, fmt.Errorf("database config %s not found", dbType)
}

// ListDatabases 列出所有数据库
func (dm *DatabaseManager) ListDatabases() map[DatabaseType]*DatabaseConfig {
	dm.mu.RLock()
	defer dm.mu.RUnlock()

	result := make(map[DatabaseType]*DatabaseConfig)
	for k, v := range dm.databases {
		result[k] = v
	}
	return result
}

// CloseAll 关闭所有数据库连接
func (dm *DatabaseManager) CloseAll(logger log.Logger) {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	for dbType, config := range dm.databases {
		if config.IsActive && config.DB != nil {
			if sqlDB, err := config.DB.DB(); err == nil {
				if err := sqlDB.Close(); err != nil {
					logger.Log(log.LevelError, "msg", "failed to close database",
						"type", dbType, "name", config.Name, "error", err)
				} else {
					logger.Log(log.LevelInfo, "msg", "database closed",
						"type", dbType, "name", config.Name)
				}
			}
			config.IsActive = false
		}
	}
}

// Data 数据层结构体
type Data struct {
	dbManager *DatabaseManager
	redis     *redis.Client
	logger    log.Logger
	sagaRepo  *SagaRepository
	// 可以添加其他仓库
}

// cleanup 清理资源
func (d *Data) cleanup() {
	// 关闭所有数据库连接
	d.dbManager.CloseAll(d.logger)

	// 关闭 Redis 连接
	if d.redis != nil {
		if err := d.redis.Close(); err != nil {
			d.logger.Log(log.LevelError, "msg", "failed to close redis", "error", err)
		} else {
			d.logger.Log(log.LevelInfo, "msg", "redis connection closed")
		}
	}

	d.logger.Log(log.LevelInfo, "msg", "all database connections closed")
}

// GetMainDB 获取主数据库
func (d *Data) GetMainDB() (*gorm.DB, error) {
	return d.dbManager.GetDatabase(MainDBType)
}

// GetSyncDB 获取同步数据库
func (d *Data) GetSyncDB() (*gorm.DB, error) {
	return d.dbManager.GetDatabase(SyncDBType)
}

// GetSagaDB 获取 Saga 数据库
func (d *Data) GetSagaDB() (*gorm.DB, error) {
	return d.dbManager.GetDatabase(SagaDBType)
}

// GetSagaRepository 获取 Saga 仓库
func (d *Data) GetSagaRepository() *SagaRepository {
	return d.sagaRepo
}

// GetDatabase 获取指定类型的数据库
func (d *Data) GetDatabase(dbType DatabaseType) (*gorm.DB, error) {
	return d.dbManager.GetDatabase(dbType)
}

// GetRedis 获取 Redis 客户端
func (d *Data) GetRedis() *redis.Client {
	return d.redis
}

// GetDBManager 获取数据库管理器
func (d *Data) GetDBManager() *DatabaseManager {
	return d.dbManager
}

// IsRedisAvailable 检查 Redis 是否可用
func (d *Data) IsRedisAvailable() bool {
	if d.redis == nil {
		return false
	}

	ctx := context.Background()
	_, err := d.redis.Ping(ctx).Result()
	return err == nil
}

// HealthCheck 健康检查
func (d *Data) HealthCheck(ctx context.Context) map[string]interface{} {
	health := make(map[string]interface{})

	// 检查数据库连接
	databases := d.dbManager.ListDatabases()
	for dbType, config := range databases {
		if config.IsActive && config.DB != nil {
			if sqlDB, err := config.DB.DB(); err == nil {
				if err := sqlDB.PingContext(ctx); err == nil {
					health[string(dbType)] = "healthy"
				} else {
					health[string(dbType)] = fmt.Sprintf("unhealthy: %v", err)
				}
			} else {
				health[string(dbType)] = "unhealthy: failed to get underlying sql.DB"
			}
		} else {
			health[string(dbType)] = "inactive"
		}
	}

	// 检查 Redis 连接
	if d.redis != nil {
		if _, err := d.redis.Ping(ctx).Result(); err == nil {
			health["redis"] = "healthy"
		} else {
			health["redis"] = fmt.Sprintf("unhealthy: %v", err)
		}
	} else {
		health["redis"] = "inactive"
	}

	return health
}

// 保持向后兼容的方法
func (d *Data) DB() *gorm.DB {
	if db, err := d.GetSyncDB(); err == nil {
		return db
	}
	return nil
}

func (d *Data) NancalDB() *gorm.DB {
	if db, err := d.GetMainDB(); err == nil {
		return db
	}
	return nil
}

// NewLocalCacheService 创建本地缓存服务
// func NewLocalCacheService(logger log.Logger) localcache.CacheRepository {
// 	return localcache.NewLocalCacheService(logger)
// }
