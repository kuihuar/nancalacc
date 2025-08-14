package data

import (
	"fmt"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
)

// DatabaseInitializer 数据库初始化器
type DatabaseInitializer struct {
	factory *DatabaseFactory
	logger  log.Logger
}

// NewDatabaseInitializer 创建数据库初始化器
func NewDatabaseInitializer(factory *DatabaseFactory, logger log.Logger) *DatabaseInitializer {
	return &DatabaseInitializer{
		factory: factory,
		logger:  logger,
	}
}

// InitializeDatabases 初始化所有数据库连接
func (di *DatabaseInitializer) InitializeDatabases() (*DatabaseManager, error) {
	dbManager := NewDatabaseManager()

	// 初始化主数据库
	if err := di.initializeMainDB(dbManager); err != nil {
		di.logger.Log(log.LevelError, "msg", "failed to initialize main database", "error", err)
		return nil, fmt.Errorf("failed to initialize main database: %w", err)
	}

	// 初始化同步数据库
	if err := di.initializeSyncDB(dbManager); err != nil {
		di.logger.Log(log.LevelError, "msg", "failed to initialize sync database", "error", err)
		return nil, fmt.Errorf("failed to initialize sync database: %w", err)
	}

	// 初始化用户数据库（可选）
	if err := di.initializeUserDB(dbManager); err != nil {
		di.logger.Log(log.LevelWarn, "msg", "failed to initialize user database", "error", err)
		// 用户数据库初始化失败不影响主流程
	}

	// 初始化日志数据库（可选）
	if err := di.initializeLogDB(dbManager); err != nil {
		di.logger.Log(log.LevelWarn, "msg", "failed to initialize log database", "error", err)
		// 日志数据库初始化失败不影响主流程
	}

	di.logger.Log(log.LevelInfo, "msg", "all databases initialized successfully")
	return dbManager, nil
}

// initializeMainDB 初始化主数据库
func (di *DatabaseInitializer) initializeMainDB(dbManager *DatabaseManager) error {
	config := di.factory.CreateMainDBConfig()

	if !config.Enable {
		di.logger.Log(log.LevelInfo, "msg", "main database is disabled")
		return nil
	}

	db, err := di.factory.CreateDatabase(MainDBType, config)
	if err != nil {
		return fmt.Errorf("failed to create main database: %w", err)
	}

	dbManager.RegisterDatabase(MainDBType, "main", db, config)
	di.logger.Log(log.LevelInfo, "msg", "main database initialized successfully")
	return nil
}

// initializeSyncDB 初始化同步数据库
func (di *DatabaseInitializer) initializeSyncDB(dbManager *DatabaseManager) error {
	config := di.factory.CreateSyncDBConfig()

	if !config.Enable {
		di.logger.Log(log.LevelInfo, "msg", "sync database is disabled")
		return nil
	}

	db, err := di.factory.CreateDatabase(SyncDBType, config)
	if err != nil {
		return fmt.Errorf("failed to create sync database: %w", err)
	}

	dbManager.RegisterDatabase(SyncDBType, "sync", db, config)
	di.logger.Log(log.LevelInfo, "msg", "sync database initialized successfully")
	return nil
}

// initializeUserDB 初始化用户数据库
func (di *DatabaseInitializer) initializeUserDB(dbManager *DatabaseManager) error {
	config := di.factory.CreateUserDBConfig()

	if !config.Enable {
		di.logger.Log(log.LevelInfo, "msg", "user database is disabled")
		return nil
	}

	db, err := di.factory.CreateDatabase(UserDBType, config)
	if err != nil {
		return fmt.Errorf("failed to create user database: %w", err)
	}

	dbManager.RegisterDatabase(UserDBType, "user", db, config)
	di.logger.Log(log.LevelInfo, "msg", "user database initialized successfully")
	return nil
}

// initializeLogDB 初始化日志数据库
func (di *DatabaseInitializer) initializeLogDB(dbManager *DatabaseManager) error {
	config := di.factory.CreateLogDBConfig()

	if !config.Enable {
		di.logger.Log(log.LevelInfo, "msg", "log database is disabled")
		return nil
	}

	db, err := di.factory.CreateDatabase(LogDBType, config)
	if err != nil {
		return fmt.Errorf("failed to create log database: %w", err)
	}

	dbManager.RegisterDatabase(LogDBType, "log", db, config)
	di.logger.Log(log.LevelInfo, "msg", "log database initialized successfully")
	return nil
}

// NewDataWithFactory 使用数据库工厂创建数据层实例
func NewDataWithFactory(
	factory *DatabaseFactory,
	redis *redis.Client,
	logger log.Logger,
) (*Data, func(), error) {
	// 创建数据库初始化器
	initializer := NewDatabaseInitializer(factory, logger)

	// 初始化所有数据库
	dbManager, err := initializer.InitializeDatabases()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to initialize databases: %w", err)
	}

	data := &Data{
		dbManager: dbManager,
		redis:     redis,
		logger:    logger,
	}

	// 返回清理函数
	cleanup := func() {
		data.cleanup()
	}

	return data, cleanup, nil
}
