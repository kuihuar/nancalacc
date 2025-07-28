package biz

import (
	"context"
	"errors"
	"fmt"
	"nancalacc/internal/auth"
	"nancalacc/internal/conf"
	"nancalacc/internal/dingtalk"
	"nancalacc/internal/wps"
	"strconv"

	//"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/open-dingtalk/dingtalk-stream-sdk-go/clientV2"
)

// GreeterUsecase is a Greeter usecase.
type AccounterIncreUsecase struct {
	repo         AccounterRepo
	dingTalkRepo dingtalk.Dingtalk
	appAuth      auth.Authenticator
	wpsSync      wps.WpsSync
	wps          wps.Wps
	bizConf      *conf.Service_Business
	log          *log.Helper
}

// NewGreeterUsecase new a Greeter usecase.
func NewAccounterIncreUsecase(repo AccounterRepo, dingTalkRepo dingtalk.Dingtalk, appAuth auth.Authenticator, wpsSync wps.WpsSync, wps wps.Wps, bizConf *conf.Service_Business, logger log.Logger) *AccounterIncreUsecase {
	return &AccounterIncreUsecase{
		repo: repo, dingTalkRepo: dingTalkRepo,
		appAuth: appAuth, wpsSync: wpsSync, wps: wps,
		bizConf: bizConf,
		log:     log.NewHelper(logger)}
}

