package biz

import (
	"context"
	"fmt"
	v1 "nancalacc/api/account/v1"
	"nancalacc/internal/auth"
	"nancalacc/internal/conf"
	"nancalacc/internal/dingtalk"
	"nancalacc/internal/wps"
	"strconv"
	"time"

	//"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
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
}

// GreeterUsecase is a Greeter usecase.
type AccounterUsecase struct {
	repo         AccounterRepo
	dingTalkRepo dingtalk.Dingtalk
	appAuth      auth.Authenticator
	wpsSync      wps.WpsSync
	wps          wps.Wps
	bizConf      *conf.Service_Business
	log          *log.Helper
}

// NewGreeterUsecase new a Greeter usecase.
func NewAccounterUsecase(repo AccounterRepo, dingTalkRepo dingtalk.Dingtalk, appAuth auth.Authenticator, wpsSync wps.WpsSync, wps wps.Wps, bizConf *conf.Service_Business, logger log.Logger) *AccounterUsecase {
	return &AccounterUsecase{repo: repo, dingTalkRepo: dingTalkRepo, appAuth: appAuth, wpsSync: wpsSync, wps: wps, bizConf: bizConf, log: log.NewHelper(logger)}
}

func (uc *AccounterUsecase) CreateSyncAccount(ctx context.Context, req *v1.CreateSyncAccountRequest) (*v1.CreateSyncAccountReply, error) {
	// return &v1.CreateSyncAccountReply{
	// 	TaskId:     "taskId",
	// 	CreateTime: timestamppb.Now(),
	// }, nil
	uc.log.WithContext(ctx).Infof("CreateSyncAccount: %v", req)

	uc.log.WithContext(ctx).Info("CreateSyncAccount.SaveCompanyCfg")
	err := uc.repo.SaveCompanyCfg(ctx, &dingtalk.DingtalkCompanyCfg{})
	uc.log.WithContext(ctx).Infof("CreateSyncAccount.SaveCompanyCfg: err: %v", err)
	if err != nil {
		return nil, err
	}
	uc.log.WithContext(ctx).Infof("CreateSyncAccount.GetAccessToken")

	// 1. 获取access_token
	accessToken, err := uc.dingTalkRepo.GetAccessToken(ctx, "code")
	uc.log.WithContext(ctx).Infof("CreateSyncAccount.GetAccessToken: accessToken: %v, err: %v", accessToken, err)
	if err != nil {
		return nil, err
	}

	taskId := time.Now().Add(time.Duration(1) * time.Second).Format("20060102150405")

	// 1. 从第三方获取部门和用户数据

	uc.log.WithContext(ctx).Infof("CreateSyncAccount.FetchDepartments")

	depts, err := uc.dingTalkRepo.FetchDepartments(ctx, accessToken)
	uc.log.WithContext(ctx).Infof("CreateSyncAccount.FetchDepartments: depts: %+v, err: %v", depts, err)
	if err != nil {
		return nil, err
	}
	for _, dept := range depts {
		uc.log.WithContext(ctx).Infof("biz.CreateSyncAccount: dept: %+v", dept)
	}

	uc.log.WithContext(ctx).Infof("CreateSyncAccount.SaveDepartments depts: %v, taskId: %v", depts, taskId)
	// 2. 数据入库
	deptCount, err := uc.repo.SaveDepartments(ctx, depts, taskId)
	uc.log.WithContext(ctx).Infof("CreateSyncAccount.SaveDepartments: deptCount: %v, err: %v", deptCount, err)
	if err != nil {
		return nil, err
	}
	var deptIds []int64
	for _, dept := range depts {
		deptIds = append(deptIds, dept.DeptID)
	}

	uc.log.WithContext(ctx).Infof("CreateSyncAccount.FetchDepartmentUsers accessToken: %v deptIds: %v", accessToken, deptIds)
	// 1. 从第三方获取用户数据
	deptUsers, err := uc.dingTalkRepo.FetchDepartmentUsers(ctx, accessToken, deptIds)
	uc.log.WithContext(ctx).Infof("CreateSyncAccount.FetchDepartmentUsers deptUsers: %v, err: %v", deptUsers, err)
	if err != nil {
		return nil, err
	}
	// 2. 数据入库
	//这里可以 将deptUsers转为model.TbLasUser,
	// SaveUsers(ctx, TbLasUser)
	uc.log.WithContext(ctx).Infof("CreateSyncAccount.SaveUsers deptUsers: %v, taskId: %v", deptUsers, taskId)
	userCount, err := uc.repo.SaveUsers(ctx, deptUsers, taskId)
	uc.log.WithContext(ctx).Infof("CreateSyncAccount.SaveUsers userCount: %v, err: %v", userCount, err)
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
	uc.log.WithContext(ctx).Infof("CreateSyncAccount.SaveDepartmentUserRelations deptUserRelations: %v, taskId: %v", deptUserRelations, taskId)
	// 3. 数据入库
	relationCount, err := uc.repo.SaveDepartmentUserRelations(ctx, deptUserRelations, taskId)
	uc.log.WithContext(ctx).Infof("CreateSyncAccount.SaveDepartmentUserRelations relationCount: %v, err: %v", relationCount, err)
	if err != nil {
		return nil, err
	}
	uc.log.WithContext(ctx).Infof("CreateSyncAccount.CallEcisaccountsyncAll taskId: %v", taskId)

	appAccessToken, err := uc.appAuth.GetAccessToken(ctx)
	if err != nil {
		return nil, err
	}
	fmt.Println("appAccessToken", appAccessToken)

	res, err := uc.wpsSync.PostEcisaccountsyncAll(ctx, appAccessToken.AccessToken, &wps.EcisaccountsyncAllRequest{
		TaskId:         taskId,
		ThirdCompanyId: uc.bizConf.ThirdCompanyId,
	})
	uc.log.WithContext(ctx).Infof("CreateSyncAccount.CallEcisaccountsyncAll res: %v, err: %v", res, err)

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
	return &v1.GetSyncAccountReply{
		Status: v1.GetSyncAccountReply_SUCCESS,
	}, nil
}
