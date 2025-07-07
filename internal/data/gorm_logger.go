package data

import (
	"context"
	"errors"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// gormLogger 适配器实现gorm.Logger接口
type gormLogger struct {
	logger        *log.Helper
	SlowThreshold time.Duration
	LogLevel      gormlogger.LogLevel
}

// NewGormLogger 创建适配Kratos的GORM日志器
func NewGormLogger(logger log.Logger, slowThreshold time.Duration) gormlogger.Interface {
	return &gormLogger{
		logger:        log.NewHelper(logger),
		SlowThreshold: slowThreshold,
		LogLevel:      gormlogger.Warn, // 默认级别
	}
}

// LogMode 设置日志级别
func (l *gormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

// Info 打印信息
func (l *gormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Info {
		l.logger.WithContext(ctx).Info(msg, data)
	}
}

// Warn 打印警告
func (l *gormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Warn {
		l.logger.WithContext(ctx).Warn(msg, data)
	}
}

// Error 打印错误
func (l *gormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Error {
		l.logger.WithContext(ctx).Error(msg, data)
	}
}

// Trace 打印SQL跟踪
func (l *gormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= gormlogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	//sql, rows := fc()
	// fields := []interface{}{
	// 	"sql", sql,
	// 	"rows", rows,
	// 	"time", elapsed,
	// }

	switch {
	case err != nil && !errors.Is(err, gorm.ErrRecordNotFound):
		//l.logger.WithContext(ctx).Errorw("gorm error", append(fields, "err", err)...)
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0:
		//l.logger.WithContext(ctx).Warnw("slow sql", append(fields, "slow", l.SlowThreshold)...)
	case l.LogLevel >= gormlogger.Info:
		//l.logger.WithContext(ctx).Debugw("gorm query", fields...)
	}
}
