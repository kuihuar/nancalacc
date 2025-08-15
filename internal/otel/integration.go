package otel

import (
	"context"

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
	service := NewService(
		config.GetTracer(),
		config.GetMeter(),
		config.GetLogger(),
	)
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
	if !i.config.Enabled || !i.config.Traces.Enabled {
		return nil
	}

	return []http.ServerOption{
		http.Middleware(
			tracing.Server(),
			i.createLoggingMiddleware(),
		),
	}
}

// CreateGRPCMiddleware 创建gRPC中间件
func (i *Integration) CreateGRPCMiddleware() []grpc.ServerOption {
	if !i.config.Enabled || !i.config.Traces.Enabled {
		return nil
	}

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
			)
			if tr, ok := transport.FromServerContext(ctx); ok {
				operation = tr.Operation()
				kind = tr.Kind().String()
			}

			// 创建span
			ctx, span := i.service.GetTracer().Start(ctx, operation)
			defer span.End()

			// 设置span属性
			span.SetAttributes(
				attribute.String("transport.kind", kind),
				attribute.String("transport.operation", operation),
			)

			// 执行请求
			reply, err = handler(ctx, req)

			// 记录错误
			if err != nil {
				span.RecordError(err)
			}

			return reply, err
		}
	}
}
