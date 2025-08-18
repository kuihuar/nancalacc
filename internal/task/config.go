// internal/task/config.go
package task

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

// JobConfig 任务配置
type JobConfig struct {
	Name        string        `json:"name"`        // 任务名称
	Spec        string        `json:"spec"`        // cron 表达式
	Description string        `json:"description"` // 任务描述
	Timeout     time.Duration `json:"timeout"`     // 执行超时时间
	RetryCount  int           `json:"retry_count"` // 重试次数
	RetryDelay  time.Duration `json:"retry_delay"` // 重试间隔
	Enabled     bool          `json:"enabled"`     // 是否启用
	MaxRetries  int           `json:"max_retries"` // 最大重试次数（包括首次执行）
	Backoff     BackoffConfig `json:"backoff"`     // 退避策略配置
}

// BackoffConfig 退避策略配置
type BackoffConfig struct {
	Type      BackoffType   `json:"type"`       // 退避类型：fixed, exponential, jitter
	BaseDelay time.Duration `json:"base_delay"` // 基础延迟时间
	MaxDelay  time.Duration `json:"max_delay"`  // 最大延迟时间
	Factor    float64       `json:"factor"`     // 指数退避因子
}

// BackoffType 退避类型
type BackoffType string

const (
	BackoffFixed       BackoffType = "fixed"       // 固定延迟
	BackoffExponential BackoffType = "exponential" // 指数退避
	BackoffJitter      BackoffType = "jitter"      // 抖动退避
)

// JobResult 任务执行结果
type JobResult struct {
	JobName    string        `json:"job_name"`
	StartTime  time.Time     `json:"start_time"`
	EndTime    time.Time     `json:"end_time"`
	Duration   time.Duration `json:"duration"`
	Success    bool          `json:"success"`
	Error      error         `json:"error,omitempty"`
	RetryCount int           `json:"retry_count"`
	IsTimeout  bool          `json:"is_timeout"`
	IsRetry    bool          `json:"is_retry"`
}

// JobExecutor 任务执行器接口
type JobExecutor interface {
	Execute(ctx context.Context) error
}

// JobExecutorFunc 任务执行器函数类型
type JobExecutorFunc func(ctx context.Context) error

// Execute 实现 JobExecutor 接口
func (f JobExecutorFunc) Execute(ctx context.Context) error {
	return f(ctx)
}

// DefaultJobConfig 创建默认任务配置
func DefaultJobConfig(name, spec string) *JobConfig {
	return &JobConfig{
		Name:        name,
		Spec:        spec,
		Description: "",
		Timeout:     30 * time.Minute, // 默认30分钟超时
		RetryCount:  3,                // 默认重试3次
		RetryDelay:  5 * time.Minute,  // 默认重试间隔5分钟
		Enabled:     true,
		MaxRetries:  3,
		Backoff: BackoffConfig{
			Type:      BackoffExponential,
			BaseDelay: 1 * time.Minute,
			MaxDelay:  30 * time.Minute,
			Factor:    2.0,
		},
	}
}

// WithTimeout 设置超时时间
func (jc *JobConfig) WithTimeout(timeout time.Duration) *JobConfig {
	jc.Timeout = timeout
	return jc
}

// WithRetry 设置重试配置
func (jc *JobConfig) WithRetry(count int, delay time.Duration) *JobConfig {
	jc.RetryCount = count
	jc.RetryDelay = delay
	jc.MaxRetries = count
	return jc
}

// WithBackoff 设置退避策略
func (jc *JobConfig) WithBackoff(backoffType BackoffType, baseDelay, maxDelay time.Duration, factor float64) *JobConfig {
	jc.Backoff = BackoffConfig{
		Type:      backoffType,
		BaseDelay: baseDelay,
		MaxDelay:  maxDelay,
		Factor:    factor,
	}
	return jc
}

// WithDescription 设置任务描述
func (jc *JobConfig) WithDescription(description string) *JobConfig {
	jc.Description = description
	return jc
}

// WithEnabled 设置是否启用
func (jc *JobConfig) WithEnabled(enabled bool) *JobConfig {
	jc.Enabled = enabled
	return jc
}

