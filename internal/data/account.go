package data

import (
	"context"
	"errors"
	"time"

	"nancalacc/internal/biz"
	"nancalacc/internal/conf"
	"nancalacc/internal/data/models"
	"nancalacc/internal/dingtalk"

	//nancalaccLog "nancalacc/internal/log"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type accounterRepo struct {
	bizConf *conf.App
	data    *Data
	log     log.Logger
}

// getSyncDB 获取同步数据库连接
func (r *accounterRepo) getSyncDB() (*gorm.DB, error) {
	return r.data.GetSyncDB()
}

var (
	Source  = "sync"
	timeout = 5 * time.Second
)

// NewAccounterRepo .
func NewAccounterRepo(data *Data, logger log.Logger) biz.AccounterRepo {
	return &accounterRepo{
		bizConf: conf.Get().GetApp(),
		data:    data,
		log:     logger,
	}
}

func (r *accounterRepo) SaveUsers(ctx context.Context, users []*dingtalk.DingtalkDeptUser, taskId string) (int, error) {

	r.log.Log(log.LevelInfo, "msg", "SaveUsers", "users_count", len(users), "taskId", taskId)

	if len(users) == 0 {
		r.log.Log(log.LevelWarn, "msg", "SaveUsers", "users is empty")
		return 0, nil
	}

	thirdCompanyID := r.bizConf.ThirdCompanyId
	platformID := r.bizConf.PlatformIds
	entities := make([]*models.TbLasUser, 0, len(users))
	invalidUsers := make([]string, 0)

	// 批量验证和转换，收集无效数据
	for _, user := range users {
		if err := dingtalk.ValidateDingTalkUser(ctx, user); err != nil {
			r.log.Log(log.LevelWarn, "msg", "invalid user skipped", "user_id", user.Userid, "err", err)
			invalidUsers = append(invalidUsers, user.Userid)
			continue
		}
		entity := models.MakeLasUser(user, thirdCompanyID, platformID, Source, taskId)
		entities = append(entities, entity)
	}

	// 记录无效数据统计
	if len(invalidUsers) > 0 {
		r.log.Log(log.LevelWarn, "msg", "invalid users found", "count", len(invalidUsers), "invalid_ids", invalidUsers)
	}

	if len(entities) == 0 {
		r.log.Log(log.LevelWarn, "msg", "no valid users to save")
		return 0, nil
	}

	db, err := r.data.GetSyncDB()
	if err != nil {
		return 0, err
	}

	// 使用 Upsert 操作避免重复键错误
	result := db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "uid"}, {Name: "task_id"}, {Name: "third_company_id"}, {Name: "platform_id"}},
		DoNothing: true,
	}).Create(&entities)

	if result.Error != nil {
		r.log.Log(log.LevelError, "msg", "SaveUsers failed", "err", result.Error)
		return 0, result.Error
	}

	r.log.Log(log.LevelInfo, "msg", "SaveUsers completed", "saved_count", int(result.RowsAffected), "total_processed", len(users))
	return int(result.RowsAffected), nil
}

func (r *accounterRepo) SaveDepartments(ctx context.Context, depts []*dingtalk.DingtalkDept, taskId string) (int, error) {

	r.log.Log(log.LevelInfo, "msg", "SaveDepartments", "depts_count", len(depts), "taskId", taskId)

	// 使用传入的 context 并设置超时
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if len(depts) == 0 {
		r.log.Log(log.LevelWarn, "msg", "SaveDepartments", "depts is empty")
		return 0, nil
	}

	entities := make([]*models.TbLasDepartment, 0, len(depts)+1) // +1 for root department

	thirdCompanyID := r.bizConf.ThirdCompanyId
	platformID := r.bizConf.PlatformIds
	companyID := r.bizConf.CompanyId

	// 先添加部门实体
	for _, dep := range depts {
		entity := models.MakeTbLasDepartment(dep, thirdCompanyID, platformID, companyID, Source, taskId)
		entities = append(entities, entity)
	}

	// 添加根部门
	rootDep := models.MakeTbLasRootDepartment(thirdCompanyID, platformID, companyID, Source, taskId)
	entities = append(entities, rootDep)

	db, err := r.data.GetSyncDB()
	if err != nil {
		return 0, err
	}

	// 使用 Upsert 操作，避免重复键错误
	result := db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "did"}, {Name: "task_id"}, {Name: "third_company_id"}, {Name: "platform_id"}},
		DoNothing: true,
	}).Create(&entities)

	if result.Error != nil {
		r.log.Log(log.LevelError, "msg", "SaveDepartments failed", "err", result.Error)
		return 0, result.Error
	}

	r.log.Log(log.LevelInfo, "msg", "SaveDepartments completed", "saved_count", int(result.RowsAffected), "total_processed", len(entities))
	return int(result.RowsAffected), nil
}

