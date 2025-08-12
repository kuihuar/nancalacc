package dingtalk

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"nancalacc/internal/auth"
	"nancalacc/internal/conf"
	"nancalacc/internal/pkg/utils"
	"nancalacc/pkg/httputil"
	"sync"
	"time"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dingtalkcontact_1_0 "github.com/alibabacloud-go/dingtalk/contact_1_0"
	dingtalkoauth2_1_0 "github.com/alibabacloud-go/dingtalk/oauth2_1_0"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/go-kratos/kratos/v2/log"
)

type dingTalkRepo struct {
	data      *conf.Auth_Dingtalk
	log       *log.Helper
	tokenAuth auth.Authenticator
	// unifiedAuthService auth.UnifiedAuthService
	dingtalkCli *dingtalkoauth2_1_0.Client

	dingtalkCliContact *dingtalkcontact_1_0.Client
}

func NewDingTalkRepo(logger log.Logger) Dingtalk {

	config := &openapi.Config{
		Protocol: tea.String("https"),
		RegionId: tea.String("central"),
	}
	client, err := dingtalkoauth2_1_0.NewClient(config)
	if err != nil {
		fmt.Printf("NewClient err: %v", err)
		//return nil, cleanup, err
		logger.Log(log.LevelError, "NewClientErr", err)
	}

	clientContact, err := dingtalkcontact_1_0.NewClient(config)
	if err != nil {
		fmt.Printf("NewClient err: %v", err)
		//return nil, cleanup, err
		logger.Log(log.LevelError, "NewClientErr", err)
	}
	tokenAuth := auth.NewDingTalkAuthenticator()

	return &dingTalkRepo{
		dingtalkCli:        client,
		tokenAuth:          tokenAuth,
		dingtalkCliContact: clientContact,
		data:               conf.Get().GetAuth().GetDingtalk(),
		log:                log.NewHelper(log.With(logger, "module", "data/dingtalk")),
	}
}

