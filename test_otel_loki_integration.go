package main

// import (
// 	"context"
// 	"nancalacc/internal/conf"
// 	"nancalacc/internal/otel"
// 	"time"

// 	"github.com/go-kratos/kratos/v2/log"
// )

// func main() {
// 	// 加载配置
// 	bc, err := conf.Load("configs/config.yaml")
// 	if err != nil {
// 		panic("failed to load config: " + err.Error())
// 	}

// 	// 初始化OpenTelemetry
// 	otelIntegration := initOpenTelemetry(bc)
// 	defer otelIntegration.Shutdown(context.Background())

// 	// 创建日志器
// 	logger := otelIntegration.CreateLogger()

// 	// 测试不同级别的日志
// 	logger.Log(log.LevelInfo, "msg", "应用启动", "service", "nancalacc")
// 	time.Sleep(1 * time.Second)

// 	logger.Log(log.LevelWarn, "msg", "警告信息", "component", "test")
// 	time.Sleep(1 * time.Second)

// 	logger.Log(log.LevelError, "msg", "错误信息", "error_code", "E001")
// 	time.Sleep(1 * time.Second)

// 	logger.Log(log.LevelDebug, "msg", "调试信息", "debug_data", "test_data")
// 	time.Sleep(1 * time.Second)

// 	// 测试结构化日志
// 	logger.Log(log.LevelInfo,
// 		"msg", "用户操作",
// 		"user_id", "12345",
// 		"action", "login",
// 		"ip", "192.168.1.100",
// 		"timestamp", time.Now().Format(time.RFC3339),
// 	)

// 	// 等待一段时间让异步推送完成
// 	time.Sleep(5 * time.Second)

// 	println("测试完成，请检查控制台输出和 Loki 中的数据")
// }

// // initOpenTelemetry 初始化OpenTelemetry
// func initOpenTelemetry(bc *conf.Bootstrap) *otel.Integration {
// 	// 创建配置适配器
// 	adapter := otel.NewConfigAdapter()
// 	config := adapter.FromBootstrap(bc)

// 	// 创建集成器
// 	integration := otel.NewIntegration(config)

// 	// 初始化
// 	if err := integration.Init(context.Background()); err != nil {
// 		panic("failed to init OpenTelemetry: " + err.Error())
// 	}

// 	return integration
// }
