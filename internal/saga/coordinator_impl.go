package saga

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"nancalacc/internal/data/models"
	"nancalacc/internal/repository/contracts"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
)

// CoordinatorImpl Saga 协调器实现
type CoordinatorImpl struct {
	repo   contracts.Repository
	logger log.Logger
	mu     sync.RWMutex
	config *Config
}

// NewCoordinator 创建新的 Saga 协调器
func NewCoordinator(repo contracts.Repository, logger log.Logger, config *Config) Coordinator {
	if config == nil {
		config = &Config{
			MaxRetries:               3,
			DefaultTimeout:           30 * time.Minute,
			RetryInterval:            5 * time.Second,
			MaxConcurrentSteps:       10,
			EventRetentionDays:       30,
			TransactionRetentionDays: 7,
			EnableMetrics:            true,
			EnableTracing:            true,
			LogLevel:                 "info",
		}
	}

	return &CoordinatorImpl{
		repo:   repo,
		logger: logger,
		config: config,
	}
}

// ==================== 事务管理 ====================

// StartTransaction 启动一个新的 Saga 事务
func (c *CoordinatorImpl) StartTransaction(ctx context.Context, name string, steps []StepDefinition) (string, error) {
	transactionID := uuid.New().String()

	c.logger.Log(log.LevelInfo, "msg", "starting saga transaction",
		"transaction_id", transactionID, "name", name, "steps_count", len(steps))

	// 创建事务记录
	transaction := &models.SagaTransaction{
		TransactionID: transactionID,
		Name:          name,
		Status:        models.SagaStatusPending,
		Progress:      0,
		StartTime:     time.Now(),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := c.repo.CreateTransaction(ctx, transaction); err != nil {
		c.logger.Log(log.LevelError, "msg", "failed to create transaction",
			"transaction_id", transactionID, "err", err)
		return "", fmt.Errorf("failed to create transaction: %w", err)
	}

	// 创建步骤记录
	for i, stepDef := range steps {
		step := &models.SagaStep{
			StepID:        stepDef.StepID,
			TransactionID: transactionID,
			StepName:      stepDef.StepName,
			Status:        models.StepStatusPending,
			RetryCount:    0,
			MaxRetries:    stepDef.MaxRetries,
			StartTime:     time.Now(),
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		// 设置默认值
		if step.MaxRetries == 0 {
			step.MaxRetries = c.config.MaxRetries
		}

		if err := c.repo.CreateStep(ctx, step); err != nil {
			c.logger.Log(log.LevelError, "msg", "failed to create step",
				"transaction_id", transactionID, "step_id", stepDef.StepID, "err", err)
			return "", fmt.Errorf("failed to create step %s: %w", stepDef.StepID, err)
		}

		// 记录步骤创建事件
		c.repo.LogEvent(ctx, transactionID, stepDef.StepID, models.EventTypeStepStarted, map[string]interface{}{
			"step_name": stepDef.StepName,
			"order":     i + 1,
		})
	}

	// 记录事务开始事件
	c.repo.LogEvent(ctx, transactionID, "", models.EventTypeSagaStarted, map[string]interface{}{
		"name":        name,
		"steps_count": len(steps),
	})

	// 异步执行事务
	go c.executeTransaction(context.Background(), transactionID)

	c.logger.Log(log.LevelInfo, "msg", "saga transaction started",
		"transaction_id", transactionID, "name", name)

	return transactionID, nil
}

// GetTransaction 获取事务信息
func (c *CoordinatorImpl) GetTransaction(ctx context.Context, transactionID string) (*TransactionInfo, error) {
	transaction, steps, err := c.repo.GetTransactionWithSteps(ctx, transactionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	// 转换为 TransactionInfo
	transactionInfo := &TransactionInfo{
		TransactionID: transaction.TransactionID,
		Name:          transaction.Name,
		Status:        transaction.Status,
		CurrentStep:   transaction.CurrentStep,
		Progress:      transaction.Progress,
		StartTime:     transaction.StartTime,
		EndTime:       transaction.EndTime,
		CreatedAt:     transaction.CreatedAt,
		UpdatedAt:     transaction.UpdatedAt,
		Steps:         make([]*StepInfo, 0, len(steps)),
	}

	// 转换步骤信息
	for _, step := range steps {
		var actionData, compensateData map[string]interface{}

		if step.ActionData != "" {
			json.Unmarshal([]byte(step.ActionData), &actionData)
		}
		if step.CompensateData != "" {
			json.Unmarshal([]byte(step.CompensateData), &compensateData)
		}

		stepInfo := &StepInfo{
			StepID:         step.StepID,
			StepName:       step.StepName,
			Status:         step.Status,
			ActionData:     actionData,
			CompensateData: compensateData,
			ErrorMessage:   step.ErrorMessage,
			RetryCount:     step.RetryCount,
			MaxRetries:     step.MaxRetries,
			StartTime:      step.StartTime,
			EndTime:        step.EndTime,
			CreatedAt:      step.CreatedAt,
			UpdatedAt:      step.UpdatedAt,
		}
		transactionInfo.Steps = append(transactionInfo.Steps, stepInfo)
	}

	return transactionInfo, nil
}

// ListTransactions 查询事务列表
func (c *CoordinatorImpl) ListTransactions(ctx context.Context, filter TransactionFilter, limit, offset int) ([]*TransactionInfo, error) {
	// 这里需要扩展 Repository 接口来支持复杂的查询
	// 暂时返回空列表，后续可以扩展
	return []*TransactionInfo{}, nil
}

// CancelTransaction 取消事务
func (c *CoordinatorImpl) CancelTransaction(ctx context.Context, transactionID string, reason string) error {
	c.logger.Log(log.LevelInfo, "msg", "cancelling transaction",
		"transaction_id", transactionID, "reason", reason)

	// 更新事务状态
	if err := c.repo.UpdateTransactionStatus(ctx, transactionID, models.SagaStatusCancelled); err != nil {
		return fmt.Errorf("failed to cancel transaction: %w", err)
	}

	// 记录取消事件
	c.repo.LogEvent(ctx, transactionID, "", models.EventTypeSagaCancelled, map[string]interface{}{
		"reason": reason,
	})

	c.logger.Log(log.LevelInfo, "msg", "transaction cancelled", "transaction_id", transactionID)
	return nil
}

// ==================== 步骤管理 ====================

// GetStep 获取步骤信息
func (c *CoordinatorImpl) GetStep(ctx context.Context, transactionID, stepID string) (*StepInfo, error) {
	step, err := c.repo.GetStep(ctx, stepID)
	if err != nil {
		return nil, fmt.Errorf("failed to get step: %w", err)
	}

	var actionData, compensateData map[string]interface{}

	if step.ActionData != "" {
		json.Unmarshal([]byte(step.ActionData), &actionData)
	}
	if step.CompensateData != "" {
		json.Unmarshal([]byte(step.CompensateData), &compensateData)
	}

	return &StepInfo{
		StepID:         step.StepID,
		StepName:       step.StepName,
		Status:         step.Status,
		ActionData:     actionData,
		CompensateData: compensateData,
		ErrorMessage:   step.ErrorMessage,
		RetryCount:     step.RetryCount,
		MaxRetries:     step.MaxRetries,
		StartTime:      step.StartTime,
		EndTime:        step.EndTime,
		CreatedAt:      step.CreatedAt,
		UpdatedAt:      step.UpdatedAt,
	}, nil
}

// ListSteps 查询步骤列表
func (c *CoordinatorImpl) ListSteps(ctx context.Context, transactionID string, filter StepFilter) ([]*StepInfo, error) {
	steps, err := c.repo.ListStepsByTransaction(ctx, transactionID)
	if err != nil {
		return nil, fmt.Errorf("failed to list steps: %w", err)
	}

	var stepInfos []*StepInfo
	for _, step := range steps {
		var actionData, compensateData map[string]interface{}

		if step.ActionData != "" {
			json.Unmarshal([]byte(step.ActionData), &actionData)
		}
		if step.CompensateData != "" {
			json.Unmarshal([]byte(step.CompensateData), &compensateData)
		}

		stepInfo := &StepInfo{
			StepID:         step.StepID,
			StepName:       step.StepName,
			Status:         step.Status,
			ActionData:     actionData,
			CompensateData: compensateData,
			ErrorMessage:   step.ErrorMessage,
			RetryCount:     step.RetryCount,
			MaxRetries:     step.MaxRetries,
			StartTime:      step.StartTime,
			EndTime:        step.EndTime,
			CreatedAt:      step.CreatedAt,
			UpdatedAt:      step.UpdatedAt,
		}
		stepInfos = append(stepInfos, stepInfo)
	}

	return stepInfos, nil
}

// RetryStep 重试步骤
func (c *CoordinatorImpl) RetryStep(ctx context.Context, transactionID, stepID string) error {
	c.logger.Log(log.LevelInfo, "msg", "retrying step",
		"transaction_id", transactionID, "step_id", stepID)

	// 获取步骤信息
	step, err := c.repo.GetStep(ctx, stepID)
	if err != nil {
		return fmt.Errorf("failed to get step: %w", err)
	}

	// 检查是否可以重试
	if step.Status != models.StepStatusFailed {
		return fmt.Errorf("step is not in failed status, cannot retry")
	}

	if step.RetryCount >= step.MaxRetries {
		return fmt.Errorf("step has reached max retry count")
	}

	// 重置步骤状态
	step.Status = models.StepStatusPending
	step.ErrorMessage = ""
	step.UpdatedAt = time.Now()

	if err := c.repo.UpdateStep(ctx, step); err != nil {
		return fmt.Errorf("failed to update step: %w", err)
	}

	// 记录重试事件
	c.repo.LogEvent(ctx, transactionID, stepID, models.EventTypeStepRetried, map[string]interface{}{
		"retry_count": step.RetryCount,
	})

	// 重新执行步骤
	go c.executeStep(context.Background(), transactionID, stepID)

	c.logger.Log(log.LevelInfo, "msg", "step retry initiated",
		"transaction_id", transactionID, "step_id", stepID)

	return nil
}

// SkipStep 跳过步骤
func (c *CoordinatorImpl) SkipStep(ctx context.Context, transactionID, stepID string, reason string) error {
	c.logger.Log(log.LevelInfo, "msg", "skipping step",
		"transaction_id", transactionID, "step_id", stepID, "reason", reason)

	// 更新步骤状态
	if err := c.repo.UpdateStepStatus(ctx, stepID, models.StepStatusSkipped); err != nil {
		return fmt.Errorf("failed to skip step: %w", err)
	}

	// 记录跳过事件
	c.repo.LogEvent(ctx, transactionID, stepID, models.EventTypeStepSkipped, map[string]interface{}{
		"reason": reason,
	})

	c.logger.Log(log.LevelInfo, "msg", "step skipped",
		"transaction_id", transactionID, "step_id", stepID)

	return nil
}

// ==================== 事件管理 ====================

// GetEvents 获取事件列表
func (c *CoordinatorImpl) GetEvents(ctx context.Context, transactionID string, filter EventFilter, limit, offset int) ([]*EventInfo, error) {
	events, err := c.repo.ListEventsByTransaction(ctx, transactionID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get events: %w", err)
	}

	var eventInfos []*EventInfo
	for _, event := range events {
		var eventData map[string]interface{}
		if event.EventData != "" {
			json.Unmarshal([]byte(event.EventData), &eventData)
		}

		eventInfo := &EventInfo{
			ID:            event.ID,
			TransactionID: event.TransactionID,
			StepID:        event.StepID,
			EventType:     event.EventType,
			EventData:     eventData,
			CreatedAt:     event.CreatedAt,
		}
		eventInfos = append(eventInfos, eventInfo)
	}

	return eventInfos, nil
}

// ==================== 监控和统计 ====================

// GetStatistics 获取统计信息
func (c *CoordinatorImpl) GetStatistics(ctx context.Context, filter StatisticsFilter) (*Statistics, error) {
	stats, err := c.repo.GetSagaStatistics(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get statistics: %w", err)
	}

	eventStats, err := c.repo.GetEventStatistics(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get event statistics: %w", err)
	}

	// 转换统计信息
	statistics := &Statistics{
		StatusDistribution:    make(map[string]int64),
		EventTypeDistribution: make(map[string]int64),
		PerformanceMetrics:    make(map[string]float64),
		CustomMetrics:         make(map[string]interface{}),
	}

	// 填充状态分布
	for status, count := range stats {
		statistics.StatusDistribution[status] = count
	}

	// 填充事件类型分布
	for eventType, count := range eventStats {
		statistics.EventTypeDistribution[eventType] = count
	}

	return statistics, nil
}

// GetHealth 获取健康状态
func (c *CoordinatorImpl) GetHealth(ctx context.Context) (*HealthInfo, error) {
	// 检查数据库连接
	if err := c.repo.Ping(ctx); err != nil {
		return &HealthInfo{
			Status:    "unhealthy",
			Message:   "database connection failed",
			Timestamp: time.Now(),
			Components: map[string]string{
				"database": "unhealthy",
			},
		}, nil
	}

	return &HealthInfo{
		Status:    "healthy",
		Message:   "all components are healthy",
		Timestamp: time.Now(),
		Components: map[string]string{
			"database": "healthy",
		},
		Metrics: map[string]interface{}{
			"uptime": time.Since(time.Now()).String(),
		},
	}, nil
}

// ==================== 配置管理 ====================

// UpdateConfig 更新配置
func (c *CoordinatorImpl) UpdateConfig(ctx context.Context, config *Config) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.config = config
	c.logger.Log(log.LevelInfo, "msg", "config updated", "config", config)
	return nil
}

// GetConfig 获取配置
func (c *CoordinatorImpl) GetConfig(ctx context.Context) (*Config, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.config, nil
}

// ==================== 内部执行方法 ====================

// executeTransaction 执行事务
func (c *CoordinatorImpl) executeTransaction(ctx context.Context, transactionID string) {
	defer func() {
		if r := recover(); r != nil {
			c.logger.Log(log.LevelError, "msg", "transaction execution panicked",
				"transaction_id", transactionID, "panic", r)
			c.handleTransactionFailure(ctx, transactionID, fmt.Errorf("transaction execution panicked: %v", r))
		}
	}()

	c.logger.Log(log.LevelInfo, "msg", "executing transaction", "transaction_id", transactionID)

	// 更新事务状态为执行中
	if err := c.repo.UpdateTransactionStatus(ctx, transactionID, models.SagaStatusInProgress); err != nil {
		c.logger.Log(log.LevelError, "msg", "failed to update transaction status",
			"transaction_id", transactionID, "err", err)
		return
	}

	// 获取所有步骤
	steps, err := c.repo.ListStepsByTransaction(ctx, transactionID)
	if err != nil {
		c.logger.Log(log.LevelError, "msg", "failed to get steps",
			"transaction_id", transactionID, "err", err)
		c.handleTransactionFailure(ctx, transactionID, err)
		return
	}

	// 按顺序执行步骤
	for i, step := range steps {
		c.logger.Log(log.LevelInfo, "msg", "executing step",
			"transaction_id", transactionID, "step_id", step.StepID, "step_name", step.StepName)

		// 更新当前步骤
		if err := c.repo.UpdateTransaction(ctx, &models.SagaTransaction{
			TransactionID: transactionID,
			CurrentStep:   step.StepID,
			Progress:      (i + 1) * 100 / len(steps),
			UpdatedAt:     time.Now(),
		}); err != nil {
			c.logger.Log(log.LevelError, "msg", "failed to update transaction progress",
				"transaction_id", transactionID, "err", err)
		}

		// 执行步骤
		if err := c.executeStep(ctx, transactionID, step.StepID); err != nil {
			c.logger.Log(log.LevelError, "msg", "step execution failed",
				"transaction_id", transactionID, "step_id", step.StepID, "err", err)
			c.handleTransactionFailure(ctx, transactionID, err)
			return
		}
	}

	// 所有步骤执行成功，完成事务
	if err := c.repo.UpdateTransactionStatus(ctx, transactionID, models.SagaStatusCompleted); err != nil {
		c.logger.Log(log.LevelError, "msg", "failed to complete transaction",
			"transaction_id", transactionID, "err", err)
		return
	}

	// 记录完成事件
	c.repo.LogEvent(ctx, transactionID, "", models.EventTypeSagaCompleted, nil)

	c.logger.Log(log.LevelInfo, "msg", "transaction completed successfully", "transaction_id", transactionID)
}

// executeStep 执行单个步骤
func (c *CoordinatorImpl) executeStep(ctx context.Context, transactionID, stepID string) error {
	step, err := c.repo.GetStep(ctx, stepID)
	if err != nil {
		return fmt.Errorf("failed to get step: %w", err)
	}

	c.logger.Log(log.LevelInfo, "msg", "executing step",
		"transaction_id", transactionID, "step_id", stepID, "step_name", step.StepName)

	// 更新步骤状态为执行中
	if err := c.repo.UpdateStepStatus(ctx, stepID, models.StepStatusInProgress); err != nil {
		return fmt.Errorf("failed to update step status: %w", err)
	}

	// 记录步骤开始事件
	c.repo.LogEvent(ctx, transactionID, stepID, models.EventTypeStepStarted, nil)

	// 执行步骤操作（这里需要从步骤定义中获取具体的操作）
	// 由于步骤定义中包含了 Action 接口，我们需要在实际使用时传入
	// 这里暂时记录一个占位事件
	c.repo.LogEvent(ctx, transactionID, stepID, models.EventTypeStepCompleted, map[string]interface{}{
		"message": "step executed successfully (placeholder)",
	})

	// 更新步骤状态为完成
	if err := c.repo.UpdateStepStatus(ctx, stepID, models.StepStatusCompleted); err != nil {
		return fmt.Errorf("failed to update step status: %w", err)
	}

	c.logger.Log(log.LevelInfo, "msg", "step completed",
		"transaction_id", transactionID, "step_id", stepID)

	return nil
}

// handleTransactionFailure 处理事务失败
func (c *CoordinatorImpl) handleTransactionFailure(ctx context.Context, transactionID string, err error) {
	c.logger.Log(log.LevelError, "msg", "transaction failed, starting compensation",
		"transaction_id", transactionID, "err", err)

	// 更新事务状态为失败
	if updateErr := c.repo.UpdateTransactionStatus(ctx, transactionID, models.SagaStatusFailed); updateErr != nil {
		c.logger.Log(log.LevelError, "msg", "failed to update transaction status to failed",
			"transaction_id", transactionID, "err", updateErr)
	}

	// 记录失败事件
	c.repo.LogEvent(ctx, transactionID, "", models.EventTypeSagaFailed, map[string]interface{}{
		"error": err.Error(),
	})

	// 执行补偿操作
	go c.compensateTransaction(context.Background(), transactionID)
}

// compensateTransaction 执行补偿操作
func (c *CoordinatorImpl) compensateTransaction(ctx context.Context, transactionID string) {
	c.logger.Log(log.LevelInfo, "msg", "starting compensation", "transaction_id", transactionID)

	// 更新状态为补偿中
	if err := c.repo.UpdateTransactionStatus(ctx, transactionID, models.SagaStatusCompensating); err != nil {
		c.logger.Log(log.LevelError, "msg", "failed to update transaction status to compensating",
			"transaction_id", transactionID, "err", err)
		return
	}

	// 获取已完成的步骤
	steps, err := c.repo.ListStepsByTransaction(ctx, transactionID)
	if err != nil {
		c.logger.Log(log.LevelError, "msg", "failed to get steps for compensation",
			"transaction_id", transactionID, "err", err)
		return
	}

	// 按相反顺序执行补偿操作
	for i := len(steps) - 1; i >= 0; i-- {
		step := steps[i]
		if step.Status == models.StepStatusCompleted {
			if err := c.compensateStep(ctx, transactionID, step.StepID); err != nil {
				c.logger.Log(log.LevelError, "msg", "compensation failed",
					"transaction_id", transactionID, "step_id", step.StepID, "err", err)
				// 补偿失败也需要记录，但不中断其他补偿操作
			}
		}
	}

	// 更新状态为已补偿
	if err := c.repo.UpdateTransactionStatus(ctx, transactionID, models.SagaStatusCompensated); err != nil {
		c.logger.Log(log.LevelError, "msg", "failed to update transaction status to compensated",
			"transaction_id", transactionID, "err", err)
	}

	c.logger.Log(log.LevelInfo, "msg", "compensation completed", "transaction_id", transactionID)
}

// compensateStep 执行步骤补偿
func (c *CoordinatorImpl) compensateStep(ctx context.Context, transactionID, stepID string) error {
	c.logger.Log(log.LevelInfo, "msg", "compensating step",
		"transaction_id", transactionID, "step_id", stepID)

	// 更新步骤状态为补偿中
	if err := c.repo.UpdateStepStatus(ctx, stepID, models.StepStatusCompensating); err != nil {
		return fmt.Errorf("failed to update step status: %w", err)
	}

	// 记录补偿开始事件
	c.repo.LogEvent(ctx, transactionID, stepID, models.EventTypeCompensationStarted, nil)

	// 执行补偿操作（这里需要从步骤定义中获取具体的补偿操作）
	// 由于步骤定义中包含了 Compensation 接口，我们需要在实际使用时传入
	// 这里暂时记录一个占位事件
	c.repo.LogEvent(ctx, transactionID, stepID, models.EventTypeCompensationCompleted, map[string]interface{}{
		"message": "step compensated successfully (placeholder)",
	})

	// 更新步骤状态为已补偿
	if err := c.repo.UpdateStepStatus(ctx, stepID, models.StepStatusCompensated); err != nil {
		return fmt.Errorf("failed to update step status: %w", err)
	}

	c.logger.Log(log.LevelInfo, "msg", "step compensated",
		"transaction_id", transactionID, "step_id", stepID)

	return nil
}
