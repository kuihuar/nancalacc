package main

// import (
// 	"context"
// 	"os"
// 	"time"

// 	"nancalacc/internal/conf"
// 	"nancalacc/internal/otel"
// 	"nancalacc/internal/service"

// 	"github.com/go-kratos/kratos/v2/log"
// 	"go.opentelemetry.io/otel/attribute"
// 	"go.opentelemetry.io/otel/codes"
// 	"go.opentelemetry.io/otel/trace"
// )

// func main() {
// 	// 加载配置
// 	bc, err := conf.Load("configs")
// 	if err != nil {
// 		log.Fatalf("Failed to load config: %v", err)
// 	}

// 	// 创建日志器
// 	logger := log.NewStdLogger(os.Stdout)
// 	helper := log.NewHelper(logger)

// 	// 初始化统一的 OpenTelemetry 服务
// 	otelService := service.NewUnifiedOTelService(bc.Otel, logger)
// 	if err := otelService.Init(context.Background()); err != nil {
// 		log.Fatalf("Failed to init OpenTelemetry: %v", err)
// 	}
// 	defer otelService.Shutdown(context.Background())

// 	// 初始化 OpenTelemetry 集成
// 	adapter := otel.NewConfigAdapter()
// 	config := adapter.FromBootstrap(bc)
// 	integration := otel.NewIntegration(config)
// 	if err := integration.Init(context.Background()); err != nil {
// 		log.Fatalf("Failed to init integration: %v", err)
// 	}
// 	defer integration.Shutdown(context.Background())

// 	// 获取 tracer
// 	tracer := otelService.GetTracer()
// 	if tracer == nil {
// 		log.Fatal("Tracer is nil")
// 	}

// 	helper.Info("🔧 OpenTelemetry 服务已初始化")
// 	helper.Infof("📊 配置信息:")
// 	helper.Infof("   - 服务名称: %s", bc.Otel.ServiceName)
// 	helper.Infof("   - 服务版本: %s", bc.Otel.ServiceVersion)
// 	helper.Infof("   - 环境: %s", bc.Otel.Environment)
// 	helper.Infof("   - 追踪启用: %v", bc.Otel.Traces.Enabled)
// 	helper.Infof("   - Jaeger启用: %v", bc.Otel.Traces.Jaeger.Enabled)
// 	helper.Infof("   - OTLP启用: %v", bc.Otel.Traces.Otlp.Enabled)
// 	helper.Infof("   - Metrics启用: %v", bc.Otel.Metrics.Enabled)
// 	helper.Infof("   - Prometheus启用: %v", bc.Otel.Metrics.Prometheus.Enabled)

// 	// 测试 1: 创建简单的 span
// 	helper.Info("🧪 测试 1: 创建简单的 span")
// 	ctx, span := tracer.Start(context.Background(), "test-simple-span")
// 	span.SetAttributes(
// 		attribute.String("test.type", "simple"),
// 		attribute.String("test.message", "这是一个简单的测试span"),
// 	)
// 	time.Sleep(100 * time.Millisecond)
// 	span.End()

// 	// 测试 2: 创建嵌套的 spans
// 	helper.Info("🧪 测试 2: 创建嵌套的 spans")
// 	ctx, parentSpan := tracer.Start(context.Background(), "test-parent-span")
// 	parentSpan.SetAttributes(
// 		attribute.String("test.type", "parent"),
// 		attribute.String("test.message", "这是父span"),
// 	)

// 	// 创建子span
// 	ctx, childSpan := tracer.Start(ctx, "test-child-span")
// 	childSpan.SetAttributes(
// 		attribute.String("test.type", "child"),
// 		attribute.String("test.message", "这是子span"),
// 	)
// 	time.Sleep(50 * time.Millisecond)
// 	childSpan.End()

// 	// 创建另一个子span
// 	ctx, childSpan2 := tracer.Start(ctx, "test-child-span-2")
// 	childSpan2.SetAttributes(
// 		attribute.String("test.type", "child2"),
// 		attribute.String("test.message", "这是第二个子span"),
// 	)
// 	time.Sleep(30 * time.Millisecond)
// 	childSpan2.End()

// 	parentSpan.End()

// 	// 测试 3: 创建带事件的 span
// 	helper.Info("🧪 测试 3: 创建带事件的 span")
// 	ctx, eventSpan := tracer.Start(context.Background(), "test-event-span")
// 	eventSpan.SetAttributes(
// 		attribute.String("test.type", "event"),
// 		attribute.String("test.message", "这是带事件的span"),
// 	)

