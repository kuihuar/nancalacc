package biz

import (
	"context"
	"errors"
	"fmt"
	v1 "nancalacc/api/account/v1"
	"nancalacc/internal/auth"
	"nancalacc/internal/conf"
	"nancalacc/internal/data/models"
	"nancalacc/internal/dingtalk"
	"nancalacc/internal/wps"
	"strconv"
	"time"

	//"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/xuri/excelize/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
// ErrUserNotFound is user not found.
// ErrUserNotFound = errors.NotFound(v1.ErrorReason_USER_NOT_FOUND.String(), "user not found")
)

// Greeter is a Greeter model.
type Accounter struct {
	Hello string
}

// GreeterRepo is a Greater repo.
type AccounterRepo interface {
	SaveUsers(ctx context.Context, users []*dingtalk.DingtalkDeptUser, taskId string) (int, error)
	SaveDepartments(ctx context.Context, depts []*dingtalk.DingtalkDept, taskId string) (int, error)
	SaveDepartmentUserRelations(ctx context.Context, relations []*dingtalk.DingtalkDeptUserRelation, taskId string) (int, error)
	SaveCompanyCfg(ctx context.Context, cfg *dingtalk.DingtalkCompanyCfg) error

	ClearAll(ctx context.Context) error

	SaveIncrementDepartments(ctx context.Context, deptsAdd, deptsDel, deptsUpd []*dingtalk.DingtalkDept) error

	SaveIncrementUsers(ctx context.Context, usersAdd, usersDel, usersUpd []*dingtalk.DingtalkDeptUser) error
	SaveIncrementDepartmentUserRelations(ctx context.Context, relationsAdd, relationsDel, relationsUpd []*dingtalk.DingtalkDeptUserRelation) error

	BatchSaveUsers(ctx context.Context, users []*models.TbLasUser) (int, error)
	BatchSaveDepts(ctx context.Context, depts []*models.TbLasDepartment) (int, error)
	BatchSaveDeptUsers(ctx context.Context, deptusers []*models.TbLasDepartmentUser) (int, error)

	CreateTask(ctx context.Context, taskName string) (int, error)
	UpdateTask(ctx context.Context, taskName, status string) error

	GetTask(ctx context.Context, taskName string) (*models.Task, error)
}

// GreeterUsecase is a Greeter usecase.
type AccounterUsecase struct {
	repo         AccounterRepo
	dingTalkRepo dingtalk.Dingtalk
	appAuth      auth.Authenticator
	wpsSync      wps.WpsSync
	wps          wps.Wps
	bizConf      *conf.Service_Business
	localCache   CacheService
	log          *log.Helper
}

var (
	prefix = "nancalacc:cache:"
)

// NewGreeterUsecase new a Greeter usecase.
func NewAccounterUsecase(repo AccounterRepo, dingTalkRepo dingtalk.Dingtalk, appAuth auth.Authenticator, wpsSync wps.WpsSync, wps wps.Wps, bizConf *conf.Service_Business, cache CacheService, logger log.Logger) *AccounterUsecase {
	return &AccounterUsecase{repo: repo, dingTalkRepo: dingTalkRepo, appAuth: appAuth, wpsSync: wpsSync, wps: wps, bizConf: bizConf, localCache: cache, log: log.NewHelper(logger)}
}

