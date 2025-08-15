package data

import (
	"errors"
	"fmt"
	"nancalacc/internal/conf"
	"nancalacc/internal/pkg/cipherutil"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DatabaseFactory 数据库连接工厂
type DatabaseFactory struct {
	config     *conf.Data
	logger     log.Logger
	otelConfig *conf.OpenTelemetry // OpenTelemetry配置
}

// NewDatabaseFactory 创建数据库工厂
func NewDatabaseFactory(config *conf.Data, logger log.Logger, otelConfig *conf.OpenTelemetry) *DatabaseFactory {
	return &DatabaseFactory{
		config:     config,
		logger:     logger,
		otelConfig: otelConfig,
	}
}

// CreateDatabase 创建数据库连接
func (df *DatabaseFactory) CreateDatabase(dbType DatabaseType, config *DatabaseConnectionConfig) (*gorm.DB, error) {
	if config == nil {
		return nil, errors.New("database config is nil")
	}

	// 获取数据库连接字符串
	dsn, err := df.getDSN(config)
	if err != nil {
		return nil, fmt.Errorf("failed to get DSN for %s: %w", dbType, err)
	}

	// 确定GORM日志级别：OpenTelemetry配置优先，数据库配置作为默认值
	var logLevel string
	if df.otelConfig != nil && df.otelConfig.Logs != nil && df.otelConfig.Logs.Gorm != nil && df.otelConfig.Logs.Gorm.LogLevel != "" {
		logLevel = df.otelConfig.Logs.Gorm.LogLevel
	} else if config.LogLevel != "" {
		logLevel = config.LogLevel
	} else {
		logLevel = "info" // 默认日志级别
	}

	// 创建使用OpenTelemetry logger的GORM配置
	gormConfig := &gorm.Config{
		Logger: NewGormLogger(df.logger, logLevel),
	}

	// 打开数据库连接
	db, err := gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		df.logger.Log(log.LevelError, "msg", "failed to open database", "type", dbType, "error", err)
		return nil, fmt.Errorf("failed to open database %s: %w", dbType, err)
	}

	// 配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB for %s: %w", dbType, err)
	}

	// 设置连接池参数
	if config.MaxOpenConns > 0 {
		sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	}
	if config.MaxIdleConns > 0 {
		sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	}
	if config.ConnMaxLifetime > 0 {
		sqlDB.SetConnMaxLifetime(config.ConnMaxLifetime)
	}

	df.logger.Log(log.LevelInfo, "msg", "database connection created",
		"type", dbType, "max_open_conns", config.MaxOpenConns,
		"max_idle_conns", config.MaxIdleConns)

	return db, nil
}

// getDSN 获取数据库连接字符串
func (df *DatabaseFactory) getDSN(config *DatabaseConnectionConfig) (string, error) {
	var dsn string

	// 开发环境直接使用配置的连接字符串
	if config.Env == "dev" || config.Env == "development" {
		dsn = config.Source
	} else {
		// 生产环境从环境变量获取加密的连接字符串
		if config.SourceKey == "" {
			return "", errors.New("source_key is required for non-dev environment")
		}

		encryptedDsn, err := conf.GetEnv(config.SourceKey)
		if err != nil {
			df.logger.Log(log.LevelError, "msg", "failed to get environment variable",
				"key", config.SourceKey, "error", err)
			return "", fmt.Errorf("failed to get environment variable %s: %w", config.SourceKey, err)
		}

		// 解密连接字符串
		appSecret := conf.Get().GetApp().GetAppSecret()
		dsn, err = cipherutil.DecryptByAes(encryptedDsn, appSecret)
		if err != nil {
			df.logger.Log(log.LevelError, "msg", "failed to decrypt DSN",
				"key", config.SourceKey, "error", err)
			return "", fmt.Errorf("failed to decrypt DSN for %s: %w", config.SourceKey, err)
		}

		if len(dsn) == 0 {
			return "", fmt.Errorf("decrypted DSN is empty for %s", config.SourceKey)
		}
	}

	// 确保连接字符串包含必要的参数
	if !strings.Contains(dsn, "parseTime=True") {
		if strings.Contains(dsn, "?") {
			dsn += "&parseTime=True"
		} else {
			dsn += "?parseTime=True"
		}
	}

	// 添加超时参数
	if !strings.Contains(dsn, "timeout=") {
		if strings.Contains(dsn, "?") {
			dsn += "&timeout=15s"
		} else {
			dsn += "?timeout=15s"
		}
	}

	// 添加字符集参数
	if !strings.Contains(dsn, "charset=") {
		if strings.Contains(dsn, "?") {
			dsn += "&charset=utf8mb4"
		} else {
			dsn += "?charset=utf8mb4"
		}
	}

	return dsn, nil
}

