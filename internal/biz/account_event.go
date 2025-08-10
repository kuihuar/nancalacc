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
func NewAccounterIncreUsecase(repo AccounterRepo, dingTalkRepo dingtalk.Dingtalk, appAuth auth.Authenticator,
	wpsSync wps.WpsSync, wps wps.Wps, bizConf *conf.Service_Business, logger log.Logger) *AccounterIncreUsecase {
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

	dingTalkAccessToken, err := uc.dingTalkRepo.GetAccessToken(ctx)
	if err != nil {
		return err
	}
	accessToken := dingTalkAccessToken.AccessToken

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
	var deptIDs []string
	tempDeptIDs := make(map[string]int64)
	for _, depInfo := range depInfos.Data.Items {
		deptIDs = append(deptIDs, depInfo.ParentID)
	}

	parentDeptInfos, err := uc.wps.BatchPostDepartments(ctx, token, wps.BatchPostDepartmentsRequest{
		DeptIDs: deptIDs,
	})
	if err != nil {
		log.Errorf("OrgDeptRemove.BatchPostDepartments err: %v", err)
		return err
	}

	for _, pdis := range parentDeptInfos.Data.Items {
		extpareId, err := strconv.ParseInt(pdis.ExDeptID, 10, 64)
		if err != nil {
			log.Errorf("OrgDeptRemove.ParseInt pdis.ExDeptID:%s, err: %v", pdis.ExDeptID, err)
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
			uc.log.Errorf("OrgDeptRemove not found parentID for DeptID: %s", dingtalkID)
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

	dingTalkAccessToken, err := uc.dingTalkRepo.GetAccessToken(ctx)
	if err != nil {
		return err
	}
	accessToken := dingTalkAccessToken.AccessToken

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

	dingTalkAccessToken, err := uc.dingTalkRepo.GetAccessToken(ctx)
	log.Infof("UserAddOrg.GetAccessToken dingTalkAccessToken: %v,userIds:%v err: %v", dingTalkAccessToken, userIds, err)
	if err != nil {
		return err
	}
	accessToken := dingTalkAccessToken.AccessToken

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

	appAccessToken, err := uc.appAuth.GetAccessToken(ctx)
	if err != nil {
		return err
	}

	wpsUsers, err := uc.FindWpsUser(ctx, userIds)

	if err != nil {
		return err
	}

	if len(wpsUsers) == 0 {
		log.Warnf("wpsUsers is empty, userIds: %v", userIds)
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
func (uc *AccounterIncreUsecase) FindWpsUser(ctx context.Context, userids []string) ([]*dingtalk.DingtalkDeptUser, error) {
	uc.log.WithContext(ctx).Infof("FindWpsUser req userids: %v", userids)
	var users []*dingtalk.DingtalkDeptUser
	appAccessToken, err := uc.appAuth.GetAccessToken(ctx)
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
			uc.log.Warnf("wpsUserInfo.Data.Items is empty, userId: %v", userId)
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
			uc.log.Warnf("wpsDeptInfo.Data.Items is empty, userId: %v", userId)
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
	uc.log.WithContext(ctx).Infof("FindWpsUser res users: %v", users)
	if len(users) == 0 {
		return nil, fmt.Errorf("wpsUsers is empty")
	}
	for _, user := range users {
		uc.log.WithContext(ctx).Infof("FindWpsUser res user: %+v", user)
	}
	return users, nil
}
func (uc *AccounterIncreUsecase) FindDingTalkUser(ctx context.Context, userids []string) ([]*dingtalk.DingtalkDeptUser, error) {
	uc.log.WithContext(ctx).Infof("FindDingTalkUser req: %v", userids)

	dingTalkAccessToken, err := uc.dingTalkRepo.GetAccessToken(ctx)
	if err != nil {
		return nil, err
	}
	accessToken := dingTalkAccessToken.AccessToken

	users, err := uc.dingTalkRepo.FetchUserDetail(ctx, accessToken, userids)
	uc.log.WithContext(ctx).Infof("FindDingTalkUser res: %+v, err:%v", users, err)

	if err != nil {
		return nil, err
	}

	return users, nil
}

// UserModifyOrg 用户信息变更（有部门变正在实现）
func (uc *AccounterIncreUsecase) UserModifyOrg(ctx context.Context, event *clientV2.GenericOpenDingTalkEvent) error {

	log := uc.log.WithContext(ctx)
	log.Infof("UserModifyOrg data: %v", event.Data)
	diffUserInfo, _ := uc.getUseInfoFromDingTalkEvent(event)
	diffUserMap := make(map[string]*dingtalk.DingtalkDeptUser, len(diffUserInfo))
	if len(diffUserInfo) > 0 {
		for _, diffUser := range diffUserInfo {
			diffUserMap[diffUser.Userid] = diffUser
			//uc.log.Infof("UserModifyOrg[基础信息变更] user: %v", user)
		}
		//uc.log.Info("UserModifyOrg[基础信息变更] modfiyUserBaseInfo:")
		//modfiyUserBaseInfo = append(modfiyUserBaseInfo, diffUserInfo...)
	}

	userIds, err := uc.getUseridsFromDingTalkEvent(event)
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
	log.Infof("UserModifyOrg userIds.size: %v", len(userIds))
	log.Infof("UserModifyOrg wpsUsers.size: %v", len(wpsUsers))
	log.Infof("UserModifyOrg dingtalkUsers.size: %v", len(dingtalkUsers))
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
					uc.log.Infof("UserModifyOrg[部门关系修改] user.Userid#deptId: %v", key1)
					delete(dingtalkUseridDeptidMap, key1)
				} else {
					//部门删除
					uc.log.Infof("UserModifyOrg[部门关系删除] user.Userid#deptId: %v", key1)
					delDepts = append(delDepts, deptId)

				}
			}

		}
		if len(dingtalkUseridDeptidMap) > 0 {
			for k, deptId := range dingtalkUseridDeptidMap {
				uc.log.Infof("UserModifyOrg[部门关系增加] user.Userid#deptId: %v", k)
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
	uc.log.Info("UserModifyOrg[部门关系新增] addRelation:")
	for i, item := range addRelation {
		uc.log.Infof("UserModifyOrg[部门关系新增] i: %d, item: %+v", i, item)
	}
	uc.log.Info("UserModifyOrg[部门关系删除] delRelation:")
	for i, item := range delRelation {
		uc.log.Infof("UserModifyOrg[部门关系删除] i: %d, item: %+v", i, item)
	}

	uc.log.Info("UserModifyOrg[部门关系修改] updRelation:")
	for i, item := range updRelation {
		uc.log.Infof("UserModifyOrg[部门关系修改] i: %d, iitem: %+v", i, item)
	}

	uc.log.Info("UserModifyOrg[基础信息变更] modfiyUserBaseInfo:")
	for i, item := range modfiyUserBaseInfo {
		uc.log.Infof("UserModifyOrg[基础信息变更] i: %d, item: %+v", i, item)
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
	appAccessToken, err := uc.appAuth.GetAccessToken(ctx)
	if err != nil {
		return err
	}
	res, err := uc.wpsSync.PostEcisaccountsyncIncrement(ctx, appAccessToken.AccessToken, &wps.EcisaccountsyncIncrementRequest{
		ThirdCompanyId: uc.bizConf.ThirdCompanyId,
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
