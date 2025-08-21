package saga

import (
	"context"
	"time"

	"nancalacc/internal/data/models"
)

// Coordinator Saga 协调器接口
// 负责协调和管理整个 Saga 事务的执行流程
type Coordinator interface {
	// ==================== 事务管理 ====================

	// StartTransaction 启动一个新的 Saga 事务
	// 参数：
	//   - ctx: 上下文
	//   - name: 事务名称
	//   - steps: 步骤定义列表
	// 返回：
	//   - transactionID: 事务ID
	//   - error: 错误信息
	StartTransaction(ctx context.Context, name string, steps []StepDefinition) (string, error)

	// GetTransaction 获取事务信息
	// 参数：
	//   - ctx: 上下文
	//   - transactionID: 事务ID
	// 返回：
	//   - *TransactionInfo: 事务信息
	//   - error: 错误信息
	GetTransaction(ctx context.Context, transactionID string) (*TransactionInfo, error)

	// ListTransactions 查询事务列表
	// 参数：
	//   - ctx: 上下文
	//   - filter: 查询过滤器
	//   - limit: 限制数量
	//   - offset: 偏移量
	// 返回：
	//   - []*TransactionInfo: 事务列表
	//   - error: 错误信息
	ListTransactions(ctx context.Context, filter TransactionFilter, limit, offset int) ([]*TransactionInfo, error)

	// CancelTransaction 取消事务
	// 参数：
	//   - ctx: 上下文
	//   - transactionID: 事务ID
	//   - reason: 取消原因
	// 返回：
	//   - error: 错误信息
	CancelTransaction(ctx context.Context, transactionID string, reason string) error

	// ==================== 步骤管理 ====================

	// GetStep 获取步骤信息
	// 参数：
	//   - ctx: 上下文
	//   - transactionID: 事务ID
	//   - stepID: 步骤ID
	// 返回：
	//   - *StepInfo: 步骤信息
	//   - error: 错误信息
	GetStep(ctx context.Context, transactionID, stepID string) (*StepInfo, error)

	// ListSteps 查询步骤列表
	// 参数：
	//   - ctx: 上下文
	//   - transactionID: 事务ID
	//   - filter: 查询过滤器
	// 返回：
	//   - []*StepInfo: 步骤列表
	//   - error: 错误信息
	ListSteps(ctx context.Context, transactionID string, filter StepFilter) ([]*StepInfo, error)

	// RetryStep 重试步骤
	// 参数：
	//   - ctx: 上下文
	//   - transactionID: 事务ID
	//   - stepID: 步骤ID
	// 返回：
	//   - error: 错误信息
	RetryStep(ctx context.Context, transactionID, stepID string) error

	// SkipStep 跳过步骤
	// 参数：
	//   - ctx: 上下文
	//   - transactionID: 事务ID
	//   - stepID: 步骤ID
	//   - reason: 跳过原因
	// 返回：
	//   - error: 错误信息
	SkipStep(ctx context.Context, transactionID, stepID string, reason string) error

	// ==================== 事件管理 ====================

	// GetEvents 获取事件列表
	// 参数：
	//   - ctx: 上下文
	//   - transactionID: 事务ID
	//   - filter: 事件过滤器
	//   - limit: 限制数量
	//   - offset: 偏移量
	// 返回：
	//   - []*EventInfo: 事件列表
	//   - error: 错误信息
	GetEvents(ctx context.Context, transactionID string, filter EventFilter, limit, offset int) ([]*EventInfo, error)

	// ==================== 监控和统计 ====================

	// GetStatistics 获取统计信息
	// 参数：
	//   - ctx: 上下文
	//   - filter: 统计过滤器
	// 返回：
	//   - *Statistics: 统计信息
	//   - error: 错误信息
	GetStatistics(ctx context.Context, filter StatisticsFilter) (*Statistics, error)

	// GetHealth 获取健康状态
	// 参数：
	//   - ctx: 上下文
	// 返回：
	//   - *HealthInfo: 健康信息
	//   - error: 错误信息
	GetHealth(ctx context.Context) (*HealthInfo, error)

	// ==================== 配置管理 ====================

	// UpdateConfig 更新配置
	// 参数：
	//   - ctx: 上下文
	//   - config: 配置信息
	// 返回：
	//   - error: 错误信息
	UpdateConfig(ctx context.Context, config *Config) error

	// GetConfig 获取配置
	// 参数：
	//   - ctx: 上下文
	// 返回：
	//   - *Config: 配置信息
	//   - error: 错误信息
	GetConfig(ctx context.Context) (*Config, error)
}

// ==================== 数据结构定义 ====================

// StepDefinition 步骤定义
type StepDefinition struct {
	StepID       string                 `json:"step_id"`      // 步骤唯一标识
	StepName     string                 `json:"step_name"`    // 步骤名称
	Action       Action                 `json:"action"`       // 正向操作
	Compensation Compensation           `json:"compensation"` // 补偿操作
	MaxRetries   int                    `json:"max_retries"`  // 最大重试次数
	Timeout      time.Duration          `json:"timeout"`      // 超时时间
	Metadata     map[string]interface{} `json:"metadata"`     // 元数据
}

// Action 正向操作接口
type Action interface {
	// Execute 执行操作
	// 参数：
	//   - ctx: 上下文
	//   - data: 输入数据
	// 返回：
	//   - result: 执行结果
	//   - error: 错误信息
	Execute(ctx context.Context, data map[string]interface{}) (map[string]interface{}, error)
}

// Compensation 补偿操作接口
type Compensation interface {
	// Compensate 执行补偿
	// 参数：
	//   - ctx: 上下文
	//   - data: 补偿数据
	// 返回：
	//   - error: 错误信息
	Compensate(ctx context.Context, data map[string]interface{}) error
}