func (r *dingTalkRepo) GetAccessToken(ctx context.Context) (*auth.AccessTokenResp, error) {

	return r.tokenAuth.GetAccessToken(ctx)
	// 	log := r.log.WithContext(ctx)
	// 	log.Info("GetAccessToken")

	// 	request := &dingtalkoauth2_1_0.GetAccessTokenRequest{
	// 		AppKey:    tea.String(r.data.AppKey),
	// 		AppSecret: tea.String(r.data.AppSecret),
	// 	}

	// 	var accessToken dingtalkoauth2_1_0.GetAccessTokenResponseBody

	// 	tryErr := func() error {
	// 		defer func() {
	// 			if r := tea.Recover(recover()); r != nil {
	// 				err := r
	// 				fmt.Printf("恢复的错误: %v\n", err)
	// 			}
	// 		}()

	// 		response, err := r.dingtalkCli.GetAccessToken(request)
	// 		if err != nil {
	// 			return err
	// 		}

	// 		accessToken = *response.Body
	// 		return nil
	// 	}()

	// 	if tryErr != nil {
	// 		// 处理错误
	// 		var sdkErr = &tea.SDKError{}
	// 		if _t, ok := tryErr.(*tea.SDKError); ok {
	// 			sdkErr = _t
	// 		} else {
	// 			sdkErr.Message = tea.String(tryErr.Error())
	// 		}

	// 		if !tea.BoolValue(util.Empty(sdkErr.Code)) && !tea.BoolValue(util.Empty(sdkErr.Message)) {
	// 			return accessToken, fmt.Errorf("获取access_token失败: [%s] %s", *sdkErr.Code, *sdkErr.Message)
	// 		}
	// 		return accessToken, fmt.Errorf("获取access_token失败: %s", *sdkErr.Message)
	// 	}

	// return accessToken, nil
}
func (r *dingTalkRepo) FetchDepartments(ctx context.Context, token string) ([]*DingtalkDept, error) {

	log := r.log.WithContext(ctx)
	log.Infof("FetchDepartments token:%s", token)

	var deptList []*DingtalkDept

	var deptIdlist []int64
	var baseDeptId int64 = 1
	// 1. 获取子部门ID列表（所有）
	deptIdsLevelOne, err := r.getDeptIds(ctx, token, baseDeptId)
	log.Infof("FetchDepartments deptIdsLevelOne: %v, err: %v", deptIdsLevelOne, err)
	if err != nil {
		return nil, err
	}

	deptIdlist = append(deptIdlist, baseDeptId)

	log.Infof("FetchDepartments deptIdlist: %v", deptIdlist)
	if len(deptIdsLevelOne) > 0 {
		log.Info("len(deptIdsLevelOne) > 0")
		deptIdlist = append(deptIdlist, deptIdsLevelOne...)
		deptIdsLeveltwo, err := r.getDeptIdsConcurrent(ctx, token, deptIdsLevelOne)

		log.Infof("FetchDepartments deptIdsLeveltwo: %v, err: %v", deptIdsLeveltwo, err)

		if err != nil {
			log.Error("getDeptIdsConcurrent failed, err: %v", err)
		}
		if len(deptIdsLeveltwo) > 0 {
			deptIdlist = append(deptIdlist, deptIdsLeveltwo...)
		}
	}

	log.Info("FetchDepartments.deptIdlist: %v", deptIdlist)
	// 2. 获取子部门详情
	deptList, err = r.FetchDeptDetails(ctx, token, deptIdlist)
	log.Infof("FetchDepartments deptList: %v, err: %v", deptList, err)
	if err != nil {
		return nil, err
	}
	return deptList, nil
}
func (r *dingTalkRepo) getDeptIds(ctx context.Context, token string, deptId int64) ([]int64, error) {

	log := r.log.WithContext(ctx)
	log.Infof("getDeptIds token:%s deptId: %v", token, deptId)

	uri := fmt.Sprintf("%s/topapi/v2/department/listsubid?access_token=%s", r.data.Endpoint, token)
	input := &ListDeptIDRequest{
		DeptID: deptId,
	}
	jsonData, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	bs, err := httputil.PostJSON(uri, jsonData, time.Second*10)
	if err != nil {
		return nil, err
	}

	//r.log.Info("FetchAccounts.deptList: %v, err: %v", string(bs), err)

	var deptIDResponse *ListDeptIDResponse
	if err = json.Unmarshal(bs, &deptIDResponse); err != nil {
		return nil, err
	}
	if deptIDResponse.Errcode != 0 {
		return nil, fmt.Errorf("钉钉API返回错误: %s, errcode: %d", deptIDResponse.Errmsg, deptIDResponse.Errcode)
	}
	if deptIDResponse.Result.DeptIDList == nil {
		return nil, fmt.Errorf("钉钉API返回错误: %s, errcode: %d", deptIDResponse.Errmsg, deptIDResponse.Errcode)
	}
	deptIdlist := deptIDResponse.Result.DeptIDList
	return deptIdlist, nil
}
func (r *dingTalkRepo) getDeptIdsConcurrent(ctx context.Context, token string, deptIds []int64) ([]int64, error) {

	uri := fmt.Sprintf("%s/topapi/v2/department/listsubid?access_token=%s", r.data.Endpoint, token)

	r.log.Info("getDeptIdsConcurrent deptIds: %v, uri: %v", deptIds, uri)
	sem := make(chan struct{}, r.data.MaxConcurrent)
	deptList := make([]int64, 0)
	var mu sync.Mutex

	var wg sync.WaitGroup

	for _, deptId := range deptIds {
		wg.Add(1)

		select {
		case sem <- struct{}{}:
		case <-ctx.Done():
			wg.Done()
			continue
		}

		// 启动goroutine处理任务
		go func(id int64) {
			defer func() {
				<-sem     // 释放信号量
				wg.Done() // 通知任务完成
			}()

			input := &ListDeptIDRequest{
				DeptID: id,
			}
			jsonData, err := json.Marshal(input)
			if err != nil {
				r.log.Errorf("getDeptIdsConcurrent.jsonData: %s, err: %v", string(jsonData), err)
				return
			}

			bs, err := httputil.PostJSON(uri, jsonData, time.Second*10)
			if err != nil {
				r.log.Errorf("getDeptIdsConcurrent.PostJSON: %s, err: %v", string(jsonData), err)
				return
			}
			var deptIDResponse *ListDeptIDResponse
			if err = json.Unmarshal(bs, &deptIDResponse); err != nil {
				r.log.Errorf("getDeptIdsConcurrent.Unmarshal: %s, err: %v", string(bs), err)
				return
			}
			if deptIDResponse.Errcode != 0 {
				r.log.Errorf("钉钉API返回错误: %s, errcode: %d", deptIDResponse.Errmsg, deptIDResponse.Errcode)
				return
			}
			if deptIDResponse.Result.DeptIDList == nil {
				r.log.Errorf("钉钉API返回错误: %s, errcode: %d", deptIDResponse.Errmsg, deptIDResponse.Errcode)
				return
			}
			deptIdlist := deptIDResponse.Result.DeptIDList

			mu.Lock()
			deptList = append(deptList, deptIdlist...)
			mu.Unlock()
		}(deptId) // 传递当前deptId值
	}
	wg.Wait()

	return deptList, nil
}

