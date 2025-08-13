package main

import (
	"context"
	"nancalacc/internal/log"
)

func main() {
	// 创建日志配置
	config := log.NewLogConfig()

	// 创建日志记录器
	logger, err := log.NewLogger(config)
	if err != nil {
		panic(err)
	}

	// 创建日志助手
	helper := log.NewLogHelper(logger)

	// 测试结构化日志
	helper.Info("测试结构化日志",
		log.NewField("user_id", "123"),
		log.NewField("action", "login"),
		log.NewField("status", "success"),
	)

	// 测试错误日志
	helper.Error("测试错误日志",
		log.NewField("error", "connection failed"),
		log.NewField("retry_count", 3),
	)

	// 测试 kratos 风格的日志
	helper.Log(1, "msg", "kratos 风格日志", "user_id", "456", "action", "logout")

	// 测试 With 方法
	ctx := context.Background()
	helperWithCtx := helper.With("request_id", "req-123", "user_agent", "test-agent")
	helperWithCtx.Log(1, "msg", "带上下文的日志", "operation", "test")

	// 测试 WithContext
	helperWithContext := helper.WithContext(ctx)
	helperWithContext.Info("带上下文的日志2", log.NewField("operation", "test2"))
}