func (uc *AccounterUsecase) CreateSyncAccount(ctx context.Context, req *v1.CreateSyncAccountRequest) (*v1.CreateSyncAccountReply, error) {
	// return &v1.CreateSyncAccountReply{
	// 	TaskId:     "taskId",
	// 	CreateTime: timestamppb.Now(),
	// }, nil
	log := uc.log.WithContext(ctx)
	log.Infof("CreateSyncAccount: %v", req)

	taskId := req.GetTaskName()

	taskCachekey := prefix + taskId

	_, ok, err := uc.localCache.Get(ctx, taskCachekey)
	if err != nil {
		return nil, err
	}
	if ok {
		return nil, status.Error(codes.AlreadyExists, "task name "+taskId+" exists")
	}
	// num, err := uc.repo.CreateTask(ctx, taskId)
	// if err != nil {
	// 	return nil, err
	// }
	// if num == 0 {
	// 	return nil, status.Error(codes.AlreadyExists, "taskId  exists")
	// }

	uc.log.WithContext(ctx).Info("CreateSyncAccount.SaveCompanyCfg")
	err = uc.repo.SaveCompanyCfg(ctx, &dingtalk.DingtalkCompanyCfg{})
	log.Infof("CreateSyncAccount.SaveCompanyCfg: err: %v", err)
	if err != nil {
		return nil, err
	}
	log.Infof("CreateSyncAccount.GetAccessToken")

	// 1. 获取access_token
	accessToken, err := uc.dingTalkRepo.GetAccessToken(ctx, "code")
	uc.log.WithContext(ctx).Infof("CreateSyncAccount.GetAccessToken: accessToken: %v, err: %v", accessToken, err)
	if err != nil {
		return nil, err
	}

	//taskId := time.Now().Add(time.Duration(1) * time.Second).Format("20060102150405")

	// 1. 从第三方获取部门和用户数据

	log.Infof("CreateSyncAccount.FetchDepartments")

	depts, err := uc.dingTalkRepo.FetchDepartments(ctx, accessToken)
	log.Infof("CreateSyncAccount.FetchDepartments: depts: %+v, err: %v", depts, err)
	if err != nil {
		return nil, err
	}
	for _, dept := range depts {
		uc.log.WithContext(ctx).Infof("biz.CreateSyncAccount: dept: %+v", dept)
	}

	log.Infof("CreateSyncAccount.SaveDepartments depts: %v, taskId: %v", depts, taskId)
	// 2. 数据入库
	deptCount, err := uc.repo.SaveDepartments(ctx, depts, taskId)
	log.Infof("CreateSyncAccount.SaveDepartments: deptCount: %v, err: %v", deptCount, err)
	if err != nil {
		return nil, err
	}
	var deptIds []int64
	for _, dept := range depts {
		deptIds = append(deptIds, dept.DeptID)
	}

	log.Infof("CreateSyncAccount.FetchDepartmentUsers accessToken: %v deptIds: %v", accessToken, deptIds)
	// 1. 从第三方获取用户数据
	deptUsers, err := uc.dingTalkRepo.FetchDepartmentUsers(ctx, accessToken, deptIds)
	log.Infof("CreateSyncAccount.FetchDepartmentUsers deptUsers: %v, err: %v", deptUsers, err)
	if err != nil {
		return nil, err
	}
	// 2. 数据入库
	//这里可以 将deptUsers转为model.TbLasUser,
	// SaveUsers(ctx, TbLasUser)
	log.Infof("CreateSyncAccount.SaveUsers deptUsers: %v, taskId: %v", deptUsers, taskId)
	userCount, err := uc.repo.SaveUsers(ctx, deptUsers, taskId)
	log.Infof("CreateSyncAccount.SaveUsers userCount: %v, err: %v", userCount, err)
	if err != nil {
		return nil, err
	}

	// 2. 关系数据入库
	var deptUserRelations []*dingtalk.DingtalkDeptUserRelation
	for _, deptUser := range deptUsers {
		order := make(map[int64]int64, 0)
		if len(deptUser.DeptOrderList) > 0 {
			for _, depIdOrder := range deptUser.DeptOrderList {
				order[depIdOrder.DeptID] = depIdOrder.DeptID
			}
		}
		for _, depId := range deptUser.DeptIDList {

			reliation := &dingtalk.DingtalkDeptUserRelation{
				Uid: deptUser.Userid,
				Did: strconv.FormatInt(depId, 10),
			}
			if order, ok := order[depId]; ok {
				reliation.Order = order
			}
			deptUserRelations = append(deptUserRelations, reliation)
		}

	}
	log.Infof("CreateSyncAccount.SaveDepartmentUserRelations deptUserRelations: %v, taskId: %v", deptUserRelations, taskId)
	// 3. 数据入库
	relationCount, err := uc.repo.SaveDepartmentUserRelations(ctx, deptUserRelations, taskId)
	uc.log.WithContext(ctx).Infof("CreateSyncAccount.SaveDepartmentUserRelations relationCount: %v, err: %v", relationCount, err)
	if err != nil {
		return nil, err
	}
	log.Infof("CreateSyncAccount.CallEcisaccountsyncAll taskId: %v", taskId)

	appAccessToken, err := uc.appAuth.GetAccessToken(ctx)
	if err != nil {
		return nil, err
	}
	log.Infof("appAccessToken", appAccessToken)

	res, err := uc.wpsSync.PostEcisaccountsyncAll(ctx, appAccessToken.AccessToken, &wps.EcisaccountsyncAllRequest{
		TaskId:         taskId,
		ThirdCompanyId: uc.bizConf.ThirdCompanyId,
	})
	log.Infof("CreateSyncAccount.CallEcisaccountsyncAll res: %v, err: %v", res, err)

	if err != nil {
		return nil, err
	}
	taskInfo := &models.Task{
		ID:          1,
		Title:       req.GetTaskName(),
		Description: req.GetTaskName(),
		Status:      "in_progress",
		CreatorID:   1,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		DueDate:     time.Now(),
		StartDate:   time.Now(),
		Progress:    30,
		ActualTime:  0,
	}
	uc.localCache.Set(ctx, taskCachekey, taskInfo, 300*time.Minute)
	return &v1.CreateSyncAccountReply{
		TaskId:     taskId,
		CreateTime: timestamppb.Now(),
	}, nil
}

