// internal/pkg/limiter/example.go
package limiter

import (
	"context"
	"fmt"
	"log"
	"time"
)

// ExampleUsage 展示限流器的使用示例
func ExampleUsage() {
	// 创建自定义配置
	config := &Config{
		CleanupInterval: 5 * time.Minute,  // 5分钟清理一次
		MaxIdleTime:     15 * time.Minute, // 15分钟空闲后清理
		MaxEntries:      1000,             // 最多1000个限流器
		DefaultRate:     5,                // 默认每秒5个请求
		DefaultBurst:    10,               // 默认突发10个请求
	}

	// 创建限流器
	limiter := NewRateLimiter(config)

	// 示例1: 使用 Allow 方法（非阻塞）
	fmt.Println("=== 示例1: 使用 Allow 方法 ===")
	for i := 0; i < 15; i++ {
		allowed := limiter.Allow("user:123", 2, 5) // 每秒2个请求，突发5个
		fmt.Printf("请求 %d: %s\n", i+1, map[bool]string{true: "允许", false: "限流"}[allowed])
		time.Sleep(100 * time.Millisecond)
	}

	// 示例2: 使用 Wait 方法（阻塞）
	fmt.Println("\n=== 示例2: 使用 Wait 方法 ===")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	for i := 0; i < 10; i++ {
		start := time.Now()
		err := limiter.Wait(ctx, "api:sync", 1, 2) // 每秒1个请求，突发2个
		duration := time.Since(start)

		if err != nil {
			fmt.Printf("请求 %d: 超时或取消\n", i+1)
		} else {
			fmt.Printf("请求 %d: 等待 %v 后执行\n", i+1, duration)
		}
	}

	// 示例3: 使用 Reserve 方法
	fmt.Println("\n=== 示例3: 使用 Reserve 方法 ===")
	for i := 0; i < 8; i++ {
		reservation := limiter.Reserve("upload:file", 0.5, 3) // 每2秒1个请求，突发3个
		delay := reservation.Delay()

		if delay > 0 {
			fmt.Printf("请求 %d: 需要等待 %v\n", i+1, delay)
			time.Sleep(delay)
		} else {
			fmt.Printf("请求 %d: 立即执行\n", i+1)
		}

		reservation.CancelAt(time.Now()) // 取消预留
	}

	// 示例4: 获取统计信息
	fmt.Println("\n=== 示例4: 统计信息 ===")
	stats := limiter.GetStats()
	for key, value := range stats {
		fmt.Printf("%s: %v\n", key, value)
	}

	// 示例5: 列出所有限流器
	fmt.Println("\n=== 示例5: 限流器列表 ===")
	keys := limiter.ListKeys()
	fmt.Printf("活跃的限流器: %v\n", keys)

	// 示例6: 获取特定限流器信息
	fmt.Println("\n=== 示例6: 限流器详情 ===")
	if info, exists := limiter.GetEntryInfo("user:123"); exists {
		for key, value := range info {
			fmt.Printf("%s: %v\n", key, value)
		}
	}

	// 示例7: 移除限流器
	fmt.Println("\n=== 示例7: 移除限流器 ===")
	removed := limiter.Remove("user:123")
	fmt.Printf("移除 user:123: %s\n", map[bool]string{true: "成功", false: "不存在"}[removed])

	// 示例8: 重置限流器
	fmt.Println("\n=== 示例8: 重置限流器 ===")
	limiter.Reset()
	fmt.Println("限流器已重置")

	// 显示最终统计
	finalStats := limiter.GetStats()
	fmt.Printf("重置后的统计: %v\n", finalStats)
}

// ExampleMiddleware 展示在HTTP中间件中的使用
func ExampleMiddleware() {
	limiter := NewRateLimiter(nil) // 使用默认配置

	// 模拟HTTP请求处理
	handleRequest := func(userID string) {
		// 使用用户ID作为限流键
		key := fmt.Sprintf("user:%s", userID)

		// 检查是否允许请求
		if limiter.Allow(key, 10, 20) { // 每秒10个请求，突发20个
			fmt.Printf("用户 %s 的请求被处理\n", userID)
			// 处理请求...
		} else {
			fmt.Printf("用户 %s 的请求被限流\n", userID)
			// 返回429状态码...
		}
	}

	// 模拟多个用户同时请求
	for i := 0; i < 3; i++ {
		userID := fmt.Sprintf("user%d", i+1)
		for j := 0; j < 15; j++ {
			handleRequest(userID)
			time.Sleep(50 * time.Millisecond)
		}
	}
}

// ExampleGRPC 展示在gRPC服务中的使用
func ExampleGRPC() {
	limiter := NewRateLimiter(&Config{
		DefaultRate:  100, // 每秒100个请求
		DefaultBurst: 200, // 突发200个请求
		MaxEntries:   5000,
	})

	// 模拟gRPC方法调用
	callGRPCMethod := func(clientIP string, method string) error {
		// 使用客户端IP和方法名作为限流键
		key := fmt.Sprintf("%s:%s", clientIP, method)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// 等待直到允许请求
		if err := limiter.Wait(ctx, key, 50, 100); err != nil {
			return fmt.Errorf("rate limit exceeded: %w", err)
		}

		// 执行gRPC方法...
		fmt.Printf("执行 gRPC 方法: %s, 客户端: %s\n", method, clientIP)
		return nil
	}

	// 模拟多个客户端调用
	clients := []string{"192.168.1.1", "192.168.1.2", "192.168.1.3"}
	methods := []string{"CreateAccount", "GetAccount", "UpdateAccount"}

	for _, client := range clients {
		for _, method := range methods {
			for i := 0; i < 5; i++ {
				if err := callGRPCMethod(client, method); err != nil {
					log.Printf("调用失败: %v", err)
				}
				time.Sleep(10 * time.Millisecond)
			}
		}
	}
}

// ExampleMonitoring 展示监控和指标收集
func ExampleMonitoring() {
	limiter := NewRateLimiter(nil)

	// 启动监控协程
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			stats := limiter.GetStats()

			// 计算限流率
			total := stats["total_requests"].(int64)
			limited := stats["limited_requests"].(int64)
			var rate float64
			if total > 0 {
				rate = float64(limited) / float64(total) * 100
			}

			fmt.Printf("监控报告 - 总请求: %d, 限流请求: %d, 限流率: %.2f%%, 活跃限流器: %d\n",
				total, limited, rate, stats["active_limiters"])
		}
	}()

	// 模拟负载
	for i := 0; i < 100; i++ {
		limiter.Allow("test:key", 5, 10)
		time.Sleep(50 * time.Millisecond)
	}

	time.Sleep(35 * time.Second) // 等待监控报告
}