func (r *dingTalkRepo) FetchDeptDetails(ctx context.Context, token string, deptIds []int64) ([]*DingtalkDept, error) {
	log := r.log.WithContext(ctx)
	log.Infof("FetchDeptDetails token:%s depIds: %v", token, deptIds)

	uriDetail := fmt.Sprintf("%s/topapi/v2/department/get?access_token=%s", r.data.Endpoint, token)

	sem := make(chan struct{}, r.data.MaxConcurrent)
	results := make(chan *DingtalkDept, len(deptIds))
	//errChan := make(chan error, 1)

	var wg sync.WaitGroup

	for _, deptId := range deptIds {
		wg.Add(1)

		select {
		case sem <- struct{}{}:
		case <-ctx.Done():
			wg.Done()
			continue
		}

		// 启动goroutine处理任务
		go func(id int64) {
			defer func() {
				<-sem     // 释放信号量
				wg.Done() // 通知任务完成
			}()

			input := &DingtalkDeptRequest{
				DeptID: id,
			}
			jsonData, err := json.Marshal(input)
			if err != nil {
				r.log.Errorf("FetchDeptDetails.jsonData: %v, err: %v", string(jsonData), err)
				//errChan <- err
				return
			}

			bs, err := httputil.PostJSON(uriDetail, jsonData, time.Second*10)
			//r.log.Infof(">>>>FetchDeptDetails.bs: %s, err: %v\n", string(bs), err)
			if err != nil {
				r.log.Errorf("FetchDeptDetails.PostJSON: %v, err: %v", string(jsonData), err)
				//errChan <- err
				return
			}
			var deptResponse *DingtalkDeptResponse
			if err = json.Unmarshal(bs, &deptResponse); err != nil {
				r.log.Errorf("FetchDeptDetails.Unmarshal: %v, err: %v", string(bs), err)
				//errChan <- err
				return
			}
			if deptResponse.Errcode != 0 {
				r.log.Errorf("FetchDeptDetails.Errcode: %v, err: %v", deptResponse.Errcode, deptResponse.Errmsg)
				//errChan <- err
				return
			}
			results <- &deptResponse.Result
		}(deptId) // 传递当前deptId值
	}
	wg.Wait()

	close(results)
	//close(errChan)
	var deptList []*DingtalkDept
	for dept := range results {
		deptList = append(deptList, dept)
	}

	return deptList, nil

}
func (r *dingTalkRepo) FetchDepartmentUsers(ctx context.Context, token string, deptIds []int64) ([]*DingtalkDeptUser, error) {
	log := r.log.WithContext(ctx)
	log.Infof("FetchDepartmentUsers: %v, %v", token, deptIds)

	// 服务端API.通讯录管理.用户管理.获取部门用户详情
	//maxConcurrent := 10
	sem := make(chan struct{}, r.data.MaxConcurrent)
	results := make(chan *DingtalkDeptUser, len(deptIds))
	// := make(chan error, 1)

	var wg sync.WaitGroup

	for _, deptId := range deptIds {
		wg.Add(1)

		select {
		case sem <- struct{}{}:
		case <-ctx.Done():
			wg.Done()
			continue
		}

		// 启动goroutine处理任务
		go func(id int64) {
			defer func() {
				<-sem     // 释放信号量
				wg.Done() // 通知任务完成
			}()

			for {
				userList, cursor, err := r.getUserListByDepId(ctx, token, id)
				if err != nil {
					r.log.Errorf("FetchDepartmentUsers.getUserListByDepId: %v, err: %v", id, err)
					//errChan <- err
					return
				}
				for _, user := range userList {
					results <- user
				}
				if cursor == 0 {
					break
				}

			}
		}(deptId)
	}
	wg.Wait()

	close(results)
	//close(errChan)
	var userList []*DingtalkDeptUser
	usersMap := make(map[string]*DingtalkDeptUser)
	for user := range results {
		// log.Infof("FetchDepartmentUsers results user: %+v", user)
		if _, ok := usersMap[user.Userid]; ok {
			usersMap[user.Userid].DeptIDList = append(usersMap[user.Userid].DeptIDList, user.DeptIDList...)
			usersMap[user.Userid].DeptIDList = utils.RemoveDuplicates(usersMap[user.Userid].DeptIDList)

		}
		usersMap[user.Userid] = user

	}
	for _, u := range usersMap {
		// log.Infof("FetchDepartmentUsers usersMap user: %+v", u)
		userList = append(userList, u)
	}
	return userList, nil
}
func (r *dingTalkRepo) getUserListByDepId(ctx context.Context, token string, deptId int64) ([]*DingtalkDeptUser, int64, error) {
	log := r.log.WithContext(ctx)
	log.Infof("getUserListByDepId token: %v,deptId: %v", token, deptId)
	// 发送post请求
	var cursor int64 = 0
	uri := fmt.Sprintf("%s/topapi/v2/user/list?access_token=%s", r.data.Endpoint, token)
	input := &ListDeptUserRequest{
		DeptID: deptId,
		Cursor: cursor,
		Size:   100,
	}
	jsonData, err := json.Marshal(input)
	if err != nil {
		return nil, 0, err
	}

	//log.Info("getUserListByDepId.uri: %v, input: %v, jsonData: %v", uri, input, string(jsonData))

	bs, err := httputil.PostJSON(uri, jsonData, time.Second*10)
	log.Info("getUserListByDepId.body: %v, err: %v", string(bs), err)
	if err != nil {
		return nil, 0, err
	}

	// 打印响应体
	//fmt.Println(string(bs))

	var userResponse ListDeptUserResponse
	if err = json.Unmarshal(bs, &userResponse); err != nil {
		return nil, 0, err
	}
	if userResponse.Errcode != 0 {
		return nil, 0, fmt.Errorf("钉钉API返回错误: %s, errcode: %v", userResponse.Errmsg, userResponse.Errcode)
	}

	var userList []*DingtalkDeptUser
	if userResponse.Result.List != nil {
		userList = make([]*DingtalkDeptUser, 0, len(userResponse.Result.List))
		for _, user := range userResponse.Result.List {
			userList = append(userList, &user)
		}
	}
	if userResponse.Result.HasMore {
		return userList, userResponse.Result.NextCursor, nil
	}
	return userList, 0, nil
}
func (r *dingTalkRepo) GetUserAccessToken(ctx context.Context, code string) (*AuthResponse, error) {

	log := r.log.WithContext(ctx)
	log.Infof("GetUserAccessToken code: %v", code)

	getUserTokenRequest := &dingtalkoauth2_1_0.GetUserTokenRequest{

		ClientId:     tea.String(r.data.AppKey),
		ClientSecret: tea.String(r.data.AppSecret),
		Code:         tea.String(code),
		//RefreshToken: tea.String("abcd"),
		GrantType: tea.String("authorization_code"),
	}

	// var accessToken string
	var response *dingtalkoauth2_1_0.GetUserTokenResponse
	var err error
	tryErr := func() (_e error) {
		defer func() {
			if er := tea.Recover(recover()); er != nil {
				_e = er
			}
		}()
		response, err = r.dingtalkCli.GetUserToken(getUserTokenRequest)
		if err != nil {
			return err
		}

		return nil
	}()

	if tryErr != nil {
		var err = &tea.SDKError{}
		if _t, ok := tryErr.(*tea.SDKError); ok {
			err = _t
		} else {
			err.Message = tea.String(tryErr.Error())
		}
		if !tea.BoolValue(util.Empty(err.Code)) && !tea.BoolValue(util.Empty(err.Message)) {
			// err 中含有 code 和 message 属性，可帮助开发定位问题
		}

	}
	if response.StatusCode != nil && *response.StatusCode != 200 {
		return nil, err
	}
	if response.Body == nil {
		return nil, err
	}
	tokenAuthResp := &AuthResponse{}
	if response.Body.AccessToken == nil {
		return nil, errors.New("response.Body.AccessToken is nil")
	}
	tokenAuthResp.AccessToken = *response.Body.AccessToken
	tokenAuthResp.RefreshToken = *response.Body.RefreshToken
	tokenAuthResp.ExpireIn = int(*response.Body.ExpireIn)

	return tokenAuthResp, nil
}
func (r *dingTalkRepo) GetUserInfo(ctx context.Context, token, unionId string) (*DingTalkUserInfo, error) {

	log := r.log.WithContext(ctx)
	log.Infof("GetUserInfo token: %v, unionId %s", token, unionId)

	getUserHeaders := &dingtalkcontact_1_0.GetUserHeaders{}
	getUserHeaders.XAcsDingtalkAccessToken = tea.String(token)
	var response *dingtalkcontact_1_0.GetUserResponse
	var err error
	tryErr := func() (_e error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				_e = r
			}
		}()
		response, err = r.dingtalkCliContact.GetUserWithOptions(tea.String(unionId), getUserHeaders, &util.RuntimeOptions{})

		r.log.WithContext(ctx).Info("response: %v, error: %v", response, err)

		if err != nil {
			return err
		}
		if response.Body == nil {
			return err
		}

		return nil
	}()

	if tryErr != nil {
		var err = &tea.SDKError{}
		if _t, ok := tryErr.(*tea.SDKError); ok {
			err = _t
		} else {
			err.Message = tea.String(tryErr.Error())
		}
		if !tea.BoolValue(util.Empty(err.Code)) && !tea.BoolValue(util.Empty(err.Message)) {
			log.Errorf("GetUserInfo error: %v", err)
			// err 中含有 code 和 message 属性，可帮助开发定位问题
		}

	}

	log.Infof("GetUserInfo response: %v", response)

	return &DingTalkUserInfo{
		UnionId: *response.Body.UnionId,
		Nick:    *response.Body.Nick,
	}, nil
}

