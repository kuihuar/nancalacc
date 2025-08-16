package main

import (
	"context"
	"fmt"
	"log"
	"nancalacc/internal/conf"
	"nancalacc/internal/otel"
	"time"

	"go.opentelemetry.io/otel/attribute"
)

func main() {
	log.Println("🚀 开始测试 OpenTelemetry 集成...")

	// 创建配置
	config := &conf.OpenTelemetry{
		Enabled:        true,
		ServiceName:    "nancalacc-test",
		ServiceVersion: "1.0.0",
		Environment:    "test",
		Traces: &conf.TracesConfig{
			Enabled: true,
			Jaeger: &conf.JaegerConfig{
				Enabled:  true,
				Endpoint: "http://192.168.1.142:14268/api/traces",
			},
		},
		Logs: &conf.LogsConfig{
			Enabled: true,
			Level:   "debug",
			Format:  "json",
		},
	}

	log.Printf("📋 配置信息:")
	log.Printf("   - 启用状态: %v", config.Enabled)
	log.Printf("   - 服务名称: %s", config.ServiceName)
	log.Printf("   - 追踪启用: %v", config.Traces.Enabled)
	log.Printf("   - Jaeger启用: %v", config.Traces.Jaeger.Enabled)
	log.Printf("   - Jaeger端点: %s", config.Traces.Jaeger.Endpoint)
	log.Printf("   - 日志启用: %v", config.Logs.Enabled)

	// 创建 OpenTelemetry 集成器
	log.Println("🔧 创建 OpenTelemetry 集成器...")
	integration := otel.NewIntegration(otel.NewConfigFromConf(config))

	// 初始化集成器
	log.Println("🚀 初始化 OpenTelemetry 集成器...")
	ctx := context.Background()
	if err := integration.Init(ctx); err != nil {
		log.Printf("❌ 初始化失败: %v", err)
		return
	}
	log.Println("✅ 初始化成功")

	// 创建中间件
	log.Println("🔧 创建 HTTP 中间件...")
	httpMiddleware := integration.CreateHTTPMiddleware()
	log.Printf("📊 HTTP 中间件数量: %d", len(httpMiddleware))

	log.Println("🔧 创建 gRPC 中间件...")
	grpcMiddleware := integration.CreateGRPCMiddleware()
	log.Printf("📊 gRPC 中间件数量: %d", len(grpcMiddleware))

	// 获取服务实例
	service := integration.GetService()
	log.Printf("🔧 获取到服务实例: %T", service)

	// 获取追踪器
	tracer := service.GetTracer()
	log.Printf("🔧 获取到追踪器: %T", tracer)

	// 创建一些测试 span
	fmt.Println("🔍 创建测试 span...")

	// 创建测试 span
	ctx, span := tracer.Start(context.Background(), "test-operation")
	defer span.End()

	// 添加一些属性
	span.SetAttributes(
		attribute.String("test.key", "test-value"),
		attribute.Int("test.number", 42),
		attribute.Bool("test.flag", true),
	)

	// 模拟一些工作
	time.Sleep(100 * time.Millisecond)

	// 创建子 span
	_, childSpan := tracer.Start(ctx, "child-operation")
	childSpan.SetAttributes(attribute.String("child.key", "child-value"))
	time.Sleep(50 * time.Millisecond)
	childSpan.End()

	// 再等待一段时间，确保数据发送
	time.Sleep(200 * time.Millisecond)

	fmt.Println("✅ 测试 span 创建完成")

	// 等待一段时间，确保数据发送到 Jaeger
	fmt.Println("⏳ 等待数据发送到 Jaeger...")
	time.Sleep(3 * time.Second)

	fmt.Println("🔄 关闭 OpenTelemetry 集成器...")
	if err := integration.Shutdown(ctx); err != nil {
		log.Printf("❌ 关闭失败: %v", err)
		return
	}
	log.Println("✅ 关闭成功")

	log.Println("🎉 测试完成！")
}
