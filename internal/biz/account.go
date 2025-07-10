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

	CallEcisaccountsyncAll(ctx context.Context, taskId string) (EcisaccountsyncResponse, error)

	ClearAll(ctx context.Context) error
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
	uc.log.WithContext(ctx).Infof("CreateSyncAccount: %v", req)
	err := uc.repo.SaveCompanyCfg(ctx, &DingtalkCompanyCfg{})
	uc.log.WithContext(ctx).Infof("biz.CreateSyncAccount: err: %v", err)
	if err != nil {
		return nil, err
	}
	uc.log.WithContext(ctx).Infof("CreateSyncAccount: %v", req)

	// 1. 获取access_token
	accessToken, err := uc.dingTalkRepo.GetAccessToken(ctx, "code")
	uc.log.WithContext(ctx).Infof("biz.CreateSyncAccount: accessToken: %v, err: %v", accessToken, err)
	if err != nil {
		return nil, err
	}

	taskId := time.Now().Add(time.Duration(1) * time.Second).Format("20060102150405")

	// 1. 从第三方获取部门和用户数据
	depts, err := uc.dingTalkRepo.FetchDepartments(ctx, accessToken)
	uc.log.WithContext(ctx).Infof("biz.CreateSyncAccount: depts: %v, err: %v", depts, err)
	if err != nil {
		return nil, err
	}
	// 2. 数据入库
	deptCount, err := uc.repo.SaveDepartments(ctx, depts, taskId)
	uc.log.WithContext(ctx).Infof("biz.CreateSyncAccount: deptCount: %v, err: %v", deptCount, err)
	if err != nil {
		return nil, err
	}
	var deptIds []int64
	for _, dept := range depts {
		deptIds = append(deptIds, dept.DeptID)
	}

	// 1. 从第三方获取用户数据
	deptUsers, err := uc.dingTalkRepo.FetchDepartmentUsers(ctx, accessToken, deptIds)
	if err != nil {
		return nil, err
	}
	// 2. 数据入库
	//这里可以 将deptUsers转为model.TbLasUser,
	// SaveUsers(ctx, TbLasUser)
	userCount, err := uc.repo.SaveUsers(ctx, deptUsers, taskId)
	uc.log.WithContext(ctx).Infof("biz.CreateSyncAccount: userCount: %v, err: %v", userCount, err)
	if err != nil {
		return nil, err
	}

	// 2. 关系数据入库
	var deptUserRelations []*DingtalkDeptUserRelation
	for _, deptUser := range deptUsers {
		for _, depId := range deptUser.DeptIDList {
			deptUserRelations = append(deptUserRelations, &DingtalkDeptUserRelation{
				Uid:   deptUser.Unionid,
				Did:   strconv.FormatInt(depId, 10),
				Order: deptUser.DeptOrder,
			})
		}

	}
	// 3. 数据入库
	relationCount, err := uc.repo.SaveDepartmentUserRelations(ctx, deptUserRelations, taskId)
	uc.log.WithContext(ctx).Infof("biz.CreateSyncAccount: relationCount: %v, err: %v", relationCount, err)
	if err != nil {
		return nil, err
	}

	_, err = uc.repo.CallEcisaccountsyncAll(ctx, taskId)
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

	userInfo, err := uc.dingTalkRepo.GetUserInfo(ctx, accessToken, "me")
	if err != nil {
		return nil, err
	}
	return &v1.GetUserInfoResponse{
		UnionId: userInfo.UnionId,
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