func (r *accounterRepo) SaveDepartmentUserRelations(ctx context.Context, relations []*dingtalk.DingtalkDeptUserRelation, taskId string) (int, error) {
	r.log.Log(log.LevelInfo, "msg", "SaveDepartmentUserRelations", "relations_count", len(relations), "taskId", taskId)

	// 使用传入的 context 并设置超时
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if len(relations) == 0 {
		r.log.Log(log.LevelWarn, "msg", "SaveDepartmentUserRelations", "relations is empty")
		return 0, nil
	}

	entities := make([]*models.TbLasDepartmentUser, 0, len(relations))

	thirdCompanyID := r.bizConf.ThirdCompanyId
	platformID := r.bizConf.PlatformIds

	for _, relation := range relations {
		entity := models.MakeTbLasDepartmentUser(relation, thirdCompanyID, platformID, "", Source, taskId)
		entities = append(entities, entity)
	}

	db, err := r.data.GetSyncDB()
	if err != nil {
		return 0, err
	}

	result := db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "did"}, {Name: "uid"}, {Name: "task_id"}, {Name: "third_company_id"}, {Name: "platform_id"}},
		DoNothing: true,
	}).Create(&entities)

	if result.Error != nil {
		r.log.Log(log.LevelError, "msg", "SaveDepartmentUserRelations failed", "err", result.Error)
		return 0, result.Error
	}

	r.log.Log(log.LevelInfo, "msg", "SaveDepartmentUserRelations completed", "saved_count", int(result.RowsAffected), "total_processed", len(relations))
	return int(result.RowsAffected), nil
}

func (r *accounterRepo) SaveCompanyCfg(ctx context.Context, input *dingtalk.DingtalkCompanyCfg) error {
	r.log.Log(log.LevelInfo, "msg", "SaveCompanyCfg", "input", input)

	now := time.Now()
	entity := &models.TbCompanyCfg{
		ThirdCompanyId: input.ThirdCompanyId,
		PlatformIds:    input.PlatformIds,
		CompanyId:      input.CompanyId,
		Status:         1,
		Ctime:          now,
		Mtime:          now,
	}

	db, err := r.data.GetSyncDB()
	if err != nil {
		return err
	}
	err = db.WithContext(ctx).Where(models.TbCompanyCfg{
		ThirdCompanyId: input.ThirdCompanyId,
		CompanyId:      input.CompanyId,
	}).FirstOrCreate(entity).Error

	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			r.log.Log(log.LevelError, "msg", "SaveCompanyCfg", "company config already exists")
		} else {
			return err
		}

	}

	return nil
}

func (r *accounterRepo) ClearAll(ctx context.Context) error {
	db, err := r.data.GetSyncDB()
	if err != nil {
		return err
	}
	err = db.WithContext(ctx).Exec("truncate table tb_company_cfg").Error
	if err != nil {
		return err
	}
	err = db.WithContext(ctx).Exec("truncate table tb_las_department").Error
	if err != nil {
		return err
	}
	err = db.WithContext(ctx).Exec("truncate table tb_las_department_user").Error
	if err != nil {
		return err
	}
	err = db.WithContext(ctx).Exec("truncate table tb_las_account").Error
	if err != nil {
		return err
	}
	return nil
}

