package biz

import (
	"context"
	v1 "nancalacc/api/account/v1"
	"nancalacc/internal/auth"
	"nancalacc/internal/dingtalk"
	"nancalacc/internal/wps"

	//"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

// GreeterUsecase is a Greeter usecase.
type AccounterUsecase struct {
	repo         AccounterRepo
	dingTalkRepo dingtalk.Dingtalk
	appAuth      auth.Authenticator
	wps          wps.Wps
	localCache   CacheService
	log          log.Logger
}

var (
	prefix = "nancalacc:cache:"
)

// NewGreeterUsecase new a Greeter usecase.
func NewAccounterUsecase(repo AccounterRepo, dingTalkRepo dingtalk.Dingtalk, wps wps.Wps, cache CacheService, logger log.Logger) *AccounterUsecase {
	appAuth := auth.NewWpsAppAuthenticator()
	return &AccounterUsecase{repo: repo, dingTalkRepo: dingTalkRepo, appAuth: appAuth, wps: wps, localCache: cache, log: logger}
}

func (uc *AccounterUsecase) CreateTask(ctx context.Context, taskName string) (int, error) {
	uc.log.Log(log.LevelInfo, "msg", "CreateTask", "taskName", taskName)
	return uc.repo.CreateTask(ctx, taskName)

}
func (uc *AccounterUsecase) GetTask(ctx context.Context, taskName string) (*v1.GetTaskReply_Task, error) {

	return &v1.GetTaskReply_Task{}, nil

}
func (uc *AccounterUsecase) UpdateTask(ctx context.Context, taskName, status string) error {
	uc.log.Log(log.LevelInfo, "msg", "UpdateTask", "taskId", taskName, "status", status)
	return uc.repo.UpdateTask(ctx, taskName, status)

}
