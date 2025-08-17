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

	r.log.Log(log.LevelInfo, "msg", "SaveUsers", "users", users, "taskId", taskId)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if len(users) == 0 {
		r.log.Log(log.LevelWarn, "msg", "SaveUsers", "users is empty")
		return 0, nil
	}
	thirdCompanyID := r.bizConf.ThirdCompanyId
	platformID := r.bizConf.PlatformIds
	entities := make([]*models.TbLasUser, 0, len(users))
	for _, user := range users {
		err := dingtalk.ValidateDingTalkUser(ctx, user)
		if err != nil {
			r.log.Log(log.LevelError, "msg", "ValidateDingTalkUser", "err", err)
			continue
		}
		entity := models.MakeLasUser(user, thirdCompanyID, platformID, Source, taskId)
		entities = append(entities, entity)
	}
	// result := r.data.db.WithContext(ctx).Clauses(clause.OnConflict{
	// 		Columns: []clause.Column{
	// 			{Name: "uid"},
	// 			{Name: "task_id"},
	// 			{Name: "platform_id"},
	// 		},
	// 		DoNothing: true,
	// 	}).Clauses(clause.OnConflict{
	// 		Columns: []clause.Column{
	// 			{Name: "account"},
	// 			{Name: "task_id"},
	// 			{Name: "third_company_id"},
	// 		},
	// 		DoNothing: true,
	// 	}).Create(&entities)
	db, err := r.data.GetSyncDB()
	if err != nil {
		return 0, err
	}
	result := db.WithContext(ctx).Create(&entities)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			r.log.Log(log.LevelError, "msg", "SaveUsers", "user already exists")
		} else {
			return 0, result.Error
		}

	}

	return int(result.RowsAffected), nil
}

func (r *accounterRepo) SaveDepartments(ctx context.Context, depts []*dingtalk.DingtalkDept, taskId string) (int, error) {

	r.log.Log(log.LevelInfo, "msg", "SaveDepartments", "depts", depts, "taskId", taskId)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	entities := make([]*models.TbLasDepartment, 0, len(depts))

	thirdCompanyID := r.bizConf.ThirdCompanyId
	platformID := r.bizConf.PlatformIds
	companyID := r.bizConf.CompanyId
	rootDep := models.MakeTbLasRootDepartment(thirdCompanyID, platformID, companyID, Source, taskId)
	for _, dep := range depts {
		entity := models.MakeTbLasDepartment(dep, thirdCompanyID, platformID, companyID, Source, taskId)
		entities = append(entities, entity)
	}
	entities = append(entities, rootDep)
	db, err := r.data.GetSyncDB()
	if err != nil {
		return 0, err
	}
	result := db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "did"}, {Name: "task_id"}, {Name: "third_company_id"}, {Name: "platform_id"}},
		DoNothing: true,
	}).Create(&entities)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			r.log.Log(log.LevelError, "msg", "SaveDepartments", "department already exists")
		} else {
			return 0, result.Error
		}

	}

	return int(result.RowsAffected), nil
}

