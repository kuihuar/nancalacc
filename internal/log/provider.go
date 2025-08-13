package log

import (
	"time"

	"nancalacc/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	gormlogger "gorm.io/gorm/logger"
)

// ProviderSet 日志提供者集合
var ProviderSet = wire.NewSet(
	NewLogger,
	NewLogHelper,
	NewGormLoggerAdapter,
	NewLogConfig,
	NewLokiClientFromConfig,
	NewLoggerFromBootstrap,
)

// NewLogConfig 创建日志配置
func NewLogConfig() *Config {
	return &Config{
		Level:      "info",
		Format:     "json",
		Output:     "stdout",
		FilePath:   "logs/app.log",
		MaxSize:    100,
		MaxBackups: 10,
		MaxAge:     30,
		Compress:   true,
		Caller:     true,
		Stacktrace: false,
		Gorm: GormConfig{
			SlowThreshold: "200ms",
			LogLevel:      "warn",
		},
		Loki: LokiConfig{
			URL:       "http://localhost:3100",
			Enable:    false,
			BatchSize: 100,
			BatchWait: time.Second,
			Timeout:   10 * time.Second,
			Labels: map[string]string{
				"service": "nancalacc",
			},
		},
	}
}

// NewLoggerFromBootstrap 从 Bootstrap 配置创建完整的日志记录器
func NewLoggerFromBootstrap(bc *conf.Bootstrap) (log.Logger, error) {
	// 从 Bootstrap 配置创建日志配置
	logConfig := &Config{
		Level:      bc.GetLogging().GetLevel(),
		Format:     bc.GetLogging().GetFormat(),
		Output:     bc.GetLogging().GetOutput(),
		FilePath:   bc.GetLogging().GetFilePath(),
		MaxSize:    int(bc.GetLogging().GetMaxSize()),
		MaxBackups: int(bc.GetLogging().GetMaxBackups()),
		MaxAge:     int(bc.GetLogging().GetMaxAge()),
		Compress:   bc.GetLogging().GetCompress(),
		Stacktrace: bc.GetLogging().GetStacktrace(),
		Loki: LokiConfig{
			Enable: bc.GetLogging().GetLoki().GetEnable(),
		},
	}

	// 创建自定义日志记录器
	customLogger, err := NewLogger(logConfig)
	if err != nil {
		return nil, err
	}

	// 添加基础字段
	serviceName := bc.App.GetName()
	loggerWithFields := customLogger.WithFields(map[string]interface{}{
		"service.name": serviceName,
	})

	// 创建 Kratos 日志适配器
	kratosLogger := NewKratosLoggerAdapter(loggerWithFields)

	return kratosLogger, nil
}

// NewLogHelper 创建日志助手
func NewLogHelper(logger Logger) *Helper {
	return &Helper{logger: logger}
}

// NewGormLoggerAdapter 创建GORM日志适配器
func NewGormLoggerAdapter(logger Logger) gormlogger.Interface {
	return NewGormLogger(logger, 200*time.Millisecond)
}

// NewLokiClientFromConfig 从配置创建Loki客户端
func NewLokiClientFromConfig(config *Config) *LokiClient {
	if !config.Loki.Enable {
		return nil
	}
	return NewLokiClient(&config.Loki)
}

// KratosLoggerAdapter Kratos日志适配器
type KratosLoggerAdapter struct {
	logger Logger
}

// NewKratosLoggerAdapter 创建Kratos日志适配器
func NewKratosLoggerAdapter(logger Logger) log.Logger {
	return &KratosLoggerAdapter{logger: logger}
}

// Log 实现Kratos日志接口
func (k *KratosLoggerAdapter) Log(level log.Level, keyvals ...interface{}) error {
	return k.logger.Log(level, keyvals...)
}

// With 添加字段
func (k *KratosLoggerAdapter) With(keyvals ...interface{}) log.Logger {
	if len(keyvals) == 0 {
		return k
	}

	fields := make(map[string]interface{})
	for i := 0; i < len(keyvals); i += 2 {
		if i+1 < len(keyvals) {
			key, ok := keyvals[i].(string)
			if ok {
				fields[key] = keyvals[i+1]
			}
		}
	}

	return &KratosLoggerAdapter{
		logger: k.logger.WithFields(fields),
	}
}
