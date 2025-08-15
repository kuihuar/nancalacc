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

	// 测试不同位置的日志调用
	testFunction1(adapter)
	testFunction2(adapter)
	testFunction3(adapter)
}

func testFunction1(logger log.Logger) {
	logger.Log(log.LevelInfo, "msg", "This is from testFunction1", "function", "testFunction1")
}

func testFunction2(logger log.Logger) {
	logger.Log(log.LevelWarn, "msg", "This is from testFunction2", "function", "testFunction2")
}

func testFunction3(logger log.Logger) {
	logger.Log(log.LevelError, "msg", "This is from testFunction3", "function", "testFunction3")
}
