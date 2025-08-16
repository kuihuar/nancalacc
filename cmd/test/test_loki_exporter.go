package main

// import (
// 	"context"
// 	"log"
// 	"os"
// 	"time"

// 	"nancalacc/internal/otel"

// 	kratoslog "github.com/go-kratos/kratos/v2/log"
// )

// func main() {
// 	// 创建日志器
// 	logger := kratoslog.NewStdLogger(os.Stdout)

// 	// 创建Loki导出器
// 	exporter := otel.NewLokiExporter("http://192.168.1.142:3100/loki/api/v1/push", logger)

// 	ctx := context.Background()

// 	// 使用当前时间
// 	now := time.Now()
// 	log.Printf("当前时间: %s", now.Format(time.RFC3339))

// 	// 测试单条日志推送
// 	err := exporter.PushLog(ctx, "info", "这是一条测试信息日志", map[string]string{
// 		"component": "test",
// 		"user_id":   "12345",
// 	})
// 	if err != nil {
// 		log.Printf("推送单条日志失败: %v", err)
// 	} else {
// 		log.Println("✅ 单条日志推送成功")
// 	}

// 	// 等待一秒
// 	time.Sleep(time.Second)

// 	// 测试错误日志推送
// 	err = exporter.PushLog(ctx, "error", "这是一条测试错误日志", map[string]string{
// 		"component":  "test",
// 		"error_type": "test_error",
// 	})
// 	if err != nil {
// 		log.Printf("推送错误日志失败: %v", err)
// 	} else {
// 		log.Println("✅ 错误日志推送成功")
// 	}

// 	// 等待一秒
// 	time.Sleep(time.Second)

// 	// 测试批量日志推送
// 	logs := []otel.LogEntry{
// 		{
// 			Level:     "info",
// 			Message:   "批量测试日志1",
// 			Service:   "nancalacc",
// 			Timestamp: time.Now(),
// 			Labels: map[string]string{
// 				"component": "batch_test",
// 				"batch_id":  "1",
// 			},
// 		},
// 		{
// 			Level:     "warn",
// 			Message:   "批量测试日志2",
// 			Service:   "nancalacc",
// 			Timestamp: time.Now().Add(time.Second),
// 			Labels: map[string]string{
// 				"component": "batch_test",
// 				"batch_id":  "1",
// 			},
// 		},
// 		{
// 			Level:     "error",
// 			Message:   "批量测试错误日志",
// 			Service:   "nancalacc",
// 			Timestamp: time.Now().Add(2 * time.Second),
// 			Labels: map[string]string{
// 				"component":  "batch_test",
// 				"batch_id":   "1",
// 				"error_type": "batch_error",
// 			},
// 		},
// 	}

// 	err = exporter.BatchPushLog(ctx, logs)
// 	if err != nil {
// 		log.Printf("批量推送日志失败: %v", err)
// 	} else {
// 		log.Println("✅ 批量日志推送成功")
// 	}

// 	log.Printf("测试完成，请检查 Loki 中的数据，时间范围: %s 到 %s",
// 		now.Format(time.RFC3339),
// 		time.Now().Format(time.RFC3339))
// }
