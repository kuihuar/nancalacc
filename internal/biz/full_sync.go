package biz

import (
	"context"
	"errors"
	"fmt"
	v1 "nancalacc/api/account/v1"
	"nancalacc/internal/auth"
	"nancalacc/internal/conf"
	"nancalacc/internal/data/models"
	"nancalacc/internal/dingtalk"
	"nancalacc/internal/wps"
	"os"
	"strconv"
	"sync"
	"time"

	//"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/xuri/excelize/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// GreeterUsecase is a Greeter usecase.
type FullSyncUsecase struct {
	repo         AccounterRepo
	dingTalkRepo dingtalk.Dingtalk
	appAuth      auth.Authenticator
	wps          wps.Wps
	bizConf      *conf.App
	localCache   CacheService
	log          log.Logger
}

// NewGreeterUsecase new a Greeter usecase.
func NewFullSyncUsecase(repo AccounterRepo, dingTalkRepo dingtalk.Dingtalk, wps wps.Wps, cache CacheService, logger log.Logger) *FullSyncUsecase {
	appAuth := auth.NewWpsAppAuthenticator()
	bizConf := conf.Get().GetApp()
	return &FullSyncUsecase{repo: repo, dingTalkRepo: dingTalkRepo,
		appAuth: appAuth, wps: wps, localCache: cache, bizConf: bizConf, log: logger}
}

func (uc *FullSyncUsecase) CreateSyncAccount(ctx context.Context, req *v1.CreateSyncAccountRequest) (*v1.CreateSyncAccountReply, error) {
	uc.log.Log(log.LevelInfo, "msg", "CreateSyncAccount", "req", req)

	taskId := req.GetTaskName()

	_, ok, err := uc.GetCacheTask(ctx, taskId)
	if err != nil {
		return nil, err
	}
	if ok {
		return nil, status.Error(codes.AlreadyExists, "task name "+taskId+" exists")
	}

	companyCfg, users, depts, deptUsers, err := uc.getFullData(ctx)
	if err != nil {
		return nil, err
	}
	err = uc.saveFullData(ctx, companyCfg, users, depts, deptUsers, taskId)
	if err != nil {
		return nil, err
	}

	err = uc.notifyFullSync(ctx, taskId)
	if err != nil {
		return nil, err
	}
	err = uc.createCacheTask(ctx, taskId, "in_progress")
	if err != nil {
		return nil, err
	}

	return &v1.CreateSyncAccountReply{
		TaskId:     taskId,
		CreateTime: timestamppb.Now(),
	}, nil
}

func (uc *FullSyncUsecase) GetSyncAccount(ctx context.Context, req *v1.GetSyncAccountRequest) (*v1.GetSyncAccountReply, error) {
	uc.log.Log(log.LevelInfo, "msg", "GetSyncAccount", "req", req)

	taskCacheInfo, ok, err := uc.localCache.Get(ctx, req.GetTaskId())
	if err != nil {
		return nil, err
	}
	if ok {
		taskInfo, ok1 := taskCacheInfo.(*models.Task)
		if ok1 {
			return &v1.GetSyncAccountReply{
				Status:     v1.GetSyncAccountReply_Status(taskInfo.Progress),
				ActualTime: int64(taskInfo.ActualTime),
				StartTime:  timestamppb.New(taskInfo.CreatedAt),
			}, nil
		}

	}
	return nil, status.Error(codes.NotFound, "task "+req.GetTaskId()+" not found")
}
func (uc *FullSyncUsecase) getFullData(ctx context.Context) (companyCfg *dingtalk.DingtalkCompanyCfg,
	users []*dingtalk.DingtalkDeptUser, depts []*dingtalk.DingtalkDept,
	deptUsers []*dingtalk.DingtalkDeptUserRelation, err error) {
	companyCft := &dingtalk.DingtalkCompanyCfg{
		ThirdCompanyId: uc.bizConf.GetThirdCompanyId(),
		PlatformIds:    uc.bizConf.GetPlatformIds(),
		CompanyId:      uc.bizConf.GetCompanyId(),
	}

	dingTalkAccessToken, err := uc.dingTalkRepo.GetAccessToken(ctx)
	uc.log.Log(log.LevelInfo, "msg", "GetAccessToken", "dingTalkAccessToken", dingTalkAccessToken, "err", err)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	accessToken := dingTalkAccessToken.AccessToken

	depts, err = uc.dingTalkRepo.FetchDepartments(ctx, accessToken)
	uc.log.Log(log.LevelInfo, "msg", "FetchDepartments", "depts", depts, "err", err)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	var deptIds []int64
	for _, dept := range depts {
		deptIds = append(deptIds, dept.DeptID)
	}
	// 1. 从第三方获取用户数据
	users, err = uc.dingTalkRepo.FetchDepartmentUsers(ctx, accessToken, deptIds)

	for _, deptUser := range users {
		uc.log.Log(log.LevelInfo, "msg", "FetchDepartmentUsers", "deptUser", deptUser)
	}
	if err != nil {
		return nil, nil, nil, nil, err
	}
	for _, deptUser := range users {
		order := make(map[int64]int64, 0)
		if len(deptUser.DeptOrderList) > 0 {
			for _, depIdOrder := range deptUser.DeptOrderList {
				order[depIdOrder.DeptID] = depIdOrder.DeptID
			}
		}
		for _, depId := range deptUser.DeptIDList {

			reliation := &dingtalk.DingtalkDeptUserRelation{
				Uid: deptUser.Userid,
				Did: strconv.FormatInt(depId, 10),
			}
			if order, ok := order[depId]; ok {
				reliation.Order = order
			}
			deptUsers = append(deptUsers, reliation)
		}

	}
	return companyCft, users, depts, deptUsers, nil
}
func (uc *FullSyncUsecase) notifyFullSync(ctx context.Context, taskId string) (err error) {
	appAccessToken, err := uc.appAuth.GetAccessToken(ctx)
	if err != nil {
		return err
	}
	resp, err := uc.wps.PostEcisaccountsyncAll(ctx, appAccessToken.AccessToken, &wps.EcisaccountsyncAllRequest{
		TaskId:         taskId,
		ThirdCompanyId: uc.bizConf.GetThirdCompanyId(),
	})
	uc.log.Log(log.LevelInfo, "msg", "PostEcisaccountsyncAll", "resp", resp, "err", err)
	return err
}
func (uc *FullSyncUsecase) saveFullData(ctx context.Context, companyCfg *dingtalk.DingtalkCompanyCfg, users []*dingtalk.DingtalkDeptUser, depts []*dingtalk.DingtalkDept,
	deptUsers []*dingtalk.DingtalkDeptUserRelation, taskId string) (err error) {
	var wg sync.WaitGroup
	errChan := make(chan error, 4)
	wg.Add(4)
	go func() {
		defer wg.Done()
		errChan <- uc.repo.SaveCompanyCfg(ctx, companyCfg)
	}()
	go func() {
		defer wg.Done()
		_, err := uc.repo.SaveDepartments(ctx, depts, taskId)
		errChan <- err
	}()
	go func() {
		defer wg.Done()
		_, err := uc.repo.SaveUsers(ctx, users, taskId)
		errChan <- err
	}()
	go func() {
		defer wg.Done()
		_, err := uc.repo.SaveDepartmentUserRelations(ctx, deptUsers, taskId)
		errChan <- err
	}()
	wg.Wait()
	close(errChan)
	for err := range errChan {
		if err != nil {
			return err
		}
	}
	return nil
}
func (uc *FullSyncUsecase) ParseExecell(ctx context.Context, taskId, filename string) (err error) {
	uc.log.Log(log.LevelInfo, "msg", "ParseExecell", "taskId", taskId, "filename", filename)

	// 更新任务状态为处理中
	if err := uc.updateTaskProgress(ctx, taskId, "in_progress", 10); err != nil {
		uc.log.Log(log.LevelError, "msg", "ParseExecell", "failed to update task progress", "err", err)
	}

	// 确保在函数结束时清理临时文件
	defer func() {
		if err := os.Remove(filename); err != nil {
			uc.log.Log(log.LevelWarn, "msg", "ParseExecell", "failed to remove temp file", "filename", filename, "err", err)
		}
	}()

	f, err := excelize.OpenFile(filename)
	if err != nil {
		uc.log.Log(log.LevelError, "msg", "ParseExecell", "failed to open excel file", "err", err)
		if updateErr := uc.updateTaskProgress(ctx, taskId, "cancelled", 0); updateErr != nil {
			uc.log.Log(log.LevelError, "msg", "ParseExecell", "failed to update task status", "err", updateErr)
		}
		return fmt.Errorf("failed to open excel file: %w", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			uc.log.Log(log.LevelError, "msg", "ParseExecell", "failed to close excel file", "err", err)
		}
	}()

	// 验证Excel文件格式
	if err := uc.validateExcelFormat(f); err != nil {
		uc.log.Log(log.LevelError, "msg", "ParseExecell", "invalid excel format", "err", err)
		if updateErr := uc.updateTaskProgress(ctx, taskId, "cancelled", 0); updateErr != nil {
			uc.log.Log(log.LevelError, "msg", "ParseExecell", "failed to update task status", "err", updateErr)
		}
		return fmt.Errorf("invalid excel format: %w", err)
	}

	// 更新进度到20%
	if err := uc.updateTaskProgress(ctx, taskId, "in_progress", 20); err != nil {
		uc.log.Log(log.LevelError, "msg", "ParseExecell", "failed to update task progress", "err", err)
	}

	// 并发处理工作表
	errChan := make(chan error, 3)
	var wg sync.WaitGroup

	sheets := f.GetSheetList()
	processSheet := map[string]bool{
		"user":            true,
		"department":      true,
		"department_user": true,
	}

	// 统计需要处理的工作表数量
	sheetCount := 0
	for _, sheet := range sheets {
		if processSheet[sheet] {
			sheetCount++
		}
	}

	if sheetCount == 0 {
		uc.log.Log(log.LevelWarn, "msg", "ParseExecell", "no valid sheets found")
		if updateErr := uc.updateTaskProgress(ctx, taskId, "cancelled", 0); updateErr != nil {
			uc.log.Log(log.LevelError, "msg", "ParseExecell", "failed to update task status", "err", updateErr)
		}
		return fmt.Errorf("no valid sheets found in excel file")
	}

	// 并发处理每个工作表
	for _, sheet := range sheets {
		if !processSheet[sheet] {
			continue
		}

		wg.Add(1)
		go func(sheetName string) {
			defer wg.Done()
			if err := uc.processSheet(ctx, taskId, f, sheetName); err != nil {
				uc.log.Log(log.LevelError, "msg", "ParseExecell", "failed to process sheet", "sheet", sheetName, "err", err)
				errChan <- fmt.Errorf("failed to process sheet %s: %w", sheetName, err)
			}
		}(sheet)
	}

	// 等待所有工作表处理完成
	wg.Wait()
	close(errChan)

	// 检查是否有错误
	for err := range errChan {
		if err != nil {
			if updateErr := uc.updateTaskProgress(ctx, taskId, "cancelled", 0); updateErr != nil {
				uc.log.Log(log.LevelError, "msg", "ParseExecell", "failed to update task status", "err", updateErr)
			}
			return err
		}
	}

	// 更新进度到80%
	if err := uc.updateTaskProgress(ctx, taskId, "in_progress", 80); err != nil {
		uc.log.Log(log.LevelError, "msg", "ParseExecell", "failed to update task progress", "err", err)
	}

	// 通知下游服务
	appAccessToken, err := uc.appAuth.GetAccessToken(ctx)
	if err != nil {
		uc.log.Log(log.LevelError, "msg", "ParseExecell", "failed to get access token", "err", err)
		if updateErr := uc.updateTaskProgress(ctx, taskId, "cancelled", 0); updateErr != nil {
			uc.log.Log(log.LevelError, "msg", "ParseExecell", "failed to update task status", "err", updateErr)
		}
		return fmt.Errorf("failed to get access token: %w", err)
	}

	_, err = uc.wps.PostEcisaccountsyncAll(ctx, appAccessToken.AccessToken, &wps.EcisaccountsyncAllRequest{
		TaskId:         taskId,
		ThirdCompanyId: uc.bizConf.GetThirdCompanyId(),
	})
	if err != nil {
		uc.log.Log(log.LevelError, "msg", "ParseExecell", "failed to notify downstream service", "err", err)
		if updateErr := uc.updateTaskProgress(ctx, taskId, "cancelled", 0); updateErr != nil {
			uc.log.Log(log.LevelError, "msg", "ParseExecell", "failed to update task status", "err", updateErr)
		}
		return fmt.Errorf("failed to notify downstream service: %w", err)
	}

	// 更新任务状态为完成
	if err := uc.updateTaskProgress(ctx, taskId, "completed", 100); err != nil {
		uc.log.Log(log.LevelError, "msg", "ParseExecell", "failed to update task progress", "err", err)
	}

	uc.log.Log(log.LevelInfo, "msg", "ParseExecell", "completed successfully", "taskId", taskId)
	return nil
}

// validateExcelFormat 验证Excel文件格式
func (uc *FullSyncUsecase) validateExcelFormat(f *excelize.File) error {
	sheets := f.GetSheetList()
	requiredSheets := map[string]bool{
		"user":            true,
		"department":      true,
		"department_user": true,
	}

	foundSheets := 0
	for _, sheet := range sheets {
		if requiredSheets[sheet] {
			foundSheets++
		}
	}

	if foundSheets == 0 {
		return fmt.Errorf("no required sheets found. Required: user, department, department_user")
	}

	return nil
}

// processSheet 处理单个工作表
func (uc *FullSyncUsecase) processSheet(ctx context.Context, taskId string, f *excelize.File, sheetName string) error {
	uc.log.Log(log.LevelInfo, "msg", "processSheet", "taskId", taskId, "sheet", sheetName)

	rows, err := f.Rows(sheetName)
	if err != nil {
		return fmt.Errorf("failed to get rows for sheet %s: %w", sheetName, err)
	}
	defer rows.Close()

	// 跳过标题行
	rows.Next()

	switch sheetName {
	case "user":
		return uc.transUser(ctx, taskId, rows)
	case "department":
		return uc.transDept(ctx, taskId, rows)
	case "department_user":
		return uc.transUserDept(ctx, taskId, rows)
	default:
		uc.log.Log(log.LevelWarn, "msg", "processSheet", "unknown sheet", "sheet", sheetName)
		return nil
	}
}

// updateTaskProgress 更新任务进度
func (uc *FullSyncUsecase) updateTaskProgress(ctx context.Context, taskId, status string, progress int) error {
	// 这里可以更新缓存中的任务状态和进度
	// 暂时使用日志记录，实际实现可以根据需要更新缓存或数据库
	uc.log.Log(log.LevelInfo, "msg", "updateTaskProgress", "taskId", taskId, "status", status, "progress", progress)
	return nil
}

func (uc *FullSyncUsecase) transUser(ctx context.Context, taskId string, rows *excelize.Rows) (err error) {
	uc.log.Log(log.LevelInfo, "msg", "transUser", "taskId", taskId)

	thirdCompanyId := uc.bizConf.GetThirdCompanyId()
	platformIds := uc.bizConf.GetPlatformIds()
	// 增加批量大小以提高性能
	users := make([]*models.TbLasUser, 0, 500)
	now := time.Now()
	processedCount := 0
	errorCount := 0

	for rows.Next() {
		row, err := rows.Columns()
		if err != nil {
			uc.log.Log(log.LevelError, "msg", "transUser", "failed to get row columns", "err", err)
			errorCount++
			continue
		}

		// 数据验证
		if len(row) < 3 {
			uc.log.Log(log.LevelWarn, "msg", "transUser", "invalid row data", "row", row, "expected_columns", 3)
			errorCount++
			continue
		}

		// 验证必填字段
		if row[0] == "" || row[1] == "" || row[2] == "" {
			uc.log.Log(log.LevelWarn, "msg", "transUser", "missing required fields", "row", row)
			errorCount++
			continue
		}

		users = append(users, &models.TbLasUser{
			TaskID:           taskId,
			ThirdCompanyID:   thirdCompanyId,
			PlatformID:       platformIds,
			Uid:              row[0],
			Account:          row[1],
			NickName:         row[2],
			EmploymentStatus: "active",
			Source:           "sync",
			Ctime:            now,
			Mtime:            now,
			CheckType:        1,
		})

		processedCount++

		// 批量保存
		if len(users) >= 500 {
			if _, err := uc.repo.BatchSaveUsers(ctx, users); err != nil {
				uc.log.Log(log.LevelError, "msg", "transUser", "failed to batch save users", "err", err)
				return fmt.Errorf("failed to batch save users: %w", err)
			}
			uc.log.Log(log.LevelInfo, "msg", "transUser", "batch saved", "count", len(users))
			users = users[:0] // 清空切片（保留底层数组，避免重新分配）
		}
	}

	// 保存剩余数据
	if len(users) > 0 {
		if _, err := uc.repo.BatchSaveUsers(ctx, users); err != nil {
			uc.log.Log(log.LevelError, "msg", "transUser", "failed to save remaining users", "err", err)
			return fmt.Errorf("failed to save remaining users: %w", err)
		}
		uc.log.Log(log.LevelInfo, "msg", "transUser", "final batch saved", "count", len(users))
	}

	if err := rows.Error(); err != nil {
		uc.log.Log(log.LevelError, "msg", "transUser", "rows error", "err", err)
		return fmt.Errorf("rows error: %w", err)
	}

	uc.log.Log(log.LevelInfo, "msg", "transUser", "completed", "processed", processedCount, "errors", errorCount)
	return nil
}
func (uc *FullSyncUsecase) transDept(ctx context.Context, taskId string, rows *excelize.Rows) (err error) {
	uc.log.Log(log.LevelInfo, "msg", "transDept", "taskId", taskId)

	thirdCompanyId := uc.bizConf.GetThirdCompanyId()
	platformIds := uc.bizConf.GetPlatformIds()
	// 增加批量大小以提高性能
	depts := make([]*models.TbLasDepartment, 0, 500)
	now := time.Now()
	processedCount := 0
	errorCount := 0

	for rows.Next() {
		row, err := rows.Columns()
		if err != nil {
			uc.log.Log(log.LevelError, "msg", "transDept", "failed to get row columns", "err", err)
			errorCount++
			continue
		}

		// 数据验证
		if len(row) < 3 {
			uc.log.Log(log.LevelWarn, "msg", "transDept", "invalid row data", "row", row, "expected_columns", 3)
			errorCount++
			continue
		}

		// 验证必填字段
		if row[0] == "" || row[2] == "" {
			uc.log.Log(log.LevelWarn, "msg", "transDept", "missing required fields", "row", row)
			errorCount++
			continue
		}

		depts = append(depts, &models.TbLasDepartment{
			TaskID:         taskId,
			ThirdCompanyID: thirdCompanyId,
			PlatformID:     platformIds,
			Did:            row[0],
			Pid:            row[1],
			Name:           row[2],
			Source:         "sync",
			Ctime:          now,
			Mtime:          now,
			CheckType:      1,
		})

		processedCount++

		// 批量保存
		if len(depts) >= 500 {
			if _, err := uc.repo.BatchSaveDepts(ctx, depts); err != nil {
				uc.log.Log(log.LevelError, "msg", "transDept", "failed to batch save departments", "err", err)
				return fmt.Errorf("failed to batch save departments: %w", err)
			}
			uc.log.Log(log.LevelInfo, "msg", "transDept", "batch saved", "count", len(depts))
			depts = depts[:0] // 清空切片
		}
	}

	// 保存剩余数据
	if len(depts) > 0 {
		if _, err := uc.repo.BatchSaveDepts(ctx, depts); err != nil {
			uc.log.Log(log.LevelError, "msg", "transDept", "failed to save remaining departments", "err", err)
			return fmt.Errorf("failed to save remaining departments: %w", err)
		}
		uc.log.Log(log.LevelInfo, "msg", "transDept", "final batch saved", "count", len(depts))
	}

	if err := rows.Error(); err != nil {
		uc.log.Log(log.LevelError, "msg", "transDept", "rows error", "err", err)
		return fmt.Errorf("rows error: %w", err)
	}

	uc.log.Log(log.LevelInfo, "msg", "transDept", "completed", "processed", processedCount, "errors", errorCount)
	return nil
}
func (uc *FullSyncUsecase) transUserDept(ctx context.Context, taskId string, rows *excelize.Rows) (err error) {
	uc.log.Log(log.LevelInfo, "msg", "transUserDept", "taskId", taskId)

	thirdCompanyId := uc.bizConf.GetThirdCompanyId()
	platformIds := uc.bizConf.GetPlatformIds()
	// 增加批量大小以提高性能
	deptusers := make([]*models.TbLasDepartmentUser, 0, 500)
	now := time.Now()
	processedCount := 0
	errorCount := 0

	for rows.Next() {
		row, err := rows.Columns()
		if err != nil {
			uc.log.Log(log.LevelError, "msg", "transUserDept", "failed to get row columns", "err", err)
			errorCount++
			continue
		}

		// 数据验证
		if len(row) < 2 {
			uc.log.Log(log.LevelWarn, "msg", "transUserDept", "invalid row data", "row", row, "expected_columns", 2)
			errorCount++
			continue
		}

		// 验证必填字段
		if row[0] == "" || row[1] == "" {
			uc.log.Log(log.LevelWarn, "msg", "transUserDept", "missing required fields", "row", row)
			errorCount++
			continue
		}

		deptusers = append(deptusers, &models.TbLasDepartmentUser{
			TaskID:         taskId,
			ThirdCompanyID: thirdCompanyId,
			PlatformID:     platformIds,
			Uid:            row[0],
			Did:            row[1],
			Ctime:          now,
			CheckType:      1,
		})

		processedCount++

		// 批量保存
		if len(deptusers) >= 500 {
			if _, err := uc.repo.BatchSaveDeptUsers(ctx, deptusers); err != nil {
				uc.log.Log(log.LevelError, "msg", "transUserDept", "failed to batch save department users", "err", err)
				return fmt.Errorf("failed to batch save department users: %w", err)
			}
			uc.log.Log(log.LevelInfo, "msg", "transUserDept", "batch saved", "count", len(deptusers))
			deptusers = deptusers[:0] // 清空切片（保留底层数组，避免重新分配）
		}
	}

	// 保存剩余数据
	if len(deptusers) > 0 {
		if _, err := uc.repo.BatchSaveDeptUsers(ctx, deptusers); err != nil {
			uc.log.Log(log.LevelError, "msg", "transUserDept", "failed to save remaining department users", "err", err)
			return fmt.Errorf("failed to save remaining department users: %w", err)
		}
		uc.log.Log(log.LevelInfo, "msg", "transUserDept", "final batch saved", "count", len(deptusers))
	}

	if err := rows.Error(); err != nil {
		uc.log.Log(log.LevelError, "msg", "transUserDept", "rows error", "err", err)
		return fmt.Errorf("rows error: %w", err)
	}

	uc.log.Log(log.LevelInfo, "msg", "transUserDept", "completed", "processed", processedCount, "errors", errorCount)
	return nil
}

func (uc *FullSyncUsecase) createCacheTask(ctx context.Context, taskName, status string) error {

	now := time.Now()
	taskInfo := &models.Task{
		ID:          1,
		Title:       taskName,
		Description: taskName,
		Status:      status,
		CreatorID:   1,
		CreatedAt:   now,
		UpdatedAt:   now,
		DueDate:     now,
		StartDate:   now,
		Progress:    30,
		ActualTime:  0,
	}
	return uc.localCache.Set(ctx, taskName, taskInfo, 300*time.Minute)
}

// func (uc *AccounterUsecase) UpdateCacheTask(ctx context.Context, taskName, status string) error {

// 	cacheKey := prefix + taskName
// 	oldTask, ok, err := uc.localCache.Get(ctx, cacheKey)
// 	if err != nil {
// 		return err
// 	}
// 	var task *models.Task
// 	var startDate time.Time
// 	now := time.Now()
// 	if ok {
// 		task, ok1 := oldTask.(*models.Task)
// 		if ok1 {
// 			startDate = task.StartDate
// 			task.ActualTime = int(now.Sub(startDate).Seconds()) + 20
// 			task.Status = status
// 			task.Progress = 100
// 			task.CompletedAt = now
// 			task.UpdatedAt = now
// 		}
// 	}

// 	if task == nil {
// 		task = &models.Task{
// 			Title:         taskName,
// 			Description:   taskName,
// 			Status:        status,
// 			Progress:      100,
// 			StartDate:     time.Now(),
// 			DueDate:       time.Now().Add(time.Minute * 30),
// 			CompletedAt:   time.Now(),
// 			CreatorID:     99,
// 			EstimatedTime: 10,
// 			ActualTime:    0,
// 		}
// 	}
// 	return uc.localCache.Set(ctx, cacheKey, task, 300*time.Minute)
// }

func (uc *FullSyncUsecase) GetCacheTask(ctx context.Context, taskName string) (*models.Task, bool, error) {

	cacheKey := prefix + taskName
	var task *models.Task
	taskInfo, ok, err := uc.localCache.Get(ctx, cacheKey)
	if err != nil {
		return nil, false, err
	}
	if !ok {
		return nil, false, nil
	}
	task, ok1 := taskInfo.(*models.Task)
	if !ok1 {
		return nil, false, errors.New("type error")
	}
	return task, true, nil

}

func (uc *FullSyncUsecase) CleanSyncAccount(ctx context.Context, taskName string, tags []string) error {

	appAccessToken, err := uc.appAuth.GetAccessToken(ctx)
	if err != nil {
		return err
	}

	if taskName == "phone" {
		users, err := uc.wps.GetCompAllUsers(ctx, appAccessToken.AccessToken, wps.GetCompAllUsersRequest{
			Recursive: true,
			PageSize:  50,
			WithTotal: true,
			Status:    []string{"active", "notactive", "disabled"},
		})
		if err != nil {
			return err
		}

		var deleteUsers []*dingtalk.DingtalkDeptUser
		for _, user := range users.Data.Items {
			uc.log.Log(log.LevelInfo, "msg", "Items", "user", user)

			for _, phone := range tags {
				if user.Phone == phone || user.LoginName == phone {
					deleteUser := &dingtalk.DingtalkDeptUser{
						Userid: user.ExUserID,
						Mobile: user.Phone,
						Name:   user.UserName,
						Email:  user.Email,
					}

					for _, dep := range user.Depts {
						depId, _ := strconv.ParseInt(dep.DeptID, 10, 64)
						deleteUser.DeptIDList = append(deleteUser.DeptIDList, depId)
					}
					deleteUsers = append(deleteUsers, deleteUser)

				}
			}

		}
		uc.log.Log(log.LevelInfo, "msg", "deleteUsers", "deleteUsers", deleteUsers)
		for i, user := range deleteUsers {
			uc.log.Log(log.LevelInfo, "msg", "deleteUsers", "i", i, "user", user)
		}

		err = uc.repo.SaveIncrementUsers(ctx, nil, deleteUsers, nil)
		if err != nil {
			uc.log.Log(log.LevelError, "msg", "OrgDeptCreate.SaveIncrementDepartments", "err", err)
			return err
		}

		res, err := uc.wps.PostEcisaccountsyncIncrement(ctx, appAccessToken.AccessToken, &wps.EcisaccountsyncIncrementRequest{
			ThirdCompanyId: uc.bizConf.GetThirdCompanyId(),
		})

		uc.log.Log(log.LevelInfo, "msg", "UserLeaveOrg.CallEcisaccountsyncIncrement", "res", res, "err", err)

		if err != nil {
			return err
		}
		if res.Code != "200" {
			uc.log.Log(log.LevelError, "msg", "UserLeaveOrg.CallEcisaccountsyncIncrement", "res", res, "err", err)
			return fmt.Errorf("code %s not 200", res.Code)
		}
	}

	if taskName == "dept" {
		//查询根部门
		rootDept, err := uc.wps.GetDepartmentRoot(ctx, appAccessToken.AccessToken, wps.GetDepartmentRootRequest{})
		if err != nil {
			return err
		}
		uc.log.Log(log.LevelInfo, "msg", "rootDept", "rootDept", rootDept)

		// 2. 查询部门下的子部门(要递归)
		allDepts, err := uc.wps.GetDeptChildren(ctx, appAccessToken.AccessToken, wps.GetDeptChildrenRequest{
			DeptID:    rootDept.Data.ID,
			Recursive: true,
			PageSize:  50,
			WithTotal: true,
		})
		if err != nil {
			return err
		}
		uc.log.Log(log.LevelInfo, "msg", "children", "allDepts", allDepts)

		var deleteDepts []*dingtalk.DingtalkDept
		//删除部门除了根部门
		for _, dept := range allDepts.Data.Items {
			if dept.ID == rootDept.Data.ID {
				continue
			}
			deptId, _ := strconv.ParseInt(dept.ExDeptID, 10, 64)

			deptDetail, err := uc.wps.BatchPostDepartments(ctx, appAccessToken.AccessToken, wps.BatchPostDepartmentsRequest{
				DeptIDs: []string{dept.ParentID},
			})
			if err != nil {
				return err
			}
			parentId, _ := strconv.ParseInt(deptDetail.Data.Items[0].ExDeptID, 10, 64)

			for _, tag := range tags {

				if tag == dept.Name {
					// 这里要找父级节点的extid
					detp := &dingtalk.DingtalkDept{
						DeptID:   deptId,
						ParentID: parentId,
						Order:    int64(dept.Order),
						Name:     dept.Name,
					}
					deleteDepts = append(deleteDepts, detp)

				}
			}

		}

		uc.log.Log(log.LevelInfo, "msg", "deleteDepts", "deleteDepts", deleteDepts)
		for i, dept := range deleteDepts {
			uc.log.Log(log.LevelInfo, "msg", "deleteDepts", "i", i, "dept", dept)
		}

		err = uc.repo.SaveIncrementDepartments(ctx, nil, deleteDepts, nil)
		if err != nil {
			uc.log.Log(log.LevelError, "msg", "OrgDeptCreate.SaveIncrementDepartments", "err", err)
			return err
		}

		res, err := uc.wps.PostEcisaccountsyncIncrement(ctx, appAccessToken.AccessToken, &wps.EcisaccountsyncIncrementRequest{
			ThirdCompanyId: uc.bizConf.GetThirdCompanyId(),
		})

		uc.log.Log(log.LevelInfo, "msg", "UserLeaveOrg.CallEcisaccountsyncIncrement", "res", res, "err", err)

		if err != nil {
			return err
		}
		if res.Code != "200" {
			uc.log.Log(log.LevelError, "msg", "UserLeaveOrg.CallEcisaccountsyncIncrement", "res", res, "err", err)
			return fmt.Errorf("code %s not 200", res.Code)
		}

	}

	return nil

}