// user_del/user_update/user_add(update_type). . auto/manual(sync_type)
func (r *accounterRepo) SaveIncrementUsers(ctx context.Context, usersAdd, usersDel, usersUpd []*dingtalk.DingtalkDeptUser) error {
	r.log.Log(log.LevelInfo, "msg", "SaveIncrementUsers", "users_add", len(usersAdd), "users_del", len(usersDel), "users_upd", len(usersUpd))

	inputlen := len(usersAdd) + len(usersDel) + len(usersUpd)
	if inputlen == 0 {
		return errors.New("empty input")
	}

	// 使用传入的 context 并设置超时
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	entities := make([]*models.TbLasUserIncrement, 0, inputlen)
	invalidUsers := make([]string, 0)

	thirdCompanyID := r.bizConf.ThirdCompanyId
	platformID := r.bizConf.PlatformIds
	companyID := r.bizConf.CompanyId

	// 批量处理用户数据
	for _, add := range usersAdd {
		if err := dingtalk.ValidateDingTalkUser(ctx, add); err != nil {
			r.log.Log(log.LevelWarn, "msg", "invalid user_add skipped", "user_id", add.Userid, "err", err)
			invalidUsers = append(invalidUsers, add.Userid)
			continue
		}
		entity := models.MakeLasUserIncrement(add, thirdCompanyID, platformID, companyID, Source, "user_add")
		entities = append(entities, entity)
	}

	for _, del := range usersDel {
		if err := dingtalk.ValidateDingTalkUser(ctx, del); err != nil {
			r.log.Log(log.LevelWarn, "msg", "invalid user_del skipped", "user_id", del.Userid, "err", err)
			invalidUsers = append(invalidUsers, del.Userid)
			continue
		}
		entity := models.MakeLasUserIncrement(del, thirdCompanyID, platformID, companyID, Source, "user_del")
		entities = append(entities, entity)
	}

	for _, upd := range usersUpd {
		if err := dingtalk.ValidateDingTalkUser(ctx, upd); err != nil {
			r.log.Log(log.LevelWarn, "msg", "invalid user_upd skipped", "user_id", upd.Userid, "err", err)
			invalidUsers = append(invalidUsers, upd.Userid)
			continue
		}
		entity := models.MakeLasUserIncrement(upd, thirdCompanyID, platformID, companyID, Source, "user_update")
		entities = append(entities, entity)
	}

	// 记录无效数据统计
	if len(invalidUsers) > 0 {
		r.log.Log(log.LevelWarn, "msg", "invalid users found", "count", len(invalidUsers), "invalid_ids", invalidUsers)
	}

	if len(entities) == 0 {
		r.log.Log(log.LevelWarn, "msg", "no valid users to save")
		return nil
	}

	db, err := r.data.GetSyncDB()
	if err != nil {
		return err
	}

	// 使用 Upsert 操作避免重复键错误
	result := db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "uid"}, {Name: "task_id"}, {Name: "third_company_id"}, {Name: "platform_id"}},
		DoNothing: true,
	}).Create(&entities)

	if result.Error != nil {
		r.log.Log(log.LevelError, "msg", "SaveIncrementUsers failed", "err", result.Error)
		return result.Error
	}

	r.log.Log(log.LevelInfo, "msg", "SaveIncrementUsers completed", "saved_count", int(result.RowsAffected), "total_processed", len(entities))
	return nil
}