// 	// 添加事件
// 	eventSpan.AddEvent("user.login", trace.WithAttributes(
// 		attribute.String("user.id", "12345"),
// 		attribute.String("user.name", "testuser"),
// 	))

// 	time.Sleep(75 * time.Millisecond)

// 	eventSpan.AddEvent("user.action", trace.WithAttributes(
// 		attribute.String("action", "click_button"),
// 		attribute.String("button.id", "submit"),
// 	))

// 	eventSpan.End()

// 	// 测试 4: 创建带错误的 span
// 	helper.Info("🧪 测试 4: 创建带错误的 span")
// 	ctx, errorSpan := tracer.Start(context.Background(), "test-error-span")
// 	errorSpan.SetAttributes(
// 		attribute.String("test.type", "error"),
// 		attribute.String("test.message", "这是带错误的span"),
// 	)

// 	// 模拟错误
// 	err = &testError{message: "这是一个测试错误"}
// 	errorSpan.RecordError(err)
// 	errorSpan.SetStatus(codes.Error, err.Error())
// 	errorSpan.End()

// 	// 测试 5: 创建长时间运行的 span
// 	helper.Info("🧪 测试 5: 创建长时间运行的 span")
// 	ctx, longSpan := tracer.Start(context.Background(), "test-long-span")
// 	longSpan.SetAttributes(
// 		attribute.String("test.type", "long"),
// 		attribute.String("test.message", "这是长时间运行的span"),
// 	)

// 	// 模拟长时间操作
// 	for i := 0; i < 3; i++ {
// 		time.Sleep(200 * time.Millisecond)
// 		longSpan.AddEvent("progress", trace.WithAttributes(
// 			attribute.Int("step", i+1),
// 			attribute.String("status", "processing"),
// 		))
// 	}

// 	longSpan.End()

// 	// 测试 6: 使用 Kratos 日志器
// 	helper.Info("🧪 测试 6: 使用 Kratos 日志器")
// 	kratosLogger := integration.CreateLogger()
// 	kratosLogger.Log(log.LevelInfo,
// 		"msg", "测试日志消息",
// 		"test_type", "kratos_logger",
// 		"timestamp", time.Now().Unix(),
// 	)

// 	kratosLogger.Log(log.LevelWarn,
// 		"msg", "测试警告消息",
// 		"test_type", "kratos_logger",
// 		"warning_code", "TEST_WARN_001",
// 	)

// 	kratosLogger.Log(log.LevelError,
// 		"msg", "测试错误消息",
// 		"test_type", "kratos_logger",
// 		"error_code", "TEST_ERROR_001",
// 	)

// 	// 测试 7: 创建多个并发 spans
// 	helper.Info("🧪 测试 7: 创建多个并发 spans")
// 	for i := 0; i < 5; i++ {
// 		go func(id int) {
// 			_, span := tracer.Start(context.Background(), "test-concurrent-span")
// 			span.SetAttributes(
// 				attribute.String("test.type", "concurrent"),
// 				attribute.Int("goroutine.id", id),
// 				attribute.String("test.message", "这是并发span"),
// 			)
// 			time.Sleep(time.Duration(100+id*50) * time.Millisecond)
// 			span.End()
// 		}(i)
// 	}

// 	// 等待所有并发操作完成
// 	time.Sleep(1 * time.Second)

// 	helper.Info("✅ 所有测试完成！")
// 	helper.Info("📋 测试总结:")
// 	helper.Info("   - 创建了 7 个不同类型的 spans")
// 	helper.Info("   - 测试了嵌套 span 结构")
// 	helper.Info("   - 测试了 span 事件")
// 	helper.Info("   - 测试了错误处理")
// 	helper.Info("   - 测试了长时间运行的操作")
// 	helper.Info("   - 测试了并发 span 创建")
// 	helper.Info("   - 测试了 Kratos 日志器集成")

// 	helper.Info("🔍 请检查以下端点以验证数据推送:")
// 	if bc.Otel.Traces.Jaeger.Enabled {
// 		helper.Infof("   - Jaeger UI: http://192.168.1.142:16686")
// 	}
// 	if bc.Otel.Traces.Otlp.Enabled {
// 		helper.Infof("   - OTLP 端点: %s", bc.Otel.Traces.Otlp.Endpoint)
// 	}
// 	if bc.Otel.Metrics.Prometheus.Enabled {
// 		helper.Infof("   - Prometheus: http://192.168.1.142:9090")
// 	}
// 	if bc.Otel.Logs.Loki.Enabled {
// 		helper.Infof("   - Loki: http://192.168.1.142:3100")
// 	}
// }

// // testError 测试错误类型
// type testError struct {
// 	message string
// }

// func (e *testError) Error() string {
// 	return e.message
// }
