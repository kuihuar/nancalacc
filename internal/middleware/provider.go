package middleware

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/google/wire"
)

// ProviderSet 中间件提供者集合
var ProviderSet = wire.NewSet(
	NewLoggingMiddleware,
)

// NewLoggingMiddleware 创建通用日志中间件
func NewLoggingMiddleware(logger log.Logger) middleware.Middleware {
	return LoggingMiddleware(logger)
}
