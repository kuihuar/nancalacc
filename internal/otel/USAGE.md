# OpenTelemetry 集成使用指南

## 概述

这个 OpenTelemetry 集成提供了完整的可观测性解决方案，包括分布式追踪、指标收集和结构化日志记录。它与 Kratos 框架完全兼容，可以无缝集成到现有项目中。

## 核心组件

### 1. Integration (集成器)
主要的集成接口，提供：
- 初始化 OpenTelemetry 服务
- 创建 Kratos 兼容的日志器
- 生成 HTTP/gRPC 中间件
- 管理服务生命周期

### 2. Service (服务)
OpenTelemetry 核心服务，管理：
- Tracer (追踪器)
- Meter (指标器)
- Logger (日志器)

### 3. Config (配置)
配置管理，支持：
- 追踪配置
- 指标配置
- 日志配置
- 资源属性

## 基本使用

### 1. 创建配置

```go
config := &conf.OpenTelemetry{
    Enabled: true,
    Traces: &conf.Traces{
        Enabled:  true,
        Endpoint: "http://localhost:14268/api/traces",
    },
    Metrics: &conf.Metrics{
        Enabled:  true,
        Endpoint: "http://localhost:14250",
    },
    Logs: &conf.Logs{
        Enabled:  true,
        Endpoint: "http://localhost:14250",
    },
}
```

### 2. 初始化集成器

```go
// 创建集成器
otelIntegration := otel.NewIntegration(config)

// 初始化
ctx := context.Background()
if err := otelIntegration.Init(ctx); err != nil {
    log.Fatalf("Failed to initialize OpenTelemetry: %v", err)
}
defer otelIntegration.Shutdown(ctx)
```

### 3. 在 Kratos 应用中使用

```go
// 创建日志器
logger := otelIntegration.CreateLogger()

// 创建服务器并添加中间件
httpSrv := http.NewServer(
    http.Address(":8000"),
    otelIntegration.CreateHTTPMiddleware()...,
)

grpcSrv := grpc.NewServer(
    grpc.Address(":9000"),
    otelIntegration.CreateGRPCMiddleware()...,
)

// 创建应用
app := kratos.New(
    kratos.Logger(logger),
    kratos.Server(httpSrv, grpcSrv),
)
```

## 在业务代码中使用

### 1. 获取追踪器和日志器

```go
// 获取追踪器
tracer := otelIntegration.GetService().GetTracer()

// 获取日志器
logger := otelIntegration.CreateLogger()
```

### 2. 创建业务服务

```go
type UserService struct {
    logger log.Logger
    tracer trace.Tracer
}

func NewUserService(logger log.Logger, tracer trace.Tracer) *UserService {
    return &UserService{
        logger: logger,
        tracer: tracer,
    }
}
```

### 3. 在方法中使用追踪和日志

```go
func (s *UserService) GetUser(ctx context.Context, userID string) (*User, error) {
    // 创建 span
    ctx, span := s.tracer.Start(ctx, "UserService.GetUser")
    defer span.End()

    // 设置 span 属性
    span.SetAttributes(
        attribute.String("user.id", userID),
        attribute.String("service.name", "UserService"),
    )

    // 记录日志
    s.logger.Log(log.LevelInfo, "msg", "Getting user", "user_id", userID)

    // 业务逻辑
    user, err := s.fetchUserFromDB(ctx, userID)
    if err != nil {
        // 记录错误
        span.RecordError(err)
        s.logger.Log(log.LevelError, "msg", "Failed to get user", "error", err.Error())
        return nil, err
    }

    // 记录成功日志
    s.logger.Log(log.LevelInfo, "msg", "User retrieved successfully", 
        "user_id", userID, 
        "user_name", user.Name)

    return user, nil
}
```

## 日志记录

### 基本日志

```go
logger.Log(log.LevelInfo, "msg", "Operation completed")
logger.Log(log.LevelError, "msg", "Operation failed", "error", err.Error())
```