func (uc *AccounterUsecase) GetSyncAccount(ctx context.Context, req *v1.GetSyncAccountRequest) (*v1.GetSyncAccountReply, error) {
	uc.log.WithContext(ctx).Infof("GetSyncAccount: %v", req)

	taskId := req.GetTaskId()

	taskCachekey := prefix + taskId

	taskCacheInfo, ok, err := uc.localCache.Get(ctx, taskCachekey)
	if err != nil {
		return nil, err
	}
	if ok {
		taskInfo, ok1 := taskCacheInfo.(*models.Task)
		if ok1 {
			return &v1.GetSyncAccountReply{
				Status:                      v1.GetSyncAccountReply_Status(taskInfo.Progress),
				UserCount:                   1,
				DepartmentCount:             1,
				UserDepartmentRelationCount: 1,
			}, nil
		}

	}
	return nil, status.Error(codes.NotFound, "task "+taskId+" not found")
}

func (uc *AccounterUsecase) CreateTask(ctx context.Context, taskName string) (int, error) {
	uc.log.WithContext(ctx).Infof("CreateTask taskName: %s", taskName)
	return uc.repo.CreateTask(ctx, taskName)

}
func (uc *AccounterUsecase) GetTask(ctx context.Context, taskName string) (*v1.GetTaskReply_Task, error) {
	uc.log.WithContext(ctx).Infof("GetTask taskName: %s", taskName)

	taskInfo, err := uc.GetCacheTask(ctx, taskName)
	if err != nil {
		return nil, err
	}

	// taskInfo := &models.Task{
	// 	ID:          1,
	// 	Title:       taskName,
	// 	Description: "desc1",
	// 	Status:      "in_progress",
	// 	CreatorID:   1,
	// 	CreatedAt:   time.Now(),
	// 	UpdatedAt:   time.Now(),
	// 	DueDate:     time.Now(),
	// 	StartDate:   time.Now(),
	// 	Progress:    30,
	// 	ActualTime:  0,
	// }
	// taskInfoJson, _ := json.Marshal(taskInfo)
	// appAccessToken, err := uc.appAuth.GetAccessToken(ctx)
	// if err != nil {
	// 	return nil, err
	// }
	// err = uc.wps.CacheSet(ctx, appAccessToken.AccessToken, taskName, string(taskInfoJson), 24*time.Hour)

	// if err != nil {
	// 	return nil, err
	// }

	// taskInfoJsonRes, err := uc.wps.CacheGet(ctx, appAccessToken.AccessToken, taskName)
	// if err != nil {
	// 	return nil, err
	// }
	// uc.log.Infof("taskInfoJsonRes: %v", taskInfoJsonRes)

	// err = uc.wps.CacheDel(ctx, appAccessToken.AccessToken, taskName)
	// if err != nil {
	// 	return nil, err
	// }

	return &v1.GetTaskReply_Task{
		Name:          taskInfo.Title,
		Status:        taskInfo.Status,
		CreateTime:    timestamppb.New(taskInfo.CreatedAt),
		StartTime:     timestamppb.New(taskInfo.StartDate),
		CompletedTime: timestamppb.New(taskInfo.CompletedAt),
		ActurlTime:    int32(taskInfo.ActualTime),
	}, nil

}
func (uc *AccounterUsecase) UpdateTask(ctx context.Context, taskName, status string) error {
	uc.log.WithContext(ctx).Infof("UpdateTask taskId: %s, status %s", taskName, status)
	return uc.repo.UpdateTask(ctx, taskName, status)

}
func (uc *AccounterUsecase) ParseExecell(ctx context.Context, taskId, filename string) (err error) {

	defer func() {
		if err != nil {
			uc.UpdateCacheTask(ctx, taskId, models.TaskStatusCancelled)
		} else {
			uc.UpdateCacheTask(ctx, taskId, models.TaskStatusCompleted)
		}
	}()
	log := uc.log.WithContext(ctx)
	log.Infof("ParseExecell taskId: %s,filename:%s", taskId, filename)

	f, err := excelize.OpenFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	processSheet := map[string]bool{
		"user":            true,
		"department":      true,
		"department_user": true,
	}
	sheets := f.GetSheetList()
	for _, sheet := range sheets {
		if _, ok := processSheet[sheet]; !ok {
			fmt.Printf("sheetname: %s\n", sheet)
			continue
		}
		rows, err := f.Rows(sheet)
		if err != nil {
			return fmt.Errorf("err: %w", err)
		}
		defer rows.Close()
		rows.Next()
		switch sheet {
		case "user":
			uc.transUser(ctx, taskId, rows)
		case "department":
			uc.transDept(ctx, taskId, rows)
		case "department_user":
			uc.transUserDept(ctx, taskId, rows)
		default:
			log.Infof("not found sheetname: %s\n", sheet)
		}

	}
	// err = uc.repo.UpdateTask(ctx, taskId, models.TaskStatusCompleted)
	// if err != nil {
	// 	return err
	// }
	appAccessToken, err := uc.appAuth.GetAccessToken(ctx)
	if err != nil {
		return err
	}

	_, err = uc.wpsSync.PostEcisaccountsyncAll(ctx, appAccessToken.AccessToken, &wps.EcisaccountsyncAllRequest{
		TaskId:         taskId,
		ThirdCompanyId: uc.bizConf.ThirdCompanyId,
	})
	return err
}