func (uc *AccounterIncreUsecase) OrgDeptCreate(ctx context.Context, event *clientV2.GenericOpenDingTalkEvent) error {
	uc.log.Infof("OrgDeptCreate: %v", event.Data)

	uc.log.WithContext(ctx).Info("")
	if event.Data == nil {
		return nil
	}

	depIds, err := uc.getDeptidsFromDingTalkEvent(event)
	if err != nil {
		return err
	}

	if len(depIds) == 0 {
		uc.log.Info("OrgDeptCreate len(depIds) eq 0")
		return nil
	}

	accessToken, err := uc.dingTalkRepo.GetAccessToken(ctx, "code")
	uc.log.WithContext(ctx).Infof("OrgDeptCreate.GetAccessToken accessToken: %v, err: %v", accessToken, err)
	if err != nil {
		return err
	}
	uc.log.WithContext(ctx).Infof("OrgDeptCreate.FetchDeptDetails accessToken: %v, depIds: %v", accessToken, depIds)
	depts, err := uc.dingTalkRepo.FetchDeptDetails(ctx, accessToken, depIds)
	uc.log.WithContext(ctx).Infof("OrgDeptCreate.FetchDeptDetails accessToken: %v, depIds: %v, err:%v", accessToken, depIds, err)
	if err != nil {
		return err
	}

	err = uc.repo.SaveIncrementDepartments(ctx, depts, nil)
	if err != nil {
		uc.log.Errorf("OrgDeptCreate.SaveIncrementDepartments err: %v", err)
		return err
	}

	appAccessToken, err := uc.appAuth.GetAccessToken(ctx)
	if err != nil {
		return err
	}

	res, err := uc.wpsSync.PostEcisaccountsyncIncrement(ctx, appAccessToken.AccessToken, &wps.EcisaccountsyncIncrementRequest{
		ThirdCompanyId: uc.bizConf.ThirdCompanyId,
	})
	if err != nil {
		return err
	}
	if res.Code != "200" {
		uc.log.Errorf("code %v, not '200'", res.Code)
		return fmt.Errorf("code %s not 200", res.Code)
	}
	return nil

}
func (uc *AccounterIncreUsecase) OrgDeptRemove(ctx context.Context, event *clientV2.GenericOpenDingTalkEvent) error {
	uc.log.Infof("OrgDeptCreate: %v", event.Data)
	if event.Data == nil {
		return nil
	}

	// 1. 已删除的部门ID列表
	depIds, err := uc.getDeptidsFromDingTalkEvent(event)
	if err != nil {
		return err
	}

	if len(depIds) == 0 {
		uc.log.Info("OrgDeptCreate len(depIds) eq 0")
		return nil
	}

	var depIdstr []string
	for _, depId := range depIds {
		depIdstr = append(depIdstr, strconv.FormatInt(depId, 10))

	}
	// appAccessToken, err := uc.appAuth.GetAccessToken(ctx)
	// if err != nil {
	// 	return err
	// }
	// token := appAccessToken.AccessToken
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTM2MTU2MTgsImNvbXBfaWQiOiIxIiwiY2xpZW50X2lkIjoiY29tLmFjYy5hc3luYyIsInRrX3R5cGUiOiJhcHAiLCJzY29wZSI6Imtzby5hY2NvdW50c3luYy5zeW5jLGtzby5jb250YWN0LnJlYWQsa3NvLmNvbnRhY3QucmVhZHdyaXRlIiwiY29tcGFueV9pZCI6MSwiY2xpZW50X3ByaW5jaXBhbF9pZCI6IjczIiwiaXNfd3BzMzY1Ijp0cnVlfQ.ZOkiwnZ6f1uW45_sq5uT_ZW3dmA6yCXuKetMaUI7mCw"

	// 2. 查询部门ID父ID
	depInfos, err := uc.wps.PostBatchDepartmentsByExDepIds(ctx, token, wps.PostBatchDepartmentsByExDepIdsRequest{
		ExDeptIDs: depIdstr,
	})

	if err != nil {
		uc.log.Errorf("OrgDeptRemove.BatchGetDepartment err: %v", err)
		return err
	}
	depInfoMap := make(map[string]*wps.WpsDepartmentItem)
	for _, depInfo := range depInfos.Data.Items {
		depInfoMap[depInfo.ID] = &depInfo
	}

	depts := make([]*dingtalk.DingtalkDept, len(depIds))
	for i, depId := range depIds {
		depts[i] = &dingtalk.DingtalkDept{
			DeptID: depId,
		}
	}
	for i, dep := range depts {
		if _, ok := depInfoMap[strconv.FormatInt(dep.DeptID, 10)]; ok {
			depts[i].ParentID = dep.ParentID
		}
	}
	// 3. 存入DB
	err = uc.repo.SaveIncrementDepartments(ctx, nil, depts)
	if err != nil {
		return err
	}
	// 4. 执行同步
	res, err := uc.wpsSync.PostEcisaccountsyncIncrement(ctx, token, &wps.EcisaccountsyncIncrementRequest{
		ThirdCompanyId: uc.bizConf.ThirdCompanyId,
	})
	if err != nil {
		return err
	}
	if res.Code != "200" {
		uc.log.Errorf("code %v, not '200'", res.Code)
		return fmt.Errorf("code %s not 200", res.Code)
	}
	return nil
}
func (uc *AccounterIncreUsecase) OrgDeptModify(ctx context.Context, event *clientV2.GenericOpenDingTalkEvent) error {
	return uc.repo.SaveIncrementDepartments(ctx, nil, nil)
}
func (uc *AccounterIncreUsecase) UserAddOrg(ctx context.Context, event *clientV2.GenericOpenDingTalkEvent) error {
	uc.log.Infof("UserAddOrg: %v", event.Data)
	if event.Data == nil {
		return nil
	}

	userIds, err := uc.getUseridsFromDingTalkEvent(event)
	if err != nil {
		return err
	}

	accessToken, err := uc.dingTalkRepo.GetAccessToken(ctx, "code")
	uc.log.WithContext(ctx).Infof("UserAddOrg.GetAccessToken accessToken: %v,userIds:%v err: %v", accessToken, userIds, err)
	if err != nil {
		return err
	}
	uc.log.WithContext(ctx).Infof("UserAddOrg.GetUserDetail userIds: %v", userIds)
	users, err := uc.dingTalkRepo.FetchUserDetail(ctx, accessToken, userIds)
	if err != nil {
		return err
	}

	err = uc.repo.SaveIncrementUsers(ctx, users, nil)
	if err != nil {
		return err
	}

	relations := generateUserDeptRelations(users)

	err = uc.repo.SaveIncrementDepartmentUserRelations(ctx, relations, nil)

	if err != nil {
		return err
	}

	appAccessToken, err := uc.appAuth.GetAccessToken(ctx)
	if err != nil {
		return err
	}

	res, err := uc.wpsSync.PostEcisaccountsyncIncrement(ctx, appAccessToken.AccessToken, &wps.EcisaccountsyncIncrementRequest{
		ThirdCompanyId: uc.bizConf.ThirdCompanyId,
	})

	uc.log.WithContext(ctx).Infof("UserLeaveOrg.CallEcisaccountsyncIncrement res: %v, err: %v", res, err)

	if err != nil {
		return err
	}
	if res.Code != "200" {
		uc.log.Errorf("code %v, not '200'", res.Code)
		return fmt.Errorf("code %s not 200", res.Code)
	}
	return nil
}

type user struct {
	DingtalkUserId string
	WpsUserId      string
	Dept           struct {
		DingtalkDeptId string
		WpsDeptid      string
	}
}

