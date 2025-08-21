# 枚举类型 (Enum) - Enumeration Types

## 概述
枚举类型用于定义一组固定的常量值，提供类型安全和代码可读性。在 Go 中通常使用 `const` 和 `iota` 来实现枚举。

## 分类

### 4.1 任务状态枚举
```go
// 任务状态
type TaskStatus int

const (
    TaskStatusPending TaskStatus = iota
    TaskStatusRunning
    TaskStatusCompleted
    TaskStatusFailed
    TaskStatusCancelled
)

// 任务状态字符串映射
var taskStatusMap = map[TaskStatus]string{
    TaskStatusPending:   "pending",
    TaskStatusRunning:   "running",
    TaskStatusCompleted: "completed",
    TaskStatusFailed:    "failed",
    TaskStatusCancelled: "cancelled",
}

// 获取任务状态字符串
func (ts TaskStatus) String() string {
    return taskStatusMap[ts]
}

// 从字符串解析任务状态
func ParseTaskStatus(s string) (TaskStatus, error) {
    for status, str := range taskStatusMap {
        if str == s {
            return status, nil
        }
    }
    return TaskStatusPending, fmt.Errorf("invalid task status: %s", s)
}
```

### 4.2 同步类型枚举
```go
// 同步类型
type SyncType int

const (
    SyncTypeFull SyncType = iota
    SyncTypeIncremental
    SyncTypeDelta
)

// 同步类型字符串映射
var syncTypeMap = map[SyncType]string{
    SyncTypeFull:        "full",
    SyncTypeIncremental: "incremental",
    SyncTypeDelta:       "delta",
}

func (st SyncType) String() string {
    return syncTypeMap[st]
}

func ParseSyncType(s string) (SyncType, error) {
    for syncType, str := range syncTypeMap {
        if str == s {
            return syncType, nil
        }
    }
    return SyncTypeFull, fmt.Errorf("invalid sync type: %s", s)
}
```

### 4.3 触发类型枚举
```go
// 触发类型
type TriggerType int

const (
    TriggerTypeManual TriggerType = iota
    TriggerTypeScheduled
    TriggerTypeEvent
    TriggerTypeAPI
)

// 触发类型字符串映射
var triggerTypeMap = map[TriggerType]string{
    TriggerTypeManual:    "manual",
    TriggerTypeScheduled: "scheduled",
    TriggerTypeEvent:     "event",
    TriggerTypeAPI:       "api",
}

func (tt TriggerType) String() string {
    return triggerTypeMap[tt]
}

func ParseTriggerType(s string) (TriggerType, error) {
    for triggerType, str := range triggerTypeMap {
        if str == s {
            return triggerType, nil
        }
    }
    return TriggerTypeManual, fmt.Errorf("invalid trigger type: %s", s)
}
```

### 4.4 平台类型枚举
```go
// 平台类型
type PlatformType int

const (
    PlatformTypeDingtalk PlatformType = iota
    PlatformTypeWechat
    PlatformTypeFeishu
    PlatformTypeLark
)

// 平台类型字符串映射
var platformTypeMap = map[PlatformType]string{
    PlatformTypeDingtalk: "dingtalk",
    PlatformTypeWechat:   "wechat",
    PlatformTypeFeishu:   "feishu",
    PlatformTypeLark:     "lark",
}

func (pt PlatformType) String() string {
    return platformTypeMap[pt]
}

func ParsePlatformType(s string) (PlatformType, error) {
    for platformType, str := range platformTypeMap {
        if str == s {
            return platformType, nil
        }
    }
    return PlatformTypeDingtalk, fmt.Errorf("invalid platform type: %s", s)
}
```

### 4.5 用户状态枚举
```go
// 用户状态
type UserStatus int

const (
    UserStatusActive UserStatus = iota
    UserStatusInactive
    UserStatusSuspended
    UserStatusDeleted
)

// 用户状态字符串映射
var userStatusMap = map[UserStatus]string{
    UserStatusActive:    "active",
    UserStatusInactive:  "inactive",
    UserStatusSuspended: "suspended",
    UserStatusDeleted:   "deleted",
}

func (us UserStatus) String() string {
    return userStatusMap[us]
}

func ParseUserStatus(s string) (UserStatus, error) {
    for userStatus, str := range userStatusMap {
        if str == s {
            return userStatus, nil
        }
    }
    return UserStatusActive, fmt.Errorf("invalid user status: %s", s)
}
```

### 4.6 Saga 状态枚举
```go
// Saga 事务状态
type SagaStatus int

const (
    SagaStatusStarted SagaStatus = iota
    SagaStatusInProgress
    SagaStatusCompleted
    SagaStatusFailed
    SagaStatusCompensating
    SagaStatusCompensated
)

// Saga 步骤状态
type StepStatus int

const (
    StepStatusPending StepStatus = iota
    StepStatusRunning
    StepStatusCompleted
    StepStatusFailed
    StepStatusCompensating
    StepStatusCompensated
    StepStatusSkipped
)
```

## 特点

1. **类型安全**: 编译时检查，避免无效值
2. **可读性**: 使用有意义的常量名
3. **维护性**: 集中管理相关常量
4. **扩展性**: 易于添加新的枚举值
5. **序列化**: 支持与字符串的相互转换

## 使用场景

- 状态管理
- 类型标识
- 配置选项
- 错误码定义
- 权限级别

## 最佳实践

1. **命名规范**: 使用清晰的前缀和描述性名称
2. **文档化**: 为每个枚举值提供注释说明
3. **验证**: 提供字符串解析和验证方法
4. **默认值**: 定义合理的默认枚举值
5. **序列化**: 实现与 JSON、数据库的序列化支持 