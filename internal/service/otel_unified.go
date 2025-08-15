package service

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"nancalacc/internal/conf"

	"go.opentelemetry.io/otel/attribute"
)

// UnifiedOTelService 统一的 OpenTelemetry 服务
type UnifiedOTelService struct {
	logger *log.Helper
	config *conf.OpenTelemetry
	tracer trace.Tracer
}

// NewUnifiedOTelService 创建统一的 OpenTelemetry 服务
func NewUnifiedOTelService(c *conf.OpenTelemetry, logger log.Logger) *UnifiedOTelService {
	return &UnifiedOTelService{
		logger: log.NewHelper(logger),
		config: c,
	}
}

// Init 初始化 OpenTelemetry
func (s *UnifiedOTelService) Init(ctx context.Context) error {
	if !s.config.Enabled {
		s.logger.Info("OpenTelemetry is disabled")
		return nil
	}

	s.logger.Info("Initializing OpenTelemetry...")

	// 1. 创建资源
	res, err := s.createResource(ctx)
	if err != nil {
		return fmt.Errorf("failed to create resource: %w", err)
	}

	// 2. 创建导出器
	exporter, err := s.createExporter(ctx)
	if err != nil {
		return fmt.Errorf("failed to create exporter: %w", err)
	}

	// 3. 创建 TracerProvider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()), // 采样所有追踪
	)

	// 4. 设置全局配置
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	// 5. 创建 tracer
	s.tracer = tp.Tracer(s.config.ServiceName)

	s.logger.Infof("OpenTelemetry initialized successfully with service: %s, version: %s, environment: %s",
		s.config.ServiceName, s.config.ServiceVersion, s.config.Environment)

	return nil
}

// createResource 创建资源
func (s *UnifiedOTelService) createResource(ctx context.Context) (*resource.Resource, error) {
	attrs := []attribute.KeyValue{
		semconv.ServiceName(s.config.ServiceName),
		semconv.ServiceVersion(s.config.ServiceVersion),
		semconv.DeploymentEnvironment(s.config.Environment),
	}

	return resource.New(ctx, resource.WithAttributes(attrs...))
}

// createExporter 创建导出器
func (s *UnifiedOTelService) createExporter(ctx context.Context) (sdktrace.SpanExporter, error) {
	// 优先级：Jaeger > OTLP > Stdout
	if s.config.Traces.Jaeger != nil && s.config.Traces.Jaeger.Enabled {
		return s.createJaegerExporter(ctx)
	}

	if s.config.Traces.Otlp != nil && s.config.Traces.Otlp.Enabled {
		return s.createOTLPExporter(ctx)
	}

	// 默认使用 stdout 导出器（用于开发环境）
	s.logger.Warn("No exporter configured, using stdout exporter for development")
	return s.createStdoutExporter()
}

// createJaegerExporter 创建 Jaeger 导出器
func (s *UnifiedOTelService) createJaegerExporter(ctx context.Context) (sdktrace.SpanExporter, error) {
	s.logger.Infof("Creating Jaeger exporter with endpoint: %s", s.config.Traces.Jaeger.Endpoint)

	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(
		jaeger.WithEndpoint(s.config.Traces.Jaeger.Endpoint),
	))
	if err != nil {
		return nil, fmt.Errorf("failed to create Jaeger exporter: %w", err)
	}

	s.logger.Info("Jaeger exporter created successfully")
	return exporter, nil
}

// createOTLPExporter 创建 OTLP 导出器
func (s *UnifiedOTelService) createOTLPExporter(ctx context.Context) (sdktrace.SpanExporter, error) {
	s.logger.Infof("Creating OTLP exporter with endpoint: %s", s.config.Traces.Otlp.Endpoint)

	timeout := time.Duration(s.config.Traces.Otlp.Timeout) * time.Second
	conn, err := grpc.DialContext(ctx, s.config.Traces.Otlp.Endpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithTimeout(timeout),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to OTLP endpoint: %w", err)
	}

	exporter, err := otlptrace.New(ctx, otlptracegrpc.NewClient(
		otlptracegrpc.WithGRPCConn(conn),
	))
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP exporter: %w", err)
	}

	s.logger.Info("OTLP exporter created successfully")
	return exporter, nil
}

// createStdoutExporter 创建标准输出导出器
func (s *UnifiedOTelService) createStdoutExporter() (sdktrace.SpanExporter, error) {
	s.logger.Info("Creating stdout exporter")

	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		return nil, fmt.Errorf("failed to create stdout exporter: %w", err)
	}

	s.logger.Info("Stdout exporter created successfully")
	return exporter, nil
}

// GetTracer 获取 tracer
func (s *UnifiedOTelService) GetTracer() trace.Tracer {
	return s.tracer
}

// Shutdown 关闭 OpenTelemetry
func (s *UnifiedOTelService) Shutdown(ctx context.Context) error {
	if !s.config.Enabled {
		return nil
	}

	s.logger.Info("Shutting down OpenTelemetry...")

	if tp, ok := otel.GetTracerProvider().(*sdktrace.TracerProvider); ok {
		if err := tp.Shutdown(ctx); err != nil {
			return fmt.Errorf("failed to shutdown tracer provider: %w", err)
		}
	}

	s.logger.Info("OpenTelemetry shutdown completed")
	return nil
}

// IsEnabled 检查是否启用
func (s *UnifiedOTelService) IsEnabled() bool {
	return s.config.Enabled
}

// GetConfig 获取配置
func (s *UnifiedOTelService) GetConfig() *conf.OpenTelemetry {
	return s.config
}
