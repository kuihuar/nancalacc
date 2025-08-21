# 接口类型 (Interface) - Interface Types

## 概述
接口类型定义了对象的行为契约，提供抽象层和依赖注入的基础。在 Go 中，接口是隐式实现的，提供了良好的解耦和测试能力。

## 分类

### 6.1 数据访问接口
```go
// 用户数据访问接口
type UserRepository interface {
    Create(ctx context.Context, user *TbLasUser) error
    GetByID(ctx context.Context, id uint) (*TbLasUser, error)
    GetByTaskID(ctx context.Context, taskID string) ([]*TbLasUser, error)
    GetByThirdCompanyID(ctx context.Context, thirdCompanyID string) ([]*TbLasUser, error)
    Update(ctx context.Context, user *TbLasUser) error
    Delete(ctx context.Context, id uint) error
    BatchCreate(ctx context.Context, users []*TbLasUser) error
    BatchUpdate(ctx context.Context, users []*TbLasUser) error
    Count(ctx context.Context, filter UserFilter) (int64, error)
}

// 部门数据访问接口
type DepartmentRepository interface {
    Create(ctx context.Context, dept *TbLasDept) error
    GetByID(ctx context.Context, id uint) (*TbLasDept, error)
    GetByTaskID(ctx context.Context, taskID string) ([]*TbLasDept, error)
    GetByParentID(ctx context.Context, parentID string) ([]*TbLasDept, error)
    Update(ctx context.Context, dept *TbLasDept) error
    Delete(ctx context.Context, id uint) error
    BatchCreate(ctx context.Context, depts []*TbLasDept) error
    GetTree(ctx context.Context, taskID string) ([]*TbLasDept, error)
}

// 任务数据访问接口
type TaskRepository interface {
    Create(ctx context.Context, task *TbLasTask) error
    GetByID(ctx context.Context, id uint) (*TbLasTask, error)
    GetByTaskID(ctx context.Context, taskID string) (*TbLasTask, error)
    Update(ctx context.Context, task *TbLasTask) error
    UpdateStatus(ctx context.Context, taskID string, status string, progress int) error
    List(ctx context.Context, filter TaskFilter) ([]*TbLasTask, error)
    Count(ctx context.Context, filter TaskFilter) (int64, error)
}
```

### 6.2 服务接口
```go
// 同步服务接口
type SyncService interface {
    CreateSyncTask(ctx context.Context, req *CreateSyncAccountRequest) (*TbLasTask, error)
    GetSyncTask(ctx context.Context, taskID string) (*TbLasTask, error)
    CancelSyncTask(ctx context.Context, taskID string) error
    ExecuteSyncTask(ctx context.Context, taskID string) error
    GetSyncProgress(ctx context.Context, taskID string) (*SyncProgress, error)
    RetryFailedTask(ctx context.Context, taskID string) error
}

// 第三方API服务接口
type ThirdPartyAPIService interface {
    GetUsers(ctx context.Context, companyID string) ([]*DingtalkDeptUser, error)
    GetDepartments(ctx context.Context, companyID string) ([]*DingtalkDept, error)
    GetUserDetails(ctx context.Context, userID string) (*DingtalkDeptUser, error)
    GetDepartmentDetails(ctx context.Context, deptID string) (*DingtalkDept, error)
    ValidateCredentials(ctx context.Context) error
    GetAccessToken(ctx context.Context) (string, error)
}

// 缓存服务接口
type CacheService interface {
    Get(ctx context.Context, key string) (interface{}, error)
    Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
    Delete(ctx context.Context, key string) error
    Exists(ctx context.Context, key string) (bool, error)
    Incr(ctx context.Context, key string) (int64, error)
    IncrBy(ctx context.Context, key string, value int64) (int64, error)
    Expire(ctx context.Context, key string, ttl time.Duration) error
    Clear(ctx context.Context, pattern string) error
}
```

