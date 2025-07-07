// internal/task/jobs.go
package task

func RegisterJobs(s *CronService) {
	// 每5秒执行一次
	// s.AddFunc("*/5 * * * * *", func() {
	// 	s.log.Info("🔥 执行任务：每5秒任务")
	// })

	// 每分钟执行一次
	s.AddFunc("0 * * * * *", func() {
		//s.log.Info("⏰ 执行任务：每分钟任务")
	})
}
