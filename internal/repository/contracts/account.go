package contracts

import (
	"context"
	"nancalacc/internal/data/models"
	"nancalacc/internal/dingtalk"
)

// AccountRepository 账户数据访问接口
type AccountRepository interface {
	// 用户相关
	SaveUsers(ctx context.Context, users []*dingtalk.DingtalkDeptUser, taskId string) (int, error)
	BatchSaveUsers(ctx context.Context, users []*models.TbLasUser) (int, error)
	BatchGetUsers(ctx context.Context, taskName, thirdCompanyId, platformId string) ([]*models.TbLasUser, error)
	SaveIncrementUsers(ctx context.Context, usersAdd, usersDel, usersUpd []*dingtalk.DingtalkDeptUser) error

	// 部门相关
	SaveDepartments(ctx context.Context, depts []*dingtalk.DingtalkDept, taskId string) (int, error)
	BatchSaveDepts(ctx context.Context, depts []*models.TbLasDepartment) (int, error)
	BatchGetDepts(ctx context.Context, taskName, thirdCompanyId, platformId string) ([]*models.TbLasDepartment, error)
	SaveIncrementDepartments(ctx context.Context, deptsAdd, deptsDel, deptsUpd []*dingtalk.DingtalkDept) error

	// 部门用户关系
	SaveDepartmentUserRelations(ctx context.Context, relations []*dingtalk.DingtalkDeptUserRelation, taskId string) (int, error)
	BatchSaveDeptUsers(ctx context.Context, deptusers []*models.TbLasDepartmentUser) (int, error)
	BatchGetDeptUsers(ctx context.Context, taskName, thirdCompanyId, platformId string) ([]*models.TbLasDepartmentUser, error)
	SaveIncrementDepartmentUserRelations(ctx context.Context, relationsAdd, relationsDel, relationsUpd []*dingtalk.DingtalkDeptUserRelation) error

	// 公司配置
	SaveCompanyCfg(ctx context.Context, cfg *dingtalk.DingtalkCompanyCfg) error

	// 清理操作
	ClearAll(ctx context.Context) error
}
