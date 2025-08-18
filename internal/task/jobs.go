// internal/task/example_with_config.go
package task

import (
	"context"
	"fmt"
	v1 "nancalacc/api/account/v1"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/robfig/cron/v3"
)

// ExampleJobWithConfig 展示如何添加包含配置的任务
func RegisterJobsConfig(s *CronService) {
	// 示例1: 添加一个简单的配置化任务
	addSimpleConfigJob(s)

	// 示例2: 添加一个带重试和超时的任务
	//addRetryConfigJob(s)

	// 示例3: 添加一个带退避策略的任务
	//addBackoffConfigJob(s)

	// 示例4: 添加一个动态配置的任务
	//addDynamicConfigJob(s)

	// 添加 CreateSyncAccount 任务
	addCreateSyncAccountJob(s)

	// 添加数据库检测任务
	addDatabaseCheckJob(s)
}

// addSimpleConfigJob 添加简单配置任务
func addSimpleConfigJob(s *CronService) {
	// 创建任务配置
	config := DefaultJobConfig("simple_task", "0 */15 * * * *") // 每15分钟执行
	config.WithTimeout(5 * time.Minute)
	config.Description = "简单的配置化任务示例"

	// 注册任务
	_, err := s.AddJobWithConfig(config, JobExecutorFunc(func(ctx context.Context) error {
		s.log.Log(log.LevelInfo, "msg", "executing simple config job",
			"job_name", config.Name,
			"timeout", config.Timeout)

		// 模拟任务执行
		time.Sleep(2 * time.Second)
		s.log.Log(log.LevelInfo, "msg", "simple config job completed")
		return nil
	}))

	if err != nil {
		s.log.Log(log.LevelError, "msg", "failed to register simple config job", "err", err)
	}
}

// addRetryConfigJob 添加带重试配置的任务
// func addRetryConfigJob(s *CronService) {
// 	// 创建任务配置
// 	config := DefaultJobConfig("retry_task", "0 */30 * * * *") // 每30分钟执行
// 	config.WithRetry(3, 2*time.Minute)                         // 重试3次，间隔2分钟
// 	config.WithTimeout(10 * time.Minute)
// 	config.Description = "带重试机制的任务示例"

// 	// 注册任务
// 	_, err := s.AddJobWithConfig(config, JobExecutorFunc(func(ctx context.Context) error {
// 		s.log.Log(log.LevelInfo, "msg", "executing retry config job",
// 			"job_name", config.Name,
// 			"retry_count", config.RetryCount,
// 			"retry_delay", config.RetryDelay)

// 		// 模拟可能失败的任务
// 		if time.Now().Second()%2 == 0 { // 偶数秒时模拟失败
// 			return fmt.Errorf("simulated failure for retry testing")
// 		}

// 		s.log.Log(log.LevelInfo, "msg", "retry config job completed successfully")
// 		return nil
// 	}))

// 	if err != nil {
// 		s.log.Log(log.LevelError, "msg", "failed to register retry config job", "err", err)
// 	}
// }

// addBackoffConfigJob 添加带退避策略的任务
// func addBackoffConfigJob(s *CronService) {
// 	// 创建任务配置
// 	config := DefaultJobConfig("backoff_task", "0 0 */2 * * *") // 每2小时执行
// 	config.Backoff = BackoffConfig{
// 		Type:      BackoffExponential,
// 		BaseDelay: 30 * time.Second,
// 		MaxDelay:  5 * time.Minute,
// 		Factor:    2.0,
// 	}
// 	config.WithRetry(5, 0) // 重试5次，使用退避策略
// 	config.WithTimeout(15 * time.Minute)
// 	config.Description = "带指数退避策略的任务示例"

// 	// 注册任务
// 	_, err := s.AddJobWithConfig(config, JobExecutorFunc(func(ctx context.Context) error {
// 		s.log.Log(log.LevelInfo, "msg", "executing backoff config job",
// 			"job_name", config.Name,
// 			"backoff_type", config.Backoff.Type,
// 			"base_delay", config.Backoff.BaseDelay)

