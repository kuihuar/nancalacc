package biz

import (
	"context"
	"errors"
	v1 "nancalacc/api/account/v1"
	"nancalacc/internal/dingtalk"

	"github.com/go-kratos/kratos/v2/log"
)

// GreeterUsecase is a Greeter usecase.
type Oauth2Usecase struct {
	dingTalkRepo dingtalk.Dingtalk
	log          log.Logger
}

// NewGreeterUsecase new a Greeter usecase.
func NewOauth2Usecase(dingTalkRepo dingtalk.Dingtalk, logger log.Logger) *Oauth2Usecase {
	return &Oauth2Usecase{dingTalkRepo: dingTalkRepo, log: logger}
}

func (uc *Oauth2Usecase) GetUserInfo(ctx context.Context, req *v1.GetUserInfoRequest) (*v1.GetUserInfoResponse, error) {

	uc.log.Log(log.LevelInfo, "msg", "GetUserInfo", "req", req)

	accessToken := req.GetAccessToken()
	if accessToken == "" {
		return nil, errors.New("access_token is empty")
	}
	var userId string
	userInfo, err := uc.dingTalkRepo.GetUserInfo(ctx, accessToken, "me")
	uc.log.Log(log.LevelInfo, "msg", "GetUserInfo.dingTalkRepo.GetUserInfo", "userInfo", userInfo, "err", err)
	if err != nil {
		uc.log.Log(log.LevelError, "msg", "GetUserInfo.dingTalkRepo.GetUserInfo", "userInfo", userInfo, "err", err)
		return nil, err
	}
	token, err := uc.dingTalkRepo.GetAccessToken(ctx)
	uc.log.Log(log.LevelInfo, "msg", "GetUserInfo.dingTalkRepo.GetAccessToken", "token", token, "err", err)
	if err != nil {
		uc.log.Log(log.LevelError, "msg", "GetUserInfo.dingTalkRepo.GetAccessToken", "token", token, "err", err)
		return nil, err
	}
	userId, err = uc.dingTalkRepo.GetUseridByUnionid(ctx, token.AccessToken, userInfo.UnionId)
	log.Infof("GetUserInfo.GetUseridByUnionid: userId: %v, err: %v", userId, err)

	if err != nil {
		uc.log.Log(log.LevelError, "msg", "GetUserInfo.GetUseridByUnionid", "userId", userId, "err", err)
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
func (uc *Oauth2Usecase) GetAccessToken(ctx context.Context, req *v1.GetAccessTokenRequest) (*v1.GetAccessTokenResponse, error) {

	uc.log.Log(log.LevelInfo, "msg", "GetAccessToken", "req", req)

	code := req.GetCode()
	if code == "" {
		return nil, errors.New("code is empty")
	}
	tokenRes, err := uc.dingTalkRepo.GetUserAccessToken(ctx, code)
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
