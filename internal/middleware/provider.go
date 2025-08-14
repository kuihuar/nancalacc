package middleware

import (
	"nancalacc/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/google/wire"
)

// ProviderSet 中间件提供者集合
var ProviderSet = wire.NewSet(
	NewLoggingMiddleware,
	NewTracingMiddleware,
)

// NewLoggingMiddleware 创建通用日志中间件
func NewLoggingMiddleware(logger log.Logger) middleware.Middleware {
	return LoggingMiddleware(logger)
}

// NewTracingMiddleware 创建OpenTelemetry追踪中间件
func NewTracingMiddleware(logger log.Logger, tracing *conf.Tracing) middleware.Middleware {
	config := &TracingConfig{
		VerboseLogging:  tracing.GetVerboseLogging(),
		LogRequestSize:  tracing.GetLogRequestSize(),
		LogResponseSize: tracing.GetLogResponseSize(),
		LogHeaders:      tracing.GetLogHeaders(),
		LogQueryParams:  tracing.GetLogQueryParams(),
	}
	return TracingMiddleware(logger, config)
}
