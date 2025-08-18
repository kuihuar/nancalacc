package service

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"time"

	v1 "nancalacc/api/account/v1"
	"nancalacc/internal/biz"
	"nancalacc/internal/conf"
	"nancalacc/internal/pkg/limiter"

	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type AccountService struct {
	v1.UnimplementedAccountServer
	accounterUsecase *biz.AccounterUsecase
	oauth2Usecase    *biz.Oauth2Usecase
	fullSyncUsecase  *biz.FullSyncUsecase
	limiter          *limiter.RateLimiter
	log              log.Logger
}

func NewAccountService(accounterUsecase *biz.AccounterUsecase, oauth2Usecase *biz.Oauth2Usecase, fullSyncUsecase *biz.FullSyncUsecase, logger log.Logger) *AccountService {
	limiter := limiter.NewRateLimiter(nil) // 使用默认配置
	return &AccountService{accounterUsecase: accounterUsecase, oauth2Usecase: oauth2Usecase, fullSyncUsecase: fullSyncUsecase, limiter: limiter, log: logger}
}

func (s *AccountService) CreateSyncAccount(ctx context.Context, req *v1.CreateSyncAccountRequest) (*v1.CreateSyncAccountReply, error) {
	s.log.Log(log.LevelInfo, "msg", "CreateSyncAccount", "req", req)
	if req.GetTaskName() != "" && len(req.GetTaskName()) != 14 {
		return nil, status.Errorf(codes.InvalidArgument, "invalid taskname: %s", req.GetTaskName())
	}
	if req.GetTaskName() == "" {
		taskId := time.Now().Add(time.Duration(1) * time.Second).Format("20060102150405")
		req.TaskName = &taskId
	}

	// 使用更合理的限流配置
	// 全局限流：每秒最多10个请求，突发20个请求
	if !s.limiter.Allow("global_sync_account", 10, 20) {
		return nil, status.Errorf(codes.ResourceExhausted, "global rate limit exceeded")
	}

	// 如果提供了任务名称，对特定任务进行额外限流
	if req.GetTaskName() != "" {
		// 特定任务限流：每秒最多2个请求，突发5个请求
		if !s.limiter.Allow("task_"+req.GetTaskName(), 2, 5) {
			return nil, status.Errorf(codes.ResourceExhausted, "task rate limit exceeded for: %s", req.GetTaskName())
		}
	}

	// 这里设置传进来的最大分钟数
	ctx, cancel := context.WithTimeout(ctx, 50*time.Minute)
	defer cancel()
	return s.fullSyncUsecase.CreateSyncAccount(ctx, req)
}
func (s *AccountService) GetSyncAccount(ctx context.Context, req *v1.GetSyncAccountRequest) (*v1.GetSyncAccountReply, error) {

	log.Infof("GetSyncAccount req: %v", req)
	globalConf := conf.Get()
	log.Infof("globalConf: %v", globalConf)
	ctx, cancel := context.WithTimeout(ctx, 50*time.Minute)
	defer cancel()
	return s.fullSyncUsecase.GetSyncAccount(ctx, req)
}
func (s *AccountService) CancelSyncTask(ctx context.Context, req *v1.CancelSyncAccountRequest) (*emptypb.Empty, error) {
	s.log.Log(log.LevelInfo, "msg", "CancelSyncTask", "req", req)
	if req.TaskId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "task_id is empty")
	}
	err := s.fullSyncUsecase.CleanSyncAccount(ctx, req.TaskId, req.GetTags())

	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
func (s *AccountService) GetUserInfo(ctx context.Context, req *v1.GetUserInfoRequest) (*v1.GetUserInfoResponse, error) {

	s.log.Log(log.LevelInfo, "msg", "GetUserInfo", "req", req)

	accessToken := req.GetAccessToken()
	if accessToken == "" {
		return nil, errors.New("access_token is empty")
	}
	userInfo, err := s.oauth2Usecase.GetUserInfo(ctx, &v1.GetUserInfoRequest{
		AccessToken: accessToken,
	})
	if err != nil {
		return nil, err
	}
	return userInfo, nil
}
func (s *AccountService) GetAccessToken(ctx context.Context, req *v1.GetAccessTokenRequest) (*v1.GetAccessTokenResponse, error) {

	s.log.Log(log.LevelInfo, "msg", "GetAccessToken", "req", req)

	code := req.GetCode()
	if code == "" {
		return nil, errors.New("code is empty")
	}
	accessTokenResp, err := s.oauth2Usecase.GetAccessToken(ctx, &v1.GetAccessTokenRequest{
		Code: code,
	})
	if err != nil {
		return nil, err
	}
	return accessTokenResp, nil
}

