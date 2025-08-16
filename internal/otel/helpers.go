package otel

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

// TraceHelper 追踪辅助函数
type TraceHelper struct {
	tracer trace.Tracer
}

// NewTraceHelper 创建追踪辅助函数
func NewTraceHelper(tracer trace.Tracer) *TraceHelper {
	return &TraceHelper{tracer: tracer}
}

// StartSpan 开始一个新的 span
func (h *TraceHelper) StartSpan(ctx context.Context, name string, attrs ...attribute.KeyValue) (context.Context, trace.Span) {
	return h.tracer.Start(ctx, name, trace.WithAttributes(attrs...))
}

// AddEvent 添加事件到当前 span
func AddEvent(ctx context.Context, name string, attrs ...attribute.KeyValue) {
	if span := trace.SpanFromContext(ctx); span != nil {
		span.AddEvent(name, trace.WithAttributes(attrs...))
	}
}

// SetAttributes 设置属性到当前 span
func SetAttributes(ctx context.Context, attrs ...attribute.KeyValue) {
	if span := trace.SpanFromContext(ctx); span != nil {
		span.SetAttributes(attrs...)
	}
}

// RecordError 记录错误到当前 span
func RecordError(ctx context.Context, err error, attrs ...attribute.KeyValue) {
	if span := trace.SpanFromContext(ctx); span != nil {
		span.RecordError(err, trace.WithAttributes(attrs...))
	}
}

// MetricsHelper 指标辅助函数
type MetricsHelper struct {
	meter metric.Meter
}

// NewMetricsHelper 创建指标辅助函数
func NewMetricsHelper(meter metric.Meter) *MetricsHelper {
	return &MetricsHelper{meter: meter}
}

// Counter 计数器
type Counter struct {
	counter metric.Int64Counter
}

// NewCounter 创建计数器
func (h *MetricsHelper) NewCounter(name, description string, unit string) *Counter {
	counter, err := h.meter.Int64Counter(name, metric.WithDescription(description), metric.WithUnit(unit))
	if err != nil {
		// 在实际应用中，你可能想要记录这个错误
		return nil
	}
	return &Counter{counter: counter}
}

// Add 增加计数器值
func (c *Counter) Add(ctx context.Context, value int64, attrs ...attribute.KeyValue) {
	if c.counter != nil {
		c.counter.Add(ctx, value, metric.WithAttributes(attrs...))
	}
}

// Histogram 直方图
type Histogram struct {
	histogram metric.Float64Histogram
}

// NewHistogram 创建直方图
func (h *MetricsHelper) NewHistogram(name, description string, unit string) *Histogram {
	histogram, err := h.meter.Float64Histogram(name, metric.WithDescription(description), metric.WithUnit(unit))
	if err != nil {
		return nil
	}
	return &Histogram{histogram: histogram}
}

// Record 记录直方图值
func (h *Histogram) Record(ctx context.Context, value float64, attrs ...attribute.KeyValue) {
	if h.histogram != nil {
		h.histogram.Record(ctx, value, metric.WithAttributes(attrs...))
	}
}

// Attribute 创建属性键值对
func Attribute(key string, value interface{}) attribute.KeyValue {
	switch v := value.(type) {
	case string:
		return attribute.String(key, v)
	case int:
		return attribute.Int(key, v)
	case int64:
		return attribute.Int64(key, v)
	case float64:
		return attribute.Float64(key, v)
	case bool:
		return attribute.Bool(key, v)
	default:
		return attribute.String(key, "unknown")
	}
}

// Gauge 仪表盘
type Gauge struct {
	gauge metric.Float64ObservableGauge
}

// NewGauge 创建仪表盘
func (h *MetricsHelper) NewGauge(name, description string, unit string, callback func(context.Context, metric.Float64Observer) error) *Gauge {
	gauge, err := h.meter.Float64ObservableGauge(name, metric.WithDescription(description), metric.WithUnit(unit), metric.WithFloat64Callback(callback))
	if err != nil {
		return nil
	}
	return &Gauge{gauge: gauge}
}

