package main

import (
	stdlog "log"

	"nancalacc/internal/otel"

	klog "github.com/go-kratos/kratos/v2/log"
	otellog "go.opentelemetry.io/otel/log"
)

func main() {
	// 创建配置 - 使用JSON格式
	config := otel.DefaultConfig()
	config.Logs.Enabled = true
	config.Logs.Level = "info"
	config.Logs.Output = "stdout"
	config.Logs.Format = "json" // 设置为JSON格式

	// 创建一个简单的OpenTelemetry Logger
	otellogger := otellog.NewNoopLoggerProvider().Logger("test")

	// 创建日志适配器
	adapter := otel.NewKratosLoggerAdapter(otellogger, config)

	// 测试不同级别的日志
	logger := klog.With(adapter,
		"service", "test",
		"version", "1.0.0",
	)

	logger.Log(klog.LevelInfo, "msg", "这是一条信息日志")
	logger.Log(klog.LevelWarn, "msg", "这是一条警告日志")
	logger.Log(klog.LevelError, "msg", "这是一条错误日志")
	logger.Log(klog.LevelDebug, "msg", "这是一条调试日志（应该被过滤）")

	// 测试结构化日志
	logger.Log(klog.LevelInfo,
		"msg", "用户登录成功",
		"user_id", "12345",
		"ip", "192.168.1.100",
		"timestamp", "2024-01-01T12:00:00Z",
		"status", "success",
		"duration_ms", 150,
	)

	// 测试复杂数据结构
	logger.Log(klog.LevelInfo,
		"msg", "API请求处理完成",
		"method", "POST",
		"path", "/api/users",
		"status_code", 200,
		"response_time", 45.67,
		"user_agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)",
		"request_id", "req-123456789",
	)

	stdlog.Println("JSON格式日志测试完成，请查看上面的JSON输出")
}
