package biz

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"nancalacc/internal/auth"
	"nancalacc/internal/conf"
	"nancalacc/internal/dingtalk"
	"nancalacc/internal/pkg/utils"
	"nancalacc/internal/repository/contracts"
	"nancalacc/internal/wps"
	"strconv"

	//"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/open-dingtalk/dingtalk-stream-sdk-go/clientV2"
)

// GreeterUsecase is a Greeter usecase.
type IncrementalSyncUsecase struct {
	repo         contracts.AccountRepository
	dingTalkRepo dingtalk.Dingtalk
	bizConf      *conf.App
	wpsAppAuth   auth.Authenticator
	wps          wps.Wps
	log          log.Logger
}

// NewGreeterUsecase new a Greeter usecase.
func NewIncrementalSyncUsecase(repo contracts.AccountRepository, dingTalkRepo dingtalk.Dingtalk, wps wps.Wps, logger log.Logger) *IncrementalSyncUsecase {
	wpsAppAuth := auth.NewWpsAppAuthenticator()
	bizConf := conf.Get().GetApp()
	return &IncrementalSyncUsecase{
		repo: repo, dingTalkRepo: dingTalkRepo, bizConf: bizConf,
		wpsAppAuth: wpsAppAuth,
		wps:        wps,
		log:        logger}
}

// OrgDeptAdd 部门新增
func (uc *IncrementalSyncUsecase) OrgDeptCreate(ctx context.Context, event *clientV2.GenericOpenDingTalkEvent) error {

	thirdCompanyId := uc.bizConf.GetThirdCompanyId()
	uc.log.Log(log.LevelInfo, "msg", "OrgDeptCreate", "data", event.Data)

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

	dingTalkAccessToken, err := uc.dingTalkRepo.GetAccessToken(ctx)
	if err != nil {
		return err
	}
	accessToken := dingTalkAccessToken.AccessToken

	uc.log.Log(log.LevelInfo, "msg", "FetchDeptDetails", "accessToken", accessToken, "depIds", depIds)
	depts, err := uc.dingTalkRepo.FetchDeptDetails(ctx, accessToken, depIds)
	uc.log.Log(log.LevelInfo, "msg", "FetchDeptDetails", "accessToken", accessToken, "depIds", depIds, "err", err)
	if err != nil {
		return err
	}

	err = uc.repo.SaveIncrementDepartments(ctx, depts, nil, nil)
	if err != nil {
		uc.log.Log(log.LevelError, "msg", "SaveIncrementDepartments", "err", err)
		return err
	}

	appAccessToken, err := uc.wpsAppAuth.GetAccessToken(ctx)
	if err != nil {
		return err
	}

	res, err := uc.wps.PostEcisaccountsyncIncrement(ctx, appAccessToken.AccessToken, &wps.EcisaccountsyncIncrementRequest{
		ThirdCompanyId: thirdCompanyId,
	})
	if err != nil {
		return err
	}
	if res.Code != "200" {
		uc.log.Log(log.LevelError, "msg", "PostEcisaccountsyncIncrement", "res", res, "err", err)
		return fmt.Errorf("code %s not 200", res.Code)
	}
	return nil

}

