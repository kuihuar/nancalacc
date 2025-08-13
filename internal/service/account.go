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
	log              *log.Helper
}

func NewAccountService(accounterUsecase *biz.AccounterUsecase, oauth2Usecase *biz.Oauth2Usecase, fullSyncUsecase *biz.FullSyncUsecase, logger log.Logger) *AccountService {
	limiter := limiter.NewRateLimiter()
	return &AccountService{accounterUsecase: accounterUsecase, oauth2Usecase: oauth2Usecase, fullSyncUsecase: fullSyncUsecase, limiter: limiter, log: log.NewHelper(logger)}
}

func (s *AccountService) CreateSyncAccount(ctx context.Context, req *v1.CreateSyncAccountRequest) (*v1.CreateSyncAccountReply, error) {
	log := s.log.WithContext(ctx)
	log.Infof("CreateSyncAccount req: %v", req)
	if req.GetTaskName() != "" && len(req.GetTaskName()) != 14 {
		return nil, status.Errorf(codes.InvalidArgument, "invalid taskname: %s", req.GetTaskName())
	}
	if req.GetTaskName() == "" {
		taskId := time.Now().Add(time.Duration(1) * time.Second).Format("20060102150405")
		req.TaskName = &taskId
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
	log := s.log.WithContext(ctx)
	log.Infof("CancelSyncTask req: %v", req)
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

	log := s.log.WithContext(ctx)
	log.Infof("GetUserInfo req: %v", req)

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

	log := s.log.WithContext(ctx)
	log.Infof("GetAccessToken req: %v", req)

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
	log := s.log.WithContext(ctx)
	log.Infof("Callback req: %v", req)

	globalConf := conf.Get()
	log.Infof("globalConf: %v", globalConf)
	return &v1.CallbackResponse{
		Status:  "success",
		Message: "success",
	}, nil
}

//	curl -X POST "http://your-server/v1/upload/excel" \
//	  -H "X-Filename: test.xlsx" \
//	  --data-binary "@/path/to/your/file.xlsx"

func (s *AccountService) UploadFile(ctx context.Context, req *v1.UploadRequest) (*v1.UploadReply, error) {

	log := s.log.WithContext(ctx)
	log.Infof("UploadFile req: %v", req)

	taskId := time.Now().Add(time.Duration(1) * time.Second).Format("20060102150405")

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
	filename := filepath.Join(tempDir, taskId+".xlsx")

	// 写入文件
	err := os.WriteFile(filename, req.GetFile(), 0644)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to write file: %v", err)
	}

	// 解析Excel文件 - 使用带超时的context
	parseCtx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Errorf("ParseExecell panic: %v", r)
			}
		}()
		s.fullSyncUsecase.ParseExecell(parseCtx, taskId, filename)
	}()
	// if err != nil {
	// 	return nil, status.Errorf(codes.Internal, "failed to parse excel: %v", err)
	// }
	s.accounterUsecase.CreateCacheTask(ctx, taskId, "")
	return &v1.UploadReply{
		//Message: "Upload success",
		Task: taskId,
	}, nil
}

func (s *AccountService) GetTask(ctx context.Context, in *v1.GetTaskRequest) (*v1.GetTaskReply, error) {

	log := s.log.WithContext(ctx)
	log.Infof("GetAccessToken req: %v", in)

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
