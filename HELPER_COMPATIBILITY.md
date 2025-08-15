# Helper 兼容性指南

## 概述

`otelIntegration.CreateHelper()` 创建了一个与 Kratos `log.Helper` 完全兼容的日志助手，它将 Kratos 的 Helper 方法转换为 OpenTelemetry 格式，并支持所有 Helper 功能。

## 兼容性特性

### 1. 完全兼容 Kratos Helper 接口

支持所有 `log.Helper` 的方法：

```go
// 基本日志方法
helper.Info(args ...interface{})
helper.Debug(args ...interface{})
helper.Warn(args ...interface{})
helper.Error(args ...interface{})
helper.Fatal(args ...interface{})

// 格式化日志方法
helper.Infof(format string, args ...interface{})
helper.Debugf(format string, args ...interface{})
helper.Warnf(format string, args ...interface{})
helper.Errorf(format string, args ...interface{})
helper.Fatalf(format string, args ...interface{})

// 带键值对的日志方法
helper.Infow(keyvals ...interface{})
helper.Debugw(keyvals ...interface{})
helper.Warnw(keyvals ...interface{})
helper.Errorw(keyvals ...interface{})
helper.Fatalw(keyvals ...interface{})

// 上下文和字段方法
helper.WithContext(ctx context.Context) *Helper
helper.WithField(key string, value interface{}) *Helper
helper.WithFields(fields map[string]interface{}) *Helper
helper.WithError(err error) *Helper

// 工具方法
helper.Enabled(level Level) bool
helper.Logger() Logger
```

### 2. 自动格式转换

Helper 方法会自动将 Kratos 格式转换为 OpenTelemetry 格式：

```go
// Kratos Helper 调用
helper.Infow("用户登录", "user_id", "123", "ip", "192.168.1.1")

// 转换为 OpenTelemetry 格式
log.Record{
    Timestamp: time.Now(),
    Severity: log.SeverityInfo,
    Body: log.StringValue("用户登录"),
    Attributes: []attribute.KeyValue{
        attribute.String("user_id", "123"),
        attribute.String("ip", "192.168.1.1"),
    },
}
```

## 使用方法

### 1. 创建 Helper

```go
// 创建 Helper
helper := otelIntegration.CreateHelper()

// 或者使用现有的日志器创建
logger := otelIntegration.CreateLogger()
helper := log.NewHelper(logger)
```

### 2. 基本日志记录

```go
// 简单日志
helper.Info("应用启动成功")
helper.Debug("调试信息")
helper.Warn("警告信息")
helper.Error("发生错误")

// 格式化日志
helper.Infof("用户 %s 登录成功", "john_doe")
helper.Errorf("连接失败: %s", "网络超时")
```

### 3. 结构化日志记录

```go
// 带键值对的日志
helper.Infow("用户操作", 
    "user_id", "123", 
    "action", "login", 
    "ip", "192.168.1.1",
    "success", true,
)

// 错误日志
helper.Errorw("数据库操作失败",
    "operation", "INSERT",
    "table", "users",
    "error", "connection timeout",
    "retry_attempt", 2,
)
```

### 4. 上下文和字段

```go
ctx := context.Background()

// 带上下文的日志
helper.WithContext(ctx).Info("带上下文的日志")

// 带字段的日志
helper.WithField("request_id", "req_456").Info("带字段的日志")

// 链式调用
helper.WithContext(ctx).
    WithField("service", "user-service").
    WithField("version", "1.0.0").
    Info("链式调用的日志记录")

// 带错误信息
err := context.DeadlineExceeded
helper.WithError(err).Error("操作超时")
```

### 5. 复杂业务场景

```go
// 订单处理日志
helper.WithContext(ctx).
    WithField("order_id", "ORD_20231201_001").
    WithField("user_id", "user_789").
    WithField("total_amount", 299.99).
    WithField("currency", "CNY").
    WithField("payment_method", "alipay").
    WithField("processing_time", "2.5s").
    WithField("status", "completed").
    WithField("items_count", 3).
    Infof("订单处理完成: %s", "支付成功")

// 性能监控日志
helper.WithContext(ctx).
    WithField("endpoint", "/api/orders").
    WithField("avg_response_time", 245.67).
    WithField("max_response_time", 1200.0).
    WithField("request_count", 1500).
    WithField("error_rate", 0.02).
    Info("API 性能统计")
```

## 日志级别转换

Helper 的日志级别会自动转换为 OpenTelemetry 级别：

| Helper 方法 | Kratos 级别 | OpenTelemetry 级别 |
|------------|------------|-------------------|
| `Debug()` | `log.LevelDebug` | `log.SeverityDebug` |
| `Info()` | `log.LevelInfo` | `log.SeverityInfo` |
| `Warn()` | `log.LevelWarn` | `log.SeverityWarn` |
| `Error()` | `log.LevelError` | `log.SeverityError` |
| `Fatal()` | `log.LevelFatal` | `log.SeverityFatal` |

## 数据类型支持

Helper 支持以下数据类型的自动转换：

- `string` → `attribute.String`
- `int` → `attribute.Int`
- `int64` → `attribute.Int64`
- `float64` → `attribute.Float64`
- `bool` → `attribute.Bool`
- `error` → `attribute.String` (错误消息)
- 其他类型 → `attribute.String` (通过 `fmt.Sprintf` 转换)

## 最佳实践

### 1. 使用结构化日志

```go
// 好的做法
helper.Infow("用户登录",
    "user_id", user.ID,
    "ip", user.IP,
    "user_agent", user.UserAgent,
    "success", true,
)

// 避免的做法
helper.Infof("用户 %s 从 %s 登录，User-Agent: %s", user.ID, user.IP, user.UserAgent)
```

### 2. 合理使用上下文

```go
// 在请求处理中使用上下文
func (s *Service) HandleRequest(ctx context.Context, req *Request) error {
    helper := s.log.WithContext(ctx)
    
    helper.Info("开始处理请求")
    // ... 处理逻辑
    helper.Info("请求处理完成")
    
    return nil
}
```

### 3. 错误处理

```go
// 记录错误信息
if err != nil {
    helper.WithError(err).
        WithField("operation", "database_query").
        WithField("table", "users").
        Error("数据库查询失败")
    return err
}
```

### 4. 性能监控

```go
start := time.Now()
// ... 执行操作
duration := time.Since(start)

helper.WithField("operation", "user_login").
    WithField("duration", duration.String()).
    WithField("duration_ms", duration.Milliseconds()).
    Info("操作完成")
```

## 测试

运行测试文件查看 Helper 兼容性：

```bash
go run test_helper_compatibility.go
```

这将演示所有 Helper 方法的使用，并展示日志如何被转换为 OpenTelemetry 格式。

## 优势

1. **完全兼容**: 支持所有 Kratos Helper 方法，无需修改现有代码
2. **自动转换**: 自动将 Kratos 格式转换为 OpenTelemetry 格式
3. **结构化日志**: 支持结构化日志记录，便于查询和分析
4. **上下文支持**: 支持上下文传递和字段添加
5. **链式调用**: 支持链式调用，代码更简洁
6. **类型安全**: 支持多种数据类型的自动转换

## 迁移指南

从现有的 Kratos Helper 迁移到 OpenTelemetry Helper：

```go
// 原有代码
helper := log.NewHelper(logger)
helper.Info("用户登录")

// 迁移后代码
helper := otelIntegration.CreateHelper()
helper.Info("用户登录") // 无需修改，完全兼容
```

这种设计使得应用可以无缝地从 Kratos 的 Helper 迁移到 OpenTelemetry 的可观测性系统，同时保持代码的一致性和可维护性。 