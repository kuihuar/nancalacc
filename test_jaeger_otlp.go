package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"nancalacc/internal/conf"
	"nancalacc/internal/service"
)

func main() {
	// 加载配置
	bc, err := conf.Load("configs")
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// 创建日志器
	logger := log.NewStdLogger(os.Stdout)
	helper := log.NewHelper(logger)

	helper.Info("🔧 开始初始化 OpenTelemetry 服务")
	helper.Infof("📊 配置信息:")
	helper.Infof("   - 服务名称: %s", bc.Otel.ServiceName)
	helper.Infof("   - 服务版本: %s", bc.Otel.ServiceVersion)
	helper.Infof("   - 环境: %s", bc.Otel.Environment)
	helper.Infof("   - 追踪启用: %v", bc.Otel.Traces.Enabled)
	helper.Infof("   - OTLP启用: %v", bc.Otel.Traces.Otlp.Enabled)
	helper.Infof("   - OTLP端点: %s", bc.Otel.Traces.Otlp.Endpoint)
	helper.Infof("   - Jaeger启用: %v", bc.Otel.Traces.Jaeger.Enabled)
	helper.Infof("   - Jaeger端点: %s", bc.Otel.Traces.Jaeger.Endpoint)
	helper.Infof("   - Metrics启用: %v", bc.Otel.Metrics.Enabled)

	// 初始化统一的 OpenTelemetry 服务
	otelService := service.NewUnifiedOTelService(bc.Otel, logger)

	// 设置更短的超时时间用于测试
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := otelService.Init(ctx); err != nil {
		helper.Warnf("⚠️  OpenTelemetry 初始化失败: %v", err)
		helper.Info("🔄 将使用 stdout 导出器进行本地测试")

		// 如果初始化失败，尝试使用 stdout 导出器
		bc.Otel.Traces.Otlp.Enabled = false
		bc.Otel.Traces.Jaeger.Enabled = false

		otelService = service.NewUnifiedOTelService(bc.Otel, logger)
		if err := otelService.Init(context.Background()); err != nil {
			helper.Errorf("❌ 连 stdout 导出器也失败了: %v", err)
			os.Exit(1)
		}
	} else {
		helper.Info("✅ OpenTelemetry 初始化成功")
	}

	defer otelService.Shutdown(context.Background())

	// 获取 tracer
	tracer := otelService.GetTracer()
	if tracer == nil {
		helper.Error("❌ Tracer 为空")
		os.Exit(1)
	}

	// 测试 1: 创建简单的 span
	helper.Info("🧪 测试 1: 创建简单的 span")
	ctx, span := tracer.Start(context.Background(), "jaeger-test-simple-span")
	span.SetAttributes(
		attribute.String("test.type", "jaeger-otlp"),
		attribute.String("test.message", "这是通过 OTLP 发送到 Jaeger 2.9 的测试span"),
		attribute.String("jaeger.version", "2.9"),
		attribute.String("protocol", "otlp"),
		attribute.String("endpoint", bc.Otel.Traces.Otlp.Endpoint),
	)
	time.Sleep(100 * time.Millisecond)
	span.End()

	// 测试 2: 创建带事件的 span
	helper.Info("🧪 测试 2: 创建带事件的 span")
	ctx, eventSpan := tracer.Start(context.Background(), "jaeger-test-event-span")
	eventSpan.SetAttributes(
		attribute.String("test.type", "jaeger-event"),
		attribute.String("test.message", "这是带事件的 Jaeger span"),
	)

	// 添加事件
	eventSpan.AddEvent("jaeger.test.start", trace.WithAttributes(
		attribute.String("test.id", "jaeger-001"),
		attribute.String("test.name", "OTLP to Jaeger"),
	))

	time.Sleep(75 * time.Millisecond)

	eventSpan.AddEvent("jaeger.test.complete", trace.WithAttributes(
		attribute.String("status", "success"),
		attribute.String("jaeger.backend", "memory"),
	))

	eventSpan.End()

	// 测试 3: 创建嵌套的 spans
	helper.Info("🧪 测试 3: 创建嵌套的 spans")
	ctx, parentSpan := tracer.Start(context.Background(), "jaeger-parent-span")
	parentSpan.SetAttributes(
		attribute.String("test.type", "jaeger-parent"),
		attribute.String("test.message", "这是父span"),
	)

	// 创建子span
	ctx, childSpan := tracer.Start(ctx, "jaeger-child-span")
	childSpan.SetAttributes(
		attribute.String("test.type", "jaeger-child"),
		attribute.String("test.message", "这是子span"),
	)
	time.Sleep(50 * time.Millisecond)
	childSpan.End()

	parentSpan.End()

	// 测试 4: 创建带错误的 span
	helper.Info("🧪 测试 4: 创建带错误的 span")
	ctx, errorSpan := tracer.Start(context.Background(), "jaeger-error-span")
	errorSpan.SetAttributes(
		attribute.String("test.type", "jaeger-error"),
		attribute.String("test.message", "这是带错误的 Jaeger span"),
	)

	// 模拟错误
	errorSpan.RecordError(fmt.Errorf("这是一个测试错误"))
	errorSpan.SetStatus(codes.Error, "测试错误")

	errorSpan.End()

	helper.Info("✅ 所有测试完成！")

	if bc.Otel.Traces.Otlp.Enabled {
		helper.Info("🌐 Jaeger UI 地址: http://192.168.1.142:16686")
		helper.Info("📊 在 Jaeger UI 中搜索服务名: nancalacc")
		helper.Info("🔗 OTLP 端点: " + bc.Otel.Traces.Otlp.Endpoint)
	} else {
		helper.Info("📝 使用 stdout 导出器，数据已输出到控制台")
		helper.Info("💡 要连接到 Jaeger 2.9，请确保:")
		helper.Info("   1. Jaeger 2.9 正在运行")
		helper.Info("   2. OTLP 端点 " + bc.Otel.Traces.Otlp.Endpoint + " 可访问")
		helper.Info("   3. 网络连接正常")
	}
}
