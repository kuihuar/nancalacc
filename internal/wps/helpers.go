package wps

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

// logAPIRequest 统一的API请求日志记录
func logAPIRequest(ctx context.Context, logger *log.Helper, operation, method, path string, input interface{}) {
	logCtx := logger.WithContext(ctx)

	logCtx.Infof("API request started: %s %s", method, path)

	if input != nil {
		logCtx.Debugf("Request input: %+v", input)
	}
}

// CodeChecker 定义检查业务码的接口
type CodeChecker[T comparable] interface {
	GetCode() T
}

// handleAPIResponse 统一的API响应处理函数，使用泛型支持int和string类型的Code字段
func handleAPIResponse[T comparable](ctx context.Context, logger *log.Helper, operation string, responseBody []byte, response CodeChecker[T], expectedCode T) error {
	start := time.Now()
	logCtx := logger.WithContext(ctx)

	// 记录响应
	logCtx.Infof("API response received: size=%d, duration_ms=%d", len(responseBody), time.Since(start).Milliseconds())

	// 解析响应
	if err := json.Unmarshal(responseBody, response); err != nil {
		logCtx.Errorf("Failed to unmarshal response for %s: %v", operation, err)
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// 检查业务错误码
	if code := response.GetCode(); code != expectedCode {
		logCtx.Errorf("Business error for %s: code=%v, expected=%v", operation, code, expectedCode)
		return fmt.Errorf("business error: code=%v", code)
	}

	logCtx.Infof("API call %s completed successfully", operation)
	return nil
}