// 		// 模拟网络请求任务
// 		if time.Now().Minute()%10 == 0 { // 每10分钟模拟一次失败
// 			return fmt.Errorf("network timeout, will retry with backoff")
// 		}

// 		s.log.Log(log.LevelInfo, "msg", "backoff config job completed successfully")
// 		return nil
// 	}))

// 	if err != nil {
// 		s.log.Log(log.LevelError, "msg", "failed to register backoff config job", "err", err)
// 	}
// }

// addDynamicConfigJob 添加动态配置任务
// func addDynamicConfigJob(s *CronService) {
// 	// 创建动态配置
// 	config := DefaultJobConfig("dynamic_task", "0 */5 * * * *") // 每5分钟执行
// 	config.Description = "动态配置任务示例"

// 	// 注册任务
// 	_, err := s.AddJobWithConfig(config, JobExecutorFunc(func(ctx context.Context) error {
// 		// 动态调整配置
// 		dynamicConfig := getDynamicConfig()

// 		s.log.Log(log.LevelInfo, "msg", "executing dynamic config job",
// 			"job_name", config.Name,
// 			"dynamic_timeout", dynamicConfig.Timeout,
// 			"dynamic_retry_count", dynamicConfig.RetryCount)

// 		// 使用动态配置执行任务
// 		return executeWithDynamicConfig(ctx, dynamicConfig)
// 	}))

// 	if err != nil {
// 		s.log.Log(log.LevelError, "msg", "failed to register dynamic config job", "err", err)
// 	}
// }

// getDynamicConfig 获取动态配置
// func getDynamicConfig() *JobConfig {
// 	// 根据当前时间或其他条件动态调整配置
// 	now := time.Now()

// 	config := DefaultJobConfig("dynamic", "")

// 	// 根据时间调整超时时间
// 	if now.Hour() >= 22 || now.Hour() <= 6 {
// 		// 夜间时段，增加超时时间
// 		config.Timeout = 20 * time.Minute
// 		config.RetryCount = 5
// 	} else {
// 		// 白天时段，正常配置
// 		config.Timeout = 10 * time.Minute
// 		config.RetryCount = 3
// 	}

// 	// 根据系统负载调整配置
// 	// if isHighLoad() {
// 	// 	config.RetryDelay = 10 * time.Minute
// 	// } else {
// 	// 	config.RetryDelay = 2 * time.Minute
// 	// }

// 	return config
// }

// executeWithDynamicConfig 使用动态配置执行任务
// func executeWithDynamicConfig(ctx context.Context, config *JobConfig) error {
// 	// 模拟任务执行
// 	time.Sleep(3 * time.Second)

// 	// 模拟偶尔的失败
// 	if time.Now().Second()%7 == 0 {
// 		return fmt.Errorf("dynamic task failed, will retry with config: timeout=%v, retry_count=%d",
// 			config.Timeout, config.RetryCount)
// 	}

// 	return nil
// }

// isHighLoad 检查系统是否高负载
// func isHighLoad() bool {
// 	// 这里可以实现实际的负载检测逻辑
// 	// 例如：检查CPU使用率、内存使用率、队列长度等
// 	return time.Now().Minute()%15 == 0 // 每15分钟模拟一次高负载
// }

// AddJobWithConfig 添加带配置的任务（需要在CronService中实现）
func (s *CronService) AddJobWithConfig(config *JobConfig, executor JobExecutor) (cron.EntryID, error) {
	// 检查任务是否已注册
	s.mu.RLock()
	if existingID, exists := s.registeredJobs[config.Name]; exists {
		s.mu.RUnlock()
		s.log.Log(log.LevelWarn, "msg", "job already registered", "name", config.Name, "existing_id", existingID)
		return existingID, nil
	}
	s.mu.RUnlock()

	// 创建带配置的任务执行器
	configExecutor := &ConfigJobExecutor{
		config:   config,
		executor: executor,
		log:      s.log,
	}

	entryID, err := s.cron.AddFunc(config.Spec, func() {
		ctx := context.Background()
		if err := configExecutor.Execute(ctx); err != nil {
			s.log.Log(log.LevelError, "msg", "job execution failed",
				"name", config.Name, "err", err)
		} else {
			s.log.Log(log.LevelInfo, "msg", "job execution completed", "name", config.Name)
		}
	})

	if err != nil {
		return 0, fmt.Errorf("failed to add job [%s]: %w", config.Name, err)
	}

	// 注册成功后，记录任务名称
	s.mu.Lock()
	s.registeredJobs[config.Name] = entryID
	s.mu.Unlock()

	s.log.Log(log.LevelInfo, "msg", "job registered with config",
		"name", config.Name,
		"spec", config.Spec,
		"timeout", config.Timeout,
		"retry_count", config.RetryCount)

	return entryID, nil
}

