package main

// import (
// 	stdlog "log"

// 	"nancalacc/internal/otel"

// 	klog "github.com/go-kratos/kratos/v2/log"
// 	otellog "go.opentelemetry.io/otel/log"
// )

// func main() {
// 	// 创建配置
// 	config := otel.DefaultConfig()
// 	config.Logs.Enabled = true
// 	config.Logs.Level = "info"
// 	config.Logs.Output = "stdout"
// 	config.Logs.Format = "console"

// 	// 创建一个简单的OpenTelemetry Logger
// 	otellogger := otellog.NewNoopLoggerProvider().Logger("test")

// 	// 创建日志适配器
// 	adapter := otel.NewKratosLoggerAdapter(otellogger, config)

// 	// 测试不同级别的日志
// 	logger := klog.With(adapter,
// 		"service", "test",
// 		"version", "1.0.0",
// 	)

// 	logger.Log(klog.LevelInfo, "msg", "这是一条信息日志")
// 	logger.Log(klog.LevelWarn, "msg", "这是一条警告日志")
// 	logger.Log(klog.LevelError, "msg", "这是一条错误日志")
// 	logger.Log(klog.LevelDebug, "msg", "这是一条调试日志（应该被过滤）")

// 	// 测试结构化日志
// 	logger.Log(klog.LevelInfo,
// 		"msg", "用户登录成功",
// 		"user_id", "12345",
// 		"ip", "192.168.1.100",
// 		"timestamp", "2024-01-01T12:00:00Z",
// 	)

// 	stdlog.Println("测试完成，请查看上面的日志输出")
// }
