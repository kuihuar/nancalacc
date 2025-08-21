# 业务实体 vs 输入输出结构体 (DTO)

## 概述

在软件架构中，业务实体和输入输出结构体是两个重要的概念，它们有不同的职责和特点。本文档详细解释它们的区别和在当前项目中的应用。

## 1. 业务实体 (Business Entity)

### 什么是业务实体？

业务实体是**领域模型的核心**，代表业务中的真实对象，具有：

- **唯一标识** - 每个实体都有唯一ID
- **业务行为/方法** - 包含业务逻辑和操作
- **业务规则** - 验证规则、状态转换等
- **生命周期** - 可以被创建、修改、删除

### 业务实体示例

```go
// 业务实体示例（理想情况）
type User struct {
    ID       UserID
    Account  Account
    Profile  UserProfile
    Status   UserStatus
}

// 业务实体的方法
func (u *User) Activate() error {
    if u.Status == UserStatusActive {
        return errors.New("user already active")
    }
    u.Status = UserStatusActive
    return nil
}

func (u *User) ChangePassword(newPassword string) error {
    if len(newPassword) < 8 {
        return errors.New("password too short")
    }
    u.Profile.Password = newPassword
    return nil
}

func (u *User) IsActive() bool {
    return u.Status == UserStatusActive
}
```

### 业务实体的特点

1. **有唯一标识** - 每个实体都有唯一ID
2. **包含业务行为** - 有业务方法和规则
3. **有生命周期** - 可以被创建、修改、删除
4. **包含业务逻辑** - 验证规则、状态转换等
5. **可能没有序列化标签** - 专注于业务逻辑

## 2. 输入输出结构体 (DTO - Data Transfer Object)

### 什么是 DTO？

输入输出结构体是**数据传输对象**，用于：

- **API 请求参数** - 接收客户端请求数据
- **API 响应结果** - 返回给客户端的数据
- **跨层数据传输** - 在不同层之间传递数据

### 当前项目中的 DTO 示例

#### API 请求结构体

```go
// api/account/v1/account.proto
message CreateSyncAccountRequest {
    TriggerType trigger_type = 1; // 触发类型
    SyncType sync_type = 2;       // 同步类型
    optional string task_name = 3; // 任务名称
}

message GetSyncAccountRequest {
    string task_id = 1;           // 要查询的任务ID
}
```

#### API 响应结构体

```go
// api/account/v1/account.proto
message CreateSyncAccountReply {
    string task_id = 1;           // 生成的任务ID
    google.protobuf.Timestamp create_time = 2; // 任务创建时间
}

message GetSyncAccountReply {
    enum Status {
        PENDING = 0;   // 待执行
        RUNNING = 1;   // 执行中
        SUCCESS = 2;   // 成功
        FAILED = 3;    // 失败
    }
    Status status = 1;
    int64 user_count = 2;
    int64 department_count = 3;
    int64 user_department_relation_count = 4;
    int64 actual_time = 5;
    google.protobuf.Timestamp start_time = 6;
    google.protobuf.Timestamp latest_sync_time = 7;
}
```

### DTO 的特点

1. **纯数据结构** - 只包含字段，没有业务方法
2. **用于数据传输** - 在 API 层和业务层之间传递数据
3. **包含序列化标签** - 如 JSON 标签、protobuf 标签
4. **生命周期短** - 仅在请求/响应过程中使用
5. **可能没有唯一标识** - 主要用于数据传输

## 3. 当前项目的实际情况

### 项目架构

```
API Request (DTO) → Service → Biz (Usecase) → Data (Model)
```

1. **API 层**: 使用 protobuf 生成的 DTO
2. **Service 层**: 处理 DTO，调用 Biz 层
3. **Biz 层**: 包含业务逻辑，使用数据模型
4. **Data 层**: 数据模型 + Repository 实现

### 当前项目中的"业务实体"

在当前项目中，`internal/data/models/` 下的结构体实际上是**数据模型**，而不是纯粹的业务实体：

```go
// internal/data/models/user.go
type TbLasUser struct {
    ID               uint      `gorm:"primaryKey;autoIncrement;column:id"`
    TaskID           string    `gorm:"not null;column:task_id"`
    ThirdCompanyID   string    `gorm:"not null;column:third_company_id"`
    PlatformID       string    `gorm:"not null;column:platform_id"`
    Uid              string    `gorm:"not null;column:uid"`
    Account          string    `gorm:"not null;column:account"`
    NickName         string    `gorm:"not null;column:nick_name"`
    Email            string    `gorm:"column:email"`
    Phone            string    `gorm:"column:phone"`
    // ... 更多字段，包含数据库标签
}
```

这个结构体：
- ✅ **有唯一标识** (ID, Uid)
- ❌ **没有业务方法**
- ❌ **没有业务规则**
- ✅ **有数据库映射标签**

### 数据流向示例

#### 创建同步账户的完整流程

