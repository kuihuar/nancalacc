package main

import (
	"log"

	"nancalacc/internal/conf"
	"nancalacc/internal/otel"

	"github.com/go-kratos/kratos/v2/log"
)

func main() {
	// 创建配置
	config := &conf.Config{
		Otel: &conf.Otel{
			Logs: &conf.Logs{
				Enabled:              true,
				Level:                "debug",
				Format:               "json",
				Output:               "stdout",
				ZapDevelopment:       true,
				ZapDisableCaller:     false,
				ZapDisableStacktrace: false,
			},
		},
	}

	// 创建logger adapter
	adapter := otel.NewKratosLoggerAdapter(config)

	// 测试日志调用
	myBusinessFunction(adapter)
}

func myBusinessFunction(logger log.Logger) {
	logger.Log(log.LevelInfo, "msg", "This should show the correct caller", "test", "caller_test")
}
