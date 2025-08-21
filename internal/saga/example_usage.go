package saga

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

// ExampleUsage 展示如何使用 Saga 协调器
func ExampleUsage() {
	// 1. 创建 Repository 实例（这里需要传入实际的数据库连接）
	// repo := data.NewSagaRepository(db, logger)

	// 2. 创建协调器
	// coordinator := NewCoordinator(repo, logger, nil)

	// 3. 定义业务步骤
	_ = []StepDefinition{
		{
			StepID:   "validate_task",
			StepName: "验证任务",
			Action: &ValidateTaskAction{
				TaskID: "task_123",
			},
			Compensation: &ValidateTaskCompensation{},
			MaxRetries:   3,
			Timeout:      5 * time.Minute,
		},
		{
			StepID:   "fetch_dingtalk_data",
			StepName: "获取钉钉数据",
			Action: &FetchDingTalkDataAction{
				CompanyID: "company_123",
			},
			Compensation: &FetchDingTalkDataCompensation{},
			MaxRetries:   3,
			Timeout:      10 * time.Minute,
		},
		{
			StepID:   "save_company_config",
			StepName: "保存公司配置",
			Action: &SaveCompanyConfigAction{
				Config: map[string]interface{}{
					"company_id": "company_123",
					"platform":   "dingtalk",
				},
			},
			Compensation: &SaveCompanyConfigCompensation{},
			MaxRetries:   2,
			Timeout:      3 * time.Minute,
		},
		{
			StepID:   "save_departments",
			StepName: "保存部门数据",
			Action: &SaveDepartmentsAction{
				Departments: []map[string]interface{}{
					{"id": "dept_1", "name": "技术部"},
					{"id": "dept_2", "name": "产品部"},
				},
			},
			Compensation: &SaveDepartmentsCompensation{},
			MaxRetries:   2,
			Timeout:      5 * time.Minute,
		},
		{
			StepID:   "save_users",
			StepName: "保存用户数据",
			Action: &SaveUsersAction{
				Users: []map[string]interface{}{
					{"id": "user_1", "name": "张三", "dept_id": "dept_1"},
					{"id": "user_2", "name": "李四", "dept_id": "dept_2"},
				},
			},
			Compensation: &SaveUsersCompensation{},
			MaxRetries:   2,
			Timeout:      5 * time.Minute,
		},
		{
			StepID:   "save_relations",
			StepName: "保存用户部门关系",
			Action: &SaveRelationsAction{
				Relations: []map[string]interface{}{
					{"user_id": "user_1", "dept_id": "dept_1"},
					{"user_id": "user_2", "dept_id": "dept_2"},
				},
			},
			Compensation: &SaveRelationsCompensation{},
			MaxRetries:   2,
			Timeout:      3 * time.Minute,
		},
		{
			StepID:   "notify_wps",
			StepName: "通知 WPS 系统",
			Action: &NotifyWPSAction{
				TaskID: "task_123",
			},
			Compensation: &NotifyWPSCompensation{},
			MaxRetries:   3,
			Timeout:      2 * time.Minute,
		},
		{
			StepID:   "update_task_status",
			StepName: "更新任务状态",
			Action: &UpdateTaskStatusAction{
				TaskID: "task_123",
				Status: "completed",
			},
			Compensation: &UpdateTaskStatusCompensation{},
			MaxRetries:   1,
			Timeout:      1 * time.Minute,
		},
	}

	// 4. 启动 Saga 事务
	// transactionID, err := coordinator.StartTransaction(context.Background(), "sync_account_task_123", steps)
	// if err != nil {
	//     log.Errorf("Failed to start saga transaction: %v", err)
	//     return
	// }

	// 5. 查询事务状态
	// status, err := coordinator.GetTransaction(context.Background(), transactionID)
	// if err != nil {
	//     log.Errorf("Failed to get transaction status: %v", err)
	//     return
	// }

	// log.Infof("Transaction status: %s, Progress: %d%%", status.Status, status.Progress)
}