func (uc *AccounterIncreUsecase) UserLeaveOrg(ctx context.Context, event *clientV2.GenericOpenDingTalkEvent) error {

	uc.log.Infof("UserLeaveOrg: %v", event.Data)
	if event.Data == nil {
		return nil
	}

	uc.log.Info("UserLeaveOrg.getUseridsFromDingTalkEvent")
	// userIds, err := uc.getUseridsFromDingTalkEvent(event)
	// if err != nil {
	// 	return err
	// }

	userIds := []string{"033014104332101118010"}
	// 2. 增量删除用户关系
	// 2.1 用户ID 换 wps 用户ID
	// 2.2 wps 用户ID 查 wps 部门ID
	// 2.3 wps 部门ID 查 部门ID

	uc.log.Info("UserLeaveOrg.GetAccessToken")
	appAccessToken, err := uc.appAuth.GetAccessToken(ctx)
	if err != nil {
		return err
	}

	fmt.Printf("userIds:%+v\n", userIds)
	// 2.1 用户ID 换 wps 用户ID
	uc.log.Info("UserLeaveOrg.PostBatchUsersByExDepIds")
	wpsUserInfo, err := uc.wps.PostBatchUsersByExDepIds(ctx, appAccessToken.AccessToken, wps.PostBatchUsersByExDepIdsRequest{
		ExUserIDs: userIds,
		Status:    []string{wps.UserStatusActive, wps.UserStatusNoActive},
	})
	uc.log.Infof("wpsUserInfo:%+v, err:%+v\n", wpsUserInfo, err)
	if err != nil {
		return err
	}

	wpsUserIds := make([]string, len(wpsUserInfo.Data.Items))
	for i, item := range wpsUserInfo.Data.Items {
		wpsUserIds[i] = item.ID
	}

	input := wps.BatchPostUsersRequest{
		UserIDs:  wpsUserIds,
		WithDept: true,
		Status:   []string{wps.UserStatusActive, wps.UserStatusNoActive},
	}
	uc.log.Infof("UserLeaveOrg.BatchPostUsers input:+%v", input)
	// 2.2 wps 用户ID 查 wps 部门ID
	wpsUserInfoWithDept, err := uc.wps.BatchPostUsers(ctx, appAccessToken.AccessToken, input)

	uc.log.Infof("UserLeaveOrg.BatchPostUser wpsUserInfoWithDept: %+v, err: %+v\n", wpsUserInfoWithDept, err)
	if err != nil {
		return err
	}
	wpsUserDeptIds := make(map[string][]string, len(wpsUserInfoWithDept.Data.Items))

	var deptIds []string
	for _, item := range wpsUserInfoWithDept.Data.Items {
		userId := item.ExUserID

		// 每个部门有多个用户
		for _, dept := range item.Depts {
			wpsUserDeptIds[dept.ID] = append(wpsUserDeptIds[dept.ID], userId)
			deptIds = append(deptIds, dept.ID)
		}

	}

	// 2.3 wps 部门ID 查 部门ID

	uc.log.Info("UserLeaveOrg.BatchPostDepartments")
	wpsDeptInfo, err := uc.wps.BatchPostDepartments(ctx, appAccessToken.AccessToken, wps.BatchPostDepartmentsRequest{
		DeptIDs: deptIds,
	})

	if err != nil {
		return err
	}
	userWithDingtalkDeptIds := make(map[string][]string, 0)

	for _, item := range wpsDeptInfo.Data.Items {
		dingtalkId := item.ExDeptID
		if _, ok := wpsUserDeptIds[item.ID]; ok {
			dingtalkUserIds := wpsUserDeptIds[item.ID]
			userWithDingtalkDeptIds[dingtalkId] = dingtalkUserIds
		}

	}

	users := make([]*dingtalk.DingtalkDeptUser, len(userIds))
	usersmap := make(map[string]*dingtalk.DingtalkDeptUser, len(userIds))
	for i, userId := range userIds {
		users[i] = &dingtalk.DingtalkDeptUser{
			Userid: userId,
		}
	}

	for dingtalkDeptId, dingtalkUserIds := range userWithDingtalkDeptIds {
		dingtalkDeptIdInt, _ := strconv.ParseInt(dingtalkDeptId, 10, 64)
		for _, dingtalkUserId := range dingtalkUserIds {
			if _, ok := usersmap[dingtalkUserId]; ok {
				usersmap[dingtalkUserId] = &dingtalk.DingtalkDeptUser{
					Userid:     dingtalkUserId,
					DeptIDList: []int64{dingtalkDeptIdInt},
				}
			} else {
				usersmap[dingtalkUserId].DeptIDList = append(usersmap[dingtalkUserId].DeptIDList, dingtalkDeptIdInt)
			}

		}
	}

	for _, user := range usersmap {
		users = append(users, user)
	}
	// 1. 增量删除用户
	uc.log.Info("UserLeaveOrg.SaveIncrementUsers")
	err = uc.repo.SaveIncrementUsers(ctx, nil, users)
	if err != nil {
		return err
	}
	relations := generateUserDeptRelations(users)

	uc.log.Info("UserLeaveOrg.SaveIncrementDepartmentUserRelations")
	err = uc.repo.SaveIncrementDepartmentUserRelations(ctx, nil, relations)

	if err != nil {
		return err
	}

	uc.log.Info("UserLeaveOrg.PostEcisaccountsyncIncrement")
	res, err := uc.wpsSync.PostEcisaccountsyncIncrement(ctx, appAccessToken.AccessToken, &wps.EcisaccountsyncIncrementRequest{
		ThirdCompanyId: uc.bizConf.ThirdCompanyId,
	})

	uc.log.WithContext(ctx).Infof("UserLeaveOrg.PostEcisaccountsyncIncrement res: %v, err: %v", res, err)

	return err
}
func (uc *AccounterIncreUsecase) UserModifyOrg(ctx context.Context, event *clientV2.GenericOpenDingTalkEvent) error {
	return uc.repo.SaveIncrementUsers(ctx, nil, nil)
}
func generateUserDeptRelations(deptUsers []*dingtalk.DingtalkDeptUser) []*dingtalk.DingtalkDeptUserRelation {
	var deptUserRelations []*dingtalk.DingtalkDeptUserRelation
	for _, deptUser := range deptUsers {
		order := int(deptUser.DeptOrder)
		if order > 0 {
			order = 1
		} else {
			order = 0
		}
		for _, depId := range deptUser.DeptIDList {

			deptUserRelations = append(deptUserRelations, &dingtalk.DingtalkDeptUserRelation{
				Uid:   deptUser.Userid,
				Did:   strconv.FormatInt(depId, 10),
				Order: order,
			})
		}

	}

	return deptUserRelations
}