func (uc *AccounterUsecase) transUser(ctx context.Context, taskId string, rows *excelize.Rows) (err error) {

	log := uc.log.WithContext(ctx)
	log.Infof("transUser taskId: %s", taskId)

	//uc.repo.UpdateTask(ctx, taskId, models.TaskStatusInProgress)
	thirdCompanyId := uc.bizConf.ThirdCompanyId
	platformIds := uc.bizConf.PlatformIds
	users := make([]*models.TbLasUser, 0, 100)
	now := time.Now()
	for rows.Next() {
		row, err := rows.Columns()
		if err != nil {
			return fmt.Errorf("err: %w", err)
		}
		log.Info(row)
		if len(row) < 3 {
			log.Warnf("row len < 3: %v", row)
			continue
		}

		users = append(users, &models.TbLasUser{
			TaskID:           taskId,
			ThirdCompanyID:   thirdCompanyId,
			PlatformID:       platformIds,
			Uid:              row[0],
			Account:          row[1],
			NickName:         row[2],
			EmploymentStatus: "active",
			Source:           "sync",
			Ctime:            now,
			Mtime:            now,
			CheckType:        1,
		})
		if len(users) >= 100 {
			if _, err := uc.repo.BatchSaveUsers(ctx, users); err != nil {
				return err
			}
			users = users[:0] // 清空切片（保留底层数组，避免重新分配）
		}
		// num := SheetDataToModel(sheet, row)
	}
	if len(users) > 0 {
		if _, err := uc.repo.BatchSaveUsers(ctx, users); err != nil {
			return err
		}
	}

	if err := rows.Error(); err != nil {
		return fmt.Errorf("err: %w", err)
	}
	return nil
}
func (uc *AccounterUsecase) transDept(ctx context.Context, taskId string, rows *excelize.Rows) (err error) {
	log := uc.log.WithContext(ctx)
	log.Infof("transDept taskId: %s", taskId)

	//uc.repo.UpdateTask(ctx, taskId, models.TaskStatusInProgress)
	thirdCompanyId := uc.bizConf.ThirdCompanyId
	platformIds := uc.bizConf.PlatformIds
	depts := make([]*models.TbLasDepartment, 0, 100)
	now := time.Now()
	for rows.Next() {
		row, err := rows.Columns()
		if err != nil {
			return fmt.Errorf("err: %w", err)
		}

		log.Info(row)
		if len(row) < 3 {
			log.Warnf("row len < 3: %v", row)
			continue
		}

		//log.Info(row)

		depts = append(depts, &models.TbLasDepartment{
			TaskID:         taskId,
			ThirdCompanyID: thirdCompanyId,
			PlatformID:     platformIds,
			Did:            row[0],
			Pid:            row[1],
			Name:           row[2],
			// Order:          row[3],
			Source:    "sync",
			Ctime:     now,
			Mtime:     now,
			CheckType: 1,
		})
		if len(depts) >= 100 {
			if _, err := uc.repo.BatchSaveDepts(ctx, depts); err != nil {
				return err
			}
			depts = depts[:0] // 清空切片
		}
	}
	if len(depts) > 0 {
		if _, err := uc.repo.BatchSaveDepts(ctx, depts); err != nil {
			return err
		}
	}

	if err := rows.Error(); err != nil {
		return fmt.Errorf("err: %w", err)
	}
	return nil
}
func (uc *AccounterUsecase) transUserDept(ctx context.Context, taskId string, rows *excelize.Rows) (err error) {
	log := uc.log.WithContext(ctx)
	log.Infof("transUserDept taskId: %s", taskId)

	//uc.repo.UpdateTask(ctx, taskId, models.TaskStatusInProgress)
	thirdCompanyId := uc.bizConf.ThirdCompanyId
	platformIds := uc.bizConf.PlatformIds
	deptusers := make([]*models.TbLasDepartmentUser, 0, 100)
	now := time.Now()
	for rows.Next() {
		row, err := rows.Columns()
		if err != nil {
			return fmt.Errorf("err: %w", err)
		}
		log.Info(row)
		if len(row) < 2 {
			log.Warnf("row len < 2: %v", row)
			continue
		}

		deptusers = append(deptusers, &models.TbLasDepartmentUser{
			TaskID:         taskId,
			ThirdCompanyID: thirdCompanyId,
			PlatformID:     platformIds,
			Uid:            row[0],
			Did:            row[1],
			Ctime:          now,
			CheckType:      1,
		})
		if len(deptusers) >= 100 {
			if _, err := uc.repo.BatchSaveDeptUsers(ctx, deptusers); err != nil {
				return err
			}
			deptusers = deptusers[:0] // 清空切片（保留底层数组，避免重新分配）
		}
	}
	if len(deptusers) > 0 {
		if _, err := uc.repo.BatchSaveDeptUsers(ctx, deptusers); err != nil {
			return err
		}
	}

	if err := rows.Error(); err != nil {
		return fmt.Errorf("err: %w", err)
	}
	return nil
}

