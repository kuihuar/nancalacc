package wps

import (
	"context"
	"encoding/json"
	"fmt"
	"nancalacc/internal/conf"
	"strings"
	"time"

	//httpwps "nancalacc/pkg/httputil/wps"

	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type wps struct {
	cfg *conf.Auth_Wpsapp
	log log.Logger
}

func NewWps(logger log.Logger) Wps {
	return &wps{
		cfg: conf.Get().GetAuth().GetWpsapp(),
		log: logger,
	}
}

// BATCH_POST_USERS_PATH        = "/v7/users/batch_read"
func (ws *wps) BatchPostUsers(ctx context.Context, accessToken string, input BatchPostUsersRequest) (BatchPostUsersResponse, error) {
	var resp BatchPostUsersResponse

	// 记录请求
	logAPIRequest(ctx, ws.log, "BatchPostUsers", "POST", BATCH_POST_USERS_PATH, input)

	// 创建请求
	wpsReq := NewWPSRequest(DOMAIN, ws.cfg.ClientId, ws.cfg.ClientSecret, WithLogger(ws.log))

	// 执行请求
	bs, err := wpsReq.PostJSON(ctx, BATCH_POST_USERS_PATH, accessToken, input)
	if err != nil {
		ws.log.Log(log.LevelError, "msg", "BatchPostUsers request failed", "err", err)
		return resp, err
	}

	// 处理响应
	if err := handleAPIResponse(ctx, ws.log, "BatchPostUsers", bs, &resp, 0); err != nil {
		return resp, err
	}

	return resp, nil
}
func (ws *wps) BatchPostDepartments(ctx context.Context, accessToken string, input BatchPostDepartmentsRequest) (BatchPostDepartmentsResponse, error) {
	ws.log.Log(log.LevelInfo, "msg", "BatchPostDepartments", "req", input)

	var resp BatchPostDepartmentsResponse

	ak := ws.cfg.ClientId
	sk := ws.cfg.ClientSecret
	wpsReq := NewWPSRequest(DOMAIN, ak, sk, WithLogger(ws.log))

	bs, err := wpsReq.PostJSON(context.Background(), BATCH_POST_DEPTS_PATH, accessToken, input)

	ws.log.Log(log.LevelInfo, "msg", "BatchPostDepartments", "res", string(bs), "err", err)
	if err != nil {
		return resp, err
	}

	err = json.Unmarshal(bs, &resp)
	if err != nil {
		return resp, err
	}
	if resp.Code != 0 {
		return resp, ErrCodeNot0
	}

	return resp, nil
}

func (ws *wps) PostBatchDepartmentsByExDepIds(ctx context.Context, accessToken string, input PostBatchDepartmentsByExDepIdsRequest) (*PostBatchDepartmentsByExDepIdsResponse, error) {
	ws.log.Log(log.LevelInfo, "msg", "PostBatchDepartmentsByExDepIds", "req", input)
	var resp *PostBatchDepartmentsByExDepIdsResponse

	ak := ws.cfg.ClientId
	sk := ws.cfg.ClientSecret
	wpsReq := NewWPSRequest(DOMAIN, ak, sk, WithLogger(ws.log))

	ws.log.Log(log.LevelInfo, "msg", "PostBatchDepartmentsByExDepIds", "uri", POST_DEPTS_BY_EXDEPTIDS_PATH, "input", input)
	bs, err := wpsReq.PostJSON(context.Background(), POST_DEPTS_BY_EXDEPTIDS_PATH, accessToken, input)

	ws.log.Log(log.LevelInfo, "msg", "PostBatchDepartmentsByExDepIds", "res", string(bs), "err", err)
	if err != nil {
		return resp, err
	}

	err = json.Unmarshal(bs, &resp)
	if err != nil {
		return resp, err
	}
	if resp.Code != 0 {
		return resp, ErrCodeNot0
	}

	return resp, nil

}

func (ws *wps) PostBatchDeleteDept(ctx context.Context, accessToken string, input PostBatchDeleteDeptRequest) (*PostBatchDeleteDeptResponse, error) {

	ws.log.Log(log.LevelInfo, "msg", "PostBatchDeleteDept", "req", input)
	var resp *PostBatchDeleteDeptResponse

	ak := ws.cfg.ClientId
	sk := ws.cfg.ClientSecret
	wpsReq := NewWPSRequest(DOMAIN, ak, sk, WithLogger(ws.log))

	ws.log.Log(log.LevelInfo, "msg", "PostBatchDeleteDept", "uri", POST_DELETE_DEPTS_PATH, "input", input)
	bs, err := wpsReq.PostJSON(context.Background(), POST_DELETE_DEPTS_PATH, accessToken, input)

	ws.log.Log(log.LevelInfo, "msg", "PostBatchDeleteDept", "res", string(bs), "err", err)
	if err != nil {
		return resp, err
	}

	err = json.Unmarshal(bs, &resp)
	if err != nil {
		return resp, err
	}
	if resp.Code != 0 {
		return resp, ErrCodeNot0
	}

	return resp, nil
}
func (ws *wps) PostBatchDeleteUser(ctx context.Context, accessToken string, input PostBatchDeleteUserRequest) (*PostBatchDeleteUserResponse, error) {
	ws.log.Log(log.LevelInfo, "msg", "PostBatchDeleteUser", "req", input)
	var resp *PostBatchDeleteUserResponse

	ak := ws.cfg.ClientId
	sk := ws.cfg.ClientSecret
	wpsReq := NewWPSRequest(DOMAIN, ak, sk, WithLogger(ws.log))

	ws.log.Log(log.LevelInfo, "msg", "PostBatchDeleteUser", "uri", POST_DELETE_USERS_PATH, "input", input)
	bs, err := wpsReq.PostJSON(context.Background(), POST_DELETE_USERS_PATH, accessToken, input)

	ws.log.Log(log.LevelInfo, "msg", "PostBatchDeleteUser", "res", string(bs), "err", err)
	if err != nil {
		return resp, err
	}

	err = json.Unmarshal(bs, &resp)
	if err != nil {
		return resp, err
	}
	if resp.Code != 0 {
		return resp, ErrCodeNot0
	}

	return resp, nil
}

func (ws *wps) PostRomoveUserIdFromDeptId(ctx context.Context, accessToken string, input PostRomoveUserIdFromDeptIdRequest) (*PostRomoveUserIdFromDeptIdResponse, error) {
	ws.log.Log(log.LevelInfo, "msg", "PostRomoveUserIdFromDeptId", "req", input)
	var resp *PostRomoveUserIdFromDeptIdResponse

	ak := ws.cfg.ClientId
	sk := ws.cfg.ClientSecret
	wpsReq := NewWPSRequest(DOMAIN, ak, sk, WithLogger(ws.log))

	uri := strings.Replace(POST_DELETE_USERID_FROM_DEPTID_PATH, "{dept_id}", input.DeptID, 1)
	uri = strings.Replace(uri, "{user_id}", input.UserID, 1)

	ws.log.Log(log.LevelInfo, "msg", "PostRomoveUserIdFromDeptId", "uri", uri, "input", nil)
	bs, err := wpsReq.PostJSON(context.Background(), uri, accessToken, nil)

	ws.log.Log(log.LevelInfo, "msg", "PostRomoveUserIdFromDeptId", "res", string(bs), "err", err)
	if err != nil {
		return resp, err
	}

	err = json.Unmarshal(bs, &resp)
	if err != nil {
		return resp, err
	}
	if resp.Code != 0 {
		return resp, ErrCodeNot0
	}

	return resp, nil
}
func (ws *wps) PostAddUserIdToDeptId(ctx context.Context, accessToken string, input PostAddUserIdToDeptIdRequest) (*PostAddUserIdToDeptIdResponse, error) {
	ws.log.Log(log.LevelInfo, "msg", "PostAddUserIdToDeptId", "req", input)
	var resp *PostAddUserIdToDeptIdResponse

	ak := ws.cfg.ClientId
	sk := ws.cfg.ClientSecret
	wpsReq := NewWPSRequest(DOMAIN, ak, sk, WithLogger(ws.log))

	uri := strings.Replace(POST_ADD_USERID_TO_DEPTID_PATH, "{dept_id}", input.DeptID, 1)
	uri = strings.Replace(uri, "{user_id}", input.UserID, 1)

	ws.log.Log(log.LevelInfo, "msg", "PostAddUserIdToDeptId", "uri", uri, "input", nil)
	bs, err := wpsReq.PostJSON(context.Background(), uri, accessToken, nil)

	ws.log.Log(log.LevelInfo, "msg", "PostAddUserIdToDeptId", "res", string(bs), "err", err)
	if err != nil {
		return resp, err
	}

	err = json.Unmarshal(bs, &resp)
	if err != nil {
		return resp, err
	}
	if resp.Code != 0 {
		return resp, ErrCodeNot0
	}

	return resp, nil
}
func (ws *wps) PostBatchUsersByExDepIds(ctx context.Context, accessToken string, input PostBatchUsersByExDepIdsRequest) (*PostBatchUsersByExDepIdsResponse, error) {

	ws.log.Log(log.LevelInfo, "msg", "PostBatchUsersByExDepIds", "req", input)

	var resp *PostBatchUsersByExDepIdsResponse

	ak := ws.cfg.ClientId
	sk := ws.cfg.ClientSecret
	wpsReq := NewWPSRequest(DOMAIN, ak, sk, WithLogger(ws.log))

	ws.log.Log(log.LevelInfo, "msg", "PostBatchUsersByExDepIds", "uri", POST_USERS_BY_EXDEPTIDS_PATH, "input", input)
	bs, err := wpsReq.PostJSON(context.Background(), POST_USERS_BY_EXDEPTIDS_PATH, accessToken, input)

	ws.log.Log(log.LevelInfo, "msg", "PostBatchUsersByExDepIds", "res", string(bs), "err", err)
	if err != nil {
		return resp, err
	}

	err = json.Unmarshal(bs, &resp)
	if err != nil {
		return resp, err
	}
	if resp.Code != 0 {
		return resp, ErrCodeNot0
	}

	return resp, nil

}

func (ws *wps) GetDepartmentRoot(ctx context.Context, accessToken string, input GetDepartmentRootRequest) (GetDepartmentRootResponse, error) {

	ws.log.Log(log.LevelInfo, "msg", "GetDepartmentRoot", "req", input)

	var resp GetDepartmentRootResponse

	ak := ws.cfg.ClientId
	sk := ws.cfg.ClientSecret
	wpsReq := NewWPSRequest(DOMAIN, ak, sk, WithLogger(ws.log))

	bs, err := wpsReq.GET(context.Background(), GET_DEPARTMENT_ROOT, accessToken, "")

	ws.log.Log(log.LevelInfo, "msg", "GetDepartmentRoot", "res", string(bs), "err", err)
	if err != nil {
		return resp, err
	}

	err = json.Unmarshal(bs, &resp)
	if err != nil {
		return resp, err
	}
	if resp.Code != 0 {
		return resp, ErrCodeNot0
	}

	return resp, nil
}

func (ws *wps) GetUserByUserId(ctx context.Context, accessToken string, input GetUserByUserIdRequest) (GetUserByUserIdResponse, error) {
	ws.log.Log(log.LevelInfo, "msg", "GetUserByUserId", "req", input)

	var resp GetUserByUserIdResponse

	ak := ws.cfg.ClientId
	sk := ws.cfg.ClientSecret
	wpsReq := NewWPSRequest(DOMAIN, ak, sk, WithLogger(ws.log))

	uri := strings.Replace(GET_USER_DEPT_BY_USERID, "{user_id}", input.UserID, 1)

	bs, err := wpsReq.GET(context.Background(), uri, accessToken, "")

	ws.log.Log(log.LevelInfo, "msg", "GetUserByUserId", "res", string(bs), "err", err)
	if err != nil {
		return resp, err
	}

	err = json.Unmarshal(bs, &resp)
	if err != nil {
		return resp, err
	}
	if resp.Code != 0 {
		return resp, ErrCodeNot0
	}

	return resp, nil

}

// TODO 内部网关的签名是不是不对
func (ws *wps) CacheSet(ctx context.Context, accessToken string, key string, value interface{}, expiration time.Duration) error {

	ws.log.Log(log.LevelInfo, "msg", "CacheSet", "req", key)
	return status.Error(codes.Unimplemented, "GetDeptByPage")

}
func (ws *wps) CacheGet(ctx context.Context, accessToken string, key string) (interface{}, error) {

	ws.log.Log(log.LevelInfo, "msg", "CacheGet", "req", key)
	return nil, status.Error(codes.Unimplemented, "GetDeptByPage")

}

func (ws *wps) CacheDel(ctx context.Context, accessToken, input string) error {
	ws.log.Log(log.LevelInfo, "msg", "CacheDel", "req", input)
	return status.Error(codes.Unimplemented, "GetDeptByPage")

}

func (ws *wps) PostUpdateDept(ctx context.Context, accessToken string, input PostUpdateDeptRequest) (*PostUpdateDeptResponse, error) {
	return nil, status.Error(codes.Unimplemented, "PostUpdateDept")
}
func (ws *wps) PostUpdateUser(ctx context.Context, accessToken string, input PostUpdateUserRequest) (*PostUpdateUserResponse, error) {
	return nil, status.Error(codes.Unimplemented, "PostUpdateUser")
}

func (ws *wps) PostBatchUserByPage(ctx context.Context, accessToken string, input PostBatchUserByPageRequest) (*PostBatchUserByPageResponse, error) {
	return nil, status.Error(codes.Unimplemented, "PostBatchUserByPage")
}

func (ws *wps) GetDeptByPage(ctx context.Context, accessToken string, input GetDeptByPageRequest) (*GetDeptByPageResponse, error) {
	return nil, status.Error(codes.Unimplemented, "GetDeptByPage")
}

func (ws *wps) GetUserDeptsByUserId(ctx context.Context, accessToken string, input GetUserDeptsByUserIdRequest) (*GetUserDeptsByUserIdResponse, error) {
	ws.log.Log(log.LevelInfo, "msg", "GetUserDeptsByUserId", "req", input)

	var resp GetUserDeptsByUserIdResponse

	ak := ws.cfg.ClientId
	sk := ws.cfg.ClientSecret
	wpsReq := NewWPSRequest(DOMAIN, ak, sk, WithLogger(ws.log))

	uri := strings.Replace(GET_USER_DEPT_BY_USERID, "{user_id}", input.UserID, 1)
	bs, err := wpsReq.GET(context.Background(), uri, accessToken, "")

	ws.log.Log(log.LevelInfo, "msg", "GetUserDeptsByUserId", "res", string(bs), "err", err)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bs, &resp)
	if err != nil {
		return nil, err
	}
	ws.log.Log(log.LevelInfo, "msg", "GetUserDeptsByUserId", "res", resp)
	if resp.Code != 0 {
		return nil, ErrCodeNot0
	}
	return &resp, nil
}