func (uc *AccounterIncreUsecase) getDeptidsFromDingTalkEvent(event *clientV2.GenericOpenDingTalkEvent) ([]int64, error) {
	uc.log.Infof("getDeptidsFromDingTalkEvent: %v", event.Data)
	if event.Data == nil {
		return nil, errors.New("getDeptidsFromDingTalkEvent event.Data is nil")
	}
	datamap := event.Data
	var depIds []int64

	deptId, exists := datamap["deptId"]

	if !exists {
		uc.log.Errorf("getDeptidsFromDingTalkEvent not deptId: %v, exists: %v", deptId, exists)
		return nil, errors.New("getDeptidsFromDingTalkEvent not deptId")
	}

	deptIdSlice, ok := deptId.([]interface{})

	if !ok {
		uc.log.Errorf("deptId not []interface{}: %v, ok: %v", deptId, ok)
		return nil, errors.New("deptId not []interface{}")
	}

	for _, item := range deptIdSlice {
		if f, ok := item.(float64); ok {
			depIds = append(depIds, int64(f))
		} else {
			uc.log.Errorf("deptId not float64: %T", item)
			return nil, errors.New("deptId not float64")
		}
	}
	return depIds, nil
}

func (uc *AccounterIncreUsecase) getUseridsFromDingTalkEvent(event *clientV2.GenericOpenDingTalkEvent) ([]string, error) {
	uc.log.Infof("getUseridsFromDingTalkEvent: %v", event.Data)
	if event.Data == nil {
		return nil, errors.New("getUseridsFromDingTalkEvent event.Data is nil")
	}
	datamap := event.Data
	var userIds []string

	userId, exists := datamap["userId"]

	if !exists {
		uc.log.Errorf("getUseridsFromDingTalkEvent not userId: %v, exists: %v", userId, exists)
		return nil, errors.New("getUseridsFromDingTalkEvent not userId")
	}

	userIdSlice, ok := userId.([]interface{})

	if !ok {
		uc.log.Errorf("userId not []interface{}: %v, ok: %v", userId, ok)
		return nil, errors.New("userId not []interface{}")
	}

	for _, item := range userIdSlice {
		if f, ok := item.(string); ok {
			userIds = append(userIds, f)
		} else {
			uc.log.Errorf("userId not string: %T", item)
			return nil, errors.New("userId not string")
		}
	}
	return userIds, nil
}