// dept_del/dept_update/dept_add/dept_move(update_type)
func (r *accounterRepo) SaveIncrementDepartments(ctx context.Context, deptsAdd, deptsDel, deptsUpd []*dingtalk.DingtalkDept) error {
	r.log.Log(log.LevelInfo, "msg", "SaveIncrementDepartments", "depts_add", len(deptsAdd), "depts_del", len(deptsDel), "depts_upd", len(deptsUpd))

	inputlen := len(deptsAdd) + len(deptsDel) + len(deptsUpd)
	if inputlen == 0 {
		return errors.New("empty input")
	}

	// 使用传入的 context 并设置超时
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	entities := make([]*models.TbLasDepartmentIncrement, 0, inputlen)
	thirdCompanyID := r.bizConf.ThirdCompanyId
	platformID := r.bizConf.PlatformIds

	// 批量处理部门数据
	for _, add := range deptsAdd {
		entity := models.MakeDepartmentIncrement(add, thirdCompanyID, platformID, "", Source, "dept_add")
		entities = append(entities, entity)
	}

	for _, del := range deptsDel {
		entity := models.MakeDepartmentIncrement(del, thirdCompanyID, platformID, "", Source, "dept_del")
		entities = append(entities, entity)
	}

	for _, upd := range deptsUpd {
		entity := models.MakeDepartmentIncrement(upd, thirdCompanyID, platformID, "", Source, "dept_update")
		entities = append(entities, entity)
	}

	if len(entities) == 0 {
		r.log.Log(log.LevelWarn, "msg", "no valid departments to save")
		return nil
	}

	db, err := r.data.GetSyncDB()
	if err != nil {
		return err
	}

	// 使用 Upsert 操作避免重复键错误
	result := db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "did"}, {Name: "task_id"}, {Name: "third_company_id"}, {Name: "platform_id"}},
		DoNothing: true,
	}).Create(&entities)

	if result.Error != nil {
		r.log.Log(log.LevelError, "msg", "SaveIncrementDepartments failed", "err", result.Error)
		return result.Error
	}

	r.log.Log(log.LevelInfo, "msg", "SaveIncrementDepartments completed", "saved_count", int(result.RowsAffected), "total_processed", len(entities))
	return nil
}

