package contracts

import (
	"context"
	"nancalacc/internal/data/models"
)

// Repository Saga 数据访问接口
type Repository interface {
	// ==================== 事务相关方法 ====================

	// CreateTransaction 创建 Saga 事务
	CreateTransaction(ctx context.Context, transaction *models.SagaTransaction) error

	// GetTransaction 根据事务ID获取 Saga 事务
	GetTransaction(ctx context.Context, transactionID string) (*models.SagaTransaction, error)

	// UpdateTransactionStatus 更新 Saga 事务状态
	UpdateTransactionStatus(ctx context.Context, transactionID string, status models.SagaStatus) error

	// UpdateTransaction 更新 Saga 事务
	UpdateTransaction(ctx context.Context, transaction *models.SagaTransaction) error

	// ListTransactionsByStatus 根据状态查询 Saga 事务列表
	ListTransactionsByStatus(ctx context.Context, status models.SagaStatus, limit, offset int) ([]*models.SagaTransaction, error)

	// ==================== 步骤相关方法 ====================

	// CreateStep 创建 Saga 步骤
	CreateStep(ctx context.Context, step *models.SagaStep) error

	// GetStep 根据步骤ID获取 Saga 步骤
	GetStep(ctx context.Context, stepID string) (*models.SagaStep, error)

	// UpdateStep 更新 Saga 步骤
	UpdateStep(ctx context.Context, step *models.SagaStep) error

	// UpdateStepStatus 更新 Saga 步骤状态
	UpdateStepStatus(ctx context.Context, stepID string, status models.StepStatus) error

	// ListStepsByTransaction 根据事务ID查询步骤列表
	ListStepsByTransaction(ctx context.Context, transactionID string) ([]*models.SagaStep, error)

	// GetPendingSteps 获取待执行的步骤
	GetPendingSteps(ctx context.Context, transactionID string) ([]*models.SagaStep, error)

	// GetFailedSteps 获取失败的步骤
	GetFailedSteps(ctx context.Context, transactionID string) ([]*models.SagaStep, error)

	// IncrementRetryCount 增加重试次数
	IncrementRetryCount(ctx context.Context, stepID string) error

	// UpdateStepError 更新步骤错误信息
	UpdateStepError(ctx context.Context, stepID string, errorMsg string) error

	// ==================== 事件相关方法 ====================

	// CreateEvent 创建 Saga 事件
	CreateEvent(ctx context.Context, event *models.SagaEvent) error

	// LogEvent 记录事件（便捷方法）
	LogEvent(ctx context.Context, transactionID, stepID string, eventType models.EventType, eventData map[string]interface{}) error

	// ListEventsByTransaction 根据事务ID查询事件列表
	ListEventsByTransaction(ctx context.Context, transactionID string, limit, offset int) ([]*models.SagaEvent, error)

	// ListEventsByType 根据事件类型查询事件列表
	ListEventsByType(ctx context.Context, eventType models.EventType, limit, offset int) ([]*models.SagaEvent, error)

	// ==================== 高级查询方法 ====================

	// GetTransactionWithSteps 获取 Saga 事务及其所有步骤
	GetTransactionWithSteps(ctx context.Context, transactionID string) (*models.SagaTransaction, []*models.SagaStep, error)

	// GetTransactionWithStepsAndEvents 获取 Saga 事务及其所有步骤和事件
	GetTransactionWithStepsAndEvents(ctx context.Context, transactionID string) (*models.SagaTransaction, []*models.SagaStep, []*models.SagaEvent, error)

	// ==================== 清理和维护方法 ====================

	// CleanupExpiredTransactions 清理过期的 Saga 事务
	CleanupExpiredTransactions(ctx context.Context) error

	// CleanupExpiredEvents 清理过期的事件
	CleanupExpiredEvents(ctx context.Context) error

	// ==================== 统计和监控方法 ====================

	// GetSagaStatistics 获取 Saga 统计信息
	GetSagaStatistics(ctx context.Context) (map[string]int64, error)

	// GetEventStatistics 获取事件统计信息
	GetEventStatistics(ctx context.Context) (map[string]int64, error)

	// Ping 健康检查
	Ping(ctx context.Context) error
}