// OrgDeptRemove 部门删除
func (uc *IncrementalSyncUsecase) OrgDeptRemove(ctx context.Context, event *clientV2.GenericOpenDingTalkEvent) error {
	uc.log.Log(log.LevelInfo, "msg", "OrgDeptRemove", "data", event.Data)

	thirdCompanyId := uc.bizConf.GetThirdCompanyId()
	if event.Data == nil {
		return nil
	}

	// 1. 已删除的部门ID列表
	depIds, err := uc.getDeptidsFromDingTalkEvent(event)
	if err != nil {
		return err
	}

	if len(depIds) == 0 {
		uc.log.Log(log.LevelWarn, "msg", "OrgDeptRemove", "len(depIds) eq 0")
		return nil
	}

	var depIdstr []string
	for _, depId := range depIds {
		depIdstr = append(depIdstr, strconv.FormatInt(depId, 10))

	}

	if len(depIdstr) == 0 {
		uc.log.Log(log.LevelWarn, "msg", "OrgDeptRemove", "len(depIdstr) eq 0")
		return errors.New("OrgDeptRemove len(depIdstr) eq 0")
	}

	appAccessToken, err := uc.wpsAppAuth.GetAccessToken(ctx)
	if err != nil {
		return err
	}
	token := appAccessToken.AccessToken

	// 2. 查询部门ID
	depInfos, err := uc.wps.PostBatchDepartmentsByExDepIds(ctx, token, wps.PostBatchDepartmentsByExDepIdsRequest{
		ExDeptIDs: depIdstr,
	})

	if err != nil {
		uc.log.Log(log.LevelError, "msg", "PostBatchDepartmentsByExDepIds", "err", err)
		return err
	}
	var deptIDs []string
	tempDeptIDs := make(map[string]int64)
	for _, depInfo := range depInfos.Data.Items {
		deptIDs = append(deptIDs, depInfo.ParentID)
	}

	uc.log.Log(log.LevelInfo, "msg", "OrgDeptRemove", "deptIDs", deptIDs)

	if len(deptIDs) == 0 {
		uc.log.Log(log.LevelWarn, "msg", "OrgDeptRemove", "len(deptIDs) eq 0")
		return errors.New("OrgDeptRemove len(deptIDs) eq 0")
	}
	parentDeptInfos, err := uc.wps.BatchPostDepartments(ctx, token, wps.BatchPostDepartmentsRequest{
		DeptIDs: deptIDs,
	})
	if err != nil {
		uc.log.Log(log.LevelError, "msg", "BatchPostDepartments", "err", err)
		return err
	}

	for _, pdis := range parentDeptInfos.Data.Items {
		extpareId, err := strconv.ParseInt(pdis.ExDeptID, 10, 64)
		if err != nil {
			uc.log.Log(log.LevelError, "msg", "ParseInt", "pdis.ExDeptID", pdis.ExDeptID, "err", err)
		}
		tempDeptIDs[pdis.ID] = extpareId
	}

	var depts []*dingtalk.DingtalkDept

	for _, depInfo := range depInfos.Data.Items {

		dingtalkID, err := strconv.ParseInt(depInfo.ExDeptID, 10, 64)
		if err != nil {
			return err
		}
		parentID, ok := tempDeptIDs[depInfo.ParentID]
		if !ok {
			uc.log.Log(log.LevelError, "msg", "OrgDeptRemove", "not found parentID for DeptID", dingtalkID)
			continue
		}
		detp := &dingtalk.DingtalkDept{
			DeptID:   dingtalkID,
			ParentID: parentID,
			Order:    int64(depInfo.Order),
			Name:     depInfo.Name,
		}

		depts = append(depts, detp)

	}

	err = uc.repo.SaveIncrementDepartments(ctx, nil, depts, nil)
	if err != nil {
		uc.log.Log(log.LevelError, "msg", "SaveIncrementDepartments", "err", err)
		return err
	}

	res, err := uc.wps.PostEcisaccountsyncIncrement(ctx, appAccessToken.AccessToken, &wps.EcisaccountsyncIncrementRequest{
		ThirdCompanyId: thirdCompanyId,
	})
	if err != nil {
		return err
	}
	if res.Code != "200" {
		uc.log.Log(log.LevelError, "msg", "PostEcisaccountsyncIncrement", "res", res, "err", err)
		return fmt.Errorf("code %s not 200", res.Code)
	}
	return nil
}

