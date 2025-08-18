package task

import (
	"nancalacc/internal/biz"
	"nancalacc/internal/data"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet 是 task 模块的 provider 集合
var ProviderSet = wire.NewSet(NewCronServiceWithJobs)

// NewCronServiceWithJobs 创建带预定义任务的定时任务服务
func NewCronServiceWithJobs(fullSyncUsecase *biz.FullSyncUsecase, data *data.Data, logger log.Logger) *CronService {
	// 创建定时任务服务
	svc := NewCronService(fullSyncUsecase, data, logger)

	// 注册预定义任务
	//RegisterJobs(svc)

	RegisterJobsConfig(svc)

	return svc
}

// startMetricsExport 启动指标导出（示例）
// func startMetricsExport(svc *CronService, mc *MetricsCollector, logger log.Logger) {
// 	log := log.NewHelper(log.With(logger, "module", "metrics_export"))

// 	ticker := time.NewTicker(5 * time.Minute) // 每5分钟导出一次指标
// 	defer ticker.Stop()

// 	// 创建一个可取消的上下文
// 	ctx, cancel := context.WithCancel(context.Background())
// 	defer cancel()

// 	for {
// 		select {
// 		case <-ticker.C:
// 			ctx := context.Background()
// 			metrics := mc.ExportMetrics(ctx)

// 			// 这里可以将指标发送到监控系统
// 			// 例如：Prometheus, InfluxDB, 或者日志系统
// 			log.Infof("exported metrics: %+v", metrics)

// 		case <-ctx.Done():
// 			log.Info("metrics export stopped")
// 			return
// 		}
// 	}
// }