// ==================== 具体的 Action 实现示例 ====================

// ValidateTaskAction 验证任务 Action
type ValidateTaskAction struct {
	TaskID string
}

func (a *ValidateTaskAction) Execute(ctx context.Context, data map[string]interface{}) (map[string]interface{}, error) {
	// 实现任务验证逻辑
	fmt.Printf("Validating task: %s\n", a.TaskID)

	// 模拟验证成功
	return map[string]interface{}{
		"task_id": a.TaskID,
		"status":  "validated",
	}, nil
}

// ValidateTaskCompensation 验证任务补偿
type ValidateTaskCompensation struct{}

func (c *ValidateTaskCompensation) Compensate(ctx context.Context, data map[string]interface{}) error {
	// 实现任务验证的补偿逻辑
	fmt.Println("Compensating task validation")
	return nil
}

// FetchDingTalkDataAction 获取钉钉数据 Action
type FetchDingTalkDataAction struct {
	CompanyID string
}

func (a *FetchDingTalkDataAction) Execute(ctx context.Context, data map[string]interface{}) (map[string]interface{}, error) {
	// 实现获取钉钉数据的逻辑
	fmt.Printf("Fetching DingTalk data for company: %s\n", a.CompanyID)

	// 模拟获取数据成功
	return map[string]interface{}{
		"company_id": a.CompanyID,
		"departments": []map[string]interface{}{
			{"id": "dept_1", "name": "技术部"},
			{"id": "dept_2", "name": "产品部"},
		},
		"users": []map[string]interface{}{
			{"id": "user_1", "name": "张三", "dept_id": "dept_1"},
			{"id": "user_2", "name": "李四", "dept_id": "dept_2"},
		},
	}, nil
}

// FetchDingTalkDataCompensation 获取钉钉数据补偿
type FetchDingTalkDataCompensation struct{}

func (c *FetchDingTalkDataCompensation) Compensate(ctx context.Context, data map[string]interface{}) error {
	// 实现获取钉钉数据的补偿逻辑
	fmt.Println("Compensating DingTalk data fetch")
	return nil
}

// SaveCompanyConfigAction 保存公司配置 Action
type SaveCompanyConfigAction struct {
	Config map[string]interface{}
}

func (a *SaveCompanyConfigAction) Execute(ctx context.Context, data map[string]interface{}) (map[string]interface{}, error) {
	// 实现保存公司配置的逻辑
	fmt.Printf("Saving company config: %+v\n", a.Config)

	// 模拟保存成功
	return map[string]interface{}{
		"config_id": "config_123",
		"status":    "saved",
	}, nil
}

// SaveCompanyConfigCompensation 保存公司配置补偿
type SaveCompanyConfigCompensation struct{}

func (c *SaveCompanyConfigCompensation) Compensate(ctx context.Context, data map[string]interface{}) error {
	// 实现保存公司配置的补偿逻辑
	fmt.Println("Compensating company config save")
	return nil
}

// SaveDepartmentsAction 保存部门数据 Action
type SaveDepartmentsAction struct {
	Departments []map[string]interface{}
}

func (a *SaveDepartmentsAction) Execute(ctx context.Context, data map[string]interface{}) (map[string]interface{}, error) {
	// 实现保存部门数据的逻辑
	fmt.Printf("Saving departments: %+v\n", a.Departments)

	// 模拟保存成功
	return map[string]interface{}{
		"saved_count": len(a.Departments),
		"status":      "saved",
	}, nil
}

// SaveDepartmentsCompensation 保存部门数据补偿
type SaveDepartmentsCompensation struct{}

func (c *SaveDepartmentsCompensation) Compensate(ctx context.Context, data map[string]interface{}) error {
	// 实现保存部门数据的补偿逻辑
	fmt.Println("Compensating departments save")
	return nil
}

