// internal/task/cron.go
package task

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/robfig/cron/v3"
)

type CronService struct {
	cron *cron.Cron
	log  *log.Helper
}

func NewCronService(logger log.Logger) *CronService {
	return &CronService{
		cron: cron.New(cron.WithSeconds()),
		log:  log.NewHelper(log.With(logger, "module", "task")),
	}
}

func (s *CronService) Start() {
	s.log.Info("启动 CronService")
	s.cron.Start()
}

func (s *CronService) Stop() {
	s.log.Info("停止 CronService")
	s.cron.Stop()
}

func (s *CronService) AddFunc(spec string, cmd func()) (cron.EntryID, error) {
	s.log.Infof("注册任务 [%s]", spec)
	return s.cron.AddFunc(spec, cmd)
}
