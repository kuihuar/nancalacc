package biz

import (
	"context"
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

	SaveIncrementDepartments(ctx context.Context, deptsAdd, deptsDel []*dingtalk.DingtalkDept) error
	SaveIncrementUsers(ctx context.Context, usersAdd, usersDel []*dingtalk.DingtalkDeptUser) error
	SaveIncrementDepartmentUserRelations(ctx context.Context, relationsAdd, relationsDel []*dingtalk.DingtalkDeptUserRelation) error

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
	redisRepo    RedisCacheRepo
	log          *log.Helper
}

// NewGreeterUsecase new a Greeter usecase.
func NewAccounterUsecase(repo AccounterRepo, dingTalkRepo dingtalk.Dingtalk, appAuth auth.Authenticator, wpsSync wps.WpsSync, wps wps.Wps, bizConf *conf.Service_Business, redisRepo RedisCacheRepo, logger log.Logger) *AccounterUsecase {
	return &AccounterUsecase{repo: repo, dingTalkRepo: dingTalkRepo, appAuth: appAuth, wpsSync: wpsSync, wps: wps, bizConf: bizConf, redisRepo: redisRepo, log: log.NewHelper(logger)}
}

func (uc *AccounterUsecase) CreateSyncAccount(ctx context.Context, req *v1.CreateSyncAccountRequest) (*v1.CreateSyncAccountReply, error) {
	// return &v1.CreateSyncAccountReply{
	// 	TaskId:     "taskId",
	// 	CreateTime: timestamppb.Now(),
	// }, nil
	log := uc.log.WithContext(ctx)
	log.Infof("CreateSyncAccount: %v", req)

	taskId := req.GetTaskName()

	num, err := uc.repo.CreateTask(ctx, taskId)
	if err != nil {
		return nil, err
	}
	if num == 0 {
		return nil, status.Error(codes.AlreadyExists, "taskId  exists")
	}

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
		order := int(deptUser.DeptOrder)
		if order > 0 {
			order = 1
		} else {
			order = 0
		}
		for _, depId := range deptUser.DeptIDList {

			deptUserRelations = append(deptUserRelations, &dingtalk.DingtalkDeptUserRelation{
				Uid:   deptUser.Userid,
				Did:   strconv.FormatInt(depId, 10),
				Order: order,
			})
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
	return &v1.CreateSyncAccountReply{
		TaskId:     taskId,
		CreateTime: timestamppb.Now(),
	}, nil
}

func (uc *AccounterUsecase) GetSyncAccount(ctx context.Context, req *v1.GetSyncAccountRequest) (*v1.GetSyncAccountReply, error) {
	uc.log.WithContext(ctx).Infof("GetSyncAccount: %v", req)

	prefix := "nancalacc:cache:"
	status := v1.GetSyncAccountReply_Status(v1.GetSyncAccountReply_SUCCESS)
	key1 := prefix + "taskId:" + req.TaskId
	uc.log.Debugf("key1: %s", key1)
	uc.redisRepo.Set(ctx, key1, status, 50*time.Minute)
	var taskStatus v1.GetSyncAccountReply_Status
	uc.redisRepo.Get(ctx, key1, &taskStatus)

	var userCount int64
	key2 := key1 + ":userCount"
	uc.log.Debugf("key2: %s", key2)
	uc.redisRepo.Get(ctx, key2, &userCount)
	uc.log.WithContext(ctx).Infof("GetSyncAccount: taskStatus: %v", taskStatus)
	return &v1.GetSyncAccountReply{
		Status:                      taskStatus,
		UserCount:                   userCount,
		DepartmentCount:             1,
		UserDepartmentRelationCount: 1,
	}, nil
}

func (uc *AccounterUsecase) CreateTask(ctx context.Context, taskName string) (int, error) {
	uc.log.WithContext(ctx).Infof("CreateTask taskName: %s", taskName)
	return uc.repo.CreateTask(ctx, taskName)

}
func (uc *AccounterUsecase) GetTask(ctx context.Context, taskName string) (*v1.GetTaskReply_Task, error) {
	uc.log.WithContext(ctx).Infof("GetTask taskName: %s", taskName)

	taskInfo, err := uc.repo.GetTask(ctx, taskName)
	if err != nil {
		return nil, err
	}

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
func (uc *AccounterUsecase) ParseExecell(ctx context.Context, taskId, filename string) error {

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
		"tb_las_user":            true,
		"tb_las_department":      true,
		"tb_las_department_user": true,
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
		case "tb_las_user":
			uc.transUser(ctx, taskId, rows)
		case "tb_las_department":
			uc.transDept(ctx, taskId, rows)
		case "tb_las_department_user":
			uc.transUserDept(ctx, taskId, rows)
		default:
			log.Infof("not found sheetname: %s\n", sheet)
		}

	}
	err = uc.repo.UpdateTask(ctx, taskId, models.TaskStatusCompleted)
	if err != nil {
		return err
	}
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

	uc.repo.UpdateTask(ctx, taskId, models.TaskStatusInProgress)
	thirdCompanyId := uc.bizConf.ThirdCompanyId
	platformIds := uc.bizConf.PlatformIds
	users := make([]*models.TbLasUser, 0, 100)
	now := time.Now()
	for rows.Next() {
		row, err := rows.Columns()
		if err != nil {
			return fmt.Errorf("err: %w", err)
		}
		//log.Info(row)

		users = append(users, &models.TbLasUser{
			TaskID:           taskId,
			ThirdCompanyID:   thirdCompanyId,
			PlatformID:       platformIds,
			Uid:              row[4],
			Account:          row[7],
			NickName:         row[8],
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

	uc.repo.UpdateTask(ctx, taskId, models.TaskStatusInProgress)
	thirdCompanyId := uc.bizConf.ThirdCompanyId
	platformIds := uc.bizConf.PlatformIds
	depts := make([]*models.TbLasDepartment, 0, 100)
	now := time.Now()
	for rows.Next() {
		row, err := rows.Columns()
		if err != nil {
			return fmt.Errorf("err: %w", err)
		}
		//log.Info(row)

		depts = append(depts, &models.TbLasDepartment{
			TaskID:         taskId,
			ThirdCompanyID: thirdCompanyId,
			PlatformID:     platformIds,
			Did:            row[1],
			Pid:            row[5],
			Name:           row[6],
			//Order:          row[7],
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

	uc.repo.UpdateTask(ctx, taskId, models.TaskStatusInProgress)
	thirdCompanyId := uc.bizConf.ThirdCompanyId
	platformIds := uc.bizConf.PlatformIds
	deptusers := make([]*models.TbLasDepartmentUser, 0, 100)
	now := time.Now()
	for rows.Next() {
		row, err := rows.Columns()
		if err != nil {
			return fmt.Errorf("err: %w", err)
		}
		//log.Info(row)

		deptusers = append(deptusers, &models.TbLasDepartmentUser{
			TaskID:         taskId,
			ThirdCompanyID: thirdCompanyId,
			PlatformID:     platformIds,
			Uid:            row[4],
			Did:            row[5],
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
