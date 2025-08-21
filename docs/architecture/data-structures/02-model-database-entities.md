# 数据模型 (Model) - Database Entities

## 概述
Model 是直接映射到数据库表结构的数据模型，通常使用 ORM 框架（如 GORM）进行数据库操作。

## 分类

### 2.1 用户相关模型
```go
// 用户表
type TbLasUser struct {
    ID               uint      `gorm:"primaryKey;autoIncrement;column:id"`
    TaskID           string    `gorm:"not null;column:task_id"`
    ThirdCompanyID   string    `gorm:"not null;column:third_company_id"`
    PlatformID       string    `gorm:"not null;column:platform_id"`
    Uid              string    `gorm:"not null;column:uid"`
    Account          string    `gorm:"not null;column:account"`
    NickName         string    `gorm:"not null;column:nick_name"`
    RealName         string    `gorm:"not null;column:real_name"`
    Mobile           string    `gorm:"not null;column:mobile"`
    Email            string    `gorm:"not null;column:email"`
    JobNumber        string    `gorm:"not null;column:job_number"`
    Position         string    `gorm:"not null;column:position"`
    WorkPlace        string    `gorm:"not null;column:work_place"`
    Avatar           string    `gorm:"not null;column:avatar"`
    HiredDate        time.Time `gorm:"not null;column:hired_date"`
    Status           int       `gorm:"not null;column:status"`
    CreatedAt        time.Time `gorm:"not null;column:created_at"`
    UpdatedAt        time.Time `gorm:"not null;column:updated_at"`
}
```

### 2.2 部门相关模型
```go
// 部门表
type TbLasDept struct {
    ID               uint      `gorm:"primaryKey;autoIncrement;column:id"`
    TaskID           string    `gorm:"not null;column:task_id"`
    ThirdCompanyID   string    `gorm:"not null;column:third_company_id"`
    PlatformID       string    `gorm:"not null;column:platform_id"`
    DeptId           string    `gorm:"not null;column:dept_id"`
    Name             string    `gorm:"not null;column:name"`
    ParentId         string    `gorm:"not null;column:parent_id"`
    Order            int       `gorm:"not null;column:order"`
    CreatedAt        time.Time `gorm:"not null;column:created_at"`
    UpdatedAt        time.Time `gorm:"not null;column:updated_at"`
}

// 部门用户关联表
type TbLasDeptUser struct {
    ID               uint      `gorm:"primaryKey;autoIncrement;column:id"`
    TaskID           string    `gorm:"not null;column:task_id"`
    ThirdCompanyID   string    `gorm:"not null;column:third_company_id"`
    PlatformID       string    `gorm:"not null;column:platform_id"`
    DeptId           string    `gorm:"not null;column:dept_id"`
    Uid              string    `gorm:"not null;column:uid"`
    CreatedAt        time.Time `gorm:"not null;column:created_at"`
    UpdatedAt        time.Time `gorm:"not null;column:updated_at"`
}
```

### 2.3 任务相关模型
```go
// 任务表
type TbLasTask struct {
    ID               uint      `gorm:"primaryKey;autoIncrement;column:id"`
    TaskID           string    `gorm:"uniqueIndex;not null;column:task_id"`
    ThirdCompanyID   string    `gorm:"not null;column:third_company_id"`
    PlatformID       string    `gorm:"not null;column:platform_id"`
    TaskName         string    `gorm:"not null;column:task_name"`
    TaskType         string    `gorm:"not null;column:task_type"`
    Status           string    `gorm:"not null;column:status"`
    Progress         int       `gorm:"not null;column:progress"`
    TotalCount       int       `gorm:"not null;column:total_count"`
    SuccessCount     int       `gorm:"not null;column:success_count"`
    ErrorCount       int       `gorm:"not null;column:error_count"`
    ErrorMessage     string    `gorm:"not null;column:error_message"`
    StartTime        time.Time `gorm:"not null;column:start_time"`
    EndTime          *time.Time `gorm:"column:end_time"`
    CreatedAt        time.Time `gorm:"not null;column:created_at"`
    UpdatedAt        time.Time `gorm:"not null;column:updated_at"`
}
```

### 2.4 Saga 事务模型
```go
// Saga 事务表
type SagaTransaction struct {
    ID          string                 `gorm:"primaryKey;column:id"`
    Name        string                 `gorm:"column:name"`
    Status      SagaStatus             `gorm:"column:status"`
    CurrentStep string                 `gorm:"column:current_step"`
    Progress    int                    `gorm:"column:progress"`
    StartTime   time.Time              `gorm:"column:start_time"`
    EndTime     *time.Time             `gorm:"column:end_time"`
    Metadata    map[string]interface{} `gorm:"column:metadata;type:json"`
    CreatedAt   time.Time              `gorm:"column:created_at"`
    UpdatedAt   time.Time              `gorm:"column:updated_at"`
}

// Saga 步骤表
type SagaStep struct {
    ID            string                 `gorm:"primaryKey;column:id"`
    TransactionID string                 `gorm:"column:transaction_id"`
    StepID        string                 `gorm:"column:step_id"`
    StepName      string                 `gorm:"column:step_name"`
    Status        StepStatus             `gorm:"column:status"`
    ActionData    map[string]interface{} `gorm:"column:action_data;type:json"`
    CompensateData map[string]interface{} `gorm:"column:compensate_data;type:json"`
    ErrorMessage  string                 `gorm:"column:error_message"`
    RetryCount    int                    `gorm:"column:retry_count"`
    MaxRetries    int                    `gorm:"column:max_retries"`
    StartTime     time.Time              `gorm:"column:start_time"`
    EndTime       *time.Time             `gorm:"column:end_time"`
    CreatedAt     time.Time              `gorm:"column:created_at"`
    UpdatedAt     time.Time              `gorm:"column:updated_at"`
}
```

## 特点

1. **数据库映射**: 直接对应数据库表结构
2. **ORM 支持**: 使用 GORM 标签进行映射配置
3. **时间戳**: 自动管理创建和更新时间
4. **索引支持**: 支持主键、唯一索引等
5. **关联关系**: 支持一对一、一对多、多对多关系

## 使用场景

- 数据库 CRUD 操作
- 数据持久化
- 数据查询和统计
- 数据迁移
- 数据备份恢复

## 最佳实践

1. **命名规范**: 使用统一的命名约定
2. **字段类型**: 选择合适的数据库字段类型
3. **索引优化**: 为查询频繁的字段添加索引
4. **约束设置**: 设置适当的非空、唯一等约束
5. **关联关系**: 明确定义表之间的关联关系 