func (ws *wps) GetDeptChildren(ctx context.Context, accessToken string, input GetDeptChildrenRequest) (*GetDeptChildrenResponse, error) {

	ws.log.Log(log.LevelInfo, "msg", "GetDeptChildren", "req", input)

	var resp GetDeptChildrenResponse

	ak := ws.cfg.ClientId
	sk := ws.cfg.ClientSecret
	wpsReq := NewWPSRequest(DOMAIN, ak, sk, WithLogger(ws.log))

	path := strings.Replace(GET_DEPT_CHILDREN, "{dept_id}", input.DeptID, 1)

	uri := fmt.Sprintf(
		"%s?recursive=%t&page_size=%d&with_total=%t",
		path,
		input.Recursive,
		input.PageSize,
		// input.PageToken,
		input.WithTotal,
	)

	if len(input.PageToken) > 0 {
		uri += fmt.Sprintf("&page_token=%s", input.PageToken)
	}

	bs, err := wpsReq.GET(context.Background(), uri, accessToken, "")

	ws.log.Log(log.LevelInfo, "msg", "GetDeptChildren", "res", string(bs), "err", err)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bs, &resp)
	if err != nil {
		return nil, err
	}
	if resp.Code != 0 {
		return nil, ErrCodeNot0
	}

	return &resp, nil

}

