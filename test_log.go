package main

import (
	"context"
	"nancalacc/internal/log"
)

func main() {
	// 1. 创建日志配置
	config := &log.Config{
		Level:      "info",
		Format:     "json",
		Output:     "stdout",
		Stacktrace: false,
	}

	// 2. 创建自定义日志记录器
	logger, err := log.NewLogger(config)
	if err != nil {
		panic(err)
	}

	// 3. 添加基础字段（模拟业务代码中的loggerWithFields）
	loggerWithFields := logger.WithFields(map[string]interface{}{
		"service.name": "test-service",
	})

	// 4. 创建 Kratos 日志适配器
	kratosLogger := log.NewKratosLoggerAdapter(loggerWithFields)

	// 5. 创建 Helper
	helper := log.NewHelper(kratosLogger)

	// 6. 测试各种调用方式
	ctx := context.Background()

	// 测试直接调用
	helper.Infof("Test direct Infof call")

	// 测试 WithContext 调用
	helper.WithContext(ctx).Infof("Test WithContext Infof call")

	// 测试直接 Info 调用
	helper.Info("Test direct Info call")

	// 测试 WithContext Info 调用
	helper.WithContext(ctx).Info("Test WithContext Info call")

	// 测试链式调用
	helper.WithContext(ctx).WithField("test_field", "value").Infof("Test chained call")

	// 测试 WithField 功能
	helper.WithField("user_id", "123").Info("Test WithField single field")

	// 测试多个字段
	helper.WithField("task_id", "task_456").WithField("status", "running").Info("Test multiple fields")

	// 测试链式调用多个字段
	helper.WithContext(ctx).WithField("request_id", "req_789").WithField("method", "POST").Infof("Test chained multiple fields")
}