func (s *AccountService) Callback(ctx context.Context, req *v1.CallbackRequest) (*v1.CallbackResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Callback not implemented")
}

//	curl -X POST "http://your-server/v1/upload/excel" \
//	  -H "X-Filename: test.xlsx" \
//	  --data-binary "@/path/to/your/file.xlsx"

func (s *AccountService) UploadFile(ctx context.Context, req *v1.UploadRequest) (*v1.UploadReply, error) {

	s.log.Log(log.LevelInfo, "msg", "UploadFile", "req", req.GetFilename())

	taskId := time.Now().Add(time.Duration(1) * time.Second).Format("20060102150405")

	s.log.Log(log.LevelInfo, "msg", "UploadFile", "taskId", taskId)

	if req.GetFile() == nil {
		return nil, status.Errorf(codes.InvalidArgument, "file is empty")
	}

	// 检查文件类型
	// if req.GetFile().GetContentType() != "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" {
	// 	return nil, status.Errorf(codes.InvalidArgument, "file type is not excel")
	// }

	// // 检查文件大小
	// if req.GetFile().GetSize() > 10*1024*1024 {
	// 	return nil, status.Errorf(codes.InvalidArgument, "file size is too large")
	// }

	// 创建临时文件
	tempDir := os.TempDir()
	tempDir = "/tmp"
	filename := filepath.Join(tempDir, taskId+".xlsx")

	// 写入文件
	err := os.WriteFile(filename, req.GetFile(), 0644)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to write file: %v", err)
	}

	// 创建缓存任务
	err = s.fullSyncUsecase.CreateCacheTask(ctx, taskId, "pending")
	if err != nil {
		// 清理临时文件
		os.Remove(filename)
		return nil, status.Errorf(codes.Internal, "failed to create cache task: %v", err)
	}
	// 使用带错误处理的异步执行
	go func() {
		defer func() {
			// 清理临时文件
			if err := os.Remove(filename); err != nil {
				s.log.Log(log.LevelWarn, "msg", "failed to remove temp file", "filename", filename, "err", err)
			}
		}()

		// 检查原始 context 状态
		// select {
		// case <-ctx.Done():
		// 	s.log.Log(log.LevelWarn, "msg", "original context already canceled before starting goroutine", "err", ctx.Err())
		// 	return
		// default:
		// 	// 继续执行
		// }

		// 使用 Background context 避免依赖可能被取消的原始 context
		parseCtx, cancel := context.WithTimeout(context.Background(), 120*time.Minute)
		defer cancel()

		s.log.Log(log.LevelInfo, "msg", "starting excel parsing", "taskId", taskId, "filename", filename)

		// 更新任务状态为进行中
		if err := s.fullSyncUsecase.UpdateCacheTask(parseCtx, taskId, "in_progress", 20); err != nil {
			s.log.Log(log.LevelError, "msg", "failed to update task status to in_progress", "err", err)
		}

		// 执行解析
		if err := s.fullSyncUsecase.ParseExecell(parseCtx, taskId, filename); err != nil {
			s.log.Log(log.LevelError, "msg", "failed to parse excel", "err", err)
			// 更新任务状态为失败
			s.fullSyncUsecase.UpdateCacheTask(parseCtx, taskId, "completed", 0)
			return
		}

		// 更新任务状态为完成
		if err := s.fullSyncUsecase.UpdateCacheTask(parseCtx, taskId, "completed", 100); err != nil {
			s.log.Log(log.LevelError, "msg", "failed to update task status to completed", "err", err)
		}

		s.log.Log(log.LevelInfo, "msg", "Excel parsing completed for task", "taskId", taskId)
	}()

	return &v1.UploadReply{
		Task: taskId,
	}, nil
}

func (s *AccountService) GetTask(ctx context.Context, in *v1.GetTaskRequest) (*v1.GetTaskReply, error) {

	s.log.Log(log.LevelInfo, "msg", "GetTask", "req", in)

	taskName := in.GetTaskName()
	if taskName == "" {
		return nil, status.Errorf(codes.InvalidArgument, "taskName empty")
	}

	if len(taskName) != 14 {
		return nil, status.Errorf(codes.InvalidArgument, "taskName invalid")
	}

	task, err := s.accounterUsecase.GetTask(ctx, taskName)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "create task failed: %v", err)
	}
	return &v1.GetTaskReply{
		Task: task,
	}, nil
}
