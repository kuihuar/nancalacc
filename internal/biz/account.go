package biz

import (
	"context"
	"errors"
	v1 "nancalacc/api/account/v1"
	"nancalacc/internal/auth"
	"nancalacc/internal/dingtalk"
	"nancalacc/internal/repository/contracts"
	"nancalacc/internal/wps"

	//"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GreeterUsecase is a Greeter usecase.
type AccounterUsecase struct {
	repo         contracts.AccountRepository
	dingTalkRepo dingtalk.Dingtalk
	appAuth      auth.Authenticator
	wps          wps.Wps
	localCache   contracts.CacheRepository
	log          log.Logger
}

// NewGreeterUsecase new a Greeter usecase.
func NewAccounterUsecase(repo contracts.AccountRepository, dingTalkRepo dingtalk.Dingtalk, wps wps.Wps, cache contracts.CacheRepository, logger log.Logger) *AccounterUsecase {
	appAuth := auth.NewWpsAppAuthenticator()
	return &AccounterUsecase{repo: repo, dingTalkRepo: dingTalkRepo, appAuth: appAuth, wps: wps, localCache: cache, log: logger}
}

func (uc *AccounterUsecase) CreateTask(ctx context.Context, taskName string) (int, error) {
	return 0, status.Errorf(codes.Unimplemented, "method CreateTask not implemented")
	//uc.log.Log(log.LevelInfo, "msg", "CreateTask", "taskName", taskName)
	//return uc.repo.CreateTask(ctx, taskName)

}
func (uc *AccounterUsecase) GetTask(ctx context.Context, taskName string) (*v1.GetTaskReply_Task, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateTask not implemented")
	//return &v1.GetTaskReply_Task{}, nil

}
func (uc *AccounterUsecase) UpdateTask(ctx context.Context, taskName, status string) error {

	return errors.New("not implemented")
	//uc.log.Log(log.LevelInfo, "msg", "UpdateTask", "taskId", taskName, "status", status)
	//return uc.repo.UpdateTask(ctx, taskName, status)

}
