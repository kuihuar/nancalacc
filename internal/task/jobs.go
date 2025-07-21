// internal/task/jobs.go
package task

func RegisterJobs(s *CronService) {
	// 每小时执行一次
	// s.AddFunc("0 3 * * * *", func() {
	// 	//s.log.Info("⏰ 执行任务：每小时任务")
	// })
	// // 0点执行一次
	// s.AddFunc("0 0 * * * *", func() {
	// 	//s.log.Info("⏰ 执行任务：每小时任务")
	// })
	// 每30分钟执行一次
	s.AddFunc("0 */30 * * * *", func() {
		//s.log.Info("⏰ 执行任务：每30分钟任务")
	})
	// 每5秒执行一次
	// s.AddFunc("10 2 * * * *", func() {
	// 	s.log.Info("🔥 执行任务: 每天2点10分0秒全量同步一次任务")
	// 	ctx, cancel := context.WithCancel(context.Background())
	// 	defer cancel()
	// 	res, err := s.accounterUsecase.CreateSyncAccount(ctx, &v1.CreateSyncAccountRequest{
	// 		TriggerType: v1.TriggerType_TRIGGER_SCHEDULED,
	// 		SyncType:    v1.SyncType_FULL,
	// 	})
	// 	s.log.Infof("CreateSyncAccount: %v, err: %v", res, err)
	// })

	// 每分钟执行一次
	// s.AddFunc("0 * * * * *", func() {
	// 	//s.log.Info("⏰ 执行任务：每分钟任务")
	// })
}
