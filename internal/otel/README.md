# OpenTelemetry 集成系统

本目录包含了完整的 OpenTelemetry 集成系统，为应用程序提供统一的观测性解决方案。

## 文件结构

### 1. `config.go` - 配置管理
- **作用**: 定义 OpenTelemetry 配置结构，管理各种后端配置
- **功能**:
  - 定义追踪、指标、日志的配置结构
  - 提供默认配置和获取各种 OpenTelemetry 组件的方法
  - 管理 Jaeger、Prometheus、Loki 等后端配置
  - 支持环境变量配置

### 2. `service.go` - 服务封装
- **作用**: 封装 OpenTelemetry 核心服务，提供统一接口
- **功能**:
  - 管理 Tracer、Meter、Logger 实例
  - 提供便捷的日志记录方法（Info、Error、Warn、Debug）
  - 统一的错误处理和资源管理
  - 支持优雅关闭

### 3. `middleware.go` - HTTP/gRPC 中间件
- **作用**: 提供 HTTP 和 gRPC 的中间件，自动添加观测性
- **功能**:
  - 自动记录请求信息、响应时间、错误等
  - 包装响应写入器以获取状态码
  - 支持自定义属性添加
  - 与 Kratos 框架集成

### 4. `helpers.go` - 辅助函数
- **作用**: 提供各种辅助函数，简化 OpenTelemetry 的使用
- **功能**:
  - 追踪辅助函数（添加属性、记录事件）
  - 指标辅助函数（记录计数器、直方图）
  - 日志辅助函数（结构化日志记录）
  - 计时操作辅助函数

### 5. `integration.go` - 集成器
- **作用**: 与 Kratos 框架集成，创建中间件选项
- **功能**:
  - 初始化服务、创建 HTTP/gRPC 中间件
  - 管理服务的生命周期
  - 提供依赖注入支持

## 主要特性

### 统一的观测性
- **分布式追踪**: 自动追踪请求链路
- **指标收集**: 性能指标和业务指标
- **结构化日志**: 统一的日志格式和级别

### 框架集成
- **Kratos 集成**: 无缝集成 Kratos 框架
- **中间件支持**: HTTP 和 gRPC 中间件
- **依赖注入**: 支持 Wire 依赖注入

### 配置灵活
- **多后端支持**: Jaeger、Prometheus、Loki
- **环境变量**: 支持环境变量配置
- **默认配置**: 提供合理的默认配置

## 使用示例

### 基本使用
```go
// 创建服务
svc := otel.NewService(cfg)

// 记录日志
svc.Info(ctx, "用户登录", log.String("user_id", "123"))

// 记录指标
svc.RecordMetric(ctx, "login_attempts", 1)

// 添加追踪
svc.AddSpanAttribute(ctx, "user_id", "123")
```

### 中间件使用
```go
// HTTP 中间件
httpMiddleware := otel.HTTPMiddleware(svc)

// gRPC 中间件
grpcMiddleware := otel.GRPCMiddleware(svc)
```

### 辅助函数使用
```go
// 计时操作
err := otel.TimeOperation(ctx, "database_query", func(ctx context.Context) error {
    return db.Query(ctx, "SELECT * FROM users")
})

// 添加追踪属性
otel.AddSpanAttribute(ctx, "operation", "user_query")
otel.AddSpanAttribute(ctx, "table", "users")
```

## 修复的问题

1. **导入问题**: 添加了正确的 noop 包导入
2. **API 调用**: 使用正确的 OpenTelemetry API
3. **日志记录**: 使用正确的日志记录方法
4. **追踪器创建**: 使用正确的追踪器提供者

## 架构优势

- **模块化设计**: 各组件职责清晰，易于维护
- **扩展性强**: 支持添加新的观测性后端
- **性能优化**: 使用 noop 实现避免性能影响
- **开发友好**: 提供丰富的辅助函数和示例

## 依赖关系

- `go.opentelemetry.io/otel`: OpenTelemetry 核心库
- `go.opentelemetry.io/otel/trace`: 追踪功能
- `go.opentelemetry.io/otel/metric`: 指标功能
- `go.opentelemetry.io/otel/log`: 日志功能
- `github.com/go-kratos/kratos/v2`: Kratos 框架

这个系统为应用程序提供了完整的观测性解决方案，帮助开发者更好地监控和调试应用程序。 