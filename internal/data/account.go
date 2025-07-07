package data

import (
	"context"

	"nancalacc/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type accounterRepo struct {
	data *Data
	log  *log.Helper
}

// NewAccounterRepo .
func NewAccounterRepo(data *Data, logger log.Logger) biz.AccounterRepo {
	return &accounterRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *accounterRepo) SaveUsers(ctx context.Context, users []*biz.DingtalkDeptUser) (int, error) {
	r.log.Infof("SaveUsers: %v", users)
	return 0, nil
}

func (r *accounterRepo) SaveDepartments(ctx context.Context, depts []*biz.DingtalkDept) (int, error) {
	return 0, nil
}

func (r *accounterRepo) SaveDepartmentUserRelations(ctx context.Context, relations []*biz.DingtalkDeptUserRelation) (int, error) {
	r.log.Infof("SaveDepartmentUserRelations: %v", relations)
	return 0, nil
}

func (r *accounterRepo) SaveCompanyCfg(ctx context.Context, cfg *biz.DingtalkCompanyCfg) error {
	r.log.Infof("SaveCompanyCfg: %v", cfg)
	return nil
}
