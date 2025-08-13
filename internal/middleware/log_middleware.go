package middleware

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
)

// LogMiddleware 日志中间件
func LogMiddleware(logger log.Logger) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			var (
				start = time.Now()
				tr    transport.Transporter
			)

			if info, ok := transport.FromServerContext(ctx); ok {
				tr = info
			}

			// 记录请求开始
			logger.Log(log.LevelInfo,
				"msg", "request started",
				"method", tr.Operation(),
				"path", tr.Endpoint(),
				"args", req,
			)

			// 处理请求
			reply, err = handler(ctx, req)

			// 记录请求结束
			logger.Log(log.LevelInfo,
				"msg", "request finished",
				"method", tr.Operation(),
				"path", tr.Endpoint(),
				"duration", time.Since(start).String(),
				"error", err,
			)

			return
		}
	}
}

// ErrorHandler 统一错误处理中间件
func ErrorHandler(logger log.Logger) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			defer func() {
				if r := recover(); r != nil {
					logger.Log(log.LevelError,
						"msg", "panic recovered",
						"panic", r,
					)
					err = &PanicError{Panic: r}
				}
			}()

			return handler(ctx, req)
		}
	}
}

// PanicError 自定义panic错误类型
type PanicError struct {
	Panic interface{}
}

func (e *PanicError) Error() string {
	return "panic occurred"
}
