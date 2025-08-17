# Nancalacc 项目文档

## 📋 项目概述

**项目名称**: Nancalacc (南卡拉克账户系统)  
**项目类型**: 企业级账户同步系统  
**技术栈**: Go + Kratos + MySQL + Redis + 钉钉API + WPS API  
**主要功能**: 钉钉组织架构数据同步、WPS系统集成、用户和部门关系管理  

## 🏗️ 系统架构

### 整体架构图
```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   钉钉 API      │    │   WPS API       │    │   用户界面      │
└─────────┬───────┘    └─────────┬───────┘    └─────────┬───────┘
          │                      │                      │
          └──────────────────────┼──────────────────────┘
                                 │
                    ┌─────────────┴─────────────┐
                    │      API 网关层           │
                    └─────────────┬─────────────┘
                                 │
                    ┌─────────────┴─────────────┐
                    │      业务逻辑层            │
                    │  ┌─────────────────────┐  │
                    │  │   AccountUsecase    │  │
                    │  │   FullSyncUsecase   │  │
                    │  │ IncrementalSyncUsecase│ │
                    │  └─────────────────────┘  │
                    └─────────────┬─────────────┘
                                 │
                    ┌─────────────┴─────────────┐
                    │      数据访问层            │
                    │  ┌─────────────────────┐  │
                    │  │     GORM + MySQL    │  │
                    │  │     Redis Cache     │  │
                    │  │   本地缓存服务      │  │
                    │  └─────────────────────┘  │
                    └───────────────────────────┘
```

### 核心模块说明

#### 1. 业务逻辑层 (internal/biz/)
- **AccountUsecase**: 基础账户操作管理
- **FullSyncUsecase**: 全量数据同步逻辑
- **IncrementalSyncUsecase**: 增量数据同步逻辑

#### 2. 数据访问层 (internal/data/)
- **MySQL**: 主数据存储，使用GORM ORM
- **Redis**: 分布式缓存和会话存储
- **本地缓存**: 高频数据的内存缓存

#### 3. 第三方集成
- **钉钉API**: 获取组织架构和用户信息
- **WPS API**: 企业级文档服务集成

## 🔧 技术实现细节

### 数据库设计
```sql
-- 主要数据表
tb_las_user              -- 用户表
tb_las_department        -- 部门表  
tb_las_department_user   -- 部门用户关系表
tb_las_company_cfg       -- 公司配置表
```

### 缓存策略
- **本地缓存**: 高频访问的任务状态和配置信息
- **Redis缓存**: 分布式缓存，支持多实例部署
- **缓存键前缀**: `nancalacc:cache:`

### 同步机制
1. **全量同步**: 定期完整同步钉钉组织架构
2. **增量同步**: 实时响应钉钉组织变更事件
3. **任务跟踪**: 记录同步任务状态和进度

## 📊 性能分析报告

### 当前性能瓶颈

#### 1. 数据库操作效率低
- **问题**: 单条插入用户和部门数据
- **影响**: 全量同步时间过长，资源利用率低
- **优化空间**: 3-5倍性能提升

#### 2. 串行处理限制
- **问题**: 部门、用户、关系保存串行执行
- **影响**: 无法充分利用系统资源
- **优化空间**: 40-60%时间减少

#### 3. 缓存策略简单
- **问题**: 单级缓存，缺乏智能失效策略
- **影响**: 缓存命中率低，数据库压力大
- **优化空间**: 60-80%查询性能提升

### 性能指标基准
```
当前性能指标:
- 全量同步时间: 约 30-60 分钟 (取决于数据量)
- 数据库写入: 约 100-500 条/秒
- 查询响应时间: 约 100-500ms
- 系统吞吐量: 约 100-200 请求/分钟
```

## 🚀 优化实施计划

### 阶段一: 核心性能优化 (1-2周)

#### 1.1 数据库批量操作优化
```go
// 优化前: 单条插入
for _, user := range users {
    db.Create(&user)
}

// 优化后: 批量插入 + 事务
func (r *accounterRepo) SaveUsersBatch(ctx context.Context, users []*models.TbLasUser) error {
    return r.data.GetSyncDB().Transaction(func(tx *gorm.DB) error {
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

#### 1.2 并发处理优化
```go
// 优化前: 串行处理
err = uc.repo.SaveDepartments(ctx, depts, taskId)
err = uc.repo.SaveUsers(ctx, users, taskId)
err = uc.repo.SaveDepartmentUserRelations(ctx, relations, taskId)