func (ws *wps) GetCompAllUsers(ctx context.Context, accessToken string, input GetCompAllUsersRequest) (*GetCompAllUsersResponse, error) {

	ws.log.Log(log.LevelInfo, "msg", "GetCompAllUsers", "req", input)

	var resp GetCompAllUsersResponse

	ak := ws.cfg.ClientId
	sk := ws.cfg.ClientSecret
	wpsReq := NewWPSRequest(DOMAIN, ak, sk, WithLogger(ws.log))

	uri := fmt.Sprintf(
		"%s?recursive=%t&page_size=%d&with_total=%t",
		GET_ALL_USER_PARH,
		input.Recursive,
		input.PageSize,
		input.WithTotal,
	)
	if len(input.Status) > 0 {
		for _, status := range input.Status {
			uri += fmt.Sprintf("&status=%s", status)
		}
	}
	if len(input.PageToken) > 0 {
		uri += fmt.Sprintf("&page_token=%s", input.PageToken)
	}
	ws.log.Log(log.LevelInfo, "msg", "GetCompAllUsers", "uri", uri)
	bs, err := wpsReq.GET(context.Background(), uri, accessToken, "")

	ws.log.Log(log.LevelInfo, "msg", "GetCompAllUsers", "res", string(bs), "err", err)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bs, &resp)
	if err != nil {
		return nil, err
	}
	if resp.Code != 0 {
		return nil, ErrCodeNot0
	}

	return &resp, nil
}

