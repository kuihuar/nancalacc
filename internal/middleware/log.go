package middleware

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
)

// LoggingMiddleware 创建日志中间件
func LoggingMiddleware(logger log.Logger) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			var (
				start     = time.Now()
				kind      string
				operation string
			)

			// 获取传输信息
			if tr, ok := transport.FromServerContext(ctx); ok {
				kind = tr.Kind().String()
				operation = tr.Operation()
			}
			// 记录请求开始日志
			log.NewHelper(logger).WithContext(ctx).Infof("request started: kind=%s, operation=%s", kind, operation)

			// 处理请求
			reply, err = handler(ctx, req)

			// 计算耗时
			duration := time.Since(start)

			// 记录响应日志
			if err != nil {
				log.NewHelper(logger).WithContext(ctx).Errorf("request failed: kind=%s, operation=%s, duration=%s, error=%v",
					kind, operation, duration, err)
			} else {
				log.NewHelper(logger).WithContext(ctx).Infof("request completed: kind=%s, operation=%s, duration=%s",
					kind, operation, duration)
			}

			return reply, err
		}
	}
}
