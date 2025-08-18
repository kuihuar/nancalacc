# 任务配置使用指南

本文档详细说明如何在系统中添加包含任务配置的定时任务。

## 概述

任务配置系统提供了以下功能：
- 任务超时控制
- 重试机制
- 退避策略（固定、指数、抖动）
- 动态配置调整
- 详细的执行结果记录

## 基本用法

### 1. 简单配置任务

```go
// 创建任务配置
config := DefaultJobConfig("simple_task", "0 */15 * * * *") // 每15分钟执行
config.WithTimeout(5 * time.Minute)
config.Description = "简单的配置化任务示例"

// 注册任务
_, err := s.AddJobWithConfig(config, JobExecutorFunc(func(ctx context.Context) error {
    // 任务执行逻辑
    s.log.Log(log.LevelInfo, "msg", "executing simple config job")
    return nil
}))
```

### 2. 带重试的任务

```go
// 创建任务配置
config := DefaultJobConfig("retry_task", "0 */30 * * * *") // 每30分钟执行
config.WithRetry(3, 2*time.Minute) // 重试3次，间隔2分钟
config.WithTimeout(10 * time.Minute)
config.Description = "带重试机制的任务示例"

// 注册任务
_, err := s.AddJobWithConfig(config, JobExecutorFunc(func(ctx context.Context) error {
    // 可能失败的任务逻辑
    if time.Now().Second()%2 == 0 {
        return fmt.Errorf("simulated failure")
    }
    return nil
}))
```

### 3. 带退避策略的任务

```go
// 创建任务配置
config := DefaultJobConfig("backoff_task", "0 0 */2 * * *") // 每2小时执行
config.Backoff = BackoffConfig{
    Type:      BackoffExponential,
    BaseDelay: 30 * time.Second,
    MaxDelay:  5 * time.Minute,
    Factor:    2.0,
}
config.WithRetry(5, 0) // 重试5次，使用退避策略
config.WithTimeout(15 * time.Minute)

// 注册任务
_, err := s.AddJobWithConfig(config, JobExecutorFunc(func(ctx context.Context) error {
    // 网络请求等可能失败的任务
    return nil
}))
```

## 配置选项详解

### JobConfig 结构

```go
type JobConfig struct {
    Name        string        `json:"name"`        // 任务名称
    Spec        string        `json:"spec"`        // cron 表达式
    Description string        `json:"description"` // 任务描述
    Timeout     time.Duration `json:"timeout"`     // 执行超时时间
    RetryCount  int           `json:"retry_count"` // 重试次数
    RetryDelay  time.Duration `json:"retry_delay"` // 重试间隔
    Enabled     bool          `json:"enabled"`     // 是否启用
    MaxRetries  int           `json:"max_retries"` // 最大重试次数
    Backoff     BackoffConfig `json:"backoff"`     // 退避策略配置
}
```

### BackoffConfig 结构

```go
type BackoffConfig struct {
    Type      BackoffType   `json:"type"`       // 退避类型
    BaseDelay time.Duration `json:"base_delay"` // 基础延迟时间
    MaxDelay  time.Duration `json:"max_delay"`  // 最大延迟时间
    Factor    float64       `json:"factor"`     // 指数退避因子
}
```

### 退避类型

- `BackoffFixed`: 固定延迟
- `BackoffExponential`: 指数退避
- `BackoffJitter`: 抖动退避

## 配置方法

### 默认配置

```go
// 创建默认配置
config := DefaultJobConfig("task_name", "0 */5 * * * *")
```

默认值：
- Timeout: 30分钟
- RetryCount: 3次
- RetryDelay: 5分钟
- MaxRetries: 3次
- Backoff: 指数退避，基础延迟1分钟，最大延迟30分钟，因子2.0

### 链式配置

```go
config := DefaultJobConfig("task_name", "0 */5 * * * *").
    WithTimeout(10 * time.Minute).
    WithRetry(5, 2*time.Minute)
```

### 自定义退避策略

