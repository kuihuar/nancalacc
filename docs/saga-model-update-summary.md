# Saga 模型更新总结

## 📋 **更新概述**

为了与 `SAGA_DISTRIBUTED_TRANSACTION_PLAN.md` 文档设计保持一致，我们对 Saga 相关的模型进行了全面更新。

## 🔄 **主要变更**

### 1. **表结构更新**

| 原表名 | 新表名 | 变更说明 |
|--------|--------|----------|
| `saga_instances` | `saga_transactions` | 重命名以更好地反映事务概念 |
| `saga_steps` | `saga_steps` | 保持原名，但字段结构更新 |
| ❌ 缺失 | `saga_events` | **新增事件表**，用于完整的审计日志 |

### 2. **字段结构更新**

#### **Saga 事务表 (`saga_transactions`)**
```sql
-- 原字段
instance_id, service_name, status, data

-- 新字段
transaction_id, name, status, current_step, progress, start_time, end_time
```

#### **Saga 步骤表 (`saga_steps`)**
```sql
-- 原字段
instance_id, step_id, step_name, order, status, compensate, request_data, response_data, error_msg, retry_count, max_retries

-- 新字段
step_id, transaction_id, step_name, status, action_data, compensate_data, error_message, retry_count, max_retries, start_time, end_time
```

#### **Saga 事件表 (`saga_events`) - 新增**
```sql
-- 新字段
transaction_id, step_id, event_type, event_data, created_at
```

### 3. **状态枚举更新**

#### **Saga 状态**
```go
// 原状态
SagaStatusPending, SagaStatusRunning, SagaStatusCompleted, SagaStatusFailed, SagaStatusCompensating, SagaStatusCompensated

// 新状态
SagaStatusPending, SagaStatusInProgress, SagaStatusCompleted, SagaStatusFailed, SagaStatusCompensating, SagaStatusCompensated
```

#### **步骤状态**
```go
// 原状态
StepStatusPending, StepStatusRunning, StepStatusCompleted, StepStatusFailed, StepStatusCompensating, StepStatusCompensated

// 新状态
StepStatusPending, StepStatusInProgress, StepStatusCompleted, StepStatusFailed, StepStatusCompensating, StepStatusCompensated
```

#### **事件类型 - 新增**
```go
EventTypeSagaStarted, EventTypeStepStarted, EventTypeStepCompleted, EventTypeStepFailed,
EventTypeCompensationStarted, EventTypeCompensationCompleted, EventTypeSagaCompleted, EventTypeSagaFailed
```

## 🏗️ **新的数据模型**

### 1. **SagaTransaction 模型**
```go
type SagaTransaction struct {
    ID            uint           `gorm:"primarykey"`
    TransactionID string         `gorm:"uniqueIndex;size:64;not null"`
    Name          string         `gorm:"size:255;not null"`
    Status        SagaStatus     `gorm:"size:20;not null;default:'pending'"`
    CurrentStep   string         `gorm:"size:64"`
    Progress      int            `gorm:"default:0"`
    StartTime     time.Time
    EndTime       *time.Time
    CreatedAt     time.Time
    UpdatedAt     time.Time
    DeletedAt     gorm.DeletedAt `gorm:"index"`
}
```

### 2. **SagaStep 模型**
```go
type SagaStep struct {
    ID             uint           `gorm:"primarykey"`
    StepID         string         `gorm:"uniqueIndex;size:64;not null"`
    TransactionID  string         `gorm:"index;size:64;not null"`
    StepName       string         `gorm:"size:255;not null"`
    Status         StepStatus     `gorm:"size:20;not null;default:'pending'"`
    ActionData     string         `gorm:"type:json"`
    CompensateData string         `gorm:"type:json"`
    ErrorMessage   string         `gorm:"type:text"`
    RetryCount     int            `gorm:"default:0"`
    MaxRetries     int            `gorm:"default:3"`
    StartTime      time.Time
    EndTime        *time.Time
    CreatedAt      time.Time
    UpdatedAt      time.Time
    DeletedAt      gorm.DeletedAt `gorm:"index"`
}
```

### 3. **SagaEvent 模型 - 新增**
```go
type SagaEvent struct {
    ID            uint           `gorm:"primarykey"`
    TransactionID string         `gorm:"index;size:64;not null"`
    StepID        string         `gorm:"size:64"`
    EventType     EventType      `gorm:"size:50;not null"`
    EventData     string         `gorm:"type:json"`
    CreatedAt     time.Time
    DeletedAt     gorm.DeletedAt `gorm:"index"`
}
```

