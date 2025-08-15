package otel

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

// Service OpenTelemetry 服务
type Service struct {
	tracer trace.Tracer
	meter  metric.Meter
	logger log.Logger
}

// NewService 创建 OpenTelemetry 服务
func NewService(tracer trace.Tracer, meter metric.Meter, logger log.Logger) *Service {
	return &Service{
		tracer: tracer,
		meter:  meter,
		logger: logger,
	}
}

// Tracer 获取追踪器
func (s *Service) Tracer() trace.Tracer {
	return s.tracer
}

// Meter 获取指标器
func (s *Service) Meter() metric.Meter {
	return s.meter
}

// Logger 获取日志器
func (s *Service) Logger() log.Logger {
	return s.logger
}

// Log 记录日志
func (s *Service) Log(ctx context.Context, level log.Severity, message string, attrs ...log.KeyValue) {
	if s.logger != nil {
		record := log.Record{}
		record.SetTimestamp(time.Now())
		record.SetSeverity(level)
		record.SetBody(log.StringValue(message))
		for _, attr := range attrs {
			record.AddAttributes(attr)
		}
		s.logger.Emit(ctx, record)
	}
}

// LogInfo 记录信息日志
func (s *Service) LogInfo(ctx context.Context, message string, attrs ...log.KeyValue) {
	s.Log(ctx, log.SeverityInfo, message, attrs...)
}

// LogError 记录错误日志
func (s *Service) LogError(ctx context.Context, message string, attrs ...log.KeyValue) {
	s.Log(ctx, log.SeverityError, message, attrs...)
}

// LogWarn 记录警告日志
func (s *Service) LogWarn(ctx context.Context, message string, attrs ...log.KeyValue) {
	s.Log(ctx, log.SeverityWarn, message, attrs...)
}

// LogDebug 记录调试日志
func (s *Service) LogDebug(ctx context.Context, message string, attrs ...log.KeyValue) {
	s.Log(ctx, log.SeverityDebug, message, attrs...)
}

// Init 初始化服务
func (s *Service) Init(ctx context.Context) error {
	// 初始化 OpenTelemetry 服务
	return nil
}

// Shutdown 关闭服务
func (s *Service) Shutdown(ctx context.Context) error {
	// 关闭 OpenTelemetry 服务
	return nil
}

// GetTracer 获取追踪器
func (s *Service) GetTracer() trace.Tracer {
	if s.tracer != nil {
		return s.tracer
	}
	return trace.NewNoopTracerProvider().Tracer("nancalacc")
}