// ConfigJobExecutor 带配置的任务执行器
type ConfigJobExecutor struct {
	config   *JobConfig
	executor JobExecutor
	log      log.Logger
}

// Execute 执行带配置的任务
func (cje *ConfigJobExecutor) Execute(ctx context.Context) error {
	startTime := time.Now()
	retryCount := 0

	for attempt := 0; attempt <= cje.config.MaxRetries; attempt++ {
		// 创建带超时的上下文
		timeoutCtx, cancel := context.WithTimeout(ctx, cje.config.Timeout)

		// 执行任务
		err := cje.executor.Execute(timeoutCtx)
		cancel()

		// 记录执行结果
		result := &JobResult{
			JobName:    cje.config.Name,
			StartTime:  startTime,
			EndTime:    time.Now(),
			Duration:   time.Since(startTime),
			Success:    err == nil,
			Error:      err,
			RetryCount: retryCount,
			IsTimeout:  timeoutCtx.Err() == context.DeadlineExceeded,
			IsRetry:    attempt > 0,
		}

		// 记录结果
		cje.logResult(result)

		// 如果成功，直接返回
		if err == nil {
			return nil
		}

		// 如果是最后一次尝试，返回错误
		if attempt == cje.config.MaxRetries {
			return fmt.Errorf("job [%s] failed after %d attempts: %w",
				cje.config.Name, cje.config.MaxRetries, err)
		}

		// 计算重试延迟
		delay := cje.calculateRetryDelay(attempt)

		cje.log.Log(log.LevelWarn, "msg", "job failed, will retry",
			"name", cje.config.Name,
			"attempt", attempt+1,
			"max_retries", cje.config.MaxRetries,
			"delay", delay,
			"error", err)

		// 等待重试
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
			retryCount++
		}
	}

	return nil
}

// calculateRetryDelay 计算重试延迟
func (cje *ConfigJobExecutor) calculateRetryDelay(attempt int) time.Duration {
	switch cje.config.Backoff.Type {
	case BackoffFixed:
		return cje.config.RetryDelay
	case BackoffExponential:
		delay := cje.config.Backoff.BaseDelay
		for i := 0; i < attempt; i++ {
			delay = time.Duration(float64(delay) * cje.config.Backoff.Factor)
			if delay > cje.config.Backoff.MaxDelay {
				delay = cje.config.Backoff.MaxDelay
				break
			}
		}
		return delay
	case BackoffJitter:
		// 实现抖动退避策略 - 基于指数退避添加随机抖动
		baseDelay := cje.config.Backoff.BaseDelay
		for i := 0; i < attempt; i++ {
			baseDelay = time.Duration(float64(baseDelay) * cje.config.Backoff.Factor)
			if baseDelay > cje.config.Backoff.MaxDelay {
				baseDelay = cje.config.Backoff.MaxDelay
				break
			}
		}
		// 添加10%的随机抖动
		jitter := time.Duration(float64(baseDelay) * 0.1)
		return baseDelay + jitter
	default:
		return cje.config.RetryDelay
	}
}

// logResult 记录执行结果
func (cje *ConfigJobExecutor) logResult(result *JobResult) {
	if result.Success {
		cje.log.Log(log.LevelInfo, "msg", "job completed successfully",
			"name", result.JobName,
			"duration", result.Duration,
			"retry_count", result.RetryCount)
	} else {
		cje.log.Log(log.LevelError, "msg", "job failed",
			"name", result.JobName,
			"duration", result.Duration,
			"retry_count", result.RetryCount,
			"is_timeout", result.IsTimeout,
			"is_retry", result.IsRetry,
			"error", result.Error)
	}
}

