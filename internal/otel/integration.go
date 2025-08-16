package otel

import (
	"context"
	stdlog "log"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"go.opentelemetry.io/otel/attribute"
	otellog "go.opentelemetry.io/otel/log"
)

// Integration OpenTelemetry 集成器
type Integration struct {
	service *Service
	config  *Config
}

// NewIntegration 创建集成器
func NewIntegration(config *Config) *Integration {
	stdlog.Printf("🔍 [DEBUG] Creating new OpenTelemetry integration")
	stdlog.Printf("🔍 [DEBUG] Config - Enabled: %v, Traces: %v, Logs: %v", config.Enabled, config.Traces.Enabled, config.Logs.Enabled)

	// 创建真正的 OpenTelemetry 服务
	service := NewService(
		config.GetTracer(),
		config.GetMeter(),
		config.GetLogger(),
	)

	stdlog.Printf("🔍 [DEBUG] OpenTelemetry service created")

	return &Integration{
		service: service,
		config:  config,
	}
}

// Init 初始化集成器
func (i *Integration) Init(ctx context.Context) error {
	return i.service.Init(ctx)
}

// Shutdown 关闭集成器
func (i *Integration) Shutdown(ctx context.Context) error {
	return i.service.Shutdown(ctx)
}

// GetService 获取服务实例
func (i *Integration) GetService() *Service {
	return i.service
}

// GetLogger 获取OpenTelemetry日志器
func (i *Integration) GetLogger() otellog.Logger {
	return i.service.Logger()
}

// CreateLogger 创建Kratos兼容的日志器
func (i *Integration) CreateLogger() log.Logger {
	return NewKratosLoggerAdapter(i.service.Logger(), i.config)
}

// CreateHTTPMiddleware 创建HTTP中间件
func (i *Integration) CreateHTTPMiddleware() []http.ServerOption {
	stdlog.Printf("🔍 [DEBUG] CreateHTTPMiddleware called, enabled: %v, traces: %v", i.config.Enabled, i.config.Traces.Enabled)

	if !i.config.Enabled || !i.config.Traces.Enabled {
		stdlog.Printf("🔍 [DEBUG] HTTP middleware disabled, returning nil")
		return nil
	}

	stdlog.Printf("🔍 [DEBUG] Creating HTTP middleware with tracing and logging")
	return []http.ServerOption{
		http.Middleware(
			tracing.Server(),
			i.createLoggingMiddleware(),
		),
	}
}

// CreateGRPCMiddleware 创建gRPC中间件
func (i *Integration) CreateGRPCMiddleware() []grpc.ServerOption {
	stdlog.Printf("🔍 [DEBUG] CreateGRPCMiddleware called, enabled: %v, traces: %v", i.config.Enabled, i.config.Traces.Enabled)

	if !i.config.Enabled || !i.config.Traces.Enabled {
		stdlog.Printf("🔍 [DEBUG] gRPC middleware disabled, returning nil")
		return nil
	}

	stdlog.Printf("🔍 [DEBUG] Creating gRPC middleware with tracing and logging")
	return []grpc.ServerOption{
		grpc.Middleware(
			tracing.Server(),
			i.createLoggingMiddleware(),
		),
	}
}

// createLoggingMiddleware 创建日志中间件
func (i *Integration) createLoggingMiddleware() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			// 获取请求信息
			var (
				operation string
				kind      string
				start     = time.Now()
			)
			if tr, ok := transport.FromServerContext(ctx); ok {
				operation = tr.Operation()
				kind = tr.Kind().String()
			}

			// 添加调试日志
			stdlog.Printf("🔍 [DEBUG] Creating span for operation: %s, kind: %s", operation, kind)

			// 创建span
			ctx, span := i.service.GetTracer().Start(ctx, operation)
			defer span.End()

			// 设置span属性
			span.SetAttributes(
				attribute.String("transport.kind", kind),
				attribute.String("transport.operation", operation),
			)

			stdlog.Printf("🔍 [DEBUG] Span created with attributes: transport.kind=%s, transport.operation=%s", kind, operation)

			// 记录请求开始日志
			if i.config.Logs.Enabled {
				stdlog.Printf("🔍 [DEBUG] Logs enabled, creating logger for request start")
				logger := i.CreateLogger()
				logger.Log(log.LevelInfo,
					"msg", "request started",
					"operation", operation,
					"kind", kind,
					"request", req,
				)
				stdlog.Printf("🔍 [DEBUG] Request start log recorded")
			} else {
				stdlog.Printf("🔍 [DEBUG] Logs disabled, skipping request start log")
			}

			// 执行请求
			reply, err = handler(ctx, req)

			// 记录请求结束日志
			if i.config.Logs.Enabled {
				stdlog.Printf("🔍 [DEBUG] Logs enabled, creating logger for request end")
				logger := i.CreateLogger()
				logger.Log(log.LevelInfo,
					"msg", "request finished",
					"operation", operation,
					"kind", kind,
					"duration", time.Since(start).String(),
					"error", err,
				)
				stdlog.Printf("🔍 [DEBUG] Request end log recorded, duration: %s", time.Since(start).String())
			} else {
				stdlog.Printf("🔍 [DEBUG] Logs disabled, skipping request end log")
			}

			// 记录错误到span
			if err != nil {
				stdlog.Printf("🔍 [DEBUG] Recording error to span: %v", err)
				span.RecordError(err)
			}

			stdlog.Printf("🔍 [DEBUG] Span ending for operation: %s", operation)
			return reply, err
		}
	}
}