### 6.3 消息队列接口
```go
// 消息生产者接口
type MessageProducer interface {
    Publish(ctx context.Context, topic string, message interface{}) error
    PublishAsync(ctx context.Context, topic string, message interface{}) error
    PublishBatch(ctx context.Context, topic string, messages []interface{}) error
    Close() error
}

// 消息消费者接口
type MessageConsumer interface {
    Subscribe(ctx context.Context, topic string, handler MessageHandler) error
    Unsubscribe(ctx context.Context, topic string) error
    Start() error
    Stop() error
    IsRunning() bool
}

// 消息处理器接口
type MessageHandler interface {
    Handle(ctx context.Context, message []byte) error
    HandleAsync(ctx context.Context, message []byte) error
    GetRetryCount() int
    GetRetryDelay() time.Duration
}
```

### 6.4 监控和日志接口
```go
// 日志接口
type Logger interface {
    Debug(ctx context.Context, msg string, fields ...Field)
    Info(ctx context.Context, msg string, fields ...Field)
    Warn(ctx context.Context, msg string, fields ...Field)
    Error(ctx context.Context, msg string, fields ...Field)
    Fatal(ctx context.Context, msg string, fields ...Field)
    WithContext(ctx context.Context) Logger
    WithFields(fields ...Field) Logger
}

// 指标收集接口
type MetricsCollector interface {
    IncrementCounter(name string, labels map[string]string)
    IncrementCounterBy(name string, value float64, labels map[string]string)
    SetGauge(name string, value float64, labels map[string]string)
    RecordHistogram(name string, value float64, labels map[string]string)
    RecordSummary(name string, value float64, labels map[string]string)
    StartTimer(name string, labels map[string]string) Timer
}

// 追踪接口
type Tracer interface {
    StartSpan(ctx context.Context, name string, opts ...SpanOption) (Span, context.Context)
    Inject(ctx context.Context, span Span, format interface{}, carrier interface{}) error
    Extract(ctx context.Context, format interface{}, carrier interface{}) (SpanContext, error)
    Close() error
}
```

### 6.5 配置管理接口
```go
// 配置提供者接口
type ConfigProvider interface {
    Get(key string) interface{}
    GetString(key string) string
    GetInt(key string) int
    GetBool(key string) bool
    GetDuration(key string) time.Duration
    GetStringSlice(key string) []string
    GetMap(key string) map[string]interface{}
    Set(key string, value interface{})
    Watch(key string, callback func(interface{})) error
    Unwatch(key string) error
}

// 配置验证器接口
type ConfigValidator interface {
    Validate(config interface{}) error
    ValidateField(field string, value interface{}) error
    AddRule(field string, rule ValidationRule) error
    RemoveRule(field string) error
}
```

### 6.6 工具接口
```go
// 加密接口
type Encryptor interface {
    Encrypt(data []byte) ([]byte, error)
    Decrypt(data []byte) ([]byte, error)
    Hash(data []byte) (string, error)
    VerifyHash(data []byte, hash string) (bool, error)
}

// 序列化接口
type Serializer interface {
    Serialize(data interface{}) ([]byte, error)
    Deserialize(data []byte, target interface{}) error
    GetContentType() string
}

// 重试接口
type RetryPolicy interface {
    ShouldRetry(err error, attempt int) bool
    GetDelay(attempt int) time.Duration
    GetMaxAttempts() int
    OnRetry(err error, attempt int)
}
```

## 特点

1. **抽象性**: 提供抽象层，隐藏实现细节
2. **可测试性**: 便于单元测试和模拟
3. **可扩展性**: 支持多种实现
4. **依赖注入**: 支持依赖注入和反转控制
5. **解耦**: 降低模块间的耦合度

## 使用场景

- 依赖注入
- 单元测试
- 插件架构
- 多实现支持
- 接口隔离

## 最佳实践

1. **接口设计**: 设计小而专注的接口
2. **命名规范**: 使用清晰的命名约定
3. **文档化**: 为接口提供清晰的文档
4. **版本控制**: 考虑接口的向后兼容性
5. **实现检查**: 在编译时检查接口实现 