```go
config := DefaultJobConfig("task_name", "0 */5 * * * *")
config.Backoff = BackoffConfig{
    Type:      BackoffExponential,
    BaseDelay: 1 * time.Minute,
    MaxDelay:  10 * time.Minute,
    Factor:    2.0,
}
```

## 动态配置

### 基于时间的动态配置

```go
func getDynamicConfig() *JobConfig {
    now := time.Now()
    config := DefaultJobConfig("dynamic", "")
    
    // 根据时间调整配置
    if now.Hour() >= 22 || now.Hour() <= 6 {
        // 夜间时段，增加超时时间
        config.Timeout = 20 * time.Minute
        config.RetryCount = 5
    } else {
        // 白天时段，正常配置
        config.Timeout = 10 * time.Minute
        config.RetryCount = 3
    }
    
    return config
}
```

### 基于系统负载的动态配置

```go
func getLoadBasedConfig() *JobConfig {
    config := DefaultJobConfig("load_based", "")
    
    if isHighLoad() {
        config.RetryDelay = 10 * time.Minute
        config.Timeout = 20 * time.Minute
    } else {
        config.RetryDelay = 2 * time.Minute
        config.Timeout = 10 * time.Minute
    }
    
    return config
}
```

## 执行结果

### JobResult 结构

```go
type JobResult struct {
    JobName    string        `json:"job_name"`
    StartTime  time.Time     `json:"start_time"`
    EndTime    time.Time     `json:"end_time"`
    Duration   time.Duration `json:"duration"`
    Success    bool          `json:"success"`
    Error      error         `json:"error,omitempty"`
    RetryCount int           `json:"retry_count"`
    IsTimeout  bool          `json:"is_timeout"`
    IsRetry    bool          `json:"is_retry"`
}
```

## 最佳实践

### 1. 任务命名

- 使用描述性的任务名称
- 遵循命名约定：`{功能}_{类型}_task`

### 2. 超时设置

- 根据任务复杂度设置合理的超时时间
- 考虑网络延迟和资源竞争

### 3. 重试策略

- 对于网络请求任务，使用指数退避
- 对于计算密集型任务，使用固定延迟
- 设置合理的最大重试次数

### 4. 错误处理

```go
_, err := s.AddJobWithConfig(config, JobExecutorFunc(func(ctx context.Context) error {
    // 检查上下文是否已取消
    if ctx.Err() != nil {
        return ctx.Err()
    }
    
    // 执行任务逻辑
    if err := doSomething(ctx); err != nil {
        // 记录详细错误信息
        s.log.Log(log.LevelError, "msg", "task execution failed", "error", err)
        return err
    }
    
    return nil
}))
```

### 5. 日志记录

- 记录任务开始和结束
- 记录重要的中间状态
- 记录错误详情和重试信息

## 示例完整代码

```go
// 在 jobs.go 中添加配置化任务
func RegisterConfigJobs(s *CronService) {
    // 数据同步任务
    syncConfig := DefaultJobConfig("data_sync", "0 0 2 * * *") // 每天凌晨2点
    syncConfig.WithTimeout(30 * time.Minute)
    syncConfig.WithRetry(3, 5*time.Minute)
    syncConfig.Description = "数据同步任务"
    
    _, err := s.AddJobWithConfig(syncConfig, JobExecutorFunc(func(ctx context.Context) error {
        s.log.Log(log.LevelInfo, "msg", "starting data sync")
        
        // 执行数据同步
        if err := s.data.SyncData(ctx); err != nil {
            s.log.Log(log.LevelError, "msg", "data sync failed", "error", err)
            return err
        }
        
        s.log.Log(log.LevelInfo, "msg", "data sync completed")
        return nil
    }))
    
    if err != nil {
        s.log.Log(log.LevelError, "msg", "failed to register data sync job", "err", err)
    }
}
```

## 注意事项

1. **任务名称唯一性**: 确保任务名称在系统中唯一
2. **资源管理**: 合理设置超时时间，避免资源泄露
3. **错误处理**: 正确处理和记录错误，便于问题排查
4. **性能考虑**: 避免在任务中执行耗时操作，考虑异步处理
5. **监控告警**: 为重要任务设置监控和告警机制 