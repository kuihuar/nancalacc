package data

import (
	"context"
	"errors"
	"strconv"
	"time"

	"nancalacc/internal/biz"
	"nancalacc/internal/conf"
	"nancalacc/internal/data/models"
	"nancalacc/internal/dingtalk"
	"nancalacc/pkg/cipherutil"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type accounterRepo struct {
	serviceConf *conf.Service
	data        *Data
	log         *log.Helper
}

var (
	Source  = "sync"
	timeout = 5 * time.Second
)

// NewAccounterRepo .
func NewAccounterRepo(serviceConf *conf.Service, data *Data, logger log.Logger) biz.AccounterRepo {
	return &accounterRepo{
		serviceConf: serviceConf,
		data:        data,
		log:         log.NewHelper(logger),
	}
}

func (r *accounterRepo) SaveUsers(ctx context.Context, users []*dingtalk.DingtalkDeptUser, taskId string) (int, error) {

	r.log.WithContext(ctx).Infof("SaveUsers users: %v, taskId :%s", users, taskId)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	r.log.Infof("SaveUsers: %v", users)
	if len(users) == 0 {
		r.log.Warn("users is empty")
		return 0, nil
	}
	entities := make([]*models.TbLasUser, 0, len(users))

	thirdCompanyID := r.serviceConf.Business.ThirdCompanyId
	platformID := r.serviceConf.Business.PlatformIds
	for _, user := range users {
		if user.Name == "" {
			user.Name = user.Userid
		}
		if user.Nickname == "" {
			user.Nickname = user.Name
		}
		var account string
		if user.Mobile != "" {
			account = user.Mobile
		} else {
			account = user.Userid
		}
		secretKey := r.serviceConf.Auth.Self.SecretKey

		email, errEmail := cipherutil.AesEncryptGcmByKey(user.Email, secretKey)

		phone, errPhone := cipherutil.AesEncryptGcmByKey(user.Mobile, secretKey)

		r.log.Infof("AesEncryptGcmByKey email: %v, ncrypt email: %v, err: %v", user.Email, email, errEmail)
		r.log.Infof("AesEncryptGcmByKey phone: %v, ncrypt phone: %v, err: %v", user.Mobile, phone, errPhone)
		entities = append(entities, &models.TbLasUser{
			TaskID:         taskId,
			ThirdCompanyID: thirdCompanyID,
			PlatformID:     platformID,
			Uid:            user.Userid,
			DefDid:         "-1",
			DefDidOrder:    0,
			Account:        account,
			NickName:       user.Nickname,
			Email:          email,
			Phone:          phone,
			Title:          user.Title,
			//Leader:         sql.NullString{String: strconv.FormatBool(account.Leader)},
			Source:           Source,
			Ctime:            time.Now(),
			Mtime:            time.Now(),
			CheckType:        1,
			EmploymentStatus: "active",
			EmploymentType:   "permanent",
			//Type:           sql.NullString{String: "dept", Valid: true},
		})
	}

	result := r.data.db.WithContext(ctx).Create(&entities)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			r.log.Error("user already exists")
		} else {
			return 0, result.Error
		}

	}

	return int(result.RowsAffected), nil
}

func (r *accounterRepo) SaveDepartments(ctx context.Context, depts []*dingtalk.DingtalkDept, taskId string) (int, error) {

	r.log.WithContext(ctx).Infof("SaveDepartments depts: %v, taskId :%s", depts, taskId)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	entities := make([]*models.TbLasDepartment, 0, len(depts))

	thirdCompanyID := r.serviceConf.Business.ThirdCompanyId
	platformID := r.serviceConf.Business.PlatformIds
	companyID := r.serviceConf.Business.CompanyId
	rootDep := &models.TbLasDepartment{
		Did:            "0",
		TaskID:         taskId,
		Name:           companyID,
		ThirdCompanyID: thirdCompanyID,
		PlatformID:     platformID,
		Pid:            "-1",
		Order:          0,
		Source:         "sync",
		Ctime:          time.Now(),
		Mtime:          time.Now(),
		CheckType:      1,
		//Type:           sql.NullString{String: "dept", Valid: true},
	}
	for _, dep := range depts {
		entities = append(entities, &models.TbLasDepartment{
			Did:            strconv.FormatInt(dep.DeptID, 10),
			TaskID:         taskId,
			Name:           dep.Name,
			ThirdCompanyID: thirdCompanyID,
			PlatformID:     platformID,
			Pid:            strconv.FormatInt(dep.ParentID, 10),
			Order:          int(dep.Order),
			Source:         "sync",
			Ctime:          time.Now(),
			Mtime:          time.Now(),
			CheckType:      1,
			//Type:           sql.NullString{String: "dept", Valid: true},
		})
	}
	entities = append(entities, rootDep)
	result := r.data.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "did"}, {Name: "task_id"}, {Name: "third_company_id"}, {Name: "platform_id"}},
		DoNothing: true,
	}).Create(&entities)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			r.log.Error("department already exists")
		} else {
			return 0, result.Error
		}

	}

	return int(result.RowsAffected), nil
}

