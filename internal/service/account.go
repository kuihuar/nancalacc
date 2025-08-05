package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	v1 "nancalacc/api/account/v1"
	"nancalacc/internal/biz"
	"nancalacc/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/xuri/excelize/v2"
)

type AccountService struct {
	v1.UnimplementedAccountServer
	accounterUsecase *biz.AccounterUsecase
	oauth2Usecase    *biz.Oauth2Usecase
	log              *log.Helper
}

func NewAccountService(accounterUsecase *biz.AccounterUsecase, oauth2Usecase *biz.Oauth2Usecase, logger log.Logger) *AccountService {
	return &AccountService{accounterUsecase: accounterUsecase, log: log.NewHelper(logger)}
}

func (s *AccountService) CreateSyncAccount(ctx context.Context, req *v1.CreateSyncAccountRequest) (*v1.CreateSyncAccountReply, error) {
	log := s.log.WithContext(ctx)
	log.Infof("CreateSyncAccount req: %v", req)
	if req.GetTaskName() != "" && len(req.GetTaskName()) != 14 {
		return nil, status.Errorf(codes.InvalidArgument, "invalid taskname: %s", req.GetTaskName())
	}
	// 这里设置传进来的最大分钟数
	ctx, cancel := context.WithTimeout(ctx, 50*time.Minute)
	defer cancel()
	return s.accounterUsecase.CreateSyncAccount(ctx, req)
}
func (s *AccountService) GetSyncAccount(ctx context.Context, req *v1.GetSyncAccountRequest) (*v1.GetSyncAccountReply, error) {

	log.Infof("GetSyncAccount req: %v", req)
	globalConf := conf.Get()
	log.Infof("globalConf: %v", globalConf)
	ctx, cancel := context.WithTimeout(ctx, 50*time.Minute)
	defer cancel()
	return s.accounterUsecase.GetSyncAccount(ctx, req)
}
func (s *AccountService) CancelSyncTask(ctx context.Context, req *v1.CancelSyncAccountRequest) (*emptypb.Empty, error) {
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
	log.Infof("UploadFile req: %v", req.Filename)
	// maxUploadSize := 4
	// if len(req.File) > maxUploadSize {
	// 	return nil, status.Errorf(codes.InvalidArgument,
	// 		"max support file size  %dMB", maxUploadSize/(1<<20))
	// }
	filename := req.Filename

	uploadDir := "/tmp"
	if filename == "" {
		filename = fmt.Sprintf("upload_%d.xlsx", time.Now().UnixNano())
	}

	// 2. 验证确实是Excel文件
	if _, err := excelize.OpenReader(bytes.NewReader(req.File)); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalied Excel file")
	}

	// 3. 创建上传目录
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return nil, status.Errorf(codes.Internal, "create folder failed: %v", err)
	}

	// 4. 生成唯一文件名

	filePath := filepath.Join(uploadDir, filename)

	// 5. 保存文件
	if err := os.WriteFile(filePath, req.GetFile(), 0644); err != nil {
		return nil, status.Errorf(codes.Internal, "safe file failed: %v", err)
	}
	taskId := time.Now().Add(time.Duration(1) * time.Second).Format("20060102150405")
	_, err := s.accounterUsecase.CreateTask(ctx, taskId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "create task failed: %v", err)
	}
	s.accounterUsecase.ParseExecell(ctx, taskId, filePath)
	return &v1.UploadReply{
		Url:  filePath,
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
