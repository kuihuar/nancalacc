# 日志系统与Loki集成指南

## 概述

本项目集成了完整的日志系统，支持多种输出方式，包括控制台、文件和Loki。Loki是一个水平可扩展、高可用的日志聚合系统，特别适合微服务架构。

## 功能特性

- **多输出支持**: 同时支持控制台、文件和Loki输出
- **结构化日志**: 支持JSON和Console格式
- **日志轮转**: 自动日志文件轮转和压缩
- **批量发送**: Loki支持批量发送以提高性能
- **标签支持**: 支持自定义标签用于日志分类
- **GORM集成**: 自动集成GORM数据库日志
- **中间件支持**: HTTP请求日志中间件

## 配置说明

### 基础配置

```yaml
logging:
  level: "info"           # 日志级别: debug, info, warn, error, fatal
  format: "json"          # 日志格式: json, console
  output: "both"          # 输出方式: stdout, file, both
  file_path: "logs/app.log"  # 日志文件路径
  max_size: 100           # 单个文件最大大小(MB)
  max_backups: 10         # 最大备份文件数
  max_age: 30             # 最大保留天数
  compress: true          # 是否压缩
  caller: true            # 是否显示调用者信息
  stacktrace: false       # 是否显示堆栈信息
```

### GORM日志配置

```yaml
logging:
  gorm:
    slow_threshold: "200ms"  # 慢查询阈值
    log_level: "warn"        # GORM日志级别
```

### Loki配置

```yaml
logging:
  loki:
    url: "http://loki:3100"  # Loki服务地址
    username: ""             # 用户名（可选）
    password: ""             # 密码（可选）
    tenant_id: "nancalacc"   # 租户ID
    enable: true             # 启用Loki
    batch_size: 100         # 批处理大小
    batch_wait: "1s"        # 批处理等待时间
    timeout: "10s"          # 请求超时时间
    labels:                 # 标签
      service: "nancalacc"
      environment: "production"
      version: "v1.0.0"
```

## 使用方法

### 1. 基础日志记录

```go
package main

import (
    "context"
    "nancalacc/internal/otel"
    "github.com/go-kratos/kratos/v2/log"
)

func main() {
    // 创建 OpenTelemetry 集成器
    config := &otel.Config{
        Enabled: true,
        Logs: otel.LogConfig{
            Enabled: true,
            Level:   "info",
        },
    }
    
    integration := otel.NewIntegration(config)
    if err := integration.Init(context.Background()); err != nil {
        panic(err)
    }
    defer integration.Shutdown(context.Background())

    // 创建日志记录器
    logger := integration.CreateLogger()

    // 记录不同级别的日志
    logger.Log(log.LevelInfo, "msg", "应用启动成功")
    logger.Log(log.LevelWarn, "msg", "配置项缺失，使用默认值")
    logger.Log(log.LevelError, "msg", "数据库连接失败", "error", err.Error())
}
```

### 2. 使用日志助手

```go
package service

import (
    "context"
    "nancalacc/internal/otel"
    "github.com/go-kratos/kratos/v2/log"
)

type UserService struct {
    logger log.Logger
}

func NewUserService(integration *otel.Integration) *UserService {
    return &UserService{
        logger: integration.CreateLogger(),
    }
}

func (s *UserService) CreateUser(user *User) error {
    s.logger.Log(log.LevelInfo, 
        "msg", "创建用户",
        "user_id", user.ID,
        "email", user.Email,
    )
    
    // 业务逻辑...
    
    s.logger.Log(log.LevelInfo, "msg", "用户创建成功", "user_id", user.ID)
    return nil
}
```

### 3. 添加上下文信息

```go
func (s *UserService) GetUser(ctx context.Context, userID string) (*User, error) {
    s.logger.Log(log.LevelInfo, 
        "msg", "获取用户信息", 
        "user_id", userID,
        "request_id", ctx.Value("request_id"),
    )
    
    // 业务逻辑...
    
    return user, nil
}
```

### 4. HTTP中间件

```go
package middleware

import (
    "context"
    "time"
    "nancalacc/internal/otel"
    "github.com/go-kratos/kratos/v2/log"
    "github.com/go-kratos/kratos/v2/middleware"
)

func LogMiddleware(integration *otel.Integration) middleware.Middleware {
    logger := integration.CreateLogger()
    
    return func(handler middleware.Handler) middleware.Handler {
        return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
            start := time.Now()
            
            reply, err = handler(ctx, req)
            
            duration := time.Since(start)
            
            logger.Log(log.LevelInfo, 
                "msg", "HTTP请求",
                "method", "POST",
                "path", "/api/users",
                "duration", duration,
                "status", "200",
            )
            
            return reply, err
        }
    }
}
```

## 部署Loki

### 使用Docker Compose

1. 启动Loki和Grafana:

```bash
docker-compose -f docker-compose.loki.yml up -d
```

2. 访问Grafana:
   - URL: http://localhost:3000
   - 用户名: admin
   - 密码: admin

3. 配置数据源:
   - Loki数据源已自动配置
   - URL: http://loki:3100

### 查询日志

在Grafana中可以使用LogQL查询语言查询日志:

```logql
# 查询所有日志
{service="nancalacc"}

# 查询错误日志
{service="nancalacc"} |= "error"

# 查询特定用户的日志
{service="nancalacc"} |~ "user_id=123"

# 查询慢查询
{service="nancalacc"} |~ "slow query"
```

## 性能优化

### 1. 批量发送

Loki客户端支持批量发送以提高性能:

```yaml
logging:
  loki:
    batch_size: 100    # 批处理大小
    batch_wait: "1s"   # 批处理等待时间
```

### 2. 标签优化

合理使用标签可以提高查询性能:

```yaml
logging:
  loki:
    labels:
      service: "nancalacc"
      environment: "production"
      version: "v1.0.0"
      component: "api"
```

### 3. 日志级别

在生产环境中使用适当的日志级别:

```yaml
logging:
  level: "info"  # 生产环境推荐使用info级别
```

## 监控和告警

### 1. 日志量监控

监控日志量变化，及时发现异常:

```logql
# 统计每分钟的日志量
sum(rate({service="nancalacc"}[1m])) by (level)
```

### 2. 错误率监控

监控错误日志比例:

```logql
# 计算错误率
sum(rate({service="nancalacc", level="error"}[5m])) / sum(rate({service="nancalacc"}[5m]))
```

### 3. 慢查询监控

监控数据库慢查询:

```logql
# 统计慢查询数量
sum(rate({service="nancalacc"} |~ "slow query" [5m]))
```

## 故障排查

### 1. 日志不显示

- 检查日志级别配置
- 确认输出方式配置正确
- 验证文件路径权限

### 2. Loki连接失败

- 检查Loki服务状态
- 验证网络连接
- 确认认证信息

### 3. 性能问题

- 调整批处理参数
- 优化标签使用
- 检查网络延迟

## 最佳实践

1. **结构化日志**: 使用JSON格式和结构化字段
2. **合理标签**: 避免过多标签，影响查询性能
3. **日志级别**: 合理使用日志级别
4. **错误处理**: 记录详细的错误信息
5. **性能监控**: 定期监控日志系统性能
6. **备份策略**: 制定日志备份和清理策略 