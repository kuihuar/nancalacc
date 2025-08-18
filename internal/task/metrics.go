// internal/task/metrics.go
package task

// import (
// 	"context"
// 	"sync"
// 	"time"

// 	"github.com/go-kratos/kratos/v2/log"
// )

// // TaskMetrics 任务指标
// type TaskMetrics struct {
// 	TotalTasks     int64
// 	RunningTasks   int64
// 	CompletedTasks int64
// 	FailedTasks    int64
// 	TotalErrors    int64
// 	AvgDuration    time.Duration
// 	LastUpdate     time.Time
// }

// // MetricsCollector 指标收集器
// type MetricsCollector struct {
// 	metrics map[string]*TaskMetrics
// 	mu      sync.RWMutex
// 	log     *log.Helper
// }

// // NewMetricsCollector 创建指标收集器
// func NewMetricsCollector(logger log.Logger) *MetricsCollector {
// 	return &MetricsCollector{
// 		metrics: make(map[string]*TaskMetrics),
// 		log:     log.NewHelper(log.With(logger, "module", "task_metrics")),
// 	}
// }

// // RecordTaskStart 记录任务开始
// func (mc *MetricsCollector) RecordTaskStart(taskName string) {
// 	mc.mu.Lock()
// 	defer mc.mu.Unlock()

// 	if mc.metrics[taskName] == nil {
// 		mc.metrics[taskName] = &TaskMetrics{}
// 	}

// 	mc.metrics[taskName].TotalTasks++
// 	mc.metrics[taskName].RunningTasks++
// 	mc.metrics[taskName].LastUpdate = time.Now()
// }

// // RecordTaskComplete 记录任务完成
// func (mc *MetricsCollector) RecordTaskComplete(taskName string, duration time.Duration, err error) {
// 	mc.mu.Lock()
// 	defer mc.mu.Unlock()

// 	if mc.metrics[taskName] == nil {
// 		mc.metrics[taskName] = &TaskMetrics{}
// 	}

// 	metrics := mc.metrics[taskName]
// 	metrics.RunningTasks--
// 	metrics.LastUpdate = time.Now()

// 	if err != nil {
// 		metrics.FailedTasks++
// 		metrics.TotalErrors++
// 		mc.log.Errorf("task [%s] failed after %v: %v", taskName, duration, err)
// 	} else {
// 		metrics.CompletedTasks++
// 		// 更新平均执行时间
// 		if metrics.CompletedTasks > 0 {
// 			totalDuration := metrics.AvgDuration * time.Duration(metrics.CompletedTasks-1)
// 			metrics.AvgDuration = (totalDuration + duration) / time.Duration(metrics.CompletedTasks)
// 		} else {
// 			metrics.AvgDuration = duration
// 		}
// 		mc.log.Infof("task [%s] completed successfully in %v", taskName, duration)
// 	}
// }

// // GetTaskMetrics 获取任务指标
// func (mc *MetricsCollector) GetTaskMetrics(taskName string) *TaskMetrics {
// 	mc.mu.RLock()
// 	defer mc.mu.RUnlock()

// 	if metrics, exists := mc.metrics[taskName]; exists {
// 		// 返回副本以避免并发修改
// 		return &TaskMetrics{
// 			TotalTasks:     metrics.TotalTasks,
// 			RunningTasks:   metrics.RunningTasks,
// 			CompletedTasks: metrics.CompletedTasks,
// 			FailedTasks:    metrics.FailedTasks,
// 			TotalErrors:    metrics.TotalErrors,
// 			AvgDuration:    metrics.AvgDuration,
// 			LastUpdate:     metrics.LastUpdate,
// 		}
// 	}
// 	return nil
// }

// // GetAllMetrics 获取所有指标
// func (mc *MetricsCollector) GetAllMetrics() map[string]*TaskMetrics {
// 	mc.mu.RLock()
// 	defer mc.mu.RUnlock()

// 	result := make(map[string]*TaskMetrics)
// 	for name, metrics := range mc.metrics {
// 		result[name] = &TaskMetrics{
// 			TotalTasks:     metrics.TotalTasks,
// 			RunningTasks:   metrics.RunningTasks,
// 			CompletedTasks: metrics.CompletedTasks,
// 			FailedTasks:    metrics.FailedTasks,
// 			TotalErrors:    metrics.TotalErrors,
// 			AvgDuration:    metrics.AvgDuration,
// 			LastUpdate:     metrics.LastUpdate,
// 		}
// 	}
// 	return result
// }

// // GetSummaryMetrics 获取汇总指标
// func (mc *MetricsCollector) GetSummaryMetrics() *TaskMetrics {
// 	mc.mu.RLock()
// 	defer mc.mu.RUnlock()

// 	summary := &TaskMetrics{
// 		LastUpdate: time.Now(),
// 	}

// 	for _, metrics := range mc.metrics {
// 		summary.TotalTasks += metrics.TotalTasks
// 		summary.RunningTasks += metrics.RunningTasks
// 		summary.CompletedTasks += metrics.CompletedTasks
// 		summary.FailedTasks += metrics.FailedTasks
// 		summary.TotalErrors += metrics.TotalErrors
// 	}

// 	// 计算总体平均执行时间
// 	if summary.CompletedTasks > 0 {
// 		var totalDuration time.Duration
// 		var totalCompleted int64
// 		for _, metrics := range mc.metrics {
// 			if metrics.CompletedTasks > 0 {
// 				totalDuration += metrics.AvgDuration * time.Duration(metrics.CompletedTasks)
// 				totalCompleted += metrics.CompletedTasks
// 			}
// 		}
// 		if totalCompleted > 0 {
// 			summary.AvgDuration = totalDuration / time.Duration(totalCompleted)
// 		}
// 	}

// 	return summary
// }

// // ResetMetrics 重置指标
// func (mc *MetricsCollector) ResetMetrics(taskName string) {
// 	mc.mu.Lock()
// 	defer mc.mu.Unlock()

// 	if taskName == "" {
// 		// 重置所有指标
// 		mc.metrics = make(map[string]*TaskMetrics)
// 		mc.log.Info("all metrics reset")
// 	} else {
// 		// 重置指定任务的指标
// 		delete(mc.metrics, taskName)
// 		mc.log.Infof("metrics for task [%s] reset", taskName)
// 	}
// }

// // ExportMetrics 导出指标（用于监控系统）
// func (mc *MetricsCollector) ExportMetrics(ctx context.Context) map[string]interface{} {
// 	summary := mc.GetSummaryMetrics()
// 	allMetrics := mc.GetAllMetrics()

// 	export := map[string]interface{}{
// 		"summary":     summary,
// 		"tasks":       allMetrics,
// 		"exported_at": time.Now().Format(time.RFC3339),
// 	}

// 	return export
// }