// addCreateSyncAccountJob 添加 CreateSyncAccount 任务
func addCreateSyncAccountJob(s *CronService) {
	// 创建任务配置
	config := DefaultJobConfig("sync_account_job", "0 */3 * * * *") // 每3分钟执行
	config.WithTimeout(50 * time.Minute)
	config.WithRetry(2, 1*time.Minute) // 重试2次，间隔1分钟
	config.Description = "定时同步账户任务"

	// 注册任务
	_, err := s.AddJobWithConfig(config, JobExecutorFunc(func(ctx context.Context) error {
		s.log.Log(log.LevelInfo, "msg", "starting scheduled sync account job, every 3 minutes")

		// 创建同步请求
		req := &v1.CreateSyncAccountRequest{
			TriggerType: v1.TriggerType_TRIGGER_SCHEDULED,
			SyncType:    v1.SyncType_FULL,
			TaskName:    &[]string{time.Now().Format("20060102150405")}[0],
		}

		// 调用 fullSyncUsecase.CreateSyncAccount
		reply, err := s.fullSyncUsecase.CreateSyncAccount(ctx, req)
		if err != nil {
			s.log.Log(log.LevelError, "msg", "failed to create sync account", "err", err)
			return err
		}

		s.log.Log(log.LevelInfo, "msg", "sync account job completed", "task_id", reply.TaskId)
		return nil
	}))

	if err != nil {
		s.log.Log(log.LevelError, "msg", "failed to register sync_account_job", "err", err)
	}
}

// addDatabaseCheckJob 添加数据库检测任务
func addDatabaseCheckJob(s *CronService) {
	// 创建任务配置
	config := DefaultJobConfig("resource_check", "0 */10 * * * *") // 每10分钟执行
	config.WithTimeout(5 * time.Minute)
	config.Description = "数据库和Redis连接检测任务"

	// 注册任务
	_, err := s.AddJobWithConfig(config, JobExecutorFunc(func(ctx context.Context) error {
		s.log.Log(log.LevelInfo, "msg", "starting resource check")

		// 检查数据库连接
		databases := s.data.GetDBManager().ListDatabases()
		for dbType, config := range databases {
			if config.IsActive && config.DB != nil {
				if sqlDB, err := config.DB.DB(); err == nil {
					// 检查连接池状态
					stats := sqlDB.Stats()
					s.log.Log(log.LevelInfo, "msg", "database connection pool stats",
						"db_type", dbType,
						"max_open_connections", stats.MaxOpenConnections,
						"open_connections", stats.OpenConnections,
						"in_use", stats.InUse,
						"idle", stats.Idle)

					// 检查连接是否可用
					if err := sqlDB.PingContext(ctx); err != nil {
						s.log.Log(log.LevelError, "msg", "database ping failed", "db_type", dbType, "err", err)
					} else {
						s.log.Log(log.LevelInfo, "msg", "database ping successful", "db_type", dbType)
					}
				}
			}
		}

		// 检查 Redis 连接
		if s.data.GetRedis() != nil {
			redisClient := s.data.GetRedis()

			// 检查 Redis 连接
			if _, err := redisClient.Ping(ctx).Result(); err != nil {
				s.log.Log(log.LevelError, "msg", "redis ping failed", "err", err)
			} else {
				s.log.Log(log.LevelInfo, "msg", "redis ping successful")

				// 获取 Redis 信息
				info, err := redisClient.Info(ctx).Result()
				if err != nil {
					s.log.Log(log.LevelWarn, "msg", "failed to get redis info", "err", err)
				} else {
					s.log.Log(log.LevelInfo, "msg", "redis info retrieved", "info_length", len(info))
				}
			}
		}

		s.log.Log(log.LevelInfo, "msg", "resource check completed")
		return nil
	}))

	if err != nil {
		s.log.Log(log.LevelError, "msg", "failed to register resource_check job", "err", err)
	}
}