// TransactionInfo 事务信息
type TransactionInfo struct {
	TransactionID string            `json:"transaction_id"` // 事务ID
	Name          string            `json:"name"`           // 事务名称
	Status        models.SagaStatus `json:"status"`         // 状态
	CurrentStep   string            `json:"current_step"`   // 当前步骤
	Progress      int               `json:"progress"`       // 进度百分比
	StartTime     time.Time         `json:"start_time"`     // 开始时间
	EndTime       *time.Time        `json:"end_time"`       // 结束时间
	CreatedAt     time.Time         `json:"created_at"`     // 创建时间
	UpdatedAt     time.Time         `json:"updated_at"`     // 更新时间
	Steps         []*StepInfo       `json:"steps"`          // 步骤列表
}

// StepInfo 步骤信息
type StepInfo struct {
	StepID         string                 `json:"step_id"`         // 步骤ID
	StepName       string                 `json:"step_name"`       // 步骤名称
	Status         models.StepStatus      `json:"status"`          // 状态
	ActionData     map[string]interface{} `json:"action_data"`     // 操作数据
	CompensateData map[string]interface{} `json:"compensate_data"` // 补偿数据
	ErrorMessage   string                 `json:"error_message"`   // 错误信息
	RetryCount     int                    `json:"retry_count"`     // 重试次数
	MaxRetries     int                    `json:"max_retries"`     // 最大重试次数
	StartTime      time.Time              `json:"start_time"`      // 开始时间
	EndTime        *time.Time             `json:"end_time"`        // 结束时间
	CreatedAt      time.Time              `json:"created_at"`      // 创建时间
	UpdatedAt      time.Time              `json:"updated_at"`      // 更新时间
}

// EventInfo 事件信息
type EventInfo struct {
	ID            uint                   `json:"id"`             // 事件ID
	TransactionID string                 `json:"transaction_id"` // 事务ID
	StepID        string                 `json:"step_id"`        // 步骤ID
	EventType     models.EventType       `json:"event_type"`     // 事件类型
	EventData     map[string]interface{} `json:"event_data"`     // 事件数据
	CreatedAt     time.Time              `json:"created_at"`     // 创建时间
}

// ==================== 过滤器定义 ====================

// TransactionFilter 事务查询过滤器
type TransactionFilter struct {
	Status    []models.SagaStatus `json:"status"`     // 状态过滤
	StartTime *time.Time          `json:"start_time"` // 开始时间过滤
	EndTime   *time.Time          `json:"end_time"`   // 结束时间过滤
	Name      string              `json:"name"`       // 名称过滤
}

// StepFilter 步骤查询过滤器
type StepFilter struct {
	Status []models.StepStatus `json:"status"` // 状态过滤
}

// EventFilter 事件查询过滤器
type EventFilter struct {
	EventType []models.EventType `json:"event_type"` // 事件类型过滤
	StartTime *time.Time         `json:"start_time"` // 开始时间过滤
	EndTime   *time.Time         `json:"end_time"`   // 结束时间过滤
}

// StatisticsFilter 统计查询过滤器
type StatisticsFilter struct {
	StartTime *time.Time `json:"start_time"` // 开始时间
	EndTime   *time.Time `json:"end_time"`   // 结束时间
}

// ==================== 统计和健康信息 ====================

// Statistics 统计信息
type Statistics struct {
	TotalTransactions       int64                  `json:"total_transactions"`       // 总事务数
	CompletedTransactions   int64                  `json:"completed_transactions"`   // 完成事务数
	FailedTransactions      int64                  `json:"failed_transactions"`      // 失败事务数
	CompensatedTransactions int64                  `json:"compensated_transactions"` // 补偿事务数
	TotalSteps              int64                  `json:"total_steps"`              // 总步骤数
	CompletedSteps          int64                  `json:"completed_steps"`          // 完成步骤数
	FailedSteps             int64                  `json:"failed_steps"`             // 失败步骤数
	TotalEvents             int64                  `json:"total_events"`             // 总事件数
	StatusDistribution      map[string]int64       `json:"status_distribution"`      // 状态分布
	EventTypeDistribution   map[string]int64       `json:"event_type_distribution"`  // 事件类型分布
	PerformanceMetrics      map[string]float64     `json:"performance_metrics"`      // 性能指标
	CustomMetrics           map[string]interface{} `json:"custom_metrics"`           // 自定义指标
}

// HealthInfo 健康信息
type HealthInfo struct {
	Status     string                 `json:"status"`     // 健康状态
	Message    string                 `json:"message"`    // 状态消息
	Timestamp  time.Time              `json:"timestamp"`  // 检查时间
	Components map[string]string      `json:"components"` // 组件状态
	Metrics    map[string]interface{} `json:"metrics"`    // 健康指标
}

// Config 配置信息
type Config struct {
	MaxRetries               int           `json:"max_retries"`                // 默认最大重试次数
	DefaultTimeout           time.Duration `json:"default_timeout"`            // 默认超时时间
	RetryInterval            time.Duration `json:"retry_interval"`             // 重试间隔
	MaxConcurrentSteps       int           `json:"max_concurrent_steps"`       // 最大并发步骤数
	EventRetentionDays       int           `json:"event_retention_days"`       // 事件保留天数
	TransactionRetentionDays int           `json:"transaction_retention_days"` // 事务保留天数
	EnableMetrics            bool          `json:"enable_metrics"`             // 是否启用指标
	EnableTracing            bool          `json:"enable_tracing"`             // 是否启用追踪
	LogLevel                 string        `json:"log_level"`                  // 日志级别
}
