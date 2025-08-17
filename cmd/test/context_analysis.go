package main

import (
	"context"
	"fmt"
	"time"
)

// 分析单个 context 的信息
func analyzeContext(ctx context.Context) {
	fmt.Println("=== Context Analysis ===")

	// 检查是否有截止时间
	if deadline, ok := ctx.Deadline(); ok {
		now := time.Now()
		remaining := deadline.Sub(now)

		fmt.Printf("✓ Has deadline: %v\n", deadline)
		fmt.Printf("✓ Remaining time: %v (%.2f seconds)\n", remaining, remaining.Seconds())

		if remaining <= 0 {
			fmt.Println("⚠️  Context has expired!")
		}
	} else {
		fmt.Println("✗ No deadline set")
	}

	// 检查是否已取消
	select {
	case <-ctx.Done():
		fmt.Printf("⚠️  Context is cancelled: %v\n", ctx.Err())
	default:
		fmt.Println("✓ Context is still active")
	}

	fmt.Println()
}

// 递归检查父级 context 链
func traceContextChain(ctx context.Context) {
	fmt.Println("=== Context Chain Trace ===")

	current := ctx
	level := 0

	for current != nil {
		fmt.Printf("Level %d:\n", level)

		// 获取 context 类型
		ctxType := fmt.Sprintf("%T", current)
		fmt.Printf("  Type: %s\n", ctxType)

		// 检查截止时间
		if deadline, ok := current.Deadline(); ok {
			now := time.Now()
			remaining := deadline.Sub(now)
			fmt.Printf("  Deadline: %v\n", deadline)
			fmt.Printf("  Remaining: %v (%.2f seconds)\n", remaining, remaining.Seconds())
		} else {
			fmt.Printf("  No deadline\n")
		}

		// 检查是否已取消
		select {
		case <-current.Done():
			fmt.Printf("  Status: CANCELLED (%v)\n", current.Err())
		default:
			fmt.Printf("  Status: ACTIVE\n")
		}

		fmt.Println()

		// 尝试获取父 context（仅对某些类型有效）
		switch v := current.(type) {
		case *context.cancelCtx:
			current = v.Context
		case *context.timerCtx:
			current = v.cancelCtx.Context
		case *context.valueCtx:
			current = v.Context
		default:
			current = nil
		}

		level++
	}
}

// 简单的超时检查函数
func checkContextTimeout(ctx context.Context) {
	deadline, ok := ctx.Deadline()
	if !ok {
		fmt.Println("Context has no deadline")
		return
	}

	now := time.Now()
	remaining := deadline.Sub(now)

	fmt.Printf("Context deadline: %v\n", deadline)
	fmt.Printf("Current time: %v\n", now)
	fmt.Printf("Remaining time: %v\n", remaining)
	fmt.Printf("Remaining seconds: %.2f\n", remaining.Seconds())
}

func main() {
	fmt.Println("Context Analysis Demo")
	fmt.Println("====================")

	// 示例1：带超时的 context
	fmt.Println("\n1. Context with timeout:")
	ctx1, cancel1 := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel1()

	analyzeContext(ctx1)

	// 示例2：嵌套的 context
	fmt.Println("\n2. Nested context:")
	ctx2, cancel2 := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel2()

	ctx3, cancel3 := context.WithTimeout(ctx2, 10*time.Second)
	defer cancel3()

	fmt.Println("Parent context:")
	analyzeContext(ctx2)

	fmt.Println("Child context:")
	analyzeContext(ctx3)

	// 示例3：已取消的 context
	fmt.Println("\n3. Cancelled context:")
	ctx4, cancel4 := context.WithTimeout(context.Background(), 5*time.Second)
	cancel4() // 立即取消

	time.Sleep(100 * time.Millisecond) // 等待取消生效
	analyzeContext(ctx4)

	// 示例4：无超时的 context
	fmt.Println("\n4. Context without timeout:")
	ctx5 := context.Background()
	analyzeContext(ctx5)

	// 示例5：完整的 context 链追踪
	fmt.Println("\n5. Context chain trace:")
	ctx6, cancel6 := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel6()

	ctx7, cancel7 := context.WithTimeout(ctx6, 15*time.Second)
	defer cancel7()

	ctx8, cancel8 := context.WithTimeout(ctx7, 5*time.Second)
	defer cancel8()

	traceContextChain(ctx8)
}
