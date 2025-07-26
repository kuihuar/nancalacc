package wps

import (
	"context"
	"nancalacc/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
)

type wps struct {
	serviceConf *conf.Service
	log         *log.Helper
}

func NewWps(serviceConf *conf.Service, logger log.Logger) Wps {
	return &wps{
		serviceConf: serviceConf,
		log:         log.NewHelper(logger),
	}
}

func (w *wps) BatchGetDepartment(ctx context.Context, accessToken string, req *BatchGetDepartmentRequest) (BatchGetDepartmentResponse, error) {
	return BatchGetDepartmentResponse{}, nil
}
