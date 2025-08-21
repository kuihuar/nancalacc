# Go-Saga 分布式事务接入方案

## 📋 项目分布式事务需求分析

###  当前业务场景分析

**nancalacc** 是一个企业级账户同步系统，主要业务流程包括：

1. **钉钉组织架构数据同步**（全量/增量）
2. **WPS 系统集成** 
3. **用户和部门关系管理**
4. **任务状态跟踪和缓存**

###  现有问题识别

从代码分析来看，当前系统存在以下关键问题：

#### 1. **数据一致性问题** (高优先级)
- **全量同步流程**：涉及多个步骤（获取钉钉数据 → 保存公司配置 → 保存部门数据 → 保存用户数据 → 保存关系数据 → 通知WPS），缺乏事务保证
- **增量同步流程**：用户加入/退出部门、部门增删改等操作，需要保证数据一致性
- **第三方系统集成**：钉钉API + WPS API调用，需要保证最终一致性

#### 2. **错误处理复杂** (高优先级)
- **部分失败处理**：某个步骤失败时，难以回滚已完成的操作
- **重试机制缺失**：网络异常、第三方服务不可用时缺乏自动重试
- **补偿操作缺失**：失败后无法自动执行补偿操作

#### 3. **并发控制不足** (中优先级)
- **并发执行风险**：虽然代码中有并发保存逻辑，但缺乏协调机制
- **资源竞争**：多个同步任务同时执行时可能产生数据竞争

#### 4. **监控困难** (中优先级)
- **执行状态不透明**：无法实时跟踪分布式操作的执行状态
- **问题排查困难**：缺乏完整的执行轨迹和审计日志

### 🏗️ Saga 模式适用性分析

#### ✅ 适合 Saga 的业务场景

1. **全量同步流程**：
   ```go
   // 当前流程（需要 Saga 保护）
   func (uc *FullSyncUsecase) CreateSyncAccount(ctx context.Context, req *v1.CreateSyncAccountRequest) {
       // 1. 验证任务
       // 2. 获取钉钉数据
       // 3. 保存公司配置
       // 4. 保存部门数据
       // 5. 保存用户数据  
       // 6. 保存关系数据
       // 7. 通知WPS系统
       // 8. 更新任务状态
   }
   ```

2. **增量同步流程**：
   ```go
   // 用户加入部门流程
   func (uc *IncrementalSyncUsecase) UserAddOrg(ctx context.Context, event *clientV2.GenericOpenDingTalkEvent) {
       // 1. 获取用户数据
       // 2. 保存用户信息
       // 3. 生成用户部门关系
       // 4. 保存关系数据
       // 5. 通知WPS系统
   }
   ```

3. **第三方系统集成**：
   - 钉钉API调用（获取组织架构数据）
   - WPS API调用（通知同步完成）
   - 需要保证两个系统的数据一致性

### 📊 现有 Saga 基础设施分析

#### ✅ 已有的基础设施

1. **数据模型** (`internal/data/models/saga.go`)：
   - `SagaTransaction`：事务记录
   - `SagaStep`：步骤记录  
   - `SagaEvent`：事件日志
   - 完整的状态枚举和事件类型

2. **数据访问层** (`internal/data/saga_repository.go`)：
   - 完整的CRUD操作
   - 事件记录功能
   - 统计和监控方法
   - 清理和维护功能

3. **协调器接口** (`internal/saga/coordinator.go`)：
   - 完整的接口定义
   - 事务管理、步骤管理、事件管理
   - 监控和统计功能

#### ❌ 缺失的关键组件

1. **协调器实现**：只有接口定义，缺少具体实现
2. **业务集成层**：现有业务代码未集成Saga模式
3. **补偿操作**：缺少具体的补偿逻辑实现
4. **监控集成**：缺少与现有监控系统的集成

### 🎯 分布式事务实施建议

#### 第一阶段：完善基础设施 (1-2周)

1. **实现 Saga 协调器**：
   ```go
   // 需要实现的核心功能
   type CoordinatorImpl struct {
       repo   Repository
       logger log.Logger
   }
   
   // 核心方法
   - StartTransaction()
   - ExecuteStep() 
   - CompensateStep()
   - GetTransactionStatus()
   ```

2. **完善数据层**：
   - 确保数据库表结构已创建
   - 验证数据访问层功能
   - 添加必要的索引优化

#### 第二阶段：业务集成 (2-3周)

1. **全量同步 Saga 化**：
   ```go
   // 将现有流程改造为 Saga 步骤
   steps := []saga.StepDefinition{
       {ID: "validate_task", Action: validateTask},
       {ID: "fetch_dingtalk_data", Action: fetchDingTalkData},
       {ID: "save_company_config", Action: saveCompanyConfig},
       {ID: "save_departments", Action: saveDepartments},
       {ID: "save_users", Action: saveUsers},
       {ID: "save_relations", Action: saveRelations},
       {ID: "notify_wps", Action: notifyWPS},
       {ID: "update_task_status", Action: updateTaskStatus},
   }
   ```

2. **增量同步 Saga 化**：
   ```go
   // 用户加入部门 Saga
   steps := []saga.StepDefinition{
       {ID: "fetch_user_data", Action: fetchUserData},
       {ID: "save_user", Action: saveUser},
       {ID: "save_relations", Action: saveRelations},
       {ID: "notify_wps", Action: notifyWPS},
   }
   ```

#### 第三阶段：补偿操作实现 (1-2周)

1. **定义补偿策略**：
   ```go
   // 每个步骤的补偿操作
   type Compensation interface {
       Compensate(ctx context.Context, data map[string]interface{}) error
   }
   ```

2. **实现具体补偿逻辑**：
   - 数据回滚策略
   - 第三方系统补偿
   - 状态恢复机制

#### 第四阶段：监控和优化 (1周)

1. **集成现有监控**：
   - 与 OpenTelemetry 集成
   - 添加 Prometheus 指标
   - 完善日志记录

2. **性能优化**：
   - 并发控制优化
   - 数据库查询优化
   - 缓存策略优化

###  预期收益

1. **数据一致性保证**：通过补偿操作确保最终一致性
2. **系统可靠性提升**：支持部分失败和自动恢复
3. **用户体验改善**：提供实时的执行状态和进度跟踪
4. **运维效率提升**：完整的执行轨迹和监控指标
5. **扩展性增强**：支持水平扩展和高并发场景

###  下一步行动建议

1. **立即开始**：实现 Saga 协调器的核心功能
2. **优先级排序**：先处理全量同步，再处理增量同步
3. **渐进式改造**：保持现有功能稳定，逐步集成 Saga 模式
4. **充分测试**：每个阶段都要进行充分的测试验证

这个分析为你的分布式事务接入提供了清晰的路线图，建议按照优先级逐步实施。 