func (uc *AccounterUsecase) CreateCacheTask(ctx context.Context, taskName, status string) error {

	cacheKey := prefix + taskName
	task := &models.Task{
		Title:         taskName,
		Description:   taskName,
		CreatedAt:     time.Now(),
		Status:        models.TaskStatusInProgress,
		Progress:      0,
		StartDate:     time.Now(),
		DueDate:       time.Now().Add(time.Minute * 30),
		CompletedAt:   time.Now(),
		CreatorID:     99,
		EstimatedTime: 10,
		ActualTime:    0,
	}
	return uc.localCache.Set(ctx, cacheKey, task, 300*time.Minute)
}
func (uc *AccounterUsecase) UpdateCacheTask(ctx context.Context, taskName, status string) error {

	cacheKey := prefix + taskName
	oldTask, ok, err := uc.localCache.Get(ctx, cacheKey)
	if err != nil {
		return err
	}
	var task *models.Task
	var startDate time.Time
	now := time.Now()
	if ok {
		task, ok1 := oldTask.(*models.Task)
		if ok1 {
			startDate = task.StartDate
			task.ActualTime = int(now.Sub(startDate).Seconds()) + 20
			task.Status = status
			task.Progress = 100
			task.CompletedAt = now
			task.UpdatedAt = now
		}
	}

	if task == nil {
		task = &models.Task{
			Title:         taskName,
			Description:   taskName,
			Status:        status,
			Progress:      100,
			StartDate:     time.Now(),
			DueDate:       time.Now().Add(time.Minute * 30),
			CompletedAt:   time.Now(),
			CreatorID:     99,
			EstimatedTime: 10,
			ActualTime:    0,
		}
	}
	return uc.localCache.Set(ctx, cacheKey, task, 300*time.Minute)
}

