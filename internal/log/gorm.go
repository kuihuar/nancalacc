package log

import (
	"context"
	"errors"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// GormLogger GORM日志适配器
type GormLogger struct {
	logger        Logger
	slowThreshold time.Duration
	logLevel      gormlogger.LogLevel
}

// NewGormLogger 创建GORM日志适配器
func NewGormLogger(logger Logger, slowThreshold time.Duration) gormlogger.Interface {
	return &GormLogger{
		logger:        logger,
		slowThreshold: slowThreshold,
		logLevel:      gormlogger.Warn,
	}
}

// LogMode 设置日志级别
func (l *GormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	newLogger := *l
	newLogger.logLevel = level
	return &newLogger
}

// Info 打印信息
func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= gormlogger.Info {
		l.logger.WithFields(map[string]interface{}{
			"type": "gorm_info",
			"msg":  msg,
			"data": data,
		}).Log(log.LevelInfo)
	}
}

// Warn 打印警告
func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= gormlogger.Warn {
		l.logger.WithFields(map[string]interface{}{
			"type": "gorm_warn",
			"msg":  msg,
			"data": data,
		}).Log(log.LevelWarn)
	}
}

// Error 打印错误
func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= gormlogger.Error {
		l.logger.WithFields(map[string]interface{}{
			"type": "gorm_error",
			"msg":  msg,
			"data": data,
		}).Log(log.LevelError)
	}
}

// Trace 打印SQL跟踪
func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.logLevel <= gormlogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	fields := map[string]interface{}{
		"type":     "gorm_trace",
		"sql":      sql,
		"rows":     rows,
		"duration": elapsed.String(),
	}

	switch {
	case err != nil && !errors.Is(err, gorm.ErrRecordNotFound):
		fields["type"] = "gorm_error"
		fields["error"] = err.Error()
		l.logger.WithFields(fields).Log(log.LevelError)
	case elapsed > l.slowThreshold && l.slowThreshold != 0:
		fields["type"] = "gorm_slow_query"
		fields["threshold"] = l.slowThreshold.String()
		l.logger.WithFields(fields).Log(log.LevelWarn)
	case l.logLevel >= gormlogger.Info:
		l.logger.WithFields(fields).Log(log.LevelInfo)
	}
}
