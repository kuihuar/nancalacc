package biz

import (
	"context"
	"errors"
	v1 "nancalacc/api/account/v1"
	"nancalacc/internal/conf"
	"nancalacc/internal/dingtalk"

	"github.com/go-kratos/kratos/v2/log"
)

// GreeterUsecase is a Greeter usecase.
type Oauth2Usecase struct {
	dingTalkRepo dingtalk.Dingtalk
	bizConf      *conf.Service_Business
	log          *log.Helper
}

// NewGreeterUsecase new a Greeter usecase.
func NewOauth2Usecase(dingTalkRepo dingtalk.Dingtalk, bizConf *conf.Service_Business, logger log.Logger) *Oauth2Usecase {
	return &Oauth2Usecase{dingTalkRepo: dingTalkRepo, bizConf: bizConf, log: log.NewHelper(logger)}
}

func (uc *Oauth2Usecase) GetUserInfo(ctx context.Context, req *v1.GetUserInfoRequest) (*v1.GetUserInfoResponse, error) {
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
func (s *Oauth2Usecase) GetAccessToken(ctx context.Context, req *v1.GetAccessTokenRequest) (*v1.GetAccessTokenResponse, error) {
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
