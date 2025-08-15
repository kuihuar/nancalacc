package main

import (
	"context"
	"nancalacc/internal/conf"
	"nancalacc/internal/otel"

	"github.com/go-kratos/kratos/v2/log"
)

func main() {
	// 加载配置
	bc, err := conf.Load("configs/config-zap.yaml")
	if err != nil {
		panic("failed to load config: " + err.Error())
	}

	// 创建配置适配器
	adapter := otel.NewConfigAdapter()
	config := adapter.FromBootstrap(bc)

	// 创建集成器
	integration := otel.NewIntegration(config)

	// 初始化
	if err := integration.Init(context.Background()); err != nil {
		panic("failed to init OpenTelemetry: " + err.Error())
	}
	defer integration.Shutdown(context.Background())

	// 创建日志器
	logger := integration.CreateLogger()

	// 测试不同级别的日志
	logger.Log(log.LevelDebug, "msg", "This is a debug message", "user_id", 123, "action", "login")
	logger.Log(log.LevelInfo, "msg", "This is an info message", "user_id", 123, "action", "login")
	logger.Log(log.LevelWarn, "msg", "This is a warning message", "user_id", 123, "action", "login")
	logger.Log(log.LevelError, "msg", "This is an error message", "user_id", 123, "action", "login", "error", "connection failed")

	// 测试结构化日志
	logger.Log(log.LevelInfo,
		"msg", "User action completed",
		"user_id", 123,
		"action", "file_upload",
		"file_size", 1024,
		"duration_ms", 150.5,
		"success", true,
	)

	// 测试错误日志
	logger.Log(log.LevelError,
		"msg", "Database connection failed",
		"error", "connection timeout",
		"retry_count", 3,
		"endpoint", "mysql://localhost:3306",
	)

	// 测试Fatal级别（注意：这会导致程序退出）
	// logger.Log(log.LevelFatal, "msg", "Critical error occurred", "error", "system failure")

	println("Zap integration test completed. Check the logs above and in logs/app.log")
}