func (uc *AccounterUsecase) GetCacheTask(ctx context.Context, taskName string) (*models.Task, error) {

	cacheKey := prefix + taskName
	var task *models.Task
	taskInfo, ok, err := uc.localCache.Get(ctx, cacheKey)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("notfound")
	}
	task, ok1 := taskInfo.(*models.Task)
	if !ok1 {
		return nil, errors.New("type error")
	}
	return task, nil

}

// 这个方法是把全量数据执插入表后，可以自已调用原生API去同步
func (uc *AccounterUsecase) ParseExecellAfter(ctx context.Context, taskId, filename string) (err error) {

	// 1. 创建部门
	uc.wps.PostCreateDept(ctx, taskId, wps.PostCreateDeptRequest{
		Name: "测试部门",
	})

	// 2. 创建部门存在的 > 更新部门
	uc.wps.PostUpdateDept(ctx, taskId, wps.PostUpdateDeptRequest{
		ExDeptID: "10000000000000000000000000000000",
		Name:     "测试部门",
	})

	// 3. 全量的的 减去 创建的， 再减去更新的， 就是删除的
	uc.wps.PostBatchDeleteDept(ctx, taskId, wps.PostBatchDeleteDeptRequest{
		DeptIDs: []string{"10000000000000000000000000000000"},
	})

	// 这里的用户包含所在部门
	uc.wps.PostCreateUser(ctx, taskId, wps.PostCreateUserRequest{
		UserName: "测试用户",
	})

	// 2. 更新用户
	uc.wps.PostUpdateUser(ctx, taskId, wps.PostUpdateUserRequest{
		ExUserID: "10000000000000000000000000000000",
		UserName: "测试部门",
	})

	// 3. 更新用户
	uc.wps.PostUpdateUser(ctx, taskId, wps.PostUpdateUserRequest{
		ExUserID: "10000000000000000000000000000000",
		UserName: "测试用户",
	})

	// 全量用户

	uc.wps.PostBatchUserByPage(ctx, taskId, wps.PostBatchUserByPageRequest{
		PageSize: 100,
		PageNum:  1,
	})

	// 5. 全量的的 减去 创建的， 再减去更新的， 就是删除的
	uc.wps.PostBatchDeleteUser(ctx, taskId, wps.PostBatchDeleteUserRequest{
		UserIDs: []string{"10000000000000000000000000000000"},
	})
	return nil
}