// OrgDeptModify 部门修改
func (uc *IncrementalSyncUsecase) OrgDeptModify(ctx context.Context, event *clientV2.GenericOpenDingTalkEvent) error {

	uc.log.Log(log.LevelInfo, "msg", "OrgDeptModify", "data", event.Data)

	thirdCompanyId := uc.bizConf.GetThirdCompanyId()
	if event.Data == nil {
		return fmt.Errorf("event.Data is nil")
	}

	depIds, err := uc.getDeptidsFromDingTalkEvent(event)
	if err != nil {
		return err
	}

	dingTalkAccessToken, err := uc.dingTalkRepo.GetAccessToken(ctx)
	if err != nil {
		return err
	}
	accessToken := dingTalkAccessToken.AccessToken

	uc.log.Log(log.LevelInfo, "msg", "FetchDeptDetails", "accessToken", accessToken, "depIds", depIds)
	depts, err := uc.dingTalkRepo.FetchDeptDetails(ctx, accessToken, depIds)
	uc.log.Log(log.LevelInfo, "msg", "FetchDeptDetails", "accessToken", accessToken, "depIds", depIds, "err", err)
	if err != nil {
		return err
	}

	err = uc.repo.SaveIncrementDepartments(ctx, nil, nil, depts)
	if err != nil {
		uc.log.Log(log.LevelError, "msg", "SaveIncrementDepartments", "err", err)
		return err
	}

	appAccessToken, err := uc.wpsAppAuth.GetAccessToken(ctx)
	if err != nil {
		return err
	}

	res, err := uc.wps.PostEcisaccountsyncIncrement(ctx, appAccessToken.AccessToken, &wps.EcisaccountsyncIncrementRequest{
		ThirdCompanyId: thirdCompanyId,
	})
	if err != nil {
		return err
	}
	if res.Code != "200" {
		uc.log.Log(log.LevelError, "msg", "PostEcisaccountsyncIncrement", "res", res, "err", err)
		return fmt.Errorf("code %s not 200", res.Code)
	}
	return nil

}

