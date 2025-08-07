package biz

import (
	"context"
	"encoding/json"
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

// OrgDeptAdd 部门新增
func (uc *AccounterIncreUsecase) OrgDeptCreate(ctx context.Context, event *clientV2.GenericOpenDingTalkEvent) error {

	log := uc.log.WithContext(ctx)
	log.Infof("OrgDeptCreate data: %v", event.Data)

	if event.Data == nil {
		return nil
	}

	depIds, err := uc.getDeptidsFromDingTalkEvent(event)
	if err != nil {
		return err
	}

	if len(depIds) == 0 {
		log.Info("OrgDeptCreate len(depIds) eq 0")
		return nil
	}

	accessToken, err := uc.dingTalkRepo.GetAccessToken(ctx, "code")
	log.Infof("OrgDeptCreate.GetAccessToken accessToken: %v, err: %v", accessToken, err)
	if err != nil {
		return err
	}
	uc.log.WithContext(ctx).Infof("OrgDeptCreate.FetchDeptDetails accessToken: %v, depIds: %v", accessToken, depIds)
	depts, err := uc.dingTalkRepo.FetchDeptDetails(ctx, accessToken, depIds)
	log.Infof("OrgDeptCreate.FetchDeptDetails accessToken: %v, depIds: %v, err:%v", accessToken, depIds, err)
	if err != nil {
		return err
	}

	err = uc.repo.SaveIncrementDepartments(ctx, depts, nil, nil)
	if err != nil {
		log.Errorf("OrgDeptCreate.SaveIncrementDepartments err: %v", err)
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
		log.Errorf("code %v, not '200'", res.Code)
		return fmt.Errorf("code %s not 200", res.Code)
	}
	return nil

}

// OrgDeptRemove 部门删除
func (uc *AccounterIncreUsecase) OrgDeptRemove(ctx context.Context, event *clientV2.GenericOpenDingTalkEvent) error {
	log := uc.log.WithContext(ctx)
	log.Infof("OrgDeptRemove data: %v", event.Data)

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
	appAccessToken, err := uc.appAuth.GetAccessToken(ctx)
	if err != nil {
		return err
	}
	token := appAccessToken.AccessToken

	// 2. 查询部门ID
	depInfos, err := uc.wps.PostBatchDepartmentsByExDepIds(ctx, token, wps.PostBatchDepartmentsByExDepIdsRequest{
		ExDeptIDs: depIdstr,
	})

	if err != nil {
		log.Errorf("OrgDeptRemove.PostBatchDepartmentsByExDepIds err: %v", err)
		return err
	}
	var depts []*dingtalk.DingtalkDept

	for _, depInfo := range depInfos.Data.Items {
		parentID, err := strconv.ParseInt(depInfo.ParentID, 10, 64)
		if err != nil {
			return err
		}
		id, err := strconv.ParseInt(depInfo.ID, 10, 64)
		if err != nil {
			return err
		}
		dingtalkID, err := strconv.ParseInt(depInfo.ExDeptID, 10, 64)
		if err != nil {
			return err
		}
		detp := &dingtalk.DingtalkDept{
			DeptID:   id,
			ParentID: parentID,
			Order:    int64(depInfo.Order),
			Name:     depInfo.Name,
		}
		detp1 := &dingtalk.DingtalkDept{
			DeptID:   dingtalkID,
			ParentID: parentID,
			Order:    int64(depInfo.Order),
			Name:     depInfo.Name,
		}
		depts = append(depts, detp, detp1)

	}

	err = uc.repo.SaveIncrementDepartments(ctx, nil, nil, depts)
	if err != nil {
		log.Errorf("OrgDeptCreate.SaveIncrementDepartments err: %v", err)
		return err
	}

	res, err := uc.wpsSync.PostEcisaccountsyncIncrement(ctx, appAccessToken.AccessToken, &wps.EcisaccountsyncIncrementRequest{
		ThirdCompanyId: uc.bizConf.ThirdCompanyId,
	})
	if err != nil {
		return err
	}
	if res.Code != "200" {
		log.Errorf("code %v, not '200'", res.Code)
		return fmt.Errorf("code %s not 200", res.Code)
	}
	return nil
}

// OrgDeptModify 部门修改
func (uc *AccounterIncreUsecase) OrgDeptModify(ctx context.Context, event *clientV2.GenericOpenDingTalkEvent) error {

	log := uc.log.WithContext(ctx)
	log.Infof("OrgDeptModify data: %v", event.Data)

	if event.Data == nil {
		return fmt.Errorf("event.Data is nil")
	}

	depIds, err := uc.getDeptidsFromDingTalkEvent(event)
	if err != nil {
		return err
	}

	accessToken, err := uc.dingTalkRepo.GetAccessToken(ctx, "code")
	log.Infof("OrgDeptCreate.GetAccessToken accessToken: %v, err: %v", accessToken, err)
	if err != nil {
		return err
	}
	uc.log.WithContext(ctx).Infof("OrgDeptCreate.FetchDeptDetails accessToken: %v, depIds: %v", accessToken, depIds)
	depts, err := uc.dingTalkRepo.FetchDeptDetails(ctx, accessToken, depIds)
	log.Infof("OrgDeptCreate.FetchDeptDetails accessToken: %v, depIds: %v, err:%v", accessToken, depIds, err)
	if err != nil {
		return err
	}

	err = uc.repo.SaveIncrementDepartments(ctx, nil, nil, depts)
	if err != nil {
		log.Errorf("OrgDeptCreate.SaveIncrementDepartments err: %v", err)
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
		log.Errorf("code %v, not '200'", res.Code)
		return fmt.Errorf("code %s not 200", res.Code)
	}
	return nil

}

// UserAddOrg 用户加入部门
// 1. 加用户
// 2. 加关系
func (uc *AccounterIncreUsecase) UserAddOrg(ctx context.Context, event *clientV2.GenericOpenDingTalkEvent) error {

	log := uc.log.WithContext(ctx)
	log.Infof("UserAddOrg data: %v", event.Data)

	if event.Data == nil {
		return nil
	}

	userIds, err := uc.getUseridsFromDingTalkEvent(event)
	if err != nil {
		return err
	}

	accessToken, err := uc.dingTalkRepo.GetAccessToken(ctx, "code")
	log.Infof("UserAddOrg.GetAccessToken accessToken: %v,userIds:%v err: %v", accessToken, userIds, err)
	if err != nil {
		return err
	}
	uc.log.WithContext(ctx).Infof("UserAddOrg.GetUserDetail userIds: %v", userIds)
	users, err := uc.dingTalkRepo.FetchUserDetail(ctx, accessToken, userIds)
	if err != nil {
		return err
	}

	err = uc.repo.SaveIncrementUsers(ctx, users, nil, nil)
	if err != nil {
		return err
	}

	relations := generateUserDeptRelations(users)

	err = uc.repo.SaveIncrementDepartmentUserRelations(ctx, relations, nil, nil)

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

	log.Infof("UserLeaveOrg.CallEcisaccountsyncIncrement res: %v, err: %v", res, err)

	if err != nil {
		return err
	}
	if res.Code != "200" {
		log.Errorf("code %v, not '200'", res.Code)
		return fmt.Errorf("code %s not 200", res.Code)
	}
	return nil
}

// UserLeaveOrg 用户退出部门
// 1. 减用户
// 2. 减关系
func (uc *AccounterIncreUsecase) UserLeaveOrg(ctx context.Context, event *clientV2.GenericOpenDingTalkEvent) error {

	log := uc.log.WithContext(ctx)
	log.Infof("UserLeaveOrg data: %v", event.Data)
	if event.Data == nil {
		return nil
	}

	userIds, err := uc.getUseridsFromDingTalkEvent(event)
	if err != nil {
		return err
	}

	dingTalkAccessToken, err := uc.dingTalkRepo.GetAccessToken(ctx, "code")
	log.Infof("UserLeaveOrg.GetAccessToken accessToken: %v,userIds:%v err: %v", dingTalkAccessToken, userIds, err)
	if err != nil {
		return err
	}
	uc.log.WithContext(ctx).Infof("UserAddOrg.GetUserDetail userIds: %v", userIds)

	appAccessToken, err := uc.appAuth.GetAccessToken(ctx)
	if err != nil {
		return err
	}

	// 这里找不到
	var users []*dingtalk.DingtalkDeptUser
	// for _, userid := range userIds {
	// 唯一的接口查不到用户的组织ID
	wpsUserInfo, err := uc.wps.PostBatchUsersByExDepIds(ctx, appAccessToken.AccessToken, wps.PostBatchUsersByExDepIdsRequest{
		ExUserIDs: userIds,
		Status:    []string{wps.UserStatusActive, wps.UserStatusNoActive, wps.UserStatusDisabled},
	})
	if err != nil {
		return err
	}
	if len(wpsUserInfo.Data.Items) == 0 {
		log.Warnf("wpsUserInfo.Data.Items is empty, userIds: %v", userIds)
		return fmt.Errorf("wpsUserInfo.Data.Items is empty")
	}

	for _, item := range wpsUserInfo.Data.Items {
		user := &dingtalk.DingtalkDeptUser{
			Userid: item.ExUserId,
			// 退出部门时，部门ID为空,估计出错
			DeptIDList: []int64{},
			Name:       item.UserName,
			Avatar:     item.Avatar,
			Email:      item.Email,
			Mobile:     item.Phone,
		}
		users = append(users, user)
	}

	err = uc.repo.SaveIncrementUsers(ctx, nil, users, nil)
	if err != nil {
		return err
	}
	relations := generateUserDeptRelations(users)

	err = uc.repo.SaveIncrementDepartmentUserRelations(ctx, nil, relations, nil)

	if err != nil {
		return err
	}

	res, err := uc.wpsSync.PostEcisaccountsyncIncrement(ctx, appAccessToken.AccessToken, &wps.EcisaccountsyncIncrementRequest{
		ThirdCompanyId: uc.bizConf.ThirdCompanyId,
	})

	log.Infof("UserLeaveOrg.CallEcisaccountsyncIncrement res: %v, err: %v", res, err)

	if err != nil {
		return err
	}
	if res.Code != "200" {
		log.Errorf("code %v, not '200'", res.Code)
		return fmt.Errorf("code %s not 200", res.Code)
	}
	return nil
}

// func (uc *AccounterIncreUsecase) UserLeaveOrgBak(ctx context.Context, event *clientV2.GenericOpenDingTalkEvent) error {

// 	log := uc.log.WithContext(ctx)
// 	log.Infof("UserLeaveOrg data: %v", event.Data)
// 	if event.Data == nil {
// 		return nil
// 	}

// 	userIds, err := uc.getUseridsFromDingTalkEvent(event)
// 	if err != nil {
// 		return err
// 	}

// 	dingTalkAccessToken, err := uc.dingTalkRepo.GetAccessToken(ctx, "code")
// 	log.Infof("UserLeaveOrg.GetAccessToken accessToken: %v,userIds:%v err: %v", dingTalkAccessToken, userIds, err)
// 	if err != nil {
// 		return err
// 	}
// 	uc.log.WithContext(ctx).Infof("UserAddOrg.GetUserDetail userIds: %v", userIds)
// 	users, err := uc.dingTalkRepo.FetchUserDetail(ctx, dingTalkAccessToken, userIds)
// 	if err != nil {
// 		return err
// 	}
// 	appAccessToken, err := uc.appAuth.GetAccessToken(ctx)
// 	if err != nil {
// 		return err
// 	}
// 	for _, user := range users {
// 		userid := user.Userid
// 		var deptIdstrs []string
// 		for _, deptid := range user.DeptIDList {
// 			deptIdstrs = append(deptIdstrs, strconv.FormatInt(deptid, 10))
// 		}
// 		wpsUserInfo, err := uc.wps.PostBatchUsersByExDepIds(ctx, appAccessToken.AccessToken, wps.PostBatchUsersByExDepIdsRequest{
// 			ExUserIDs: []string{userid},
// 			Status:    []string{wps.UserStatusActive, wps.UserStatusNoActive, wps.UserStatusDisabled},
// 		})
// 		if err != nil {
// 			return err
// 		}
// 		if len(wpsUserInfo.Data.Items) == 0 {
// 			log.Warnf("wpsUserInfo.Data.Items is empty, userid: %v", userid)
// 			continue
// 		}
// 		wpsUserid := wpsUserInfo.Data.Items[0].ID
// 		wpsDeptInfo, err := uc.wps.PostBatchDepartmentsByExDepIds(ctx, appAccessToken.AccessToken, wps.PostBatchDepartmentsByExDepIdsRequest{
// 			ExDeptIDs: deptIdstrs,
// 		})
// 		if err != nil {
// 			return err
// 		}
// 		if len(wpsDeptInfo.Data.Items) == 0 {
// 			log.Warnf("wpsDeptInfo.Data.Items is empty, deptIdstrs: %v", deptIdstrs)
// 			continue
// 		}
// 		for _, dept := range wpsDeptInfo.Data.Items {
// 			wpsDetpId := dept.ID
// 			res, err := uc.wps.PostRomoveUserIdFromDeptId(ctx, appAccessToken.AccessToken, wps.PostRomoveUserIdFromDeptIdRequest{
// 				UserID: wpsUserid,
// 				DeptID: wpsDetpId,
// 			})
// 			if err != nil {
// 				return err
// 			}
// 			if res.Code != 0 {
// 				log.Errorf("code %v, not 0", res.Code)
// 				return fmt.Errorf("code %d not 0", res.Code)
// 			}
// 		}
// 	}

// 	return nil
// }

// UserModifyOrg 用户信息变更（部门变更， 版本没法模拟，暂时未实现）
func (uc *AccounterIncreUsecase) UserModifyOrg(ctx context.Context, event *clientV2.GenericOpenDingTalkEvent) error {

	log := uc.log.WithContext(ctx)
	log.Infof("UserModifyOrg data: %v", event.Data)

	users, err := uc.getUseInfoFromDingTalkEvent(event)
	if err != nil {
		return err
	}
	log.Infof("UserModifyOrg event user to deptuser: %v", users)
	err = uc.repo.SaveIncrementUsers(ctx, nil, nil, users)
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

	log.Infof("UserLeaveOrg.CallEcisaccountsyncIncrement res: %v, err: %v", res, err)

	if err != nil {
		return err
	}
	if res.Code != "200" {
		log.Errorf("code %v, not '200'", res.Code)
		return fmt.Errorf("code %s not 200", res.Code)
	}
	return nil
}
func generateUserDeptRelations(deptUsers []*dingtalk.DingtalkDeptUser) []*dingtalk.DingtalkDeptUserRelation {
	var deptUserRelations []*dingtalk.DingtalkDeptUserRelation
	for _, deptUser := range deptUsers {

		order := make(map[int64]int64, 0)
		if len(deptUser.DeptOrderList) > 0 {
			for _, depIdOrder := range deptUser.DeptOrderList {
				order[depIdOrder.DeptID] = depIdOrder.DeptID
			}
		}

		for _, depId := range deptUser.DeptIDList {
			relation := &dingtalk.DingtalkDeptUserRelation{
				Uid: deptUser.Userid,
				Did: strconv.FormatInt(depId, 10),
				// Order: order,
			}
			if order, ok := order[depId]; ok {
				relation.Order = order
			}
			deptUserRelations = append(deptUserRelations, relation)
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

// map[
//
//	diffInfo:[
//		map[
//			curr:map[email:ian@googla.om hiredDate:2025-08-07 jobNumber:20 name:Ianmodity remark:me telephone: workPlace:北京]
//			prev:map[email:ian@googla.om hiredDate:2025-08-07 jobNumber:20 name:Ian remark:me telephone: workPlace:北京]
//			userid:03301410433273270
//		]
//	]
//	eventId:ebb4c3f1284e45f680ac50ec55b5c5d8
//	optStaffId:manager331
//	timeStamp:1754553836642
//	userId:[03301410433273270]
//
// ]
func (uc *AccounterIncreUsecase) getUseInfoFromDingTalkEvent(event *clientV2.GenericOpenDingTalkEvent) ([]*dingtalk.DingtalkDeptUser, error) {
	uc.log.Infof("getUseInfoFromDingTalkEvent: %v", event.Data)
	data := event.Data

	var userInfos []*dingtalk.DingtalkDeptUser
	jsonData, err := json.Marshal(data)
	uc.log.Infof("getUseInfoFromDingTalkEvent Marshal: %v, err:%v", jsonData, err)
	if err != nil {
		return nil, fmt.Errorf("marshal error: %v", err)
	}

	var modifyInfo dingtalk.UserModifyOrgEventData
	err = json.Unmarshal(jsonData, &modifyInfo)
	uc.log.Infof("getUseInfoFromDingTalkEvent Unmarshal err: %v, modifyInfo: %v", err, modifyInfo)
	if err != nil {
		return nil, fmt.Errorf("unmarshal error: %v", err)
	}
	for _, modifyInfo := range modifyInfo.DiffInfo {
		userInfo := &dingtalk.DingtalkDeptUser{
			Userid:    modifyInfo.Userid,
			Name:      modifyInfo.Curr.Name,
			Email:     modifyInfo.Curr.Email,
			WorkPlace: modifyInfo.Curr.WorkPlace,
			JobNumber: modifyInfo.Curr.JobNumber,
			Mobile:    modifyInfo.Curr.Telephone,
			Remark:    modifyInfo.Curr.Remark,
		}
		userInfos = append(userInfos, userInfo)
	}

	return userInfos, nil
}
