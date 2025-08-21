package models

import (
	"time"

	"gorm.io/gorm"
)

// SagaTransaction Saga 事务（对应文档中的 saga_transactions 表）
type SagaTransaction struct {
	ID            uint           `gorm:"primarykey" json:"id"`
	TransactionID string         `gorm:"uniqueIndex;size:64;not null" json:"transaction_id"` // Saga 事务唯一标识
	Name          string         `gorm:"size:255;not null" json:"name"`                      // 事务名称
	Status        SagaStatus     `gorm:"size:20;not null;default:'pending'" json:"status"`   // Saga 状态
	CurrentStep   string         `gorm:"size:64" json:"current_step"`                        // 当前步骤
	Progress      int            `gorm:"default:0" json:"progress"`                          // 进度百分比
	StartTime     time.Time      `json:"start_time"`                                         // 开始时间
	EndTime       *time.Time     `json:"end_time,omitempty"`                                 // 结束时间
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

// SagaStep Saga 步骤（对应文档中的 saga_steps 表）
type SagaStep struct {
	ID             uint           `gorm:"primarykey" json:"id"`
	StepID         string         `gorm:"uniqueIndex;size:64;not null" json:"step_id"`      // 步骤唯一标识
	TransactionID  string         `gorm:"index;size:64;not null" json:"transaction_id"`     // 关联的 Saga 事务ID
	StepName       string         `gorm:"size:255;not null" json:"step_name"`               // 步骤名称
	Status         StepStatus     `gorm:"size:20;not null;default:'pending'" json:"status"` // 步骤状态
	ActionData     string         `gorm:"type:json" json:"action_data"`                     // 操作数据（JSON格式）
	CompensateData string         `gorm:"type:json" json:"compensate_data"`                 // 补偿数据（JSON格式）
	ErrorMessage   string         `gorm:"type:text" json:"error_message"`                   // 错误信息
	RetryCount     int            `gorm:"default:0" json:"retry_count"`                     // 重试次数
	MaxRetries     int            `gorm:"default:3" json:"max_retries"`                     // 最大重试次数
	StartTime      time.Time      `json:"start_time"`                                       // 开始时间
	EndTime        *time.Time     `json:"end_time,omitempty"`                               // 结束时间
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

// SagaEvent Saga 事件（对应文档中的 saga_events 表）
type SagaEvent struct {
	ID            uint           `gorm:"primarykey" json:"id"`
	TransactionID string         `gorm:"index;size:64;not null" json:"transaction_id"` // 关联的 Saga 事务ID
	StepID        string         `gorm:"size:64" json:"step_id"`                       // 关联的步骤ID（可选）
	EventType     EventType      `gorm:"size:50;not null" json:"event_type"`           // 事件类型
	EventData     string         `gorm:"type:json" json:"event_data"`                  // 事件数据（JSON格式）
	CreatedAt     time.Time      `json:"created_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

// SagaStatus Saga 状态枚举
type SagaStatus string

const (
	SagaStatusPending      SagaStatus = "pending"      // 待执行
	SagaStatusInProgress   SagaStatus = "in_progress"  // 执行中
	SagaStatusCompleted    SagaStatus = "completed"    // 已完成
	SagaStatusFailed       SagaStatus = "failed"       // 执行失败
	SagaStatusCompensating SagaStatus = "compensating" // 补偿中
	SagaStatusCompensated  SagaStatus = "compensated"  // 已补偿
	SagaStatusCancelled    SagaStatus = "cancelled"    // 已取消
)

// StepStatus 步骤状态枚举
type StepStatus string

const (
	StepStatusPending      StepStatus = "pending"      // 待执行
	StepStatusInProgress   StepStatus = "in_progress"  // 执行中
	StepStatusCompleted    StepStatus = "completed"    // 已完成
	StepStatusFailed       StepStatus = "failed"       // 执行失败
	StepStatusCompensating StepStatus = "compensating" // 补偿中
	StepStatusCompensated  StepStatus = "compensated"  // 已补偿
	StepStatusSkipped      StepStatus = "skipped"      // 已跳过
)

// EventType 事件类型枚举
type EventType string

const (
	EventTypeSagaStarted           EventType = "saga_started"           // Saga 开始
	EventTypeStepStarted           EventType = "step_started"           // 步骤开始
	EventTypeStepCompleted         EventType = "step_completed"         // 步骤完成
	EventTypeStepFailed            EventType = "step_failed"            // 步骤失败
	EventTypeStepRetried           EventType = "step_retried"           // 步骤重试
	EventTypeStepSkipped           EventType = "step_skipped"           // 步骤跳过
	EventTypeCompensationStarted   EventType = "compensation_started"   // 补偿开始
	EventTypeCompensationCompleted EventType = "compensation_completed" // 补偿完成
	EventTypeSagaCompleted         EventType = "saga_completed"         // Saga 完成
	EventTypeSagaFailed            EventType = "saga_failed"            // Saga 失败
	EventTypeSagaCancelled         EventType = "saga_cancelled"         // Saga 取消
)

// TableName 指定表名
func (SagaTransaction) TableName() string {
	return "saga_transactions"
}

func (SagaStep) TableName() string {
	return "saga_steps"
}

func (SagaEvent) TableName() string {
	return "saga_events"
}
