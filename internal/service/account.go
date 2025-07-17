package service

import (
	"context"
	"errors"
	"time"

	v1 "nancalacc/api/account/v1"
	"nancalacc/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/emptypb"
)

type AccountService struct {
	v1.UnimplementedAccountServer
	accounterUsecase *biz.AccounterUsecase
	log              *log.Helper
}

func NewAccountService(accounterUsecase *biz.AccounterUsecase, logger log.Logger) *AccountService {
	return &AccountService{accounterUsecase: accounterUsecase, log: log.NewHelper(logger)}
}

func (s *AccountService) CreateSyncAccount(ctx context.Context, req *v1.CreateSyncAccountRequest) (*v1.CreateSyncAccountReply, error) {
	s.log.Infof("CreateSyncAccount req: %v", req)
	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()
	return s.accounterUsecase.CreateSyncAccount(ctx, req)
}
func (s *AccountService) GetSyncAccount(ctx context.Context, req *v1.GetSyncAccountRequest) (*v1.GetSyncAccountReply, error) {
	return &v1.GetSyncAccountReply{}, nil
}
func (s *AccountService) CancelSyncTask(ctx context.Context, req *v1.CancelSyncAccountRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}
func (s *AccountService) GetUserInfo(ctx context.Context, req *v1.GetUserInfoRequest) (*v1.GetUserInfoResponse, error) {
	s.log.Infof("GetUserInfo req: %v", req)
	accessToken := req.GetAccessToken()
	if accessToken == "" {
		return nil, errors.New("access_token is empty")
	}
	userInfo, err := s.accounterUsecase.GetUserInfo(ctx, &v1.GetUserInfoRequest{
		AccessToken: accessToken,
	})
	if err != nil {
		return nil, err
	}
	return userInfo, nil
}
func (s *AccountService) GetAccessToken(ctx context.Context, req *v1.GetAccessTokenRequest) (*v1.GetAccessTokenResponse, error) {
	s.log.Infof("GetAccessToken req: %v", req)
	code := req.GetCode()
	if code == "" {
		return nil, errors.New("code is empty")
	}
	accessTokenResp, err := s.accounterUsecase.GetAccessToken(ctx, &v1.GetAccessTokenRequest{
		Code: code,
	})
	if err != nil {
		return nil, err
	}
	return accessTokenResp, nil
}

func (s *AccountService) Callback(ctx context.Context, req *v1.CallbackRequest) (*v1.CallbackResponse, error) {
	s.log.Infof("Callback req: %v", req.Code)
	return &v1.CallbackResponse{
		Status:  "success",
		Message: "success",
	}, nil
}
