# 业务逻辑优化计划

## 项目概述
nancalacc 是一个企业级账户同步系统，主要功能包括：
- 钉钉组织架构数据同步（全量/增量）
- WPS 系统集成
- 用户和部门关系管理
- 任务状态跟踪和缓存

## 当前架构分析

### 核心业务模块
1. **AccountUsecase**: 基础账户操作
2. **FullSyncUsecase**: 全量同步逻辑
    - 1. 验证是否提交过uc.localCache，否则返回错误
    - 2. 获取token并从第三方获取部门和用户数据
    - 3. 数据入库，包括SaveCompanyCfg公司配置入库，SaveDepartments部门入库，SaveUsers用户入库，SaveDepartmentUserRelations关系入
    - 4. 调用wps接口，通知调用 PostEcisaccountsyncAll 开始同步
    - 5. 更新任务状态uc.localCache ，包括任务状态，进度，开始时间，结束时间，实际时间
    - 6. 返回结果，包括任务id

其中3里面的的三项可以是并发，但需要保证数据一致性，
2里面已经有并发，但需要保证数据一致性，
确保可以入库完成时，再有通知调用
3. **IncrementalSyncUsecase**: 增量同步逻辑
4. **数据层**: GORM + MySQL + Redis

### 技术栈
- **框架**: Kratos (Go)
- **数据库**: MySQL + Redis
- **ORM**: GORM
- **缓存**: Redis + 本地缓存
- **第三方集成**: 钉钉 API + WPS API

## 性能瓶颈分析

### 1. 数据库操作优化

#### 问题识别
- **批量插入效率低**: `SaveUsers` 和 `SaveDepartments` 使用单条插入
- **分页查询性能**: `BatchGetDeptUsers` 和 `BatchGetUsers` 使用 OFFSET 分页
- **事务管理缺失**: 缺乏批量操作的事务控制
- **索引优化不足**: 可能缺少关键字段的复合索引

#### 优化方案
```go
// 优化前：单条插入
for _, user := range users {
    db.Create(&user)
}

// 优化后：批量插入 + 事务
func (r *accounterRepo) SaveUsersBatch(ctx context.Context, users []*models.TbLasUser) error {
    return r.data.GetSyncDB().Transaction(func(tx *gorm.DB) error {
        // 使用批量插入，每批1000条
        batchSize := 1000
        for i := 0; i < len(users); i += batchSize {
            end := i + batchSize
            if end > len(users) {
                end = len(users)
            }
            if err := tx.CreateInBatches(users[i:end], batchSize).Error; err != nil {
                return err
            }
        }
        return nil
    })
}
```

### 2. 并发处理优化

#### 问题识别
- **串行处理**: 全量同步中部门、用户、关系保存是串行的
- **缺乏并发控制**: 没有利用 Go 的并发特性
- **资源竞争**: 多个操作可能竞争数据库连接

#### 优化方案
```go
// 优化前：串行处理
err = uc.repo.SaveCompanyCfg(ctx, companyCfg)
err = uc.repo.SaveDepartments(ctx, depts, taskId)
err = uc.repo.SaveUsers(ctx, users, taskId)
err = uc.repo.SaveDepartmentUserRelations(ctx, relations, taskId)

// 优化后：并发处理 + 错误聚合
func (uc *FullSyncUsecase) saveDataConcurrently(ctx context.Context, depts, users, relations) error {
    var wg sync.WaitGroup
    errChan := make(chan error, 3)
    
    // 并发保存部门
    wg.Add(1)
    go func() {
        defer wg.Done()
        if err := uc.repo.SaveDepartments(ctx, depts, taskId); err != nil {
            errChan <- fmt.Errorf("save departments failed: %w", err)
        }
    }()
    
    // 并发保存用户
    wg.Add(1)
    go func() {
        defer wg.Done()
        if err := uc.repo.SaveUsers(ctx, users, taskId); err != nil {
            errChan <- fmt.Errorf("save users failed: %w", err)
        }
    }()
    
    // 并发保存关系
    wg.Add(1)
    go func() {
        defer wg.Done()
        if err := uc.repo.SaveDepartmentUserRelations(ctx, relations, taskId); err != nil {
            errChan <- fmt.Errorf("save relations failed: %w", err)
        }
    }()
    
    wg.Wait()
    close(errChan)
    
    // 收集所有错误
    var errors []error
    for err := range errChan {
        errors = append(errors, err)
    }
    
    if len(errors) > 0 {
        return fmt.Errorf("concurrent operations failed: %v", errors)
    }
    
    return nil
}
```

### 3. 缓存策略优化

#### 问题识别
- **缓存穿透**: 频繁查询不存在的任务
- **缓存雪崩**: 大量缓存同时过期
- **缓存一致性**: 缓存与数据库数据不一致

#### 优化方案
```go
// 优化前：简单缓存
func (uc *AccounterUsecase) GetCacheTask(ctx context.Context, taskName string) (*models.Task, error) {
    cacheKey := prefix + taskName
    task, ok, err := uc.localCache.Get(ctx, cacheKey)
    // ... 简单返回
}

// 优化后：多级缓存 + 布隆过滤器
type CacheStrategy struct {
    localCache  CacheService
    redisCache  CacheService
    bloomFilter *bloom.BloomFilter
}

func (cs *CacheStrategy) GetTask(ctx context.Context, taskName string) (*models.Task, error) {
    // 1. 布隆过滤器快速判断
    if !cs.bloomFilter.Test([]byte(taskName)) {
        return nil, errors.New("task not found")
    }
    
    // 2. 本地缓存
    if task, ok, err := cs.localCache.Get(ctx, taskName); ok && err == nil {
        return task.(*models.Task), nil
    }
    
    // 3. Redis 缓存
    if task, ok, err := cs.redisCache.Get(ctx, taskName); ok && err == nil {
        // 回填本地缓存
        cs.localCache.Set(ctx, taskName, task, 5*time.Minute)
        return task.(*models.Task), nil
    }
    
    // 4. 数据库查询
    task, err := cs.repo.GetTask(ctx, taskName)
    if err != nil {
        return nil, err
    }
    
    // 5. 更新缓存
    cs.redisCache.Set(ctx, taskName, task, 30*time.Minute)
    cs.localCache.Set(ctx, taskName, task, 5*time.Minute)
    
    return task, nil
}
```