func (ws *wps) PostCreateDept(ctx context.Context, accessToken string, input PostCreateDeptRequest) (*PostCreateDeptResponse, error) {

	ws.log.Log(log.LevelInfo, "msg", "PostCreateDept", "req", input)

	var resp *PostCreateDeptResponse

	ak := ws.cfg.ClientId
	sk := ws.cfg.ClientSecret
	wpsReq := NewWPSRequest(DOMAIN, ak, sk, WithLogger(ws.log))

	ws.log.Log(log.LevelInfo, "msg", "PostCreateDept", "uri", POST_CREATE_DEPT_PATH, "input", input)
	bs, err := wpsReq.PostJSON(context.Background(), POST_CREATE_DEPT_PATH, accessToken, input)

	ws.log.Log(log.LevelInfo, "msg", "PostCreateDept", "res", string(bs), "err", err)
	if err != nil {
		return resp, err
	}

	err = json.Unmarshal(bs, &resp)
	if err != nil {
		return resp, err
	}
	if resp.Code != 0 {
		return resp, ErrCodeNot0
	}

	return resp, nil

}

func (ws *wps) PostCreateUser(ctx context.Context, accessToken string, input PostCreateUserRequest) (*PostCreateUserResponse, error) {
	ws.log.Log(log.LevelInfo, "msg", "PostCreateUser", "req", input)

	var resp *PostCreateUserResponse

	// input := &EcisaccountsyncIncrementRequest{
	// 	ThirdCompanyId: thirdCompanyId,
	// }
	ak := ws.cfg.ClientId
	sk := ws.cfg.ClientSecret
	wpsReq := NewWPSRequest(DOMAIN, ak, sk, WithLogger(ws.log))

	ws.log.Log(log.LevelInfo, "msg", "PostCreateUser", "uri", POST_CREATE_USER_PATH, "input", input)
	bs, err := wpsReq.PostJSON(context.Background(), POST_CREATE_USER_PATH, accessToken, input)

	ws.log.Log(log.LevelInfo, "msg", "PostCreateUser", "res", string(bs), "err", err)
	if err != nil {
		return resp, err
	}

	err = json.Unmarshal(bs, &resp)
	if err != nil {
		return resp, err
	}
	if resp.Code != 0 {
		return resp, ErrCodeNot0
	}

	return resp, nil
}

