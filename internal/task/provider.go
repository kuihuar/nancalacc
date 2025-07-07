package task

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewCronServiceWithJobs,
)

func NewCronServiceWithJobs(logger log.Logger) *CronService {
	svc := NewCronService(logger)
	RegisterJobs(svc)
	return svc
}