### 结构化日志

```go
logger.Log(log.LevelInfo, 
    "msg", "User operation", 
    "user_id", "123", 
    "operation", "create",
    "duration_ms", 150,
    "status", "success",
)
```

### 日志级别

- `log.LevelDebug`: 调试信息
- `log.LevelInfo`: 一般信息
- `log.LevelWarn`: 警告信息
- `log.LevelError`: 错误信息
- `log.LevelFatal`: 致命错误

## 分布式追踪

### 创建 Span

```go
ctx, span := tracer.Start(ctx, "operation.name")
defer span.End()
```

### 设置属性

```go
span.SetAttributes(
    attribute.String("user.id", userID),
    attribute.Int("request.size", len(data)),
    attribute.Bool("cache.hit", true),
    attribute.Float64("duration.ms", 123.45),
)
```

### 记录事件

```go
span.AddEvent("user.authenticated", trace.WithAttributes(
    attribute.String("user.id", userID),
    attribute.String("auth.method", "jwt"),
))
```

### 记录错误

```go
if err != nil {
    span.RecordError(err)
    return nil, err
}
```

## 指标收集

### 创建指标

```go
meter := otelIntegration.GetService().GetMeter()

// 计数器
requestCounter, _ := meter.Int64Counter("http.requests.total")
requestCounter.Add(ctx, 1, metric.WithAttributes(
    attribute.String("method", "GET"),
    attribute.String("path", "/users"),
))

// 直方图
requestDuration, _ := meter.Float64Histogram("http.request.duration")
requestDuration.Record(ctx, 123.45, metric.WithAttributes(
    attribute.String("method", "GET"),
    attribute.String("path", "/users"),
))
```

## 配置选项

### 追踪配置

```go
Traces: &conf.Traces{
    Enabled:  true,
    Endpoint: "http://localhost:14268/api/traces",
    Sampler: &conf.Sampler{
        Type:  "always_on",
        Param: 1.0,
    },
}
```

### 指标配置

```go
Metrics: &conf.Metrics{
    Enabled:  true,
    Endpoint: "http://localhost:14250",
    Interval: "10s",
}
```

### 日志配置

```go
Logs: &conf.Logs{
    Enabled:  true,
    Endpoint: "http://localhost:14250",
    Level:    "info",
}
```

## 生产环境建议

### 1. 采样配置

```go
Sampler: &conf.Sampler{
    Type:  "probabilistic",
    Param: 0.1, // 10% 采样率
}
```

### 2. 批量导出

```go
Traces: &conf.Traces{
    Enabled:  true,
    Endpoint: "http://otel-collector:14268/api/traces",
    Batch: &conf.Batch{
        MaxQueueSize: 1000,
        MaxExportBatchSize: 512,
        ExportTimeout: "30s",
    },
}
```

### 3. 资源属性

```go
Resource: &conf.Resource{
    Attributes: []*conf.KeyValue{
        {Key: "service.name", Value: "my-service"},
        {Key: "service.version", Value: "1.0.0"},
        {Key: "deployment.environment", Value: "production"},
    },
}
```

## 故障排除

### 常见问题

1. **连接失败**
   - 检查端点配置
   - 验证网络连接
   - 确认服务可用性

2. **数据丢失**
   - 检查缓冲区配置
   - 验证网络稳定性
   - 调整批量大小

3. **性能问题**
   - 调整采样率
   - 优化批量配置
   - 监控资源使用

### 调试模式

```go
config := &conf.OpenTelemetry{
    Enabled: true,
    Debug:   true,
    // ... 其他配置
}
```

## 示例项目

完整的示例项目位于 `cmd/otel_example/` 目录，包含：

- 完整的应用示例
- 业务服务示例
- 配置文件
- 测试脚本
- 详细文档

运行示例：

```bash
cd cmd/otel_example
go run .
``` 