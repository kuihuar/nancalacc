package middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

// TracingConfig 追踪中间件配置
type TracingConfig struct {
	// 是否启用详细日志
	VerboseLogging bool
	// 是否记录请求体大小
	LogRequestSize bool
	// 是否记录响应体大小
	LogResponseSize bool
	// 是否记录请求头
	LogHeaders bool
	// 是否记录查询参数
	LogQueryParams bool
}

// DefaultTracingConfig 默认配置
func DefaultTracingConfig() *TracingConfig {
	return &TracingConfig{
		VerboseLogging:  false, // 默认不启用详细日志
		LogRequestSize:  false,
		LogResponseSize: false,
		LogHeaders:      false,
		LogQueryParams:  false,
	}
}

// TracingMiddleware 创建OpenTelemetry追踪中间件
func TracingMiddleware(logger log.Logger, config *TracingConfig) middleware.Middleware {
	if config == nil {
		config = DefaultTracingConfig()
	}

	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			var (
				startTime = time.Now()
				tr        = otel.GetTracerProvider().Tracer("nancalacc")
				prop      = otel.GetTextMapPropagator()
			)

			// 从请求中提取追踪上下文
			var span trace.Span
			if httpReq, ok := req.(*khttp.Request); ok {
				// 从HTTP请求头中提取追踪信息
				ctx = prop.Extract(ctx, propagation.HeaderCarrier(httpReq.Header))
			}

			// 创建新的span
			ctx, span = tr.Start(ctx, "request",
				trace.WithSpanKind(trace.SpanKindServer),
			)
			defer span.End()

			// 获取trace ID和span ID
			spanCtx := span.SpanContext()
			traceID := spanCtx.TraceID().String()
			spanID := spanCtx.SpanID().String()

			// 生成或提取request ID
			requestID := extractRequestID(ctx, req)
			if requestID == "" {
				requestID = generateRequestID()
			}

			// 将追踪信息添加到上下文
			ctx = context.WithValue(ctx, "trace_id", traceID)
			ctx = context.WithValue(ctx, "span_id", spanID)
			ctx = context.WithValue(ctx, "request_id", requestID)

			// 获取请求信息
			method := getMethod(ctx, req)
			path := getPath(ctx, req)
			operationName := fmt.Sprintf("%s %s", method, path)

			// 记录简化的请求开始日志
			if config.VerboseLogging {
				log.NewHelper(logger).WithContext(ctx).Infof("Request started: trace_id=%s, span_id=%s, request_id=%s, method=%s, path=%s",
					traceID, spanID, requestID, method, path)
			} else {
				log.NewHelper(logger).WithContext(ctx).Infof("Request: trace_id=%s, span_id=%s, name=%s",
					traceID, spanID, operationName)
			}

			// 执行实际的处理器
			reply, err = handler(ctx, req)

			// 记录简化的请求结束日志
			duration := time.Since(startTime)
			if config.VerboseLogging {
				log.NewHelper(logger).WithContext(ctx).Infof("Request completed: trace_id=%s, span_id=%s, request_id=%s, method=%s, path=%s, duration=%s, error=%v",
					traceID, spanID, requestID, method, path, duration, err)
			} else {
				log.NewHelper(logger).WithContext(ctx).Infof("Response: trace_id=%s, span_id=%s, name=%s, duration=%s, error=%v",
					traceID, spanID, operationName, duration, err)
			}

			// 设置span属性（保持详细信息用于追踪系统）
			span.SetAttributes(
				attribute.String("request_id", requestID),
				attribute.String("method", method),
				attribute.String("path", path),
				attribute.Int64("duration_ms", duration.Milliseconds()),
			)

			if err != nil {
				span.RecordError(err)
			}

			return reply, err
		}
	}
}

// extractRequestID 从上下文或请求头中提取request ID
func extractRequestID(ctx context.Context, req interface{}) string {
	// 首先从上下文中获取
	if requestID, ok := ctx.Value("request_id").(string); ok && requestID != "" {
		return requestID
	}

	// 从HTTP请求头中获取
	if httpReq, ok := req.(*khttp.Request); ok {
		if requestID := httpReq.Header.Get("X-Request-ID"); requestID != "" {
			return requestID
		}
	}

	return ""
}

// generateRequestID 生成新的request ID
func generateRequestID() string {
	return fmt.Sprintf("req_%d", time.Now().UnixNano())
}

// getMethod 获取请求方法
func getMethod(ctx context.Context, req interface{}) string {
	if httpReq, ok := req.(*khttp.Request); ok {
		return httpReq.Method
	}
	return "unknown"
}

// getPath 获取请求路径
func getPath(ctx context.Context, req interface{}) string {
	if httpReq, ok := req.(*khttp.Request); ok {
		return httpReq.URL.Path
	}
	return "unknown"
}

// GetTraceID 从上下文中获取trace ID
func GetTraceID(ctx context.Context) string {
	if traceID, ok := ctx.Value("trace_id").(string); ok {
		return traceID
	}
	return ""
}

// GetSpanID 从上下文中获取span ID
func GetSpanID(ctx context.Context) string {
	if spanID, ok := ctx.Value("span_id").(string); ok {
		return spanID
	}
	return ""
}

// GetRequestID 从上下文中获取request ID
func GetRequestID(ctx context.Context) string {
	if requestID, ok := ctx.Value("request_id").(string); ok {
		return requestID
	}
	return ""
}
