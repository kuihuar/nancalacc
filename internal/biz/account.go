package biz

import (
	"context"
	"errors"
	v1 "nancalacc/api/account/v1"
	"nancalacc/internal/auth"
	"nancalacc/internal/data/models"
	"nancalacc/internal/dingtalk"
	"nancalacc/internal/wps"
	"time"

	//"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// GreeterUsecase is a Greeter usecase.
type AccounterUsecase struct {
	repo         AccounterRepo
	dingTalkRepo dingtalk.Dingtalk
	appAuth      auth.Authenticator
	wps          wps.Wps
	localCache   CacheService
	log          *log.Helper
}

var (
	prefix = "nancalacc:cache:"
)

// NewGreeterUsecase new a Greeter usecase.
func NewAccounterUsecase(repo AccounterRepo, dingTalkRepo dingtalk.Dingtalk, wps wps.Wps, cache CacheService, logger log.Logger) *AccounterUsecase {
	appAuth := auth.NewWpsAppAuthenticator()
	return &AccounterUsecase{repo: repo, dingTalkRepo: dingTalkRepo, appAuth: appAuth, wps: wps, localCache: cache, log: log.NewHelper(logger)}
}

func (uc *AccounterUsecase) CreateTask(ctx context.Context, taskName string) (int, error) {
	uc.log.WithContext(ctx).Infof("CreateTask taskName: %s", taskName)
	return uc.repo.CreateTask(ctx, taskName)

}
func (uc *AccounterUsecase) GetTask(ctx context.Context, taskName string) (*v1.GetTaskReply_Task, error) {
	uc.log.WithContext(ctx).Infof("GetTask taskName: %s", taskName)

	taskInfo, err := uc.GetCacheTask(ctx, taskName)
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

func (uc *AccounterUsecase) CreateCacheTask(ctx context.Context, taskName, status string) error {

	cacheKey := prefix + taskName
	task := &models.Task{
		Title:         taskName,
		Description:   taskName,
		CreatedAt:     time.Now(),
		Status:        models.TaskStatusInProgress,
		Progress:      0,
		StartDate:     time.Now(),
		DueDate:       time.Now().Add(time.Minute * 30),
		CompletedAt:   time.Now(),
		CreatorID:     99,
		EstimatedTime: 10,
		ActualTime:    0,
	}
	return uc.localCache.Set(ctx, cacheKey, task, 300*time.Minute)
}
func (uc *AccounterUsecase) UpdateCacheTask(ctx context.Context, taskName, status string) error {

	cacheKey := prefix + taskName
	oldTask, ok, err := uc.localCache.Get(ctx, cacheKey)
	if err != nil {
		return err
	}
	var task *models.Task
	var startDate time.Time
	now := time.Now()
	if ok {
		task, ok1 := oldTask.(*models.Task)
		if ok1 {
			startDate = task.StartDate
			task.ActualTime = int(now.Sub(startDate).Seconds()) + 20
			task.Status = status
			task.Progress = 100
			task.CompletedAt = now
			task.UpdatedAt = now
		}
	}

	if task == nil {
		task = &models.Task{
			Title:         taskName,
			Description:   taskName,
			Status:        status,
			Progress:      100,
			StartDate:     time.Now(),
			DueDate:       time.Now().Add(time.Minute * 30),
			CompletedAt:   time.Now(),
			CreatorID:     99,
			EstimatedTime: 10,
			ActualTime:    0,
		}
	}
	return uc.localCache.Set(ctx, cacheKey, task, 300*time.Minute)
}

func (uc *AccounterUsecase) GetCacheTask(ctx context.Context, taskName string) (*models.Task, error) {

	cacheKey := prefix + taskName
	var task *models.Task
	taskInfo, ok, err := uc.localCache.Get(ctx, cacheKey)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("notfound")
	}
	task, ok1 := taskInfo.(*models.Task)
	if !ok1 {
		return nil, errors.New("type error")
	}
	return task, nil

}