func (r *dingTalkRepo) GetUseridByUnionid(ctx context.Context, token, unionid string) (string, error) {

	log := r.log.WithContext(ctx)

	log.Info("GetUseridByUnionid token: %v,unionid %v", token, unionid)
	uri := fmt.Sprintf("%s/topapi/user/getbyunionid?access_token=%s", r.data.Endpoint, token)
	input := &DingTalkUseridByUnionidRequest{
		Unionid: unionid,
	}
	jsonData, err := json.Marshal(input)
	if err != nil {
		return "", err
	}

	bs, err := httputil.PostJSON(uri, jsonData, time.Second*10)
	if err != nil {
		return "", err
	}

	r.log.Info("GetUseridByUnionid: %v, err: %v", string(bs), err)

	var getUseridByUnionidResponse *DingTalkUseridByUnionidResponse
	if err = json.Unmarshal(bs, &getUseridByUnionidResponse); err != nil {
		return "", err
	}
	if getUseridByUnionidResponse.Errcode != 0 {
		return "", fmt.Errorf("钉钉API返回错误: %s, errcode: %d", getUseridByUnionidResponse.Errmsg, getUseridByUnionidResponse.Errcode)
	}
	if getUseridByUnionidResponse.Result.Userid == "" {
		return "", fmt.Errorf("钉钉API返回错误 Result: %+v, Result.Userid: %s", getUseridByUnionidResponse.Result, getUseridByUnionidResponse.Result.Userid)
	}
	return getUseridByUnionidResponse.Result.Userid, nil
}