### 4. 错误处理和重试机制

#### 问题识别
- **缺乏重试机制**: API 调用失败后直接返回错误
- **错误分类不清晰**: 没有区分可重试和不可重试的错误
- **缺乏熔断器**: 第三方服务异常时没有保护机制

#### 优化方案
```go
// 重试配置
type RetryConfig struct {
    MaxAttempts int
    Backoff     time.Duration
    MaxBackoff  time.Duration
}

// 熔断器配置
type CircuitBreakerConfig struct {
    Threshold   int
    Timeout     time.Duration
    HalfOpen   time.Duration
}

// 带重试和熔断的 API 调用
func (uc *FullSyncUsecase) callWithRetryAndCircuitBreaker(ctx context.Context, operation func() error) error {
    // 熔断器检查
    if uc.circuitBreaker.IsOpen() {
        return errors.New("circuit breaker is open")
    }
    
    var lastErr error
    for attempt := 1; attempt <= uc.retryConfig.MaxAttempts; attempt++ {
        if err := operation(); err != nil {
            lastErr = err
            
            // 判断是否可重试
            if !isRetryableError(err) {
                return err
            }
            
            // 指数退避
            backoff := time.Duration(attempt) * uc.retryConfig.Backoff
            if backoff > uc.retryConfig.MaxBackoff {
                backoff = uc.retryConfig.MaxBackoff
            }
            
            select {
            case <-ctx.Done():
                return ctx.Err()
            case <-time.After(backoff):
                continue
            }
        }
        
        // 成功，重置熔断器
        uc.circuitBreaker.OnSuccess()
        return nil
    }
    
    // 达到最大重试次数，触发熔断器
    uc.circuitBreaker.OnFailure()
    return fmt.Errorf("operation failed after %d attempts: %w", uc.retryConfig.MaxAttempts, lastErr)
}
```

### 5. 监控和指标收集

#### 问题识别
- **缺乏性能指标**: 没有关键操作的耗时统计
- **缺乏业务指标**: 没有同步成功率、数据量等统计
- **缺乏告警机制**: 异常情况没有及时通知

#### 优化方案
```go
// 性能指标收集
type Metrics struct {
    syncDuration    prometheus.Histogram
    syncSuccess     prometheus.Counter
    syncFailure     prometheus.Counter
    dataSize        prometheus.Histogram
    cacheHitRate    prometheus.Gauge
}

// 在关键操作中添加指标收集
func (uc *FullSyncUsecase) CreateSyncAccount(ctx context.Context, req *v1.CreateSyncAccountRequest) (*v1.CreateSyncAccountReply, error) {
    start := time.Now()
    defer func() {
        uc.metrics.syncDuration.Observe(time.Since(start).Seconds())
    }()
    
    // ... 业务逻辑
    
    if err != nil {
        uc.metrics.syncFailure.Inc()
        return nil, err
    }
    
    uc.metrics.syncSuccess.Inc()
    return reply, nil
}
```

## 实施优先级

### 高优先级 (立即实施)
1. **数据库批量操作优化**: 提升数据插入性能 3-5倍
2. **并发处理优化**: 减少全量同步时间 40-60%
3. **基础缓存优化**: 减少数据库查询压力

### 中优先级 (1-2周内)
1. **重试和熔断机制**: 提升系统稳定性
2. **分页查询优化**: 改善大数据量查询性能
3. **监控指标收集**: 为后续优化提供数据支撑

### 低优先级 (长期规划)
1. **多级缓存策略**: 进一步减少延迟
2. **数据库索引优化**: 根据实际查询模式优化
3. **异步处理**: 将非关键操作异步化

## 预期效果

### 性能提升
- **全量同步时间**: 从当前时间减少 50-70%
- **数据库写入性能**: 提升 3-5倍
- **查询响应时间**: 减少 60-80%
- **系统吞吐量**: 提升 2-3倍

### 稳定性提升
- **错误恢复能力**: 自动重试和熔断保护
- **资源利用率**: 更好的并发控制和资源管理
- **监控告警**: 及时发现问题并处理

### 可维护性提升
- **代码结构**: 更清晰的错误处理和重试逻辑
- **性能监控**: 详细的性能指标和告警
- **配置管理**: 可配置的重试和熔断参数

## 风险评估

### 技术风险
- **并发控制复杂性**: 需要仔细处理竞态条件
- **缓存一致性**: 多级缓存可能导致数据不一致
- **数据库连接池**: 并发增加可能影响连接池配置

### 缓解措施
- **充分测试**: 单元测试和集成测试覆盖
- **渐进式部署**: 分阶段实施，监控效果
- **回滚计划**: 准备快速回滚方案

## 下一步行动

1. **代码审查**: 与团队讨论优化方案
2. **性能基准测试**: 建立当前性能基准
3. **分阶段实施**: 按优先级逐步实施优化
4. **效果验证**: 每个阶段完成后验证效果
5. **持续优化**: 根据实际运行情况持续改进

---
*最后更新: 2024年1月* 