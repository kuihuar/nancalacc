package service

import (
	"context"
	"errors"

	pb "nancalacc/api/account/v1"
	v1 "nancalacc/api/account/v1"
	"nancalacc/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/emptypb"
)

type AccountService struct {
	pb.UnimplementedAccountServer
	accounterUsecase *biz.AccounterUsecase
	log              *log.Helper
}

func NewAccountService(accounterUsecase *biz.AccounterUsecase, logger log.Logger) *AccountService {
	return &AccountService{accounterUsecase: accounterUsecase, log: log.NewHelper(logger)}
}

func (s *AccountService) CreateSyncAccount(ctx context.Context, req *pb.CreateSyncAccountRequest) (*pb.CreateSyncAccountReply, error) {
	s.log.Infof("CreateSyncAccount req: %v", req)
	_, err := s.accounterUsecase.CreateSyncAccount(ctx, req)
	if err != nil {
		return nil, err
	}
	return &pb.CreateSyncAccountReply{
		TaskId: "10",
	}, nil
}
func (s *AccountService) GetSyncAccount(ctx context.Context, req *pb.GetSyncAccountRequest) (*pb.GetSyncAccountReply, error) {
	return &pb.GetSyncAccountReply{}, nil
}
func (s *AccountService) CancelSyncTask(ctx context.Context, req *pb.CancelSyncAccountRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}
func (s *AccountService) GetUserInfo(ctx context.Context, req *pb.GetUserInfoRequest) (*pb.GetUserInfoResponse, error) {
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
func (s *AccountService) GetAccessToken(ctx context.Context, req *pb.GetAccessTokenRequest) (*pb.GetAccessTokenResponse, error) {
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