func (r *dingTalkRepo) FetchUserDetail(ctx context.Context, token string, userIds []string) ([]*DingtalkDeptUser, error) {
	log := r.log.WithContext(ctx)

	log.Infof("FetchUserDetail token: %s, userIds: %v", token, userIds)
	uri := fmt.Sprintf("%s/topapi/v2/user/get?access_token=%s", r.data.Endpoint, token)

	log.Info("FetchUserDetail deptIds: %v, uri: %v", userIds, uri)
	sem := make(chan struct{}, r.data.MaxConcurrent)
	userList := make([]*DingtalkDeptUser, 0)
	var mu sync.Mutex

	var wg sync.WaitGroup

	for _, userId := range userIds {
		wg.Add(1)

		select {
		case sem <- struct{}{}:
		case <-ctx.Done():
			wg.Done()
			continue
		}

		// 启动goroutine处理任务
		go func(id string) {
			defer func() {
				<-sem     // 释放信号量
				wg.Done() // 通知任务完成
			}()

			input := &DingTalkUserDetailRequest{
				Userid: id,
			}
			jsonData, err := json.Marshal(input)
			if err != nil {
				r.log.Errorf("GetUserDetail.jsonData: %v, err: %v", string(jsonData), err)
				return
			}

			bs, err := httputil.PostJSON(uri, jsonData, time.Second*10)
			r.log.Infof(">>>>>>>>>GetUserDetail.PostJSON: %v, err: %v\n", string(bs), err)
			if err != nil {
				r.log.Errorf("GetUserDetail.PostJSON: %v, err: %v", string(bs), err)
				return
			}
			var userDetail *DingTalkUserDetailResponse
			if err = json.Unmarshal(bs, &userDetail); err != nil {
				r.log.Errorf("GetUserDetail.Unmarshal: %v, err: %v", string(bs), err)
				return
			}
			if userDetail.Errcode != 0 {
				r.log.Errorf("钉钉API返回错误: %s, errcode: %d", userDetail.Errmsg, userDetail.Errcode)
				return
			}
			user := userDetail.Result
			//r.log.Info("GetUserDetail user: %v", user)
			mu.Lock()
			userList = append(userList, &user)
			mu.Unlock()
		}(userId)
	}
	wg.Wait()

	return userList, nil
}