func (r *accounterRepo) SaveIncrementDepartmentUserRelations(ctx context.Context, relationsAdd, relationsDel, relationsUpd []*dingtalk.DingtalkDeptUserRelation) error {
	r.log.Log(log.LevelInfo, "msg", "SaveIncrementDepartmentUserRelations", "relations_add", len(relationsAdd), "relations_del", len(relationsDel), "relations_upd", len(relationsUpd))

	inputlen := len(relationsAdd) + len(relationsDel) + len(relationsUpd)
	if inputlen == 0 {
		return errors.New("empty input")
	}

	entities := make([]*models.TbLasDepartmentUserIncrement, 0, inputlen)

	thirdCompanyID := r.bizConf.ThirdCompanyId
	platformID := r.bizConf.PlatformIds

	// 批量处理部门用户关系数据
	for _, add := range relationsAdd {
		entity := models.MmakeTbLasDepartmentUserIncrement(add, thirdCompanyID, platformID, "", Source, "user_dept_add")
		entities = append(entities, entity)
	}

	for _, del := range relationsDel {
		entity := models.MmakeTbLasDepartmentUserIncrement(del, thirdCompanyID, platformID, "", Source, "user_dept_del")
		entities = append(entities, entity)
	}

	for _, upd := range relationsUpd {
		entity := models.MmakeTbLasDepartmentUserIncrement(upd, thirdCompanyID, platformID, "", Source, "user_dept_update")
		entities = append(entities, entity)
	}

	if len(entities) == 0 {
		r.log.Log(log.LevelWarn, "msg", "no valid relations to save")
		return nil
	}

	db, err := r.data.GetSyncDB()
	if err != nil {
		return err
	}

	// 使用 Upsert 操作避免重复键错误
	result := db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "did"}, {Name: "uid"}, {Name: "task_id"}, {Name: "third_company_id"}, {Name: "platform_id"}},
		DoNothing: true,
	}).Create(&entities)

	if result.Error != nil {
		r.log.Log(log.LevelError, "msg", "SaveIncrementDepartmentUserRelations failed", "err", result.Error)
		return result.Error
	}

	r.log.Log(log.LevelInfo, "msg", "SaveIncrementDepartmentUserRelations completed", "saved_count", int(result.RowsAffected), "total_processed", len(entities))
	return nil
}
func (r *accounterRepo) BatchSaveUsers(ctx context.Context, users []*models.TbLasUser) (int, error) {

	r.log.Log(log.LevelInfo, "msg", "BatchSaveUsers", "users_count", len(users))

	if len(users) == 0 {
		r.log.Log(log.LevelWarn, "msg", "BatchSaveUsers", "users is empty")
		return 0, nil
	}

	// 检查 context 是否已取消
	select {
	case <-ctx.Done():
		r.log.Log(log.LevelError, "msg", "BatchSaveUsers", "context canceled before database operation", "err", ctx.Err())
		return 0, ctx.Err()
	default:
		// 继续执行
	}

	// 记录 context 的截止时间（如果有的话）
	if deadline, ok := ctx.Deadline(); ok {
		r.log.Log(log.LevelInfo, "msg", "BatchSaveUsers", "context deadline", "deadline", deadline, "time_until_deadline", time.Until(deadline))
	}

	db, err := r.data.GetSyncDB()
	if err != nil {
		return 0, err
	}

	result := db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "uid"}, {Name: "task_id"}, {Name: "third_company_id"}, {Name: "platform_id"}},
		DoNothing: true,
	}).Create(users)

	if result.Error != nil {
		// 检查是否是 context 相关错误
		if errors.Is(result.Error, context.DeadlineExceeded) {
			r.log.Log(log.LevelError, "msg", "BatchSaveUsers", "context deadline exceeded", "err", result.Error)
		} else if errors.Is(result.Error, context.Canceled) {
			r.log.Log(log.LevelError, "msg", "BatchSaveUsers", "context canceled during database operation", "err", result.Error)
		} else {
			r.log.Log(log.LevelError, "msg", "BatchSaveUsers failed", "err", result.Error)
		}
		return 0, result.Error
	}

	r.log.Log(log.LevelInfo, "msg", "BatchSaveUsers completed", "saved_count", int(result.RowsAffected), "total_processed", len(users))
	return int(result.RowsAffected), nil
}
func (r *accounterRepo) BatchSaveDepts(ctx context.Context, depts []*models.TbLasDepartment) (int, error) {

	r.log.Log(log.LevelInfo, "msg", "BatchSaveDepts", "depts_count", len(depts))

	if len(depts) == 0 {
		r.log.Log(log.LevelWarn, "msg", "BatchSaveDepts", "depts is empty")
		return 0, nil
	}

	// 检查 context 是否已取消
	select {
	case <-ctx.Done():
		r.log.Log(log.LevelError, "msg", "BatchSaveDepts", "context canceled before database operation", "err", ctx.Err())
		return 0, ctx.Err()
	default:
		// 继续执行
	}

	// 记录 context 的截止时间（如果有的话）
	if deadline, ok := ctx.Deadline(); ok {
		r.log.Log(log.LevelInfo, "msg", "BatchSaveDepts", "context deadline", "deadline", deadline, "time_until_deadline", time.Until(deadline))
	}

	db, err := r.data.GetSyncDB()
	if err != nil {
		return 0, err
	}

	// 使用 Upsert 操作避免重复键错误
	result := db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "did"}, {Name: "task_id"}, {Name: "third_company_id"}, {Name: "platform_id"}},
		DoNothing: true,
	}).Create(depts)

	if result.Error != nil {
		// 检查是否是 context 相关错误
		if errors.Is(result.Error, context.DeadlineExceeded) {
			r.log.Log(log.LevelError, "msg", "BatchSaveDepts", "context deadline exceeded", "err", result.Error)
		} else if errors.Is(result.Error, context.Canceled) {
			r.log.Log(log.LevelError, "msg", "BatchSaveDepts", "context canceled during database operation", "err", result.Error)
		} else {
			r.log.Log(log.LevelError, "msg", "BatchSaveDepts failed", "err", result.Error)
		}
		return 0, result.Error
	}

	r.log.Log(log.LevelInfo, "msg", "BatchSaveDepts completed", "saved_count", int(result.RowsAffected), "total_processed", len(depts))
	return int(result.RowsAffected), nil
}
func (r *accounterRepo) BatchSaveDeptUsers(ctx context.Context, usersdepts []*models.TbLasDepartmentUser) (int, error) {

	r.log.Log(log.LevelInfo, "msg", "BatchSaveDeptUsers", "usersdepts_count", len(usersdepts))

	// 直接使用传入的 ctx，继承父级的超时设置

	if len(usersdepts) == 0 {
		r.log.Log(log.LevelWarn, "msg", "BatchSaveDeptUsers", "usersdepts is empty")
		return 0, nil
	}

	// 检查 context 是否已取消
	select {
	case <-ctx.Done():
		r.log.Log(log.LevelError, "msg", "BatchSaveDeptUsers", "context canceled before database operation", "err", ctx.Err())
		return 0, ctx.Err()
	default:
		// 继续执行
	}

	// 记录 context 的截止时间（如果有的话）
	if deadline, ok := ctx.Deadline(); ok {
		r.log.Log(log.LevelInfo, "msg", "BatchSaveDeptUsers", "context deadline", "deadline", deadline, "time_until_deadline", time.Until(deadline))
	}

	db, err := r.data.GetSyncDB()
	if err != nil {
		return 0, err
	}

	// 使用 Upsert 操作避免重复键错误
	result := db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "did"}, {Name: "uid"}, {Name: "task_id"}, {Name: "third_company_id"}, {Name: "platform_id"}},
		DoNothing: true,
	}).Create(usersdepts)

	if result.Error != nil {
		// 检查是否是 context 相关错误
		if errors.Is(result.Error, context.DeadlineExceeded) {
			r.log.Log(log.LevelError, "msg", "BatchSaveDeptUsers", "context deadline exceeded", "err", result.Error)
		} else if errors.Is(result.Error, context.Canceled) {
			r.log.Log(log.LevelError, "msg", "BatchSaveDeptUsers", "context canceled during database operation", "err", result.Error)
		} else {
			r.log.Log(log.LevelError, "msg", "BatchSaveDeptUsers failed", "err", result.Error)
		}
		return 0, result.Error
	}

	r.log.Log(log.LevelInfo, "msg", "BatchSaveDeptUsers completed", "saved_count", int(result.RowsAffected), "total_processed", len(usersdepts))
	return int(result.RowsAffected), nil
}