func (r *accounterRepo) SaveDepartmentUserRelations(ctx context.Context, relations []*dingtalk.DingtalkDeptUserRelation, taskId string) (int, error) {
	r.log.WithContext(ctx).Infof("SaveDepartmentUserRelations relations: %v, taskId :%s", relations, taskId)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	entities := make([]*models.TbLasDepartmentUser, 0, len(relations))

	thirdCompanyID := r.serviceConf.Business.ThirdCompanyId
	platformID := r.serviceConf.Business.PlatformIds

	for _, relation := range relations {
		entities = append(entities, &models.TbLasDepartmentUser{
			Did:            relation.Did,
			TaskID:         taskId,
			ThirdCompanyID: thirdCompanyID,
			PlatformID:     platformID,
			Uid:            relation.Uid,
			Ctime:          time.Now(),
			Order:          relation.Order,
			CheckType:      1,
		})
	}

	result := r.data.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "did"}, {Name: "uid"}, {Name: "task_id"}, {Name: "third_company_id"}, {Name: "platform_id"}},
		DoNothing: true,
	}).Create(&entities)

	if result.Error != nil {
		return 0, result.Error
	}

	return int(result.RowsAffected), nil
}

func (r *accounterRepo) SaveCompanyCfg(ctx context.Context, cfg *dingtalk.DingtalkCompanyCfg) error {
	r.log.Infof("SaveCompanyCfg: %v", cfg)

	thirdCompanyID := r.serviceConf.Business.ThirdCompanyId
	platformID := r.serviceConf.Business.PlatformIds
	companyID := r.serviceConf.Business.CompanyId
	entity := &models.TbCompanyCfg{
		ThirdCompanyId: thirdCompanyID,
		PlatformIds:    platformID,
		CompanyId:      companyID,
		Status:         1,
		Ctime:          time.Now(),
		Mtime:          time.Now(),
	}

	err := r.data.db.WithContext(ctx).Where(models.TbCompanyCfg{
		ThirdCompanyId: thirdCompanyID,
		CompanyId:      companyID,
	}).FirstOrCreate(entity).Error

	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			r.log.Error("company config already exists")
		} else {
			return err
		}

	}

	return nil
}

// func (r *accounterRepo) CallEcisaccountsyncAll(ctx context.Context, taskId string) (biz.EcisaccountsyncAllResponse, error) {
// 	r.log.Infof("CallEcisaccountsyncAll: %v", taskId)

// 	path := r.serviceConf.Business.EcisaccountsyncUrl

// 	// path := "http://encs-pri-proxy-gateway/ecisaccountsync/api/sync/all"
// 	var resp biz.EcisaccountsyncAllResponse

// 	thirdCompanyID := r.serviceConf.Business.ThirdCompanyId

// 	collectCost := "1100000"
// 	uri := fmt.Sprintf("%s?taskId=%s&thirdCompanyId=%s&collectCost=%s", path, taskId, thirdCompanyID, collectCost)

// 	r.log.Infof("CallEcisaccountsyncAll uri: %s", uri)
// 	bs, err := httputil.PostJSON(uri, nil, time.Second*10)
// 	r.log.Infof("CallEcisaccountsyncAll.Post output: bs:%s, err:%w", string(bs), err)

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
	err := r.data.db.WithContext(ctx).Exec("truncate table tb_company_cfg").Error
	if err != nil {
		return err
	}
	err = r.data.db.WithContext(ctx).Exec("truncate table tb_las_department").Error
	if err != nil {
		return err
	}
	err = r.data.db.WithContext(ctx).Exec("truncate table tb_las_department_user").Error
	if err != nil {
		return err
	}
	err = r.data.db.WithContext(ctx).Exec("truncate table tb_las_account").Error
	if err != nil {
		return err
	}
	return nil
}

