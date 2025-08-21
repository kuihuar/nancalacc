package biz

import (
	"context"
	"nancalacc/internal/data/models"
	"nancalacc/internal/dingtalk"
	"time"
)

// AccounterRepo 定义了账户同步相关的数据访问接口
type AccounterRepo interface {
	// 全量同步相关方法
	SaveUsers(ctx context.Context, users []*dingtalk.DingtalkDeptUser, taskId string) (int, error)
	SaveDepartments(ctx context.Context, depts []*dingtalk.DingtalkDept, taskId string) (int, error)
	SaveDepartmentUserRelations(ctx context.Context, relations []*dingtalk.DingtalkDeptUserRelation, taskId string) (int, error)
	SaveCompanyCfg(ctx context.Context, cfg *dingtalk.DingtalkCompanyCfg) error

	// 增量同步相关方法
	SaveIncrementUsers(ctx context.Context, usersAdd, usersDel, usersUpd []*dingtalk.DingtalkDeptUser) error
	SaveIncrementDepartments(ctx context.Context, deptsAdd, deptsDel, deptsUpd []*dingtalk.DingtalkDept) error
	SaveIncrementDepartmentUserRelations(ctx context.Context, relationsAdd, relationsDel, relationsUpd []*dingtalk.DingtalkDeptUserRelation) error

	// 批量操作
	BatchSaveUsers(ctx context.Context, users []*models.TbLasUser) (int, error)
	BatchSaveDepts(ctx context.Context, depts []*models.TbLasDepartment) (int, error)
	BatchSaveDeptUsers(ctx context.Context, deptusers []*models.TbLasDepartmentUser) (int, error)

	// 任务管理
	CreateTask(ctx context.Context, taskName string) (int, error)
	UpdateTask(ctx context.Context, taskName, status string) error
	GetTask(ctx context.Context, taskName string) (*models.Task, error)

	// 查询相关
	BatchGetDeptUsers(ctx context.Context, taskName, thirdCompanyId, platformId string) ([]*models.TbLasDepartmentUser, error)
	BatchGetUsers(ctx context.Context, taskName, thirdCompanyId, platformId string) ([]*models.TbLasUser, error)
	BatchGetDepts(ctx context.Context, taskName, thirdCompanyId, platformId string) ([]*models.TbLasDepartment, error)

	// 清理操作
	ClearAll(ctx context.Context) error
}

// CacheService 定义了缓存服务接口
type CacheService interface {
	Get(ctx context.Context, key string) (interface{}, bool, error)
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Del(ctx context.Context, key string) error
}