func (ws *wps) GetUsersSearch(ctx context.Context, accessToken string, input GetUsersSearchRequest) (*GetUsersSearchResponse, error) {
	ws.log.Log(log.LevelInfo, "msg", "GetUsersSearch", "req", input)

	var resp *GetUsersSearchResponse

	ak := ws.cfg.ClientId
	sk := ws.cfg.ClientSecret
	wpsReq := NewWPSRequest(DOMAIN, ak, sk, WithLogger(ws.log))

	ws.log.Log(log.LevelInfo, "msg", "GetUsersSearch", "uri", GET_USERS_SEARCH_PATH, "input", input)
	uri := fmt.Sprintf("%s?keyword=%s&page_size=%d", GET_USERS_SEARCH_PATH, input.Keyword, input.PageSize)

	if input.PageToken != "" {
		uri += "&page_token=true"
		uri += fmt.Sprintf("&page_token=%s", input.PageToken)
	}
	if len(input.Status) > 0 {
		for _, status := range input.Status {
			uri += fmt.Sprintf("&status=%s", status)
		}
	}
	if input.SearchFieldConfigEnabled {
		uri += "&search_field_config_enabled=true"
	}

	if len(input.SearchSource) > 0 {
		for _, searchSource := range input.SearchSource {
			uri += fmt.Sprintf("&search_source=%s", searchSource)
		}
	}
	bs, err := wpsReq.GET(context.Background(), uri, accessToken, "")

	ws.log.Log(log.LevelInfo, "msg", "GetUsersSearch", "res", string(bs), "err", err)
	if err != nil {
		return resp, err
	}

	err = json.Unmarshal(bs, &resp)
	if err != nil {
		return resp, err
	}
	if resp.Code != 0 {
		return resp, ErrCodeNot0
	}

	return resp, nil
}