## 🔧 **仓库方法更新**

### 1. **新增方法**
- `CreateTransaction()` - 创建事务
- `GetTransaction()` - 获取事务
- `UpdateTransactionStatus()` - 更新事务状态
- `CreateEvent()` - 创建事件
- `LogEvent()` - 记录事件
- `ListEventsByTransaction()` - 查询事务事件
- `GetEventStatistics()` - 获取事件统计

### 2. **向后兼容方法**
为了确保现有代码不受影响，保留了原有的方法名：
- `CreateInstance()` → 内部调用 `CreateTransaction()`
- `GetInstance()` → 内部调用 `GetTransaction()`
- `UpdateInstanceStatus()` → 内部调用 `UpdateTransactionStatus()`

## 📊 **事件表的作用**

### 1. **完整的审计日志**
```go
// 记录 Saga 开始事件
sagaRepo.LogEvent(ctx, transactionID, "", models.EventTypeSagaStarted, nil)

// 记录步骤开始事件
sagaRepo.LogEvent(ctx, transactionID, stepID, models.EventTypeStepStarted, map[string]interface{}{
    "step_name": "validate_user",
})

// 记录步骤完成事件
sagaRepo.LogEvent(ctx, transactionID, stepID, models.EventTypeStepCompleted, map[string]interface{}{
    "duration": "2.5s",
    "result": "success",
})
```

### 2. **监控和调试**
- 完整的执行轨迹
- 性能分析
- 错误追踪
- 业务分析

### 3. **合规要求**
- 操作审计
- 数据追溯
- 合规报告

## 🚀 **使用示例**

### 1. **创建 Saga 事务**
```go
transaction := &models.SagaTransaction{
    TransactionID: "saga_001",
    Name:          "sync_account",
    Status:        models.SagaStatusPending,
    Progress:      0,
    StartTime:     time.Now(),
}

err := sagaRepo.CreateTransaction(ctx, transaction)
```

### 2. **创建 Saga 步骤**
```go
step := &models.SagaStep{
    StepID:        "step_001",
    TransactionID: "saga_001",
    StepName:      "validate_user",
    Status:        models.StepStatusPending,
    MaxRetries:    3,
    StartTime:     time.Now(),
}

err := sagaRepo.CreateStep(ctx, step)
```

### 3. **记录事件**
```go
err := sagaRepo.LogEvent(ctx, "saga_001", "step_001", models.EventTypeStepStarted, map[string]interface{}{
    "user_id": "123",
    "action":  "validation",
})
```

## 🔄 **迁移策略**

### 1. **数据库迁移**
```sql
-- 创建新表
CREATE TABLE saga_transactions (...);
CREATE TABLE saga_events (...);

-- 迁移旧数据（可选）
INSERT INTO saga_transactions 
SELECT id, instance_id, service_name, status, NULL, 0, created_at, NULL, created_at, updated_at, deleted_at 
FROM saga_instances;
```

### 2. **代码迁移**
- 逐步替换旧的方法调用
- 利用向后兼容方法平滑过渡
- 添加事件记录功能

### 3. **测试验证**
- 单元测试更新
- 集成测试验证
- 性能测试确认

## ✅ **与文档的一致性**

现在我们的实现与 `SAGA_DISTRIBUTED_TRANSACTION_PLAN.md` 文档完全一致：

1. ✅ **表结构一致**：`saga_transactions`, `saga_steps`, `saga_events`
2. ✅ **字段设计一致**：所有字段名称和类型匹配
3. ✅ **状态枚举一致**：状态值和命名规范匹配
4. ✅ **事件系统一致**：完整的事件记录和审计功能
5. ✅ **方法接口一致**：仓库方法覆盖文档中的所有功能

## 🎯 **总结**

通过这次更新，我们实现了：

1. **完整的事件系统**：支持完整的审计日志和监控
2. **更好的数据模型**：更清晰的字段命名和结构
3. **向后兼容性**：确保现有代码不受影响
4. **文档一致性**：实现与设计文档完全匹配
5. **扩展性**：为未来的功能扩展奠定基础

这个更新为 nancalacc 项目提供了一个完整、可靠、可扩展的 Saga 分布式事务解决方案。 