// 优化后: 并发处理
func (uc *FullSyncUsecase) saveDataConcurrently(ctx context.Context, depts, users, relations) error {
    var wg sync.WaitGroup
    errChan := make(chan error, 3)
    
    // 并发保存部门、用户、关系
    wg.Add(3)
    go func() { defer wg.Done(); errChan <- uc.repo.SaveDepartments(ctx, depts, taskId) }()
    go func() { defer wg.Done(); errChan <- uc.repo.SaveUsers(ctx, users, taskId) }()
    go func() { defer wg.Done(); errChan <- uc.repo.SaveDepartmentUserRelations(ctx, relations, taskId) }()
    
    wg.Wait()
    close(errChan)
    
    // 收集错误
    for err := range errChan {
        if err != nil {
            return err
        }
    }
    return nil
}
```

### 阶段二: 稳定性提升 (2-3周)

#### 2.1 重试和熔断机制
```go
type RetryConfig struct {
    MaxAttempts int
    Backoff     time.Duration
    MaxBackoff  time.Duration
}

type CircuitBreakerConfig struct {
    Threshold   int
    Timeout     time.Duration
    HalfOpen   time.Duration
}

func (uc *FullSyncUsecase) callWithRetryAndCircuitBreaker(ctx context.Context, operation func() error) error {
    if uc.circuitBreaker.IsOpen() {
        return errors.New("circuit breaker is open")
    }
    
    for attempt := 1; attempt <= uc.retryConfig.MaxAttempts; attempt++ {
        if err := operation(); err != nil {
            if !isRetryableError(err) {
                return err
            }
            
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
        
        uc.circuitBreaker.OnSuccess()
        return nil
    }
    
    uc.circuitBreaker.OnFailure()
    return fmt.Errorf("operation failed after %d attempts", uc.retryConfig.MaxAttempts)
}
```

#### 2.2 缓存策略优化
```go
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

### 阶段三: 监控和运维 (3-4周)

#### 3.1 性能指标收集
```go
type Metrics struct {
    syncDuration    prometheus.Histogram
    syncSuccess     prometheus.Counter
    syncFailure     prometheus.Counter
    dataSize        prometheus.Histogram
    cacheHitRate    prometheus.Gauge
}

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

## 📈 预期优化效果

### 性能提升指标
| 指标 | 优化前 | 优化后 | 提升幅度 |
|------|--------|--------|----------|
| 全量同步时间 | 30-60分钟 | 15-25分钟 | 50-70% |
| 数据库写入性能 | 100-500条/秒 | 300-1500条/秒 | 3-5倍 |
| 查询响应时间 | 100-500ms | 40-100ms | 60-80% |
| 系统吞吐量 | 100-200请求/分钟 | 200-400请求/分钟 | 2-3倍 |

### 稳定性提升
- **错误恢复能力**: 自动重试和熔断保护
- **资源利用率**: 更好的并发控制和资源管理
- **监控告警**: 及时发现问题并处理

## 🛠️ 开发环境配置

### 环境要求
- **Go版本**: 1.19+
- **MySQL版本**: 8.0+
- **Redis版本**: 6.0+
- **操作系统**: Linux/macOS/Windows

### 本地开发设置
```bash
# 1. 克隆项目
git clone <repository-url>
cd nancalacc_optimization

# 2. 安装依赖
go mod download

# 3. 配置环境变量
cp .env.example .env
# 编辑 .env 文件，配置数据库连接等

# 4. 启动依赖服务
docker-compose up -d

# 5. 运行项目
go run cmd/nancalacc/main.go
```

### 配置文件说明
```yaml
# configs/config.yaml
server:
  http:
    addr: 0.0.0.0:8000
    timeout: 1s
  grpc:
    addr: 0.0.0.0:9000
    timeout: 1s

data:
  database:
    driver: mysql
    source: "user:password@tcp(127.0.0.1:3306)/nancalacc?charset=utf8mb4&parseTime=True&loc=Local"
  redis:
    addr: 127.0.0.1:6379
    password: ""
    db: 0

app:
  third_company_id: "your_company_id"
  platform_ids: "your_platform_ids"
  company_id: "your_company_id"
```

## 🧪 测试策略

### 单元测试
```bash
# 运行所有测试
go test ./...

# 运行特定包的测试
go test ./internal/biz/...

# 运行测试并显示覆盖率
go test -cover ./...
```

### 性能测试
```bash
# 使用 wrk 进行 HTTP 性能测试
wrk -t12 -c400 -d30s http://localhost:8000/health

# 使用 go test 进行基准测试
go test -bench=. ./internal/biz/...
```

### 集成测试
```bash
# 运行集成测试
go test -tags=integration ./...

# 使用 Docker Compose 进行端到端测试
docker-compose -f docker-compose.test.yml up --abort-on-container-exit
```

## 📝 代码规范

### Go 代码规范
- 遵循 [Effective Go](https://golang.org/doc/effective_go.html) 规范
- 使用 `gofmt` 格式化代码
- 使用 `golint` 检查代码质量
- 使用 `go vet` 检查潜在问题

### 提交规范
```bash
# 提交前检查
go fmt ./...
go vet ./...
go test ./...

# 提交信息格式
feat: 添加用户批量导入功能
fix: 修复缓存失效问题
docs: 更新API文档
style: 格式化代码
refactor: 重构同步逻辑
test: 添加性能测试用例
chore: 更新依赖版本
```

## 🔍 故障排查指南

### 常见问题及解决方案

#### 1. 数据库连接问题
```bash
# 检查数据库连接
mysql -h localhost -u user -p -e "SELECT 1"

# 检查连接池配置
show variables like 'max_connections';
show status like 'Threads_connected';
```

#### 2. Redis 连接问题
```bash
# 检查 Redis 连接
redis-cli ping

# 检查 Redis 内存使用
redis-cli info memory
```

#### 3. 性能问题排查
```bash
# 使用 pprof 进行性能分析
go tool pprof http://localhost:8000/debug/pprof/profile

# 使用 trace 进行追踪分析
go tool trace trace.out
```

## 📚 相关资源

### 官方文档
- [Kratos 框架文档](https://go-kratos.dev/)
- [GORM 文档](https://gorm.io/)
- [钉钉开放平台文档](https://open.dingtalk.com/)
- [WPS 开放平台文档](https://open.wps.cn/)

### 技术博客
- [Go 性能优化实践](https://blog.golang.org/profiling-go-programs)
- [MySQL 性能调优](https://dev.mysql.com/doc/refman/8.0/en/optimization.html)
- [Redis 最佳实践](https://redis.io/topics/optimization)

### 工具推荐
- **性能分析**: pprof, trace, go-torch
- **代码质量**: golint, go vet, staticcheck
- **测试工具**: testify, gomock, sqlmock
- **监控工具**: Prometheus, Grafana, Jaeger

## 📅 项目时间线

### 2024年1月
- [x] 项目初始化和代码分析
- [x] 性能瓶颈识别
- [x] 优化计划制定

### 2024年2月
- [ ] 阶段一：核心性能优化
- [ ] 数据库批量操作优化
- [ ] 并发处理优化

### 2024年3月
- [ ] 阶段二：稳定性提升
- [ ] 重试和熔断机制
- [ ] 缓存策略优化

### 2024年4月
- [ ] 阶段三：监控和运维
- [ ] 性能指标收集
- [ ] 告警机制实现

## 👥 团队分工

### 开发团队
- **架构师**: 负责整体架构设计和优化方案
- **后端开发**: 负责核心业务逻辑优化
- **数据库专家**: 负责数据库性能优化
- **运维工程师**: 负责监控和部署优化

### 协作方式
- **代码审查**: 所有代码变更需要至少一名团队成员审查
- **定期同步**: 每周进行进度同步和问题讨论
- **文档更新**: 及时更新技术文档和操作手册

---

**文档维护**: 开发团队  
**最后更新**: 2024年1月  
**版本**: v1.0.0  

---

*如有问题或建议，请联系开发团队或提交 Issue。* 