// SaveUsersAction 保存用户数据 Action
type SaveUsersAction struct {
	Users []map[string]interface{}
}

func (a *SaveUsersAction) Execute(ctx context.Context, data map[string]interface{}) (map[string]interface{}, error) {
	// 实现保存用户数据的逻辑
	fmt.Printf("Saving users: %+v\n", a.Users)

	// 模拟保存成功
	return map[string]interface{}{
		"saved_count": len(a.Users),
		"status":      "saved",
	}, nil
}

// SaveUsersCompensation 保存用户数据补偿
type SaveUsersCompensation struct{}

func (c *SaveUsersCompensation) Compensate(ctx context.Context, data map[string]interface{}) error {
	// 实现保存用户数据的补偿逻辑
	fmt.Println("Compensating users save")
	return nil
}

// SaveRelationsAction 保存用户部门关系 Action
type SaveRelationsAction struct {
	Relations []map[string]interface{}
}

func (a *SaveRelationsAction) Execute(ctx context.Context, data map[string]interface{}) (map[string]interface{}, error) {
	// 实现保存用户部门关系的逻辑
	fmt.Printf("Saving relations: %+v\n", a.Relations)

	// 模拟保存成功
	return map[string]interface{}{
		"saved_count": len(a.Relations),
		"status":      "saved",
	}, nil
}

// SaveRelationsCompensation 保存用户部门关系补偿
type SaveRelationsCompensation struct{}

func (c *SaveRelationsCompensation) Compensate(ctx context.Context, data map[string]interface{}) error {
	// 实现保存用户部门关系的补偿逻辑
	fmt.Println("Compensating relations save")
	return nil
}

// NotifyWPSAction 通知 WPS 系统 Action
type NotifyWPSAction struct {
	TaskID string
}

func (a *NotifyWPSAction) Execute(ctx context.Context, data map[string]interface{}) (map[string]interface{}, error) {
	// 实现通知 WPS 系统的逻辑
	fmt.Printf("Notifying WPS system for task: %s\n", a.TaskID)

	// 模拟通知成功
	return map[string]interface{}{
		"task_id": a.TaskID,
		"status":  "notified",
	}, nil
}

// NotifyWPSCompensation 通知 WPS 系统补偿
type NotifyWPSCompensation struct{}

func (c *NotifyWPSCompensation) Compensate(ctx context.Context, data map[string]interface{}) error {
	// 实现通知 WPS 系统的补偿逻辑
	fmt.Println("Compensating WPS notification")
	return nil
}

// UpdateTaskStatusAction 更新任务状态 Action
type UpdateTaskStatusAction struct {
	TaskID string
	Status string
}

func (a *UpdateTaskStatusAction) Execute(ctx context.Context, data map[string]interface{}) (map[string]interface{}, error) {
	// 实现更新任务状态的逻辑
	fmt.Printf("Updating task status: %s -> %s\n", a.TaskID, a.Status)

	// 模拟更新成功
	return map[string]interface{}{
		"task_id": a.TaskID,
		"status":  a.Status,
	}, nil
}

// UpdateTaskStatusCompensation 更新任务状态补偿
type UpdateTaskStatusCompensation struct{}

func (c *UpdateTaskStatusCompensation) Compensate(ctx context.Context, data map[string]interface{}) error {
	// 实现更新任务状态的补偿逻辑
	fmt.Println("Compensating task status update")
	return nil
}

// ==================== 业务集成示例 ====================

// SyncAccountService 账户同步服务示例
type SyncAccountService struct {
	coordinator Coordinator
	logger      log.Logger
}

// NewSyncAccountService 创建账户同步服务
func NewSyncAccountService(coordinator Coordinator, logger log.Logger) *SyncAccountService {
	return &SyncAccountService{
		coordinator: coordinator,
		logger:      logger,
	}
}

