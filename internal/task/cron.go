// internal/task/cron.go
package task

import (
	"nancalacc/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/robfig/cron/v3"
)

type CronService struct {
	cron             *cron.Cron
	log              *log.Helper
	accounterUsecase *biz.AccounterUsecase
}

func NewCronService(accounterUsecase *biz.AccounterUsecase, logger log.Logger) *CronService {
	return &CronService{
		cron:             cron.New(cron.WithSeconds()),
		log:              log.NewHelper(log.With(logger, "module", "task")),
		accounterUsecase: accounterUsecase,
	}
}

func (s *CronService) Start() {
	s.log.Info("cron service starting...")
	s.cron.Start()
	s.log.Info("cron service start success")
}

func (s *CronService) Stop() {
	s.log.Info("cron service stoping...")
	s.cron.Stop()
}

func (s *CronService) AddFunc(spec string, cmd func()) (cron.EntryID, error) {
	s.log.Infof("add task [%s]", spec)
	return s.cron.AddFunc(spec, cmd)
}
