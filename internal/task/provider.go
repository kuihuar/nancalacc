package task

import (
	"nancalacc/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewCronServiceWithJobs,
)

func NewCronServiceWithJobs(accounterUsecase *biz.AccounterUsecase, logger log.Logger) *CronService {
	svc := NewCronService(accounterUsecase, logger)
	RegisterJobs(svc)
	return svc
}