func (ws *wps) GetContactPermission(ctx context.Context, accessToken string, input GetContactPermissionRequest) (*GetContactPermissionResponse, error) {
	ws.log.Log(log.LevelInfo, "msg", "GetContactPermission", "req", input)

	var resp *GetContactPermissionResponse

	ak := ws.cfg.ClientId
	sk := ws.cfg.ClientSecret
	wpsReq := NewWPSRequest(DOMAIN, ak, sk, WithLogger(ws.log))

	ws.log.Log(log.LevelInfo, "msg", "GetContactPermission", "uri", GET_CONTACT_PERMISSION_PATH, "input", input)
	var uri string
	for _, status := range input.Scopes {
		uri += fmt.Sprintf("&scope=%s", status)
	}
	bs, err := wpsReq.GET(context.Background(), GET_CONTACT_PERMISSION_PATH+"?"+uri, accessToken, "")

	ws.log.Log(log.LevelInfo, "msg", "GetContactPermission", "res", string(bs), "err", err)
	if err != nil {
		return resp, err
	}

	err = json.Unmarshal(bs, &resp)
	if err != nil {
		return resp, err
	}
	if resp.Code != 0 {
		return resp, ErrCodeNot0
	}

	return resp, nil
}

// 内部网关
func (ws *wps) GetObjUploadUrl(ctx context.Context, accessToken string, input GetObjUploadUrlRequest) (*GetObjUploadUrlResponse, error) {
	ws.log.Log(log.LevelInfo, "msg", "GetObjUploadUrl", "req", input)

	var resp *GetObjUploadUrlResponse

	ak := ws.cfg.ClientId
	sk := ws.cfg.ClientSecret
	wpsReq := NewWPSRequest(DOMAIN, ak, sk, WithLogger(ws.log))

	// http: //encs-pri-cams-engine/{c}/asyncacc/v1/task

	uri := "http://119.3.173.229/api/cams/sdk/api/v1/wps3/presigned_upload"
	uri = "http://119.3.173.229/api/cams/sdk/api/v1/wps3/presigned_upload"

	ws.log.Log(log.LevelInfo, "msg", "GetObjUploadUrl", "uri", uri, "input", "")
	bs, err := wpsReq.GET(context.Background(), uri, accessToken, "")

	ws.log.Log(log.LevelInfo, "msg", "GetObjUploadUrl", "res", string(bs), "err", err)
	if err != nil {
		return resp, err
	}

	err = json.Unmarshal(bs, &resp)
	if err != nil {
		return resp, err
	}
	if resp.Code != 0 {
		return resp, ErrCodeNot0
	}

	return resp, nil
}