func (r *accounterRepo) CreateTask(ctx context.Context, taskName string) (int, error) {
	r.log.Log(log.LevelInfo, "msg", "CreateTask", "name", taskName)

	// 使用传入的 context 并设置超时
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	task := &models.Task{
		Title:         taskName,
		Description:   taskName,
		CreatedAt:     time.Now(),
		Status:        models.TaskStatusPending,
		Progress:      0,
		StartDate:     time.Now(),
		DueDate:       time.Now().Add(time.Minute * 30),
		CompletedAt:   time.Now(),
		CreatorID:     99,
		EstimatedTime: 10,
		ActualTime:    0,
	}

	db, err := r.data.GetMainDB()
	if err != nil {
		return 0, err
	}
	result := db.WithContext(ctx).Where("title=?", taskName).FirstOrCreate(task)

	if result.Error != nil {
		// 处理其他错误
		return 0, result.Error
	}

	if result.RowsAffected > 0 {
		return 1, nil
	} else {
		return 0, nil
	}
}

func (r *accounterRepo) UpdateTask(ctx context.Context, taskName, status string) error {

	r.log.Log(log.LevelInfo, "msg", "UpdateTask", "taskName", taskName, "status", status)
	// pending/in_progress/completed/cancelled

	var task models.Task
	db, err := r.data.GetMainDB()
	if err != nil {
		return err
	}
	if err := db.Model(&models.Task{}).WithContext(ctx).Where("title=?", taskName).Find(&task).Error; err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			r.log.Log(log.LevelError, "msg", "UpdateTask", "查询超时")
		}
		return err
	}
	now := time.Now()
	task.ActualTime = int(now.Sub(task.StartDate).Seconds())
	task.UpdatedAt = now
	task.Status = status
	if status == models.TaskStatusCompleted || status == models.TaskStatusCancelled {
		task.CompletedAt = now
	}

	result := db.WithContext(ctx).Updates(task)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			r.log.Log(log.LevelError, "msg", "UpdateTask", "task already exists")
		} else {
			return result.Error
		}

	}

	return nil
}

func (r *accounterRepo) GetTask(ctx context.Context, taskName string) (*models.Task, error) {

	r.log.Log(log.LevelInfo, "msg", "GetTask", "name", taskName)
	task := &models.Task{}
	db, err := r.data.GetMainDB()
	if err != nil {
		return nil, err
	}
	result := db.WithContext(ctx).Where("title=?", taskName).Find(task)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, errors.New("notfound")
	}
	return task, nil
}