func (r *accounterRepo) SaveIncrementUsers(ctx context.Context, usersAdd, usersDel []*dingtalk.DingtalkDeptUser) error {

	r.log.WithContext(ctx).Infof("SaveIncrementUsers usersAdd: %v", usersAdd)
	r.log.WithContext(ctx).Infof("SaveIncrementUsers usersDel: %v", usersDel)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	entities := make([]*models.TbLasUserIncrement, 0, len(usersAdd)+len(usersDel))

	thirdCompanyID := r.serviceConf.Business.ThirdCompanyId
	platformID := r.serviceConf.Business.PlatformIds

	for _, user := range usersAdd {
		if user.Name == "" {
			user.Name = user.Userid
		}
		if user.Nickname == "" {
			user.Nickname = user.Name
		}
		var account string
		if user.Mobile != "" {
			account = user.Mobile
		} else {
			account = user.Userid
		}

		secretKey := r.serviceConf.Auth.Self.SecretKey

		email, errEmail := cipherutil.AesEncryptGcmByKey(user.Email, secretKey)

		phone, errPhone := cipherutil.AesEncryptGcmByKey(user.Mobile, secretKey)

		r.log.Infof("AesEncryptGcmByKey email: %v, ncrypt email: %v, err: %v", user.Email, email, errEmail)
		r.log.Infof("AesEncryptGcmByKey phone: %v, ncrypt phone: %v, err: %v", user.Mobile, phone, errPhone)

		entities = append(entities, &models.TbLasUserIncrement{
			ThirdCompanyID: thirdCompanyID,
			PlatformID:     platformID,
			Uid:            user.Userid,
			DefDid:         "-1",
			DefDidOrder:    0,
			Account:        account,
			NickName:       user.Nickname,
			Email:          email,
			Phone:          phone,
			Title:          user.Title,
			//Leader:         sql.NullString{String: strconv.FormatBool(account.Leader)},
			Source:           Source,
			Ctime:            time.Now(),
			Mtime:            time.Now(),
			EmploymentStatus: "active",
			EmploymentType:   "permanent",
			UpdateType:       "user_add",
			SyncType:         "auto",
			SyncTime:         time.Now(),
			Status:           0,
			//Type:           sql.NullString{String: "dept", Valid: true},
		})
	}
	for _, user := range usersDel {
		entities = append(entities, &models.TbLasUserIncrement{
			ThirdCompanyID: thirdCompanyID,
			PlatformID:     platformID,
			Uid:            user.Userid,
			DefDid:         "-1",
			DefDidOrder:    0,
			Account:        user.Userid,
			NickName:       user.Nickname,
			//Email:          "email",
			//Phone:          "phone",
			Title: user.Title,
			//Leader:         sql.NullString{String: strconv.FormatBool(account.Leader)},
			Source:           Source,
			Ctime:            time.Now(),
			Mtime:            time.Now(),
			EmploymentStatus: "active",
			UpdateType:       "user_del",
			SyncType:         "auto",
			SyncTime:         time.Now(),
			Status:           0,
			//Type:           sql.NullString{String: "dept", Valid: true},
		})
	}
	result := r.data.db.WithContext(ctx).Create(&entities)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			r.log.Error("user already exists")
		} else {
			return result.Error
		}

	}

	return nil
}
func (r *accounterRepo) SaveIncrementDepartments(ctx context.Context, deptsAdd, deptsDel []*dingtalk.DingtalkDept) error {

	r.log.WithContext(ctx).Infof("SaveIncrementDepartments deptsAdd: %v", deptsAdd)
	r.log.WithContext(ctx).Infof("SaveIncrementDepartments deptsDel: %v", deptsDel)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	entities := make([]*models.TbLasDepartmentIncrement, 0, len(deptsAdd)+len(deptsDel))

	thirdCompanyID := r.serviceConf.Business.ThirdCompanyId
	platformID := r.serviceConf.Business.PlatformIds

	for _, dep := range deptsDel {
		entities = append(entities, &models.TbLasDepartmentIncrement{
			Did:            strconv.FormatInt(dep.DeptID, 10),
			Name:           dep.Name,
			ThirdCompanyID: thirdCompanyID,
			PlatformID:     platformID,
			// 删除的时候这个父id, 怎么传
			Pid:        "",
			Order:      int32(dep.Order),
			Source:     "sync",
			Ctime:      time.Now(),
			Mtime:      time.Now(),
			UpdateType: "dept_del", //dept_move
			SyncTime:   time.Now(),
			SyncType:   "auto",
			Status:     0,
		})
	}

	for _, dep := range deptsAdd {
		entities = append(entities, &models.TbLasDepartmentIncrement{
			Did:            strconv.FormatInt(dep.DeptID, 10),
			Name:           dep.Name,
			ThirdCompanyID: thirdCompanyID,
			PlatformID:     platformID,
			Pid:            strconv.FormatInt(dep.ParentID, 10),
			Order:          int32(dep.Order),
			Source:         "sync",
			Ctime:          time.Now(),
			Mtime:          time.Now(),
			UpdateType:     "dept_add",
			SyncTime:       time.Now(),
			SyncType:       "auto",
			Status:         0,
		})
	}
	result := r.data.db.WithContext(ctx).Create(&entities)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			r.log.Error("user already exists")
		} else {
			return result.Error
		}

	}

	return nil
}