// SyncAccount 同步账户（使用 Saga 模式）
func (s *SyncAccountService) SyncAccount(ctx context.Context, taskID, companyID string) (string, error) {
	s.logger.Log(log.LevelInfo, "msg", "starting account sync", "task_id", taskID, "company_id", companyID)

	// 定义 Saga 步骤
	steps := []StepDefinition{
		{
			StepID:   "validate_task",
			StepName: "验证任务",
			Action: &ValidateTaskAction{
				TaskID: taskID,
			},
			Compensation: &ValidateTaskCompensation{},
			MaxRetries:   3,
			Timeout:      5 * time.Minute,
		},
		{
			StepID:   "fetch_dingtalk_data",
			StepName: "获取钉钉数据",
			Action: &FetchDingTalkDataAction{
				CompanyID: companyID,
			},
			Compensation: &FetchDingTalkDataCompensation{},
			MaxRetries:   3,
			Timeout:      10 * time.Minute,
		},
		{
			StepID:   "save_company_config",
			StepName: "保存公司配置",
			Action: &SaveCompanyConfigAction{
				Config: map[string]interface{}{
					"company_id": companyID,
					"platform":   "dingtalk",
				},
			},
			Compensation: &SaveCompanyConfigCompensation{},
			MaxRetries:   2,
			Timeout:      3 * time.Minute,
		},
		{
			StepID:   "save_departments",
			StepName: "保存部门数据",
			Action: &SaveDepartmentsAction{
				Departments: []map[string]interface{}{
					{"id": "dept_1", "name": "技术部"},
					{"id": "dept_2", "name": "产品部"},
				},
			},
			Compensation: &SaveDepartmentsCompensation{},
			MaxRetries:   2,
			Timeout:      5 * time.Minute,
		},
		{
			StepID:   "save_users",
			StepName: "保存用户数据",
			Action: &SaveUsersAction{
				Users: []map[string]interface{}{
					{"id": "user_1", "name": "张三", "dept_id": "dept_1"},
					{"id": "user_2", "name": "李四", "dept_id": "dept_2"},
				},
			},
			Compensation: &SaveUsersCompensation{},
			MaxRetries:   2,
			Timeout:      5 * time.Minute,
		},
		{
			StepID:   "save_relations",
			StepName: "保存用户部门关系",
			Action: &SaveRelationsAction{
				Relations: []map[string]interface{}{
					{"user_id": "user_1", "dept_id": "dept_1"},
					{"user_id": "user_2", "dept_id": "dept_2"},
				},
			},
			Compensation: &SaveRelationsCompensation{},
			MaxRetries:   2,
			Timeout:      3 * time.Minute,
		},
		{
			StepID:   "notify_wps",
			StepName: "通知 WPS 系统",
			Action: &NotifyWPSAction{
				TaskID: taskID,
			},
			Compensation: &NotifyWPSCompensation{},
			MaxRetries:   3,
			Timeout:      2 * time.Minute,
		},
		{
			StepID:   "update_task_status",
			StepName: "更新任务状态",
			Action: &UpdateTaskStatusAction{
				TaskID: taskID,
				Status: "completed",
			},
			Compensation: &UpdateTaskStatusCompensation{},
			MaxRetries:   1,
			Timeout:      1 * time.Minute,
		},
	}

	// 启动 Saga 事务
	transactionID, err := s.coordinator.StartTransaction(ctx, fmt.Sprintf("sync_account_%s", taskID), steps)
	if err != nil {
		s.logger.Log(log.LevelError, "msg", "failed to start saga transaction", "err", err)
		return "", err
	}

	s.logger.Log(log.LevelInfo, "msg", "saga transaction started", "transaction_id", transactionID)
	return transactionID, nil
}

// GetSyncStatus 获取同步状态
func (s *SyncAccountService) GetSyncStatus(ctx context.Context, transactionID string) (*TransactionInfo, error) {
	return s.coordinator.GetTransaction(ctx, transactionID)
}