// LogHelper 日志辅助函数
type LogHelper struct {
	logger log.Logger
}

// NewLogHelper 创建日志辅助函数
func NewLogHelper(logger log.Logger) *LogHelper {
	return &LogHelper{logger: logger}
}

// Info 记录信息日志
func (h *LogHelper) Info(ctx context.Context, msg string, attrs ...attribute.KeyValue) {
	// 记录日志
	if h.logger != nil {
		record := log.Record{}
		record.SetTimestamp(time.Now())
		record.SetSeverity(log.SeverityInfo)
		record.SetBody(log.StringValue(msg))
		for _, attr := range attrs {
			record.AddAttributes(log.KeyValueFromAttribute(attr))
		}
		h.logger.Emit(ctx, record)
	}
}

// Error 记录错误日志
func (h *LogHelper) Error(ctx context.Context, msg string, attrs ...attribute.KeyValue) {
	if h.logger != nil {
		record := log.Record{}
		record.SetTimestamp(time.Now())
		record.SetSeverity(log.SeverityError)
		record.SetBody(log.StringValue(msg))
		for _, attr := range attrs {
			record.AddAttributes(log.KeyValueFromAttribute(attr))
		}
		h.logger.Emit(ctx, record)
	}
}

// Debug 记录调试日志
func (h *LogHelper) Debug(ctx context.Context, msg string, attrs ...attribute.KeyValue) {
	if h.logger != nil {
		record := log.Record{}
		record.SetTimestamp(time.Now())
		record.SetSeverity(log.SeverityDebug)
		record.SetBody(log.StringValue(msg))
		for _, attr := range attrs {
			record.AddAttributes(log.KeyValueFromAttribute(attr))
		}
		h.logger.Emit(ctx, record)
	}
}

// Warn 记录警告日志
func (h *LogHelper) Warn(ctx context.Context, msg string, attrs ...attribute.KeyValue) {
	if h.logger != nil {
		record := log.Record{}
		record.SetTimestamp(time.Now())
		record.SetSeverity(log.SeverityWarn)
		record.SetBody(log.StringValue(msg))
		for _, attr := range attrs {
			record.AddAttributes(log.KeyValueFromAttribute(attr))
		}
		h.logger.Emit(ctx, record)
	}
}

// TimeOperation 计时操作辅助函数
func TimeOperation(ctx context.Context, name string, operation func(context.Context) error) error {
	start := time.Now()

	// 创建 span
	ctx, span := trace.SpanFromContext(ctx).TracerProvider().Tracer("nancalacc").Start(ctx, name)
	defer span.End()

	// 执行操作
	err := operation(ctx)

	// 记录执行时间
	duration := time.Since(start)
	span.SetAttributes(attribute.Int64("duration_ms", duration.Milliseconds()))

	// 记录指标 - 使用全局 meter
	if meter := otel.GetMeterProvider().Meter("nancalacc"); meter != nil {
		if histogram, _ := meter.Float64Histogram("operation_duration_seconds"); histogram != nil {
			histogram.Record(ctx, duration.Seconds(), metric.WithAttributes(attribute.String("operation", name)))
		}
	}

	return err
}

// CommonAttributes 常用属性
var CommonAttributes = struct {
	ServiceName    attribute.Key
	ServiceVersion attribute.Key
	Environment    attribute.Key
	Operation      attribute.Key
	Status         attribute.Key
	Error          attribute.Key
	Duration       attribute.Key
}{
	ServiceName:    attribute.Key("service.name"),
	ServiceVersion: attribute.Key("service.version"),
	Environment:    attribute.Key("environment"),
	Operation:      attribute.Key("operation"),
	Status:         attribute.Key("status"),
	Error:          attribute.Key("error"),
	Duration:       attribute.Key("duration"),
}
