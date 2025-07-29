package tracer

import (
	"context"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

type TracerManager struct {
	tracerProvider *trace.TracerProvider
}

func NewTracerManager() *TracerManager {
	return &TracerManager{}
}

// 初始化 Tracer
func (tm *TracerManager) Init(env, name string) error {
	// io.Discard
	// os.Stdout

	exporter, err := stdouttrace.New(stdouttrace.WithWriter(os.Stdout))
	if err != nil {
		return err // 返回错误而不是直接退出
	}

	tm.tracerProvider = trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("name"),
			attribute.String("environment", "env"),
		)),
	)
	otel.SetTracerProvider(tm.tracerProvider)

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))
	return nil
}

// 优雅关闭 Tracer
func (tm *TracerManager) Shutdown() error {
	if tm.tracerProvider != nil {
		return tm.tracerProvider.Shutdown(context.Background())
	}
	return nil
}