// DatabaseConnectionConfig 数据库连接配置
type DatabaseConnectionConfig struct {
	Source          string        // 连接字符串（开发环境使用）
	SourceKey       string        // 环境变量键名（生产环境使用）
	Env             string        // 环境：dev, prod, test
	MaxOpenConns    int           // 最大打开连接数
	MaxIdleConns    int           // 最大空闲连接数
	ConnMaxLifetime time.Duration // 连接最大生命周期
	LogLevel        string        // 日志级别：silent, error, warn, info
	Enable          bool          // 是否启用
}

// NewDatabaseConnectionConfig 创建数据库连接配置
func NewDatabaseConnectionConfig() *DatabaseConnectionConfig {
	return &DatabaseConnectionConfig{
		MaxOpenConns:    20,
		MaxIdleConns:    10,
		ConnMaxLifetime: 6 * time.Hour,
		LogLevel:        "info",
		Enable:          true,
	}
}

// CreateMainDBConfig 创建主数据库配置
func (df *DatabaseFactory) CreateMainDBConfig() *DatabaseConnectionConfig {
	config := NewDatabaseConnectionConfig()

	if df.config.Database != nil {
		config.Source = df.config.Database.Source
		config.SourceKey = df.config.Database.SourceKey
		config.Env = df.config.Database.Env
		config.MaxOpenConns = int(df.config.Database.MaxOpenConns)
		config.MaxIdleConns = int(df.config.Database.MaxIdleConns)
		config.Enable = df.config.Database.Enable

		if df.config.Database.ConnMaxLifetime != "" {
			if duration, err := time.ParseDuration(df.config.Database.ConnMaxLifetime); err == nil {
				config.ConnMaxLifetime = duration
			}
		}
	}

	return config
}

// CreateSyncDBConfig 创建同步数据库配置
func (df *DatabaseFactory) CreateSyncDBConfig() *DatabaseConnectionConfig {
	config := NewDatabaseConnectionConfig()

	if df.config.DatabaseSync != nil {
		config.Source = df.config.DatabaseSync.Source
		config.SourceKey = df.config.DatabaseSync.SourceKey
		config.Env = df.config.DatabaseSync.Env
		config.MaxOpenConns = int(df.config.DatabaseSync.MaxOpenConns)
		config.MaxIdleConns = int(df.config.DatabaseSync.MaxIdleConns)

		if df.config.DatabaseSync.ConnMaxLifetime != "" {
			if duration, err := time.ParseDuration(df.config.DatabaseSync.ConnMaxLifetime); err == nil {
				config.ConnMaxLifetime = duration
			}
		}
	}

	return config
}

// CreateUserDBConfig 创建用户数据库配置（示例）
func (df *DatabaseFactory) CreateUserDBConfig() *DatabaseConnectionConfig {
	config := NewDatabaseConnectionConfig()
	// 这里可以根据实际需求配置用户数据库
	// 可以从配置文件、环境变量等获取配置
	return config
}

// CreateLogDBConfig 创建日志数据库配置（示例）
func (df *DatabaseFactory) CreateLogDBConfig() *DatabaseConnectionConfig {
	config := NewDatabaseConnectionConfig()
	// 这里可以根据实际需求配置日志数据库
	return config
}