func (r *accounterRepo) SaveIncrementDepartmentUserRelations(ctx context.Context, relationsAdd, relationsDel []*dingtalk.DingtalkDeptUserRelation) error {

	r.log.WithContext(ctx).Infof("SaveIncrementDepartmentUserRelations relationsAdd: %v", relationsAdd)
	r.log.WithContext(ctx).Infof("SaveIncrementDepartmentUserRelations relationsDel: %v", relationsDel)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	entities := make([]*models.TbLasDepartmentUserIncrement, 0, len(relationsAdd)+len(relationsDel))

	thirdCompanyID := r.serviceConf.Business.ThirdCompanyId
	platformID := r.serviceConf.Business.PlatformIds

	for _, relation := range relationsAdd {
		entities = append(entities, &models.TbLasDepartmentUserIncrement{
			Did:            relation.Did,
			ThirdCompanyID: thirdCompanyID,
			PlatformID:     platformID,
			Uid:            relation.Uid,
			Ctime:          time.Now(),
			Order:          1,
			UpdateType:     "user_dept_add",
			SyncType:       "auto",
			SyncTime:       time.Now(),
			Status:         0,
		})
	}

	for _, relation := range relationsDel {
		entities = append(entities, &models.TbLasDepartmentUserIncrement{
			Did:            relation.Did,
			ThirdCompanyID: thirdCompanyID,
			PlatformID:     platformID,
			Uid:            relation.Uid,
			Ctime:          time.Now(),
			Order:          1,
			UpdateType:     "user_dept_del",
			SyncType:       "auto",
			SyncTime:       time.Now(),
			Status:         0,
		})
	}
	result := r.data.db.WithContext(ctx).Create(&entities)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			r.log.Infof("relation already exists: %v")
		} else {
			return result.Error
		}

	}

	return nil
}
func (r *accounterRepo) BatchSaveUsers(ctx context.Context, users []*models.TbLasUser) (int, error) {

	r.log.WithContext(ctx).Infof("BatchSaveUsers users: %v", users)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if len(users) == 0 {
		r.log.Infof("users is empty")
		return 0, nil
	}

	// for _, user := range users {
	// 	r.log.Warn(user)
	// }

	result := r.data.nancalDB.WithContext(ctx).Create(users)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			r.log.Error("user already exists")
		} else {
			return 0, result.Error
		}
	}

	return int(result.RowsAffected), nil
}
func (r *accounterRepo) BatchSaveDepts(ctx context.Context, depts []*models.TbLasDepartment) (int, error) {

	r.log.WithContext(ctx).Infof("BatchSaveDepts depts: %v", depts)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if len(depts) == 0 {
		r.log.Warn("users is empty")
		return 0, nil
	}
	// for _, dept := range depts {
	// 	r.log.Info(dept)
	// }

	result := r.data.nancalDB.WithContext(ctx).Create(depts)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			r.log.Error("dept already exists")
		} else {
			return 0, result.Error
		}

	}
	return int(result.RowsAffected), nil
}
func (r *accounterRepo) BatchSaveDeptUsers(ctx context.Context, usersdepts []*models.TbLasDepartmentUser) (int, error) {

	r.log.WithContext(ctx).Infof("BatchSaveDeptUsers usersdepts: %v", usersdepts)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if len(usersdepts) == 0 {
		r.log.Warn("users is empty")
		return 0, nil
	}
	// for _, userdept := range usersdepts {
	// 	r.log.Info(userdept)
	// }
	result := r.data.nancalDB.WithContext(ctx).Create(usersdepts)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			r.log.Error("deptuser already exists")
		} else {
			return 0, result.Error
		}

	}
	return int(result.RowsAffected), nil
}

func (r *accounterRepo) CreateTask(ctx context.Context, taskName string) (int, error) {
	log := r.log.WithContext(ctx)
	log.Infof("CreateTask name: %s", taskName)

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

	result := r.data.nancalDB.WithContext(ctx).Where("title=?", taskName).FirstOrCreate(task)

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

	log := r.log.WithContext(ctx)
	log.Infof("UpdateTask taskName: %s, status: %s", taskName, status)
	// pending/in_progress/completed/cancelled

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var task models.Task
	if err := r.data.nancalDB.Model(&models.Task{}).WithContext(ctx).Where("title=?", taskName).Find(&task).Error; err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Error("查询超时")
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

	result := r.data.db.WithContext(ctx).Updates(task)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			log.Infof("task already exists")
		} else {
			return result.Error
		}

	}

	return nil
}

func (r *accounterRepo) GetTask(ctx context.Context, taskName string) (*models.Task, error) {

	r.log.WithContext(ctx).Infof("CreateTask name: %s", taskName)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	task := &models.Task{}
	result := r.data.nancalDB.WithContext(ctx).Where("title=?", taskName).Find(task)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, errors.New("notfound")
	}
	return task, nil
}
