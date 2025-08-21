package saga

import (
	"context"
	"time"
)

type SagaCoordinator interface {
	StartSaga(ctx context.Context, sagaID string, steps []SagaStep) error

	ExecuteStep(ctx context.Context, sagaID string, stepID string) error

	CompenstateStep(ctx context.Context, sagaID string, stepID string) error

	GetSagaStatus(ctx context.Context, sagaID string) (*SagaStatus, error)
}

type RetryPolicy struct {
	MaxAttempts int           `json:"max_attempts"`
	Backoff     time.Duration `json:"backoff"`
}

type StepStatusType string

const (
	StepStatusTypePending   StepStatusType = "pending"
	StepStatusTypeRunning   StepStatusType = "running"
	StepStatusTypeCompleted StepStatusType = "completed"
	StepStatusTypeFailed    StepStatusType = "failed"
)

type SagaError struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

type SagaStatusType string

type StepStatus struct {
	Status      StepStatusType `json:"status"`
	Attempts    int            `json:"attempts"`
	LastAttempt time.Time      `json:"last_attempt"`
	Error       *SagaError     `json:"error,omitempty"`
}

const (
	SagaStatusTypePending   SagaStatusType = "pending"
	SagaStatusTypeRunning   SagaStatusType = "running"
	SagaStatusTypeCompleted SagaStatusType = "completed"
	SagaStatusTypeFailed    SagaStatusType = "failed"
)

type Steps map[string]SagaStep

type SagaStep struct {
	ID          string                          `json:"id"`
	Name        string                          `json:"name"`
	Action      func(ctx context.Context) error `json:"-"`
	Compensate  func(ctx context.Context) error `json:"-"`
	RetryPolicy *RetryPolicy                    `json:"retry_policy"`
	Timeout     time.Duration                   `json:"timeout"`
	DependsOn   []string                        `json:"depends_on"`
	Metadata    map[string]interface{}          `json:"metadata"`
}
type SagaStatus struct {
	SagaID      string                `json:"saga_id"`
	Status      SagaStatusType        `json:"status"`
	CurrentStep string                `json:"current_step"`
	Progress    int                   `json:"progress"`
	StartTime   time.Time             `json:"start_time"`
	EndTime     *time.Time            `json:"end_time,omitempty"`
	Steps       map[string]StepStatus `json:"steps"`
	Error       *SagaError            `json:"error,omitempty"`
}
