package mysql

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"nancalacc/internal/data/models"

	"gorm.io/gorm"
)

// SagaRepository Saga 数据访问层
type SagaRepository struct {
	db *gorm.DB
}

// NewSagaRepository 创建 Saga 仓库实例
func NewSagaRepository(db *gorm.DB) *SagaRepository {
	return &SagaRepository{db: db}
}

// ==================== Saga 事务相关方法 ====================

// CreateTransaction 创建 Saga 事务
func (r *SagaRepository) CreateTransaction(ctx context.Context, transaction *models.SagaTransaction) error {
	return r.db.WithContext(ctx).Create(transaction).Error
}

// GetTransaction 根据事务ID获取 Saga 事务
func (r *SagaRepository) GetTransaction(ctx context.Context, transactionID string) (*models.SagaTransaction, error) {
	var transaction models.SagaTransaction
	err := r.db.WithContext(ctx).Where("transaction_id = ?", transactionID).First(&transaction).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

// UpdateTransactionStatus 更新 Saga 事务状态
func (r *SagaRepository) UpdateTransactionStatus(ctx context.Context, transactionID string, status models.SagaStatus) error {
	updates := map[string]interface{}{
		"status":     status,
		"updated_at": time.Now(),
	}

	if status == models.SagaStatusCompleted || status == models.SagaStatusFailed || status == models.SagaStatusCompensated {
		now := time.Now()
		updates["end_time"] = &now
	}

	return r.db.WithContext(ctx).Model(&models.SagaTransaction{}).
		Where("transaction_id = ?", transactionID).
		Updates(updates).Error
}

// UpdateTransaction 更新 Saga 事务
func (r *SagaRepository) UpdateTransaction(ctx context.Context, transaction *models.SagaTransaction) error {
	return r.db.WithContext(ctx).Save(transaction).Error
}

// ListTransactionsByStatus 根据状态查询 Saga 事务列表
func (r *SagaRepository) ListTransactionsByStatus(ctx context.Context, status models.SagaStatus, limit, offset int) ([]*models.SagaTransaction, error) {
	var transactions []*models.SagaTransaction
	err := r.db.WithContext(ctx).
		Where("status = ?", status).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&transactions).Error
	return transactions, err
}

// ==================== Saga 步骤相关方法 ====================

// CreateStep 创建 Saga 步骤
func (r *SagaRepository) CreateStep(ctx context.Context, step *models.SagaStep) error {
	return r.db.WithContext(ctx).Create(step).Error
}

// GetStep 根据步骤ID获取 Saga 步骤
func (r *SagaRepository) GetStep(ctx context.Context, stepID string) (*models.SagaStep, error) {
	var step models.SagaStep
	err := r.db.WithContext(ctx).Where("step_id = ?", stepID).First(&step).Error
	if err != nil {
		return nil, err
	}
	return &step, nil
}

// UpdateStep 更新 Saga 步骤
func (r *SagaRepository) UpdateStep(ctx context.Context, step *models.SagaStep) error {
	return r.db.WithContext(ctx).Save(step).Error
}

// UpdateStepStatus 更新 Saga 步骤状态
func (r *SagaRepository) UpdateStepStatus(ctx context.Context, stepID string, status models.StepStatus) error {
	updates := map[string]interface{}{
		"status":     status,
		"updated_at": time.Now(),
	}

	if status == models.StepStatusCompleted || status == models.StepStatusFailed || status == models.StepStatusCompensated {
		now := time.Now()
		updates["end_time"] = &now
	}

	return r.db.WithContext(ctx).Model(&models.SagaStep{}).
		Where("step_id = ?", stepID).
		Updates(updates).Error
}

// ListStepsByTransaction 根据事务ID查询步骤列表
func (r *SagaRepository) ListStepsByTransaction(ctx context.Context, transactionID string) ([]*models.SagaStep, error) {
	var steps []*models.SagaStep
	err := r.db.WithContext(ctx).
		Where("transaction_id = ?", transactionID).
		Order("start_time ASC").
		Find(&steps).Error
	return steps, err
}

