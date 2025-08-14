package middleware

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
)

// LoggingMiddleware 创建通用日志中间件
func LoggingMiddleware(logger log.Logger) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			var (
				startTime = time.Now()
				method    = getMethod(ctx, req)
				path      = getPath(ctx, req)
			)

			// 获取追踪信息
			traceID := GetTraceID(ctx)
			spanID := GetSpanID(ctx)
			requestID := GetRequestID(ctx)

			// 记录请求开始日志
			log.NewHelper(logger).WithContext(ctx).Infof("Request started: trace_id=%s, span_id=%s, request_id=%s, method=%s, path=%s",
				traceID, spanID, requestID, method, path)

			// 执行实际的处理器
			reply, err = handler(ctx, req)

			// 记录请求结束日志
			duration := time.Since(startTime)
			log.NewHelper(logger).WithContext(ctx).Infof("Request completed: trace_id=%s, span_id=%s, request_id=%s, method=%s, path=%s, duration=%s, error=%v",
				traceID, spanID, requestID, method, path, duration, err)

			return reply, err
		}
	}
}