// ExecuteWithRetry 带重试机制的任务执行
func (jc *JobConfig) ExecuteWithRetry(ctx context.Context, executor JobExecutor, logger log.Logger) *JobResult {
	result := &JobResult{
		JobName:   jc.Name,
		StartTime: time.Now(),
	}

	var lastErr error
	for attempt := 0; attempt <= jc.MaxRetries; attempt++ {
		// 创建带超时的上下文
		execCtx, cancel := context.WithTimeout(ctx, jc.Timeout)

		// 执行任务
		err := executor.Execute(execCtx)
		cancel()

		// 记录执行时间
		result.EndTime = time.Now()
		result.Duration = result.EndTime.Sub(result.StartTime)
		result.RetryCount = attempt

		if err == nil {
			// 任务执行成功
			result.Success = true
			logger.Log(log.LevelInfo, "msg", "job executed successfully",
				"job_name", jc.Name,
				"attempt", attempt+1,
				"duration", result.Duration)
			return result
		}

		// 任务执行失败
		lastErr = err
		result.Error = err
		result.IsRetry = attempt < jc.MaxRetries

		// 检查是否超时
		if execCtx.Err() == context.DeadlineExceeded {
			result.IsTimeout = true
			logger.Log(log.LevelError, "msg", "job execution timeout",
				"job_name", jc.Name,
				"attempt", attempt+1,
				"timeout", jc.Timeout)
		} else {
			logger.Log(log.LevelError, "msg", "job execution failed",
				"job_name", jc.Name,
				"attempt", attempt+1,
				"error", err)
		}

		// 如果还有重试机会，等待后重试
		if attempt < jc.MaxRetries {
			delay := jc.calculateDelay(attempt)
			logger.Log(log.LevelInfo, "msg", "retrying job",
				"job_name", jc.Name,
				"attempt", attempt+1,
				"next_attempt", attempt+2,
				"delay", delay)

			select {
			case <-ctx.Done():
				// 父上下文被取消
				result.Error = ctx.Err()
				return result
			case <-time.After(delay):
				// 等待重试延迟
				continue
			}
		}
	}

	// 所有重试都失败了
	logger.Log(log.LevelError, "msg", "job execution failed after all retries",
		"job_name", jc.Name,
		"total_attempts", jc.MaxRetries+1,
		"final_error", lastErr)

	return result
}

// calculateDelay 计算重试延迟时间
func (jc *JobConfig) calculateDelay(attempt int) time.Duration {
	switch jc.Backoff.Type {
	case BackoffFixed:
		return jc.RetryDelay

	case BackoffExponential:
		delay := jc.Backoff.BaseDelay
		for i := 0; i < attempt; i++ {
			delay = time.Duration(float64(delay) * jc.Backoff.Factor)
			if delay > jc.Backoff.MaxDelay {
				delay = jc.Backoff.MaxDelay
				break
			}
		}
		return delay

	case BackoffJitter:
		// 指数退避 + 随机抖动
		baseDelay := jc.Backoff.BaseDelay
		for i := 0; i < attempt; i++ {
			baseDelay = time.Duration(float64(baseDelay) * jc.Backoff.Factor)
			if baseDelay > jc.Backoff.MaxDelay {
				baseDelay = jc.Backoff.MaxDelay
				break
			}
		}
		// 添加 ±25% 的随机抖动
		jitter := time.Duration(float64(baseDelay) * 0.25)
		return baseDelay + time.Duration(float64(jitter)*(0.5-0.5*float64(time.Now().UnixNano()%2)))

	default:
		return jc.RetryDelay
	}
}

// Validate 验证任务配置
func (jc *JobConfig) Validate() error {
	if jc.Name == "" {
		return fmt.Errorf("job name cannot be empty")
	}
	if jc.Spec == "" {
		return fmt.Errorf("cron spec cannot be empty")
	}
	if jc.Timeout <= 0 {
		return fmt.Errorf("timeout must be positive")
	}
	if jc.RetryCount < 0 {
		return fmt.Errorf("retry count cannot be negative")
	}
	if jc.MaxRetries < jc.RetryCount {
		return fmt.Errorf("max retries cannot be less than retry count")
	}
	return nil
}