func (r *accounterRepo) SaveDepartmentUserRelations(ctx context.Context, relations []*dingtalk.DingtalkDeptUserRelation, taskId string) (int, error) {
	r.log.Log(log.LevelInfo, "msg", "SaveDepartmentUserRelations", "relations.size", len(relations), "taskId", taskId)

	for _, relation := range relations {
		r.log.Log(log.LevelInfo, "msg", "SaveDepartmentUserRelations", "relation", relation)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

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
		return 0, result.Error
	}

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

// func (r *accounterRepo) CallEcisaccountsyncAll(ctx context.Context, taskId string) (biz.EcisaccountsyncAllResponse, error) {
// 	r.log.Log(log.LevelInfo, "msg", "CallEcisaccountsyncAll", "taskId", taskId)

// 	path := r.serviceConf.Business.EcisaccountsyncUrl

// 	// path := "http://encs-pri-proxy-gateway/ecisaccountsync/api/sync/all"
// 	var resp biz.EcisaccountsyncAllResponse

// 	thirdCompanyID := r.serviceConf.Business.ThirdCompanyId

// 	collectCost := "1100000"
// 	uri := fmt.Sprintf("%s?taskId=%s&thirdCompanyId=%s&collectCost=%s", path, taskId, thirdCompanyID, collectCost)

// 	r.log.Log(log.LevelInfo, "msg", "CallEcisaccountsyncAll", "uri", uri)
// 	bs, err := httputil.PostJSON(uri, nil, time.Second*10)
// 	r.log.Log(log.LevelInfo, "msg", "CallEcisaccountsyncAll.Post", "output", "bs", string(bs), "err", err)

// 	if err != nil {
// 		return resp, err
// 	}
// 	err = json.Unmarshal(bs, &resp)
// 	if err != nil {
// 		return resp, fmt.Errorf("Unmarshal err: %w", err)
// 	}
// 	if resp.Code != "200" {
// 		return resp, fmt.Errorf("code not 200: %s", resp.Code)
// 	}

// 	return resp, nil

// }

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
	r.log.Log(log.LevelInfo, "msg", "SaveIncrementUsers")
	inputlen := len(usersAdd) + len(usersDel) + len(usersUpd)
	if inputlen == 0 {
		return errors.New("empty input")
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	entities := make([]*models.TbLasUserIncrement, 0, inputlen)

	thirdCompanyID := r.bizConf.ThirdCompanyId
	platformID := r.bizConf.PlatformIds
	companyID := r.bizConf.CompanyId

	// user_del/user_update/user_add
	for _, add := range usersAdd {
		err := dingtalk.ValidateDingTalkUser(ctx, add)
		if err != nil {
			r.log.Log(log.LevelError, "msg", "ValidateDingTalkUser", "err", err)
			continue
		}
		entity := models.MakeLasUserIncrement(add, thirdCompanyID, platformID, companyID, Source, "user_add")
		entities = append(entities, entity)
	}
	for _, del := range usersDel {
		err := dingtalk.ValidateDingTalkUser(ctx, del)
		if err != nil {
			r.log.Log(log.LevelError, "msg", "ValidateDingTalkUser", "err", err)
			continue
		}
		entity := models.MakeLasUserIncrement(del, thirdCompanyID, platformID, companyID, Source, "user_del")
		// entity := r.makeLasUserIncrement(user, "user_del")
		entities = append(entities, entity)
	}
	for _, upd := range usersUpd {
		err := dingtalk.ValidateDingTalkUser(ctx, upd)
		if err != nil {
			r.log.Log(log.LevelError, "msg", "ValidateDingTalkUser", "err", err)
			continue
		}
		// entity := r.makeLasUserIncrement(user, "user_update")
		entity := models.MakeLasUserIncrement(upd, thirdCompanyID, platformID, companyID, Source, "user_update")
		entities = append(entities, entity)
	}

	r.log.Log(log.LevelInfo, "msg", "SaveIncrementUsers", "entities", entities)

	for i, item := range entities {
		r.log.Log(log.LevelInfo, "msg", "SaveIncrementUsers", "entities", "i", i, "item", item)
	}

	db, err := r.data.GetSyncDB()
	if err != nil {
		return err
	}
	result := db.WithContext(ctx).Create(&entities)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			r.log.Log(log.LevelError, "msg", "SaveIncrementUsers", "user already exists")
		} else {
			return result.Error
		}
	}
	return nil
}

// dept_del/dept_update/dept_add/dept_move(update_type)
func (r *accounterRepo) SaveIncrementDepartments(ctx context.Context, deptsAdd, deptsDel, deptsUpd []*dingtalk.DingtalkDept) error {
	r.log.Log(log.LevelInfo, "msg", "SaveIncrementDepartments")
	inputlen := len(deptsAdd) + len(deptsDel) + len(deptsUpd)
	if inputlen == 0 {
		return errors.New("empty input")
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// dept_del/dept_update/dept_add
	entities := make([]*models.TbLasDepartmentIncrement, 0, inputlen)
	thirdCompanyID := r.bizConf.ThirdCompanyId
	platformID := r.bizConf.PlatformIds

	for _, add := range deptsAdd {
		entity := models.MakeDepartmentIncrement(add, thirdCompanyID, platformID, "", Source, "dept_add")
		entities = append(entities, entity)
	}
	for _, del := range deptsDel {
		// entity := r.makeDepartmentIncrement(dep, "dept_del")
		entity := models.MakeDepartmentIncrement(del, thirdCompanyID, platformID, "", Source, "dept_del")
		entities = append(entities, entity)
	}

	for _, upd := range deptsUpd {
		// entity := r.makeDepartmentIncrement(dep, "dept_update")
		entity := models.MakeDepartmentIncrement(upd, thirdCompanyID, platformID, "", Source, "dept_update")
		entities = append(entities, entity)
	}

	r.log.Log(log.LevelInfo, "msg", "SaveIncrementDepartments", "entities", entities)

	for i, item := range entities {
		r.log.Log(log.LevelInfo, "msg", "SaveIncrementDepartments", "entities", "i", i, "item", item)
	}

	db, err := r.data.GetSyncDB()
	if err != nil {
		return err
	}
	result := db.WithContext(ctx).Create(&entities)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			r.log.Log(log.LevelError, "msg", "SaveIncrementDepartments", "user already exists")
		} else {
			return result.Error
		}

	}

	return nil
}

