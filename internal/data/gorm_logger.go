package data

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// GormLogger 实现GORM logger接口，使用OpenTelemetry logger
type GormLogger struct {
	logger   log.Logger
	logLevel gormlogger.LogLevel
}

// NewGormLogger 创建新的GORM logger
func NewGormLogger(logger log.Logger, level string) gormlogger.Interface {
	var logLevel gormlogger.LogLevel
	switch level {
	case "silent":
		logLevel = gormlogger.Silent
	case "error":
		logLevel = gormlogger.Error
	case "warn":
		logLevel = gormlogger.Warn
	case "info":
		logLevel = gormlogger.Info
	default:
		logLevel = gormlogger.Info
	}

	return &GormLogger{
		logger:   logger,
		logLevel: logLevel,
	}
}

// LogMode 设置日志级别
func (l *GormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	newLogger := *l
	newLogger.logLevel = level
	return &newLogger
}

// Info 记录信息日志
func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= gormlogger.Info {
		l.logger.Log(log.LevelInfo, "msg", fmt.Sprintf(msg, data...))
	}
}

// Warn 记录警告日志
func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= gormlogger.Warn {
		l.logger.Log(log.LevelWarn, "msg", fmt.Sprintf(msg, data...))
	}
}

// Error 记录错误日志
func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= gormlogger.Error {
		l.logger.Log(log.LevelError, "msg", fmt.Sprintf(msg, data...))
	}
}

// Trace 记录SQL跟踪日志
func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.logLevel <= gormlogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	// 记录慢查询
	if elapsed > time.Second && l.logLevel >= gormlogger.Warn {
		l.logger.Log(log.LevelWarn, "msg", "slow sql query",
			"sql", sql,
			"rows", rows,
			"elapsed", elapsed.String())
		return
	}

	// 记录错误
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) && l.logLevel >= gormlogger.Error {
		l.logger.Log(log.LevelError, "msg", "sql error",
			"sql", sql,
			"rows", rows,
			"elapsed", elapsed.String(),
			"error", err)
		return
	}

	// 记录普通SQL日志
	if l.logLevel >= gormlogger.Info {
		l.logger.Log(log.LevelInfo, "msg", "sql query",
			"sql", sql,
			"rows", rows,
			"elapsed", elapsed.String())
	}
}