// GetPendingSteps 获取待执行的步骤
func (r *SagaRepository) GetPendingSteps(ctx context.Context, transactionID string) ([]*models.SagaStep, error) {
	var steps []*models.SagaStep
	err := r.db.WithContext(ctx).
		Where("transaction_id = ? AND status = ?", transactionID, models.StepStatusPending).
		Order("start_time ASC").
		Find(&steps).Error
	return steps, err
}

// GetFailedSteps 获取失败的步骤
func (r *SagaRepository) GetFailedSteps(ctx context.Context, transactionID string) ([]*models.SagaStep, error) {
	var steps []*models.SagaStep
	err := r.db.WithContext(ctx).
		Where("transaction_id = ? AND status = ?", transactionID, models.StepStatusFailed).
		Order("start_time ASC").
		Find(&steps).Error
	return steps, err
}

// IncrementRetryCount 增加重试次数
func (r *SagaRepository) IncrementRetryCount(ctx context.Context, stepID string) error {
	return r.db.WithContext(ctx).Model(&models.SagaStep{}).
		Where("step_id = ?", stepID).
		UpdateColumn("retry_count", gorm.Expr("retry_count + 1")).Error
}

// UpdateStepError 更新步骤错误信息
func (r *SagaRepository) UpdateStepError(ctx context.Context, stepID string, errorMsg string) error {
	return r.db.WithContext(ctx).Model(&models.SagaStep{}).
		Where("step_id = ?", stepID).
		Update("error_message", errorMsg).Error
}

// ==================== Saga 事件相关方法 ====================

// CreateEvent 创建 Saga 事件
func (r *SagaRepository) CreateEvent(ctx context.Context, event *models.SagaEvent) error {
	return r.db.WithContext(ctx).Create(event).Error
}

// LogEvent 记录事件（便捷方法）
func (r *SagaRepository) LogEvent(ctx context.Context, transactionID, stepID string, eventType models.EventType, eventData map[string]interface{}) error {
	dataJSON, _ := json.Marshal(eventData)

	event := &models.SagaEvent{
		TransactionID: transactionID,
		StepID:        stepID,
		EventType:     eventType,
		EventData:     string(dataJSON),
		CreatedAt:     time.Now(),
	}

	return r.CreateEvent(ctx, event)
}

// ListEventsByTransaction 根据事务ID查询事件列表
func (r *SagaRepository) ListEventsByTransaction(ctx context.Context, transactionID string, limit, offset int) ([]*models.SagaEvent, error) {
	var events []*models.SagaEvent
	err := r.db.WithContext(ctx).
		Where("transaction_id = ?", transactionID).
		Order("created_at ASC").
		Limit(limit).
		Offset(offset).
		Find(&events).Error
	return events, err
}

// ListEventsByType 根据事件类型查询事件列表
func (r *SagaRepository) ListEventsByType(ctx context.Context, eventType models.EventType, limit, offset int) ([]*models.SagaEvent, error) {
	var events []*models.SagaEvent
	err := r.db.WithContext(ctx).
		Where("event_type = ?", eventType).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&events).Error
	return events, err
}

// ==================== 高级查询方法 ====================

// GetTransactionWithSteps 获取 Saga 事务及其所有步骤
func (r *SagaRepository) GetTransactionWithSteps(ctx context.Context, transactionID string) (*models.SagaTransaction, []*models.SagaStep, error) {
	// 获取事务
	transaction, err := r.GetTransaction(ctx, transactionID)
	if err != nil {
		return nil, nil, err
	}

	// 获取步骤
	steps, err := r.ListStepsByTransaction(ctx, transactionID)
	if err != nil {
		return nil, nil, err
	}

	return transaction, steps, nil
}

