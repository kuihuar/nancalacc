# Zap Logger 集成指南

本项目已经集成了 Uber 的 Zap 日志库，提供了高性能的结构化日志记录功能。

## 特性

- **高性能**: Zap 是 Go 生态中最快的日志库之一
- **结构化日志**: 支持 JSON 和 Console 格式
- **灵活配置**: 支持多种输出方式和配置选项
- **向后兼容**: 与现有的 OpenTelemetry 和 Kratos 日志系统完全兼容
- **异步支持**: 可以通过配置实现异步日志记录

## 配置

### 基本配置

在配置文件中添加以下 zap 相关配置：

```yaml
otel:
  logs:
    enabled: true
    level: "info"
    format: "json"
    output: "both"  # stdout, file, both
    file_path: "logs/app.log"
    
    # Zap配置
    use_zap: true                    # 启用zap logger
    zap_development: false           # 生产模式
    zap_disable_caller: false        # 启用调用者信息
    zap_disable_stacktrace: false    # 启用堆栈跟踪
    zap_encoding: "json"             # 编码格式: json, console
```

### 高级配置

```yaml
otel:
  logs:
    # 基本配置
    use_zap: true
    zap_development: false
    
    # 字段名自定义
    zap_time_key: "timestamp"
    zap_level_key: "level"
    zap_name_key: "logger"
    zap_caller_key: "caller"
    zap_function_key: "func"
    zap_message_key: "message"
    zap_stacktrace_key: "stacktrace"
    
    # 输出配置
    output: "both"
    file_path: "logs/app.log"
    max_size: 100    # MB
    max_backups: 3
    max_age: 28      # days
    compress: true
```

## 使用方法

### 1. 基本日志记录

```go
logger := integration.CreateLogger()

// 信息日志
logger.Log(log.LevelInfo, "msg", "User logged in", "user_id", 123)

// 错误日志
logger.Log(log.LevelError, "msg", "Database connection failed", "error", "timeout")

// 调试日志
logger.Log(log.LevelDebug, "msg", "Processing request", "request_id", "abc123")
```

### 2. 结构化日志

```go
logger.Log(log.LevelInfo, 
    "msg", "API request completed",
    "method", "POST",
    "path", "/api/users",
    "status_code", 200,
    "duration_ms", 150.5,
    "user_id", 123,
    "ip", "192.168.1.1",
)
```

### 3. 错误日志

```go
err := someFunction()
if err != nil {
    logger.Log(log.LevelError,
        "msg", "Operation failed",
        "error", err.Error(),
        "operation", "user_creation",
        "user_id", 123,
    )
}
```

## 输出格式

### JSON 格式 (默认)

```json
{
  "level": "INFO",
  "timestamp": "2025-08-15T23:21:29.790+0800",
  "caller": "otel/logger_adapter.go:478",
  "func": "nancalacc/internal/otel.(*KratosLoggerAdapter).logWithZap",
  "message": "User action completed",
  "user_id": 123,
  "action": "file_upload",
  "file_size": 1024,
  "duration_ms": 150.5,
  "success": true
}
```

### Console 格式

```
2025-08-15T23:21:29.790+0800    INFO    User action completed    {"user_id": 123, "action": "file_upload", "file_size": 1024, "duration_ms": 150.5, "success": true}
```

## 性能优势

1. **零分配**: Zap 在大多数情况下实现零内存分配
2. **结构化**: 支持结构化字段，便于日志分析
3. **异步**: 可以通过配置实现异步日志记录
4. **级别过滤**: 高效的日志级别过滤

## 与现有系统的兼容性

- ✅ 完全兼容 Kratos 日志接口
- ✅ 兼容 OpenTelemetry 日志系统
- ✅ 支持现有的日志配置
- ✅ 无需修改现有代码

## 配置示例

### 开发环境

```yaml
otel:
  logs:
    use_zap: true
    zap_development: true
    zap_encoding: "console"
    output: "stdout"
    level: "debug"
```

### 生产环境

```yaml
otel:
  logs:
    use_zap: true
    zap_development: false
    zap_encoding: "json"
    output: "both"
    file_path: "logs/app.log"
    level: "info"
    max_size: 100
    max_backups: 5
    compress: true
```

### 高性能配置

```yaml
otel:
  logs:
    use_zap: true
    zap_development: false
    zap_disable_caller: true
    zap_disable_stacktrace: true
    zap_encoding: "json"
    output: "file"
    file_path: "logs/app.log"
    level: "warn"
```

## 测试

运行测试文件验证集成：

```bash
go run test_zap_integration.go
```

这将输出结构化日志到控制台和文件。

## 注意事项

1. **文件权限**: 确保应用有权限写入日志文件目录
2. **磁盘空间**: 定期清理日志文件，避免磁盘空间不足
3. **性能**: 在生产环境中，建议禁用调用者信息和堆栈跟踪以提高性能
4. **日志轮转**: 使用 `lumberjack` 进行日志轮转，避免单个文件过大 