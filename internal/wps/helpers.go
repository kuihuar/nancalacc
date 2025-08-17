package wps

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

// logAPIRequest 统一的API请求日志记录
func logAPIRequest(ctx context.Context, logger log.Logger, operation, method, path string, input interface{}) {
	logger.Log(log.LevelInfo, "msg", "API request started", "method", method, "path", path)

	if input != nil {
		logger.Log(log.LevelInfo, "msg", "Request input", "input", input)
	}
}

// CodeChecker 定义检查业务码的接口
type CodeChecker[T comparable] interface {
	GetCode() T
}

// handleAPIResponse 统一的API响应处理函数，使用泛型支持int和string类型的Code字段
func handleAPIResponse[T comparable](ctx context.Context, logger log.Logger, operation string, responseBody []byte, response CodeChecker[T], expectedCode T) error {
	start := time.Now()

	// 记录响应
	logger.Log(log.LevelInfo, "msg", "API response received", "size", len(responseBody), "duration_ms", time.Since(start).Milliseconds())

	// 解析响应
	if err := json.Unmarshal(responseBody, response); err != nil {
		logger.Log(log.LevelError, "msg", "Failed to unmarshal response", "operation", operation, "err", err)
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// 检查业务错误码
	if code := response.GetCode(); code != expectedCode {
		logger.Log(log.LevelError, "msg", "Business error", "operation", operation, "code", code, "expected", expectedCode)
		return fmt.Errorf("business error: code=%v", code)
	}

	logger.Log(log.LevelInfo, "msg", "API call completed", "operation", operation)
	return nil
}
