package service

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"nancalacc/internal/conf"
)

// ProviderSet 是otel服务的提供者集合

// OTelService OpenTelemetry服务
type OTelService struct {
	logger *log.Helper
	config *conf.OpenTelemetry
	tracer trace.Tracer
}

// NewOTelService 创建OpenTelemetry服务
func NewOTelService(c *conf.OpenTelemetry, logger log.Logger) *OTelService {
	return &OTelService{
		logger: log.NewHelper(logger),
		config: c,
	}
}

// Init 初始化OpenTelemetry
func (s *OTelService) Init(ctx context.Context) error {
	if !s.config.Enabled {
		s.logger.Info("OpenTelemetry is disabled")
		return nil
	}

	// 创建资源
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(s.config.ServiceName),
			semconv.ServiceVersion(s.config.ServiceVersion),
			semconv.DeploymentEnvironment(s.config.Environment),
		),
	)
	if err != nil {
		return err
	}

	// 创建导出器
	var exporter sdktrace.SpanExporter
	if s.config.Traces.Jaeger != nil && s.config.Traces.Jaeger.Enabled {
		exporter, err = jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(s.config.Traces.Jaeger.Endpoint)))
		if err != nil {
			return err
		}
		s.logger.Infof("Jaeger exporter initialized: %s", s.config.Traces.Jaeger.Endpoint)
	} else if s.config.Traces.Otlp != nil && s.config.Traces.Otlp.Enabled {
		timeout := time.Duration(s.config.Traces.Otlp.Timeout) * time.Second
		conn, err := grpc.DialContext(ctx, s.config.Traces.Otlp.Endpoint,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithBlock(),
			grpc.WithTimeout(timeout),
		)
		if err != nil {
			return err
		}
		exporter, err = otlptrace.New(ctx, otlptracegrpc.NewClient(otlptracegrpc.WithGRPCConn(conn)))
		if err != nil {
			return err
		}
		s.logger.Infof("OTLP exporter initialized: %s", s.config.Traces.Otlp.Endpoint)
	} else {
		s.logger.Warn("No exporter configured, using noop exporter")
		return nil
	}

	// 创建TracerProvider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)

	// 设置全局TracerProvider
	otel.SetTracerProvider(tp)

	// 设置全局Propagator
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	// 创建tracer
	s.tracer = tp.Tracer(s.config.ServiceName)

	s.logger.Info("OpenTelemetry initialized successfully")
	return nil
}

// GetTracer 获取tracer
func (s *OTelService) GetTracer() trace.Tracer {
	return s.tracer
}

// Shutdown 关闭OpenTelemetry
func (s *OTelService) Shutdown(ctx context.Context) error {
	if s.config.Enabled {
		if tp, ok := otel.GetTracerProvider().(*sdktrace.TracerProvider); ok {
			return tp.Shutdown(ctx)
		}
	}
	return nil
}