// UserAddOrg 用户加入部门
// 1. 加用户
// 2. 加关系
func (uc *IncrementalSyncUsecase) UserAddOrg(ctx context.Context, event *clientV2.GenericOpenDingTalkEvent) error {

	uc.log.Log(log.LevelInfo, "msg", "UserAddOrg", "data", event.Data)

	thirdCompanyId := uc.bizConf.GetThirdCompanyId()
	if event.Data == nil {
		return nil
	}

	userIds, err := uc.getUseridsFromDingTalkEvent(event)
	if err != nil {
		return err
	}

	dingTalkAccessToken, err := uc.dingTalkRepo.GetAccessToken(ctx)
	uc.log.Log(log.LevelInfo, "msg", "GetAccessToken", "dingTalkAccessToken", dingTalkAccessToken, "userIds", userIds, "err", err)
	if err != nil {
		return err
	}
	accessToken := dingTalkAccessToken.AccessToken

	uc.log.Log(log.LevelInfo, "msg", "GetUserDetail", "userIds", userIds)
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

	appAccessToken, err := uc.wpsAppAuth.GetAccessToken(ctx)
	if err != nil {
		return err
	}

	res, err := uc.wps.PostEcisaccountsyncIncrement(ctx, appAccessToken.AccessToken, &wps.EcisaccountsyncIncrementRequest{
		ThirdCompanyId: thirdCompanyId,
	})

	log.Infof("UserAddOrg.CallEcisaccountsyncIncrement res: %v, err: %v", res, err)

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
// 2. 减关系 //未自测
func (uc *IncrementalSyncUsecase) UserLeaveOrg(ctx context.Context, event *clientV2.GenericOpenDingTalkEvent) error {

	thirdCompanyId := uc.bizConf.GetThirdCompanyId()
	uc.log.Log(log.LevelInfo, "msg", "UserLeaveOrg", "data", event.Data)
	if event.Data == nil {
		return nil
	}

	userIds, err := uc.getUseridsFromDingTalkEvent(event)
	if err != nil {
		return err
	}

	appAccessToken, err := uc.wpsAppAuth.GetAccessToken(ctx)
	if err != nil {
		return err
	}

	wpsUsers, err := uc.FindWpsUser(ctx, userIds)

	if err != nil {
		return err
	}

	if len(wpsUsers) == 0 {
		uc.log.Log(log.LevelWarn, "msg", "UserLeaveOrg", "wpsUsers is empty, userIds", userIds)
		return fmt.Errorf("wpsUsers is empty")
	}
	err = uc.repo.SaveIncrementUsers(ctx, nil, wpsUsers, nil)
	if err != nil {
		return err
	}
	relations := generateUserDeptRelations(wpsUsers)

	err = uc.repo.SaveIncrementDepartmentUserRelations(ctx, nil, relations, nil)

	if err != nil {
		return err
	}

	res, err := uc.wps.PostEcisaccountsyncIncrement(ctx, appAccessToken.AccessToken, &wps.EcisaccountsyncIncrementRequest{
		ThirdCompanyId: thirdCompanyId,
	})

	uc.log.Log(log.LevelInfo, "msg", "CallEcisaccountsyncIncrement", "res", res, "err", err)

	if err != nil {
		return err
	}
	if res.Code != "200" {
		uc.log.Log(log.LevelError, "msg", "CallEcisaccountsyncIncrement", "res", res, "err", err)
		return fmt.Errorf("code %s not 200", res.Code)
	}
	return nil
}
func (uc *IncrementalSyncUsecase) FindWpsUser(ctx context.Context, userids []string) ([]*dingtalk.DingtalkDeptUser, error) {
	uc.log.Log(log.LevelInfo, "msg", "FindWpsUser", "req", "userids", userids)
	var users []*dingtalk.DingtalkDeptUser
	appAccessToken, err := uc.wpsAppAuth.GetAccessToken(ctx)
	if err != nil {
		return nil, err
	}
	for _, userId := range userids {
		wpsUserInfo, err := uc.wps.PostBatchUsersByExDepIds(ctx, appAccessToken.AccessToken, wps.PostBatchUsersByExDepIdsRequest{
			ExUserIDs: []string{userId},
			Status:    []string{wps.UserStatusActive, wps.UserStatusNoActive, wps.UserStatusDisabled},
		})
		if err != nil {
			return nil, err
		}
		if len(wpsUserInfo.Data.Items) == 0 {
			uc.log.Log(log.LevelWarn, "msg", "FindWpsUser", "wpsUserInfo.Data.Items is empty, userId", userId)
			continue
		}

		wpsUser := wpsUserInfo.Data.Items[0]
		wpsUserid := wpsUser.ID

		wpsDeptInfo, err := uc.wps.GetUserDeptsByUserId(ctx, appAccessToken.AccessToken, wps.GetUserDeptsByUserIdRequest{
			UserID: wpsUserid,
		})
		if err != nil {
			return nil, err
		}
		if len(wpsDeptInfo.Data.Items) == 0 {
			uc.log.Log(log.LevelWarn, "msg", "FindWpsUser", "wpsDeptInfo.Data.Items is empty, userId", userId)
			continue
		}
		user := &dingtalk.DingtalkDeptUser{
			Userid: wpsUser.ExUserId,
			Name:   wpsUser.UserName,
			Email:  wpsUser.Email,
			Mobile: wpsUser.Phone,
		}
		for _, item := range wpsDeptInfo.Data.Items {
			deptId, err := strconv.ParseInt(item.ExDeptID, 10, 64)
			if err != nil {
				return nil, err
			}
			user.DeptIDList = append(user.DeptIDList, deptId)
		}
		users = append(users, user)
	}
	uc.log.Log(log.LevelInfo, "msg", "FindWpsUser", "res", "users", users)
	if len(users) == 0 {
		return nil, fmt.Errorf("wpsUsers is empty")
	}
	for _, user := range users {
		uc.log.Log(log.LevelInfo, "msg", "FindWpsUser", "res", "user", user)
	}
	return users, nil
}
func (uc *IncrementalSyncUsecase) FindDingTalkUser(ctx context.Context, userids []string) ([]*dingtalk.DingtalkDeptUser, error) {
	uc.log.Log(log.LevelInfo, "msg", "FindDingTalkUser", "req", "userids", userids)

	dingTalkAccessToken, err := uc.dingTalkRepo.GetAccessToken(ctx)
	if err != nil {
		return nil, err
	}
	accessToken := dingTalkAccessToken.AccessToken

	users, err := uc.dingTalkRepo.FetchUserDetail(ctx, accessToken, userids)
	uc.log.Log(log.LevelInfo, "msg", "FindDingTalkUser", "res", "users", users, "err", err)

	if err != nil {
		return nil, err
	}

	return users, nil
}

// UserModifyOrg 用户信息变更（有部门变正在实现）
func (uc *IncrementalSyncUsecase) UserModifyOrg(ctx context.Context, event *clientV2.GenericOpenDingTalkEvent) error {

	thirdCompanyId := uc.bizConf.GetThirdCompanyId()
	uc.log.Log(log.LevelInfo, "msg", "UserModifyOrg", "data", event.Data)
	diffUserInfo, _ := uc.getUseInfoFromDingTalkEvent(event)
	diffUserMap := make(map[string]*dingtalk.DingtalkDeptUser, len(diffUserInfo))
	if len(diffUserInfo) > 0 {
		for _, diffUser := range diffUserInfo {
			diffUserMap[diffUser.Userid] = diffUser
			uc.log.Log(log.LevelInfo, "msg", "UserModifyOrg", "diffUser", diffUser)
		}
	}

	userIds, err := uc.getUseridsFromDingTalkEvent(event)

	userIds = utils.RemoveDuplicates(userIds)
	if err != nil {
		return err
	}
	wpsUsersMap := make(map[string]*dingtalk.DingtalkDeptUser)
	wpsUsers, err := uc.FindWpsUser(ctx, userIds)
	if err != nil {
		return err
	}
	if len(wpsUsers) > 0 {
		for _, user := range wpsUsers {
			wpsUsersMap[user.Userid] = user
		}
	}

	for _, u := range wpsUsers {
		log.Infof("UserModifyOrg wpsuser: %+v", u)

		for _, deptId := range u.DeptIDList {
			log.Infof("UserModifyOrg wpsuser deptId: %+v", deptId)
		}
	}

	dingtalkUsers, err := uc.FindDingTalkUser(ctx, userIds)

	if err != nil {
		return err
	}
	uc.log.Log(log.LevelInfo, "msg", "UserModifyOrg", "userIds.size", len(userIds))
	uc.log.Log(log.LevelInfo, "msg", "UserModifyOrg", "wpsUsers.size", len(wpsUsers))
	uc.log.Log(log.LevelInfo, "msg", "UserModifyOrg", "dingtalkUsers.size", len(dingtalkUsers))
	var modfiyUserBaseInfo []*dingtalk.DingtalkDeptUser
	var delRelation []*dingtalk.DingtalkDeptUserRelation
	var addRelation []*dingtalk.DingtalkDeptUserRelation
	var updRelation []*dingtalk.DingtalkDeptUserRelation

	for _, dingtalkUser := range dingtalkUsers { //4个
		var delDepts []int64
		var addDepts []int64
		var updDepts []int64
		finalUser := dingtalkUser
		dingtalkUseridDeptidMap := make(map[string]int64)
		for _, deptId := range dingtalkUser.DeptIDList {
			key := dingtalkUser.Userid + "#" + strconv.FormatInt(deptId, 10)
			dingtalkUseridDeptidMap[key] = deptId
		}
		if wpsUser, ok := wpsUsersMap[dingtalkUser.Userid]; ok { //先找到用户

			for _, deptId := range wpsUser.DeptIDList {
				key1 := wpsUser.Userid + "#" + strconv.FormatInt(deptId, 10)

				if _, ok := dingtalkUseridDeptidMap[key1]; ok {
					//部门修改
					updDepts = append(updDepts, deptId)
					uc.log.Log(log.LevelInfo, "msg", "UserModifyOrg", "部门关系修改", "user.Userid#deptId", key1)
					delete(dingtalkUseridDeptidMap, key1)
				} else {
					//部门删除
					uc.log.Log(log.LevelInfo, "msg", "UserModifyOrg", "部门关系删除", "user.Userid#deptId", key1)
					delDepts = append(delDepts, deptId)

				}
			}

		}
		if len(dingtalkUseridDeptidMap) > 0 {
			for k, deptId := range dingtalkUseridDeptidMap {
				uc.log.Log(log.LevelInfo, "msg", "UserModifyOrg", "部门关系增加", "user.Userid#deptId", k)
				addDepts = append(addDepts, deptId)
			}

		}

		if len(addDepts) > 0 {
			adduser := finalUser
			adduser.DeptIDList = addDepts
			addRelation = append(addRelation, generateUserDeptRelations([]*dingtalk.DingtalkDeptUser{adduser})...)

		}
		if len(delDepts) > 0 {
			deluser := finalUser
			deluser.DeptIDList = delDepts
			delRelation = append(delRelation, generateUserDeptRelations([]*dingtalk.DingtalkDeptUser{deluser})...)
		}
		if len(updDepts) > 0 {
			upduser := finalUser
			upduser.DeptIDList = updDepts
			updRelation = append(updRelation, generateUserDeptRelations([]*dingtalk.DingtalkDeptUser{upduser})...)
		}

		if _, ok := diffUserMap[finalUser.Userid]; ok {
			modfiyUserBaseInfo = append(modfiyUserBaseInfo, finalUser)
		}
	}
	uc.log.Log(log.LevelInfo, "msg", "UserModifyOrg", "部门关系新增", "addRelation", addRelation)
	for i, item := range addRelation {
		uc.log.Log(log.LevelInfo, "msg", "UserModifyOrg", "部门关系新增", "i", i, "item", item)
	}
	uc.log.Log(log.LevelInfo, "msg", "UserModifyOrg", "部门关系删除", "delRelation", delRelation)
	for i, item := range delRelation {
		uc.log.Log(log.LevelInfo, "msg", "UserModifyOrg", "部门关系删除", "i", i, "item", item)
	}

	uc.log.Log(log.LevelInfo, "msg", "UserModifyOrg", "部门关系修改", "updRelation", updRelation)
	for i, item := range updRelation {
		uc.log.Log(log.LevelInfo, "msg", "UserModifyOrg", "部门关系修改", "i", i, "item", item)
	}

	uc.log.Log(log.LevelInfo, "msg", "UserModifyOrg", "基础信息变更", "modfiyUserBaseInfo", modfiyUserBaseInfo)
	for i, item := range modfiyUserBaseInfo {
		uc.log.Log(log.LevelInfo, "msg", "UserModifyOrg", "基础信息变更", "i", i, "item", item)
	}
	if len(modfiyUserBaseInfo) > 0 {
		err = uc.repo.SaveIncrementUsers(ctx, nil, nil, modfiyUserBaseInfo)
		if err != nil {
			return err
		}
	}

	if len(addRelation)+len(delRelation)+len(updRelation) > 0 {
		err = uc.repo.SaveIncrementDepartmentUserRelations(ctx, addRelation, delRelation, updRelation)
		if err != nil {
			return err
		}
	}
	log.Infof("UserModifyOrg.CallEcisaccountsyncIncrement test %s", event.Data)
	//return nil
	appAccessToken, err := uc.wpsAppAuth.GetAccessToken(ctx)
	if err != nil {
		return err
	}
	res, err := uc.wps.PostEcisaccountsyncIncrement(ctx, appAccessToken.AccessToken, &wps.EcisaccountsyncIncrementRequest{
		ThirdCompanyId: thirdCompanyId,
	})

	log.Infof("UserModifyOrg.CallEcisaccountsyncIncrement res: %v, err: %v", res, err)

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

func (uc *IncrementalSyncUsecase) getDeptidsFromDingTalkEvent(event *clientV2.GenericOpenDingTalkEvent) ([]int64, error) {
	uc.log.Log(log.LevelInfo, "msg", "getDeptidsFromDingTalkEvent", "event.Data", event.Data)

	if event.Data == nil {
		return nil, errors.New("getDeptidsFromDingTalkEvent event.Data is nil")
	}
	datamap := event.Data
	var depIds []int64

	deptId, exists := datamap["deptId"]

	if !exists {
		uc.log.Log(log.LevelError, "msg", "getDeptidsFromDingTalkEvent", "not deptId", deptId, "exists", exists)
		return nil, errors.New("getDeptidsFromDingTalkEvent not deptId")
	}

	deptIdSlice, ok := deptId.([]interface{})

	if !ok {
		uc.log.Log(log.LevelError, "msg", "getDeptidsFromDingTalkEvent", "deptId not []interface{}", deptId, "ok", ok)
		return nil, errors.New("deptId not []interface{}")
	}

	for _, item := range deptIdSlice {
		if f, ok := item.(float64); ok {
			depIds = append(depIds, int64(f))
		} else {
			uc.log.Log(log.LevelError, "msg", "getDeptidsFromDingTalkEvent", "deptId not float64", item)
			return nil, errors.New("deptId not float64")
		}
	}
	return depIds, nil
}

func (uc *IncrementalSyncUsecase) getUseridsFromDingTalkEvent(event *clientV2.GenericOpenDingTalkEvent) ([]string, error) {
	uc.log.Log(log.LevelInfo, "msg", "getUseridsFromDingTalkEvent", "event.Data", event.Data)
	if event.Data == nil {
		return nil, errors.New("getUseridsFromDingTalkEvent event.Data is nil")
	}
	datamap := event.Data
	var userIds []string

	userId, exists := datamap["userId"]

	if !exists {
		uc.log.Log(log.LevelError, "msg", "getUseridsFromDingTalkEvent", "not userId", userId, "exists", exists)
		return nil, errors.New("getUseridsFromDingTalkEvent not userId")
	}

	userIdSlice, ok := userId.([]interface{})

	if !ok {
		uc.log.Log(log.LevelError, "msg", "getUseridsFromDingTalkEvent", "userId not []interface{}", userId, "ok", ok)
		return nil, errors.New("userId not []interface{}")
	}

	for _, item := range userIdSlice {
		if f, ok := item.(string); ok {
			userIds = append(userIds, f)
		} else {
			uc.log.Log(log.LevelError, "msg", "getUseridsFromDingTalkEvent", "userId not string", item)
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
func (uc *IncrementalSyncUsecase) getUseInfoFromDingTalkEvent(event *clientV2.GenericOpenDingTalkEvent) ([]*dingtalk.DingtalkDeptUser, error) {
	uc.log.Log(log.LevelInfo, "msg", "getUseInfoFromDingTalkEvent", "event.Data", event.Data)
	data := event.Data

	var userInfos []*dingtalk.DingtalkDeptUser
	jsonData, err := json.Marshal(data)
	uc.log.Log(log.LevelInfo, "msg", "getUseInfoFromDingTalkEvent", "Marshal", "jsonData", jsonData, "err", err)
	if err != nil {
		return nil, fmt.Errorf("marshal error: %v", err)
	}

	var modifyInfo dingtalk.UserModifyOrgEventData
	err = json.Unmarshal(jsonData, &modifyInfo)
	uc.log.Log(log.LevelInfo, "msg", "getUseInfoFromDingTalkEvent", "Unmarshal", "err", err, "modifyInfo", modifyInfo)
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
