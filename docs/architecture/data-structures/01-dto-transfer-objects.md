# 数据传输对象 (DTO) - Data Transfer Objects

## 概述
DTO 是用于在不同层之间传输数据的对象，主要用于 API 接口、服务间通信等场景。

## 分类

### 1.1 Proto 生成的 DTO
从 Protocol Buffers 定义自动生成的数据传输对象。

```go
// 创建同步账户请求
type CreateSyncAccountRequest struct {
    TriggerType TriggerType `protobuf:"varint,1,opt,name=trigger_type,json=triggerType,proto3,enum=api.account.v1.TriggerType" json:"trigger_type,omitempty"`
    SyncType    SyncType    `protobuf:"varint,2,opt,name=sync_type,json=syncType,proto3,enum=api.account.v1.SyncType" json:"sync_type,omitempty"`
    TaskName    *string     `protobuf:"bytes,3,opt,name=task_name,json=taskName,proto3,oneof" json:"task_name,omitempty"`
}

// 获取同步账户请求
type GetSyncAccountRequest struct {
    TaskId string `protobuf:"bytes,1,opt,name=task_id,json=taskId,proto3" json:"task_id,omitempty"`
}

// 取消同步任务请求
type CancelSyncAccountRequest struct {
    TaskId string `protobuf:"bytes,1,opt,name=task_id,json=taskId,proto3" json:"task_id,omitempty"`
}
```

### 1.2 第三方 API 的 DTO
用于与外部系统通信的数据传输对象。

```go
// 钉钉账户同步请求
type EcisaccountsyncAllRequest struct {
    TaskId         string `json:"task_id"`
    ThirdCompanyId string `json:"third_company_id"`
    CollectCost    int    `json:"collect_cost"`
}

// 钉钉部门用户响应
type DingtalkDeptUser struct {
    Userid   string `json:"userid"`
    Name     string `json:"name"`
    Mobile   string `json:"mobile"`
    Email    string `json:"email"`
    Jobnumber string `json:"jobnumber"`
    DeptIds  []int  `json:"dept_ids"`
}

// 钉钉部门响应
type DingtalkDept struct {
    DeptId   int    `json:"dept_id"`
    Name     string `json:"name"`
    ParentId int    `json:"parent_id"`
    Order    int    `json:"order"`
}
```

### 1.3 内部服务间 DTO
用于微服务内部通信的数据传输对象。

```go
// 任务状态更新
type TaskStatusUpdate struct {
    TaskID     string    `json:"task_id"`
    Status     string    `json:"status"`
    Progress   int       `json:"progress"`
    UpdatedAt  time.Time `json:"updated_at"`
    ErrorMsg   string    `json:"error_msg,omitempty"`
}

// 同步结果
type SyncResult struct {
    Success     bool   `json:"success"`
    TotalCount  int    `json:"total_count"`
    SuccessCount int   `json:"success_count"`
    ErrorCount  int    `json:"error_count"`
    Message     string `json:"message,omitempty"`
}
```

## 特点

1. **数据封装**: 将多个相关字段封装在一个对象中
2. **版本兼容**: 支持向前和向后兼容
3. **序列化友好**: 支持 JSON、Protobuf 等序列化格式
4. **验证支持**: 可以包含数据验证规则
5. **文档化**: 通过注释和标签提供清晰的文档

## 使用场景

- API 接口的请求和响应
- 微服务间的数据交换
- 第三方系统集成
- 数据导入导出
- 缓存数据结构

## 最佳实践

1. **保持简单**: DTO 应该只包含传输所需的数据
2. **明确命名**: 使用清晰的命名约定
3. **版本控制**: 考虑 API 版本兼容性
4. **验证数据**: 在边界处进行数据验证
5. **文档化**: 提供清晰的字段说明 