func (r *accounterRepo) SaveIncrementDepartmentUserRelations(ctx context.Context, relationsAdd, relationsDel, relationsUpd []*dingtalk.DingtalkDeptUserRelation) error {
	r.log.Log(log.LevelInfo, "msg", "SaveIncrementDepartmentUserRelations")

	inputlen := len(relationsAdd) + len(relationsDel) + len(relationsUpd)
	if inputlen == 0 {
		return errors.New("empty input")
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	r.log.Log(log.LevelInfo, "msg", "SaveIncrementDepartmentUserRelations", "input", "relationsAdd", "relationsDel", relationsDel, "relationsUpd", relationsUpd)

	for i, item := range relationsAdd {
		r.log.Log(log.LevelInfo, "msg", "SaveIncrementDepartmentUserRelations", "input", "relationsAdd", "i", i, "item", item)
	}
	for i, item := range relationsDel {
		r.log.Log(log.LevelInfo, "msg", "SaveIncrementDepartmentUserRelations", "input", "relationsDel", "i", i, "item", item)
	}
	for i, item := range relationsUpd {
		r.log.Log(log.LevelInfo, "msg", "SaveIncrementDepartmentUserRelations", "input", "relationsUpd", "i", i, "item", item)
	}

	entities := make([]*models.TbLasDepartmentUserIncrement, 0, inputlen)

	thirdCompanyID := r.bizConf.ThirdCompanyId
	platformID := r.bizConf.PlatformIds
	// user_dept_add/user_dept_del/user_dept_update/user_dept_move
	for _, add := range relationsAdd {
		entity := models.MmakeTbLasDepartmentUserIncrement(add, thirdCompanyID, platformID, "", Source, "user_dept_add")
		entities = append(entities, entity)
	}

	for _, del := range relationsDel {
		// entity := r.makeDeptUserRelatins(relation, "user_dept_del")
		entity := models.MmakeTbLasDepartmentUserIncrement(del, thirdCompanyID, platformID, "", Source, "user_dept_del")
		entities = append(entities, entity)
	}
	for _, upd := range relationsUpd {
		// entity := r.makeDeptUserRelatins(relation, "user_dept_update")
		entity := models.MmakeTbLasDepartmentUserIncrement(upd, thirdCompanyID, platformID, "", Source, "user_dept_update")
		entities = append(entities, entity)
	}

	r.log.Log(log.LevelInfo, "msg", "SaveIncrementDepartmentUserRelations", "entities", entities)

	for i, item := range entities {
		r.log.Log(log.LevelInfo, "msg", "SaveIncrementDepartmentUserRelations", "entities", "i", i, "item", item)
	}

	db, err := r.data.GetSyncDB()
	if err != nil {
		return err
	}
	result := db.WithContext(ctx).Create(&entities)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			r.log.Log(log.LevelError, "msg", "SaveIncrementDepartmentUserRelations", "relation already exists")
		} else {
			return result.Error
		}

	}

	return nil
}
func (r *accounterRepo) BatchSaveUsers(ctx context.Context, users []*models.TbLasUser) (int, error) {

	r.log.Log(log.LevelInfo, "msg", "BatchSaveUsers", "users", users)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if len(users) == 0 {
		r.log.Log(log.LevelWarn, "msg", "BatchSaveUsers", "users is empty")
		return 0, nil
	}

	// for _, user := range users {
	// 	r.log.Warn(user)
	// }

	db, err := r.data.GetSyncDB()
	if err != nil {
		return 0, err
	}
	result := db.WithContext(ctx).Create(users)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			r.log.Log(log.LevelError, "msg", "BatchSaveUsers", "user already exists")
		} else {
			return 0, result.Error
		}
	}

	return int(result.RowsAffected), nil
}
func (r *accounterRepo) BatchSaveDepts(ctx context.Context, depts []*models.TbLasDepartment) (int, error) {

	r.log.Log(log.LevelInfo, "msg", "BatchSaveDepts", "depts", depts)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if len(depts) == 0 {
		r.log.Log(log.LevelWarn, "msg", "BatchSaveDepts", "depts is empty")
		return 0, nil
	}
	// for _, dept := range depts {
	// 	r.log.Info(dept)
	// }

	db, err := r.data.GetSyncDB()
	if err != nil {
		return 0, err
	}
	result := db.WithContext(ctx).Create(depts)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			r.log.Log(log.LevelError, "msg", "BatchSaveDepts", "dept already exists")
		} else {
			return 0, result.Error
		}

	}
	return int(result.RowsAffected), nil
}
func (r *accounterRepo) BatchSaveDeptUsers(ctx context.Context, usersdepts []*models.TbLasDepartmentUser) (int, error) {

	r.log.Log(log.LevelInfo, "msg", "BatchSaveDeptUsers", "usersdepts", usersdepts)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if len(usersdepts) == 0 {
		r.log.Log(log.LevelWarn, "msg", "BatchSaveDeptUsers", "usersdepts is empty")
		return 0, nil
	}
	// for _, userdept := range usersdepts {
	// 	r.log.Info(userdept)
	// }
	db, err := r.data.GetSyncDB()
	if err != nil {
		return 0, err
	}
	result := db.WithContext(ctx).Create(usersdepts)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			r.log.Log(log.LevelError, "msg", "BatchSaveDeptUsers", "deptuser already exists")
		} else {
			return 0, result.Error
		}

	}
	return int(result.RowsAffected), nil
}

func (r *accounterRepo) CreateTask(ctx context.Context, taskName string) (int, error) {
	r.log.Log(log.LevelInfo, "msg", "CreateTask", "name", taskName)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
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

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

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

	r.log.Log(log.LevelInfo, "msg", "CreateTask", "name", taskName)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
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