// GetTransactionWithStepsAndEvents 获取 Saga 事务及其所有步骤和事件
func (r *SagaRepository) GetTransactionWithStepsAndEvents(ctx context.Context, transactionID string) (*models.SagaTransaction, []*models.SagaStep, []*models.SagaEvent, error) {
	// 获取事务和步骤
	transaction, steps, err := r.GetTransactionWithSteps(ctx, transactionID)
	if err != nil {
		return nil, nil, nil, err
	}

	// 获取事件
	events, err := r.ListEventsByTransaction(ctx, transactionID, 1000, 0) // 获取所有事件
	if err != nil {
		return nil, nil, nil, err
	}

	return transaction, steps, events, nil
}

// ==================== 清理和维护方法 ====================

// CleanupExpiredTransactions 清理过期的 Saga 事务（保留7天）
func (r *SagaRepository) CleanupExpiredTransactions(ctx context.Context) error {
	expiredTime := time.Now().AddDate(0, 0, -7)
	return r.db.WithContext(ctx).
		Where("status IN (?, ?, ?) AND updated_at < ?",
			models.SagaStatusCompleted, models.SagaStatusCompensated, models.SagaStatusFailed, expiredTime).
		Delete(&models.SagaTransaction{}).Error
}

// CleanupExpiredEvents 清理过期的事件（保留30天）
func (r *SagaRepository) CleanupExpiredEvents(ctx context.Context) error {
	expiredTime := time.Now().AddDate(0, 0, -30)
	return r.db.WithContext(ctx).
		Where("created_at < ?", expiredTime).
		Delete(&models.SagaEvent{}).Error
}

// ==================== 统计和监控方法 ====================

// GetSagaStatistics 获取 Saga 统计信息
func (r *SagaRepository) GetSagaStatistics(ctx context.Context) (map[string]int64, error) {
	stats := make(map[string]int64)

	// 统计各状态的事务数量
	statuses := []models.SagaStatus{
		models.SagaStatusPending,
		models.SagaStatusInProgress,
		models.SagaStatusCompleted,
		models.SagaStatusFailed,
		models.SagaStatusCompensating,
		models.SagaStatusCompensated,
	}

	for _, status := range statuses {
		var count int64
		err := r.db.WithContext(ctx).Model(&models.SagaTransaction{}).
			Where("status = ?", status).
			Count(&count).Error
		if err != nil {
			return nil, fmt.Errorf("failed to count status %s: %w", status, err)
		}
		stats[string(status)] = count
	}

	// 统计事件数量
	var eventCount int64
	err := r.db.WithContext(ctx).Model(&models.SagaEvent{}).Count(&eventCount).Error
	if err != nil {
		return nil, fmt.Errorf("failed to count events: %w", err)
	}
	stats["total_events"] = eventCount

	return stats, nil
}

// GetEventStatistics 获取事件统计信息
func (r *SagaRepository) GetEventStatistics(ctx context.Context) (map[string]int64, error) {
	stats := make(map[string]int64)

	// 统计各类型事件数量
	eventTypes := []models.EventType{
		models.EventTypeSagaStarted,
		models.EventTypeStepStarted,
		models.EventTypeStepCompleted,
		models.EventTypeStepFailed,
		models.EventTypeCompensationStarted,
		models.EventTypeCompensationCompleted,
		models.EventTypeSagaCompleted,
		models.EventTypeSagaFailed,
	}

	for _, eventType := range eventTypes {
		var count int64
		err := r.db.WithContext(ctx).Model(&models.SagaEvent{}).
			Where("event_type = ?", eventType).
			Count(&count).Error
		if err != nil {
			return nil, fmt.Errorf("failed to count event type %s: %w", eventType, err)
		}
		stats[string(eventType)] = count
	}

	return stats, nil
}

func (r *SagaRepository) Ping(ctx context.Context) error {
	return r.db.WithContext(ctx).Exec("SELECT 1").Error
}
