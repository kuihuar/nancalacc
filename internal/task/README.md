# 定时任务模块 (internal/task)

## 📋 概述

这是一个基于 `github.com/robfig/cron/v3` 的优化定时任务模块，提供了完整的任务管理、监控、错误处理和指标收集功能。

## ✨ 主要特性

### 🔧 核心功能
- **任务调度**: 支持标准的 cron 表达式，秒级精度
- **错误处理**: 完整的 panic 恢复和错误处理机制
- **超时控制**: 支持任务执行超时控制
- **上下文支持**: 支持 context 取消和超时
- **任务监控**: 实时监控任务执行状态和性能指标

### 📊 监控指标
- 任务执行次数统计
- 成功/失败率统计
- 平均执行时间
- 当前运行任务数
- 错误计数和最后错误信息

### 🛡️ 可靠性
- 并发安全的任务管理
- 优雅的启动和停止
- 详细的日志记录
- 任务状态持久化（内存）

## 🏗️ 架构设计

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   CronService   │    │ MetricsCollector│    │   JobConfig     │
│                 │    │                 │    │                 │
│ - 任务调度      │    │ - 指标收集      │    │ - 任务配置      │
│ - 错误处理      │    │ - 性能统计      │    │ - 超时设置      │
│ - 状态管理      │    │ - 数据导出      │    │ - 重试策略      │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         └───────────────────────┼───────────────────────┘
                                 │
                    ┌─────────────────┐
                    │   TaskInfo      │
                    │                 │
                    │ - 任务元数据    │
                    │ - 执行统计      │
                    │ - 状态信息      │
                    └─────────────────┘
```

## 🚀 使用方法

### 1. 基本使用

```go
// 创建定时任务服务
cronService := NewCronService(accounterUsecase, logger)

// 启动服务
cronService.Start()
defer cronService.Stop()

// 添加简单任务
entryID, err := cronService.AddFunc("simple_task", "0 */5 * * * *", func() {
    log.Info("执行简单任务")
})
```

### 2. 带上下文的任务

```go
// 添加带上下文的任务
entryID, err := cronService.AddFuncWithContext("context_task", "0 */10 * * * *", func(ctx context.Context) error {
    // 检查上下文是否被取消
    select {
    case <-ctx.Done():
        return ctx.Err()
    default:
        // 执行任务逻辑
        return nil
    }
})
```

### 3. 任务配置

```go
// 定义任务配置
jobConfig := JobConfig{
    Name:        "sync_account",
    Spec:        "0 0 2 * * *", // 每天凌晨2点执行
    Description: "每日账户同步",
    Timeout:     30 * time.Minute,
    RetryCount:  3,
    RetryDelay:  5 * time.Minute,
}

// 注册任务
registerJob(cronService, jobConfig)
```

### 4. 监控指标

```go
// 获取任务指标
taskInfo, exists := cronService.GetTaskInfo(entryID)
if exists {
    log.Infof("任务执行次数: %d", taskInfo.RunCount)
    log.Infof("任务失败次数: %d", taskInfo.ErrorCount)
    log.Infof("平均执行时间: %v", taskInfo.AvgDuration)
}

// 获取所有任务指标
allTasks := cronService.GetAllTasks()
for name, info := range allTasks {
    log.Infof("任务 [%s]: 执行 %d 次, 失败 %d 次", name, info.RunCount, info.ErrorCount)
}
```

## 📝 Cron 表达式格式

支持标准的 cron 表达式格式：

```
秒 分 时 日 月 星期
* * * * * *
│ │ │ │ │ │
│ │ │ │ │ └── 星期 (0-7, 0和7都表示星期日)
│ │ │ │ └──── 月 (1-12)
│ │ │ └────── 日 (1-31)
│ │ └──────── 时 (0-23)
│ └────────── 分 (0-59)
└──────────── 秒 (0-59)
```

### 常用表达式示例

```go
"0 0 * * * *"     // 每小时执行
"0 */30 * * * *"  // 每30分钟执行
"0 0 2 * * *"     // 每天凌晨2点执行
"0 0 9 * * 1"     // 每周一上午9点执行
"0 0 12 1 * *"    // 每月1号中午12点执行
```

## 🔍 监控和调试

### 1. 日志级别

模块使用结构化日志，包含以下字段：
- `module`: 模块名称
- `task`: 任务名称
- `duration`: 执行时间
- `error`: 错误信息

### 2. 指标导出

```go
// 导出所有指标
metrics := metricsCollector.ExportMetrics(ctx)
log.Infof("指标数据: %+v", metrics)
```

### 3. 健康检查

```go
// 检查任务服务状态
tasks := cronService.GetAllTasks()
for name, info := range tasks {
    if info.ErrorCount > 10 {
        log.Warnf("任务 [%s] 失败次数过多: %d", name, info.ErrorCount)
    }
}
```

## ⚡ 性能优化

### 1. 内存管理
- 使用 sync.RWMutex 保证并发安全
- 定期清理过期的任务信息
- 避免内存泄漏

### 2. 执行优化
- 任务超时控制防止长时间阻塞
- 并发执行支持
- 错误重试机制

### 3. 监控优化
- 异步指标收集
- 批量指标导出
- 可配置的监控间隔

## 🔧 配置选项

### 1. 任务配置

```go
type JobConfig struct {
    Name        string        // 任务名称
    Spec        string        // cron 表达式
    Description string        // 任务描述
    Timeout     time.Duration // 执行超时
    RetryCount  int           // 重试次数
    RetryDelay  time.Duration // 重试间隔
}
```

### 2. 服务配置

```go
// 创建带配置的 cron 服务
cron.New(cron.WithSeconds(), cron.WithChain(cron.Recover(cron.DefaultLogger)))
```

## 🚨 错误处理

### 1. 常见错误

```go
// cron 表达式错误
if err != nil {
    log.Errorf("invalid cron spec: %v", err)
    return
}

// 任务执行错误
if err != nil {
    log.Errorf("task execution failed: %v", err)
    // 记录错误指标
    metricsCollector.RecordTaskComplete(taskName, duration, err)
}
```

### 2. 恢复机制

```go
// panic 恢复
defer func() {
    if r := recover(); r != nil {
        log.Errorf("task panicked: %v", r)
        // 记录错误并恢复
    }
}()
```

## 📈 扩展建议

### 1. 分布式支持
- 集成 Redis 或 etcd 实现分布式锁
- 支持多实例任务协调
- 任务分片执行

### 2. 持久化存储
- 任务配置持久化到数据库
- 执行历史记录
- 指标数据存储

### 3. 监控集成
- Prometheus 指标导出
- Grafana 仪表板
- 告警机制

### 4. 高级功能
- 动态任务配置
- 任务依赖关系
- 资源限制控制

## 🤝 贡献指南

1. 遵循 Go 代码规范
2. 添加单元测试
3. 更新文档
4. 提交前运行测试

## 📄 许可证

本项目采用 MIT 许可证。 