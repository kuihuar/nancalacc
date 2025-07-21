package biz

import (
	"context"
	"errors"
	v1 "nancalacc/api/account/v1"
	"nancalacc/internal/conf"
	"strconv"
	"time"

	//"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/open-dingtalk/dingtalk-stream-sdk-go/clientV2"
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

type AccounterConf struct {
	Env            string `json:"env"`
	LogLevel       string `json:"log_level"`
	AccessKey      string `json:"access_key"`
	SecretKey      string `json:"secret_key"`
	ThirdCompanyID string `json:"third_company_id"`
	PlatformIDs    string `json:"platform_ids"`
	CompanyID      string `json:"company_id"`
}

func NewAccounterConf(c *conf.ServiceConf) *AccounterConf {
	conf := &AccounterConf{
		Env:            c.Env,
		LogLevel:       c.LogLevel,
		AccessKey:      c.AccessKey,
		SecretKey:      c.SecretKey,
		ThirdCompanyID: c.ThirdCompanyId,
		PlatformIDs:    c.PlatformIds,
		CompanyID:      c.CompanyId,
	}
	return conf
}

// GreeterRepo is a Greater repo.
type AccounterRepo interface {
	SaveUsers(ctx context.Context, users []*DingtalkDeptUser, taskId string) (int, error)
	SaveDepartments(ctx context.Context, depts []*DingtalkDept, taskId string) (int, error)
	SaveDepartmentUserRelations(ctx context.Context, relations []*DingtalkDeptUserRelation, taskId string) (int, error)
	SaveCompanyCfg(ctx context.Context, cfg *DingtalkCompanyCfg) error

	CallEcisaccountsyncAll(ctx context.Context, taskId string) (EcisaccountsyncAllResponse, error)

	ClearAll(ctx context.Context) error

	SaveIncrementDepartments(ctx context.Context, deptsAdd, deptsDel []*DingtalkDept) error
	SaveIncrementUsers(ctx context.Context, usersAdd, usersDel []*DingtalkDeptUser) error
	SaveIncrementDepartmentUserRelations(ctx context.Context, relationsAdd, relationsDel []*DingtalkDeptUserRelation) error

	CallEcisaccountsyncIncrement(ctx context.Context, thirdCompanyId string) (EcisaccountsyncIncrementResponse, error)
}

// GreeterUsecase is a Greeter usecase.
type AccounterUsecase struct {
	repo         AccounterRepo
	dingTalkRepo DingTalkRepo
	log          *log.Helper
}

// NewGreeterUsecase new a Greeter usecase.
func NewAccounterUsecase(repo AccounterRepo, dingTalkRepo DingTalkRepo, logger log.Logger) *AccounterUsecase {
	return &AccounterUsecase{repo: repo, dingTalkRepo: dingTalkRepo, log: log.NewHelper(logger)}
}

func (uc *AccounterUsecase) CreateSyncAccount(ctx context.Context, req *v1.CreateSyncAccountRequest) (*v1.CreateSyncAccountReply, error) {
	// return &v1.CreateSyncAccountReply{
	// 	TaskId:     "taskId",
	// 	CreateTime: timestamppb.Now(),
	// }, nil

	uc.log.WithContext(ctx).Infof("CreateSyncAccount: %v", req)

	uc.log.WithContext(ctx).Info("CreateSyncAccount.SaveCompanyCfg")
	err := uc.repo.SaveCompanyCfg(ctx, &DingtalkCompanyCfg{})
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
	var deptUserRelations []*DingtalkDeptUserRelation
	for _, deptUser := range deptUsers {
		order := int(deptUser.DeptOrder)
		if order > 0 {
			order = 1
		} else {
			order = 0
		}
		for _, depId := range deptUser.DeptIDList {

			deptUserRelations = append(deptUserRelations, &DingtalkDeptUserRelation{
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

	res, err := uc.repo.CallEcisaccountsyncAll(ctx, taskId)
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

func (uc *AccounterUsecase) GetUserInfo(ctx context.Context, req *v1.GetUserInfoRequest) (*v1.GetUserInfoResponse, error) {
	uc.log.WithContext(ctx).Infof("GetUserInfo: %v", req)
	accessToken := req.GetAccessToken()
	if accessToken == "" {
		return nil, errors.New("access_token is empty")
	}
	var userId string
	userInfo, err := uc.dingTalkRepo.GetUserInfo(ctx, accessToken, "me")
	uc.log.WithContext(ctx).Infof("GetUserInfo.dingTalkRepo.GetUserInfo: %v, err:%v", userInfo, err)
	if err != nil {
		uc.log.WithContext(ctx).Errorf("GetUserInfo.dingTalkRepo.GetUserInfo: %v, err:%v", userInfo, err)
		return nil, err
	}
	token, err := uc.dingTalkRepo.GetAccessToken(ctx, "code")
	uc.log.WithContext(ctx).Infof("GetUserInfo.dingTalkRepo.GetAccessToken: token: %v, err: %v", token, err)
	if err != nil {
		uc.log.WithContext(ctx).Error("GetUserInfo.dingTalkRepo.GetAccessToken: token: %v, err: %v", token, err)
		return nil, err
	}
	userId, err = uc.dingTalkRepo.GetUseridByUnionid(ctx, token, userInfo.UnionId)
	uc.log.WithContext(ctx).Infof("GetUserInfo.GetUseridByUnionid: userId: %v, err: %v", userId, err)

	if err != nil {
		uc.log.WithContext(ctx).Error("GetUserInfo.GetUseridByUnionid: userId: %v, err: %v", userId, err)
		return nil, err
	}

	return &v1.GetUserInfoResponse{
		UnionId: userInfo.UnionId,
		UserId:  userId,
		Name:    userInfo.Nick,
		Email:   userInfo.Email,
		Avatar:  userInfo.AvatarUrl,
	}, nil
}

// https://login.dingtalk.com/oauth2/challenge.htm?
// client_id=dinglz1setxqhrpp7aa0
// &redirect_uri=http://119.3.173.229/cloud/login/api/v1/oauth/code/login?auth_type=oauth
// &platform_id=1
// &response_type=code
// &state=6c938a3e11174f67bf40b2d7d679dbe1
func (s *AccounterUsecase) GetAccessToken(ctx context.Context, req *v1.GetAccessTokenRequest) (*v1.GetAccessTokenResponse, error) {
	s.log.WithContext(ctx).Infof("GetAccessToken: %v", req)
	code := req.GetCode()
	if code == "" {
		return nil, errors.New("code is empty")
	}
	tokenRes, err := s.dingTalkRepo.GetUserAccessToken(ctx, code)
	if err != nil {
		return nil, err
	}
	return &v1.GetAccessTokenResponse{
		AccessToken:  tokenRes.AccessToken,
		RefreshToken: tokenRes.RefreshToken,
		ExpiresIn:    int64(tokenRes.ExpireIn),
		//RefreshToken: tokenRes.RefreshToken,
	}, nil
}

func (uc *AccounterUsecase) OrgDeptCreate(ctx context.Context, event *clientV2.GenericOpenDingTalkEvent) error {
	uc.log.Infof("OrgDeptCreate: %v", event.Data)
	if event.Data == nil {
		return nil
	}

	depIds, err := uc.getDeptidsFromDingTalkEvent(event)
	if err != nil {
		return err
	}

	if len(depIds) == 0 {
		uc.log.Info("OrgDeptCreate len(depIds) eq 0")
		return nil
	}

	accessToken, err := uc.dingTalkRepo.GetAccessToken(ctx, "code")
	uc.log.WithContext(ctx).Infof("OrgDeptCreate.GetAccessToken accessToken: %v, err: %v", accessToken, err)
	if err != nil {
		return err
	}
	uc.log.WithContext(ctx).Infof("OrgDeptCreate.FetchDeptDetails accessToken: %v, depIds: %v", accessToken, depIds)
	depts, err := uc.dingTalkRepo.FetchDeptDetails(ctx, accessToken, depIds)
	uc.log.WithContext(ctx).Infof("OrgDeptCreate.FetchDeptDetails accessToken: %v, depIds: %v, err:%v", accessToken, depIds, err)
	if err != nil {
		return err
	}

	err = uc.repo.SaveIncrementDepartments(ctx, depts, nil)
	if err != nil {
		uc.log.Errorf("OrgDeptCreate.SaveIncrementDepartments err: %v", err)
		return err
	}

	res, err := uc.repo.CallEcisaccountsyncIncrement(ctx, "taskId")

	uc.log.WithContext(ctx).Infof("UserLeaveOrg.CallEcisaccountsyncIncrement res: %v, err: %v", res, err)
	if err != nil {
		uc.log.Errorf("UserLeaveOrg.CallEcisaccountsyncIncrement res: %v, err: %v", res, err)
		//return err
	}

	return nil

}
func (uc *AccounterUsecase) OrgDeptModify(ctx context.Context, event *clientV2.GenericOpenDingTalkEvent) error {
	return uc.repo.SaveIncrementDepartments(ctx, nil, nil)
}
func (uc *AccounterUsecase) OrgDeptRemove(ctx context.Context, event *clientV2.GenericOpenDingTalkEvent) error {
	uc.log.Infof("OrgDeptCreate: %v", event.Data)
	if event.Data == nil {
		return nil
	}

	depIds, err := uc.getDeptidsFromDingTalkEvent(event)
	if err != nil {
		return err
	}

	if len(depIds) == 0 {
		uc.log.Info("OrgDeptCreate len(depIds) eq 0")
		return nil
	}
	depts := make([]*DingtalkDept, len(depIds))

	// TODO 从 wps 获取 ParentID, 才可以提交
	for i, depId := range depIds {
		depts[i] = &DingtalkDept{
			DeptID: depId,
		}
	}

	// accessToken, err := uc.dingTalkRepo.GetAccessToken(ctx, "code")
	// uc.log.WithContext(ctx).Infof("OrgDeptCreate.GetAccessToken accessToken: %v, err: %v", accessToken, err)
	// if err != nil {
	// 	return err
	// }
	// uc.log.WithContext(ctx).Infof("OrgDeptCreate.FetchDeptDetails accessToken: %v, depIds: %v", accessToken, depIds)
	// depts, err := uc.dingTalkRepo.FetchDeptDetails(ctx, accessToken, depIds)
	// uc.log.WithContext(ctx).Infof("OrgDeptCreate.FetchDeptDetails accessToken: %v, depIds: %v, err:%v", accessToken, depIds, err)
	// if err != nil {
	// 	return err
	// }

	res, err := uc.repo.CallEcisaccountsyncIncrement(ctx, "taskId")

	uc.log.WithContext(ctx).Infof("UserLeaveOrg.CallEcisaccountsyncIncrement res: %v, err: %v", res, err)
	if err != nil {
		uc.log.Errorf("UserLeaveOrg.CallEcisaccountsyncIncrement res: %v, err: %v", res, err)
		//return err
	}

	return uc.repo.SaveIncrementDepartments(ctx, nil, depts)
}
func (uc *AccounterUsecase) UserAddOrg(ctx context.Context, event *clientV2.GenericOpenDingTalkEvent) error {
	uc.log.Infof("UserAddOrg: %v", event.Data)
	if event.Data == nil {
		return nil
	}

	userIds, err := uc.getUseridsFromDingTalkEvent(event)
	if err != nil {
		return err
	}

	accessToken, err := uc.dingTalkRepo.GetAccessToken(ctx, "code")
	uc.log.WithContext(ctx).Infof("UserAddOrg.GetAccessToken accessToken: %v,userIds:%v err: %v", accessToken, userIds, err)
	if err != nil {
		return err
	}
	uc.log.WithContext(ctx).Infof("UserAddOrg.GetUserDetail userIds: %v", userIds)
	users, err := uc.dingTalkRepo.FetchUserDetail(ctx, accessToken, userIds)
	if err != nil {
		return err
	}

	err = uc.repo.SaveIncrementUsers(ctx, users, nil)
	if err != nil {
		return err
	}

	relations := generateUserDeptRelations(users)

	err = uc.repo.SaveIncrementDepartmentUserRelations(ctx, relations, nil)

	if err != nil {
		return err
	}

	res, err := uc.repo.CallEcisaccountsyncIncrement(ctx, "taskId")

	uc.log.WithContext(ctx).Infof("UserLeaveOrg.CallEcisaccountsyncIncrement res: %v, err: %v", res, err)
	if err != nil {
		uc.log.Errorf("UserLeaveOrg.CallEcisaccountsyncIncrement res: %v, err: %v", res, err)
		//return err
	}

	return nil
}
func (uc *AccounterUsecase) UserModifyOrg(ctx context.Context, event *clientV2.GenericOpenDingTalkEvent) error {
	return uc.repo.SaveIncrementUsers(ctx, nil, nil)
}
func (uc *AccounterUsecase) UserLeaveOrg(ctx context.Context, event *clientV2.GenericOpenDingTalkEvent) error {

	uc.log.Infof("UserLeaveOrg: %v", event.Data)
	if event.Data == nil {
		return nil
	}

	userIds, err := uc.getUseridsFromDingTalkEvent(event)
	if err != nil {
		return err
	}
	users := make([]*DingtalkDeptUser, len(userIds))
	for i, userId := range userIds {
		users[i] = &DingtalkDeptUser{
			Userid: userId,
		}
	}
	// accessToken, err := uc.dingTalkRepo.GetAccessToken(ctx, "code")
	// uc.log.WithContext(ctx).Infof("CreateSyncAccount.GetAccessToken accessToken: %v,userIds:%v err: %v", accessToken, userIds, err)
	// if err != nil {
	// 	return err
	// }
	// uc.log.WithContext(ctx).Infof("CreateSyncAccount.GetAccessToken: accessToken: %v, err: %v", accessToken, err)
	// users, err := uc.dingTalkRepo.FetchUserDetail(ctx, accessToken, userIds)
	// if err != nil {
	// 	return err
	// }
	err = uc.repo.SaveIncrementUsers(ctx, nil, users)
	if err != nil {
		return err
	}
	// relations := generateUserDeptRelations(users)

	// err = uc.repo.SaveIncrementDepartmentUserRelations(ctx, nil, relations)

	// if err != nil {
	// 	return err
	// }

	res, err := uc.repo.CallEcisaccountsyncIncrement(ctx, "taskId")

	uc.log.WithContext(ctx).Infof("UserLeaveOrg.CallEcisaccountsyncIncrement res: %v, err: %v", res, err)
	if err != nil {
		uc.log.Errorf("UserLeaveOrg.CallEcisaccountsyncIncrement res: %v, err: %v", res, err)
		//return err
	}
	return nil
}

func generateUserDeptRelations(deptUsers []*DingtalkDeptUser) []*DingtalkDeptUserRelation {
	var deptUserRelations []*DingtalkDeptUserRelation
	for _, deptUser := range deptUsers {
		order := int(deptUser.DeptOrder)
		if order > 0 {
			order = 1
		} else {
			order = 0
		}
		for _, depId := range deptUser.DeptIDList {

			deptUserRelations = append(deptUserRelations, &DingtalkDeptUserRelation{
				Uid:   deptUser.Userid,
				Did:   strconv.FormatInt(depId, 10),
				Order: order,
			})
		}

	}

	return deptUserRelations
}

func (uc *AccounterUsecase) getDeptidsFromDingTalkEvent(event *clientV2.GenericOpenDingTalkEvent) ([]int64, error) {
	uc.log.Infof("getDeptidsFromDingTalkEvent: %v", event.Data)
	if event.Data == nil {
		return nil, errors.New("getDeptidsFromDingTalkEvent event.Data is nil")
	}
	datamap := event.Data
	var depIds []int64

	deptId, exists := datamap["deptId"]

	if !exists {
		uc.log.Errorf("getDeptidsFromDingTalkEvent not deptId: %v, exists: %v", deptId, exists)
		return nil, errors.New("getDeptidsFromDingTalkEvent not deptId")
	}

	deptIdSlice, ok := deptId.([]interface{})

	if !ok {
		uc.log.Errorf("deptId not []interface{}: %v, exists: %v", deptId, exists)
		return nil, errors.New("deptId not []interface{}")
	}

	for _, item := range deptIdSlice {
		if f, ok := item.(float64); ok {
			depIds = append(depIds, int64(f))
		} else {
			uc.log.Errorf("deptId not float64: %T", item)
			return nil, errors.New("deptId not float64")
		}
	}
	return depIds, nil
}

func (uc *AccounterUsecase) getUseridsFromDingTalkEvent(event *clientV2.GenericOpenDingTalkEvent) ([]string, error) {
	uc.log.Infof("getUseridsFromDingTalkEvent: %v", event.Data)
	if event.Data == nil {
		return nil, errors.New("getUseridsFromDingTalkEvent event.Data is nil")
	}
	datamap := event.Data
	var userIds []string

	userId, exists := datamap["userId"]

	if !exists {
		uc.log.Errorf("getUseridsFromDingTalkEvent not userId: %v, exists: %v", userId, exists)
		return nil, errors.New("getUseridsFromDingTalkEvent not userId")
	}

	userIdSlice, ok := userId.([]interface{})

	if !ok {
		uc.log.Errorf("deptId not []interface{}: %v, exists: %v", userId, exists)
		return nil, errors.New("userId not []interface{}")
	}

	for _, item := range userIdSlice {
		if f, ok := item.(string); ok {
			userIds = append(userIds, f)
		} else {
			uc.log.Errorf("userId not string: %T", item)
			return nil, errors.New("userId not string")
		}
	}
	return userIds, nil
}