```go
// 1. API 请求 (DTO)
CreateSyncAccountRequest {
    trigger_type: TRIGGER_MANUAL,
    sync_type: FULL,
    task_name: "20231201120000"
}

// 2. Service 层处理
func (s *AccountService) CreateSyncAccount(ctx context.Context, req *v1.CreateSyncAccountRequest) (*v1.CreateSyncAccountReply, error) {
    // 参数验证
    if req.GetTaskName() != "" && len(req.GetTaskName()) != 14 {
        return nil, status.Errorf(codes.InvalidArgument, "invalid taskname: %s", req.GetTaskName())
    }
    
    // 限流控制
    if !s.limiter.Allow("global_sync_account", 10, 20) {
        return nil, status.Errorf(codes.ResourceExhausted, "global rate limit exceeded")
    }
    
    // 调用 Biz 层
    return s.fullSyncUsecase.CreateSyncAccount(ctx, req)
}

// 3. Biz 层业务逻辑
func (uc *FullSyncUsecase) CreateSyncAccount(ctx context.Context, req *v1.CreateSyncAccountRequest) (*v1.CreateSyncAccountReply, error) {
    taskId := req.GetTaskName()
    
    // 业务逻辑处理
    companyCfg, users, depts, deptUsers, err := uc.getFullData(ctx)
    if err != nil {
        return nil, err
    }
    
    // 保存数据
    err = uc.saveFullData(ctx, companyCfg, users, depts, deptUsers, taskId)
    if err != nil {
        return nil, err
    }
    
    // 通知外部系统
    err = uc.notifyFullSync(ctx, taskId)
    if err != nil {
        return nil, err
    }
    
    return &v1.CreateSyncAccountReply{
        TaskId:     taskId,
        CreateTime: timestamppb.Now(),
    }, nil
}

// 4. Data 层数据操作
func (r *AccountRepository) SaveUsers(ctx context.Context, users []*dingtalk.DingtalkDeptUser, taskId string) (int, error) {
    // 将 DTO 转换为数据模型
    // 保存到数据库
    return len(users), nil
}
```

## 4. 对比总结

### 业务实体 vs 输入输出结构体

| 特征 | 业务实体 | 输入输出结构体 (DTO) |
|------|----------|---------------------|
| **用途** | 业务逻辑的核心对象 | 数据传输 |
| **方法** | 包含业务方法 | 纯数据结构 |
| **生命周期** | 长期存在 | 短期使用 |
| **标识** | 有唯一标识 | 可能没有 |
| **业务规则** | 包含业务规则 | 不包含 |
| **序列化标签** | 可能没有 | 有序列化标签 |
| **职责** | 业务逻辑 | 数据传输 |

### 当前项目的设计特点

1. **输入输出结构体**: `CreateSyncAccountRequest`, `CreateSyncAccountReply` 等
2. **数据模型**: `TbLasUser`, `TbLasDepartment` 等
3. **业务逻辑**: 在 `Usecase` 中实现
4. **没有纯粹的业务实体**: 业务逻辑和数据模型混合

## 5. 最佳实践建议

### 何时使用业务实体

- 当需要封装复杂的业务逻辑时
- 当对象有明确的生命周期和状态时
- 当需要确保业务规则的一致性时

### 何时使用 DTO

- 在 API 层接收和返回数据时
- 在不同层之间传递数据时
- 当需要序列化/反序列化时

### 当前项目的改进方向

1. **分离关注点**: 可以考虑将业务逻辑从数据模型中分离出来
2. **创建真正的业务实体**: 在 Biz 层定义纯粹的业务实体
3. **使用转换器**: 在数据模型和业务实体之间进行转换

```go
// 建议的改进示例
// 1. 业务实体
type User struct {
    ID       UserID
    Account  Account
    Profile  UserProfile
    Status   UserStatus
}

// 2. 数据模型
type TbLasUser struct {
    ID       uint   `gorm:"primaryKey;autoIncrement;column:id"`
    Uid      string `gorm:"not null;column:uid"`
    Account  string `gorm:"not null;column:account"`
    // ... 数据库字段
}

// 3. 转换器
func (u *User) ToModel() *TbLasUser {
    return &TbLasUser{
        Uid:     u.ID.String(),
        Account: u.Account.String(),
        // ...
    }
}

func (m *TbLasUser) ToEntity() *User {
    return &User{
        ID:      UserID(m.Uid),
        Account: Account(m.Account),
        // ...
    }
}
```

## 6. 总结

业务实体和输入输出结构体在软件架构中扮演不同的角色：

- **业务实体**专注于业务逻辑和规则
- **DTO**专注于数据传输和序列化

当前项目采用了一种实用的混合设计，既保证了功能的完整性，又避免了过度设计。这种设计在中小型项目中是有效的，但随着项目的增长，可以考虑进一步分离关注点，创建更纯粹的业务实体。 