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
	log *log.Helper
}

func NewWps(logger log.Logger) Wps {
	return &wps{
		cfg: conf.Get().GetAuth().GetWpsapp(),
		log: log.NewHelper(logger),
	}
}

// BATCH_POST_USERS_PATH        = "/v7/users/batch_read"
func (ws *wps) BatchPostUsers(ctx context.Context, accessToken string, input BatchPostUsersRequest) (BatchPostUsersResponse, error) {
	log := ws.log.WithContext(ctx)
	log.Infof("PostBatchUsersByExDepIds req %v", input)

	var resp BatchPostUsersResponse

	ak := ws.cfg.ClientId
	sk := ws.cfg.ClientSecret
	wpsReq := NewWPSRequest(DOMAIN, ak, sk, WithLogger(ws.log))

	log.Infof("BatchPostUsers uri: %s, input: %+v\n", BATCH_POST_USERS_PATH, input)
	bs, err := wpsReq.PostJSON(context.Background(), BATCH_POST_USERS_PATH, accessToken, input)
	log.Infof("BatchPostUsers: %s, err: %+v\n", string(bs), err)
	if err != nil {
		return resp, err
	}

	err = json.Unmarshal(bs, &resp)
	if err != nil {
		return resp, err
	}
	if resp.Code != 200 {
		return resp, ErrCodeNot200
	}

	return resp, nil
}
func (ws *wps) BatchPostDepartments(ctx context.Context, accessToken string, input BatchPostDepartmentsRequest) (BatchPostDepartmentsResponse, error) {
	log := ws.log.WithContext(ctx)
	log.Infof("BatchPostDepartments req %v", input)

	var resp BatchPostDepartmentsResponse

	ak := ws.cfg.ClientId
	sk := ws.cfg.ClientSecret
	wpsReq := NewWPSRequest(DOMAIN, ak, sk, WithLogger(ws.log))

	bs, err := wpsReq.PostJSON(context.Background(), BATCH_POST_DEPTS_PATH, accessToken, input)

	log.Infof("BatchPostDepartments: %s, err: %+v\n", string(bs), err)
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
	log := ws.log.WithContext(ctx)
	log.Infof("PostBatchDepartmentsByExDepIds req %v", input)
	var resp *PostBatchDepartmentsByExDepIdsResponse

	ak := ws.cfg.ClientId
	sk := ws.cfg.ClientSecret
	wpsReq := NewWPSRequest(DOMAIN, ak, sk)

	log.Infof("PostBatchDepartmentsByExDepIds uri: %s, input: %+v\n", POST_DEPTS_BY_EXDEPTIDS_PATH, input)
	bs, err := wpsReq.PostJSON(context.Background(), POST_DEPTS_BY_EXDEPTIDS_PATH, accessToken, input)

	log.Infof("PostBatchDepartmentsByExDepIds: %s, err: %+v\n", string(bs), err)
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

	log := ws.log.WithContext(ctx)
	log.Infof("PostBatchDeleteDept req %v", input)
	var resp *PostBatchDeleteDeptResponse

	ak := ws.cfg.ClientId
	sk := ws.cfg.ClientSecret
	wpsReq := NewWPSRequest(DOMAIN, ak, sk)

	log.Infof("PostBatchDeleteDept uri: %s, input: %+v\n", POST_DELETE_DEPTS_PATH, input)
	bs, err := wpsReq.PostJSON(context.Background(), POST_DELETE_DEPTS_PATH, accessToken, input)

	log.Infof("PostBatchDeleteDept: %s, err: %+v\n", string(bs), err)
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
	log := ws.log.WithContext(ctx)
	log.Infof("PostBatchDeleteUser req %v", input)
	var resp *PostBatchDeleteUserResponse

	ak := ws.cfg.ClientId
	sk := ws.cfg.ClientSecret
	wpsReq := NewWPSRequest(DOMAIN, ak, sk)

	log.Infof("PostBatchDeleteUser uri: %s, input: %+v\n", POST_DELETE_USERS_PATH, input)
	bs, err := wpsReq.PostJSON(context.Background(), POST_DELETE_USERS_PATH, accessToken, input)

	log.Infof("PostBatchDeleteUser: %s, err: %+v\n", string(bs), err)
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
	log := ws.log.WithContext(ctx)
	log.Infof("PostRomoveUserIdFromDeptId req %v", input)
	var resp *PostRomoveUserIdFromDeptIdResponse

	ak := ws.cfg.ClientId
	sk := ws.cfg.ClientSecret
	wpsReq := NewWPSRequest(DOMAIN, ak, sk)

	uri := strings.Replace(POST_DELETE_USERID_FROM_DEPTID_PATH, "{dept_id}", input.DeptID, 1)
	uri = strings.Replace(uri, "{user_id}", input.UserID, 1)

	log.Infof("PostRomoveUserIdFromDeptId uri: %s, input: %+v\n", uri, nil)
	bs, err := wpsReq.PostJSON(context.Background(), uri, accessToken, nil)

	log.Infof("PostRomoveUserIdFromDeptId: %s, err: %+v\n", string(bs), err)
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
	log := ws.log.WithContext(ctx)
	log.Infof("PostAddUserIdToDeptId req %v", input)
	var resp *PostAddUserIdToDeptIdResponse

	ak := ws.cfg.ClientId
	sk := ws.cfg.ClientSecret
	wpsReq := NewWPSRequest(DOMAIN, ak, sk)

	uri := strings.Replace(POST_ADD_USERID_TO_DEPTID_PATH, "{dept_id}", input.DeptID, 1)
	uri = strings.Replace(uri, "{user_id}", input.UserID, 1)

	log.Infof("PostAddUserIdToDeptId uri: %s, input: %+v\n", uri, nil)
	bs, err := wpsReq.PostJSON(context.Background(), uri, accessToken, nil)

	log.Infof("PostAddUserIdToDeptId: %s, err: %+v\n", string(bs), err)
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

	log := ws.log.WithContext(ctx)
	log.Infof("PostBatchUsersByExDepIds req %v", input)

	var resp *PostBatchUsersByExDepIdsResponse

	// input := &EcisaccountsyncIncrementRequest{
	// 	ThirdCompanyId: thirdCompanyId,
	// }
	ak := ws.cfg.ClientId
	sk := ws.cfg.ClientSecret
	wpsReq := NewWPSRequest(DOMAIN, ak, sk, WithLogger(ws.log))

	log.Infof("PostBatchUsersByExDepIds uri: %s, input: %+v\n", POST_USERS_BY_EXDEPTIDS_PATH, input)
	bs, err := wpsReq.PostJSON(context.Background(), POST_USERS_BY_EXDEPTIDS_PATH, accessToken, input)

	log.Infof("PostBatchUsersByExDepIds res: %s, err: %+v\n", string(bs), err)
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

	log := ws.log.WithContext(ctx)
	log.Infof("GetDepartmentRoot req %v", input)

	var resp GetDepartmentRootResponse

	// input := &EcisaccountsyncIncrementRequest{
	// 	ThirdCompanyId: thirdCompanyId,
	// }
	ak := ws.cfg.ClientId
	sk := ws.cfg.ClientSecret
	wpsReq := NewWPSRequest(DOMAIN, ak, sk, WithLogger(ws.log))

	bs, err := wpsReq.GET(context.Background(), GET_DEPARTMENT_ROOT, accessToken, "")

	log.Infof("GetDepartmentRoot: %s, err: %+v\n", string(bs), err)
	if err != nil {
		return resp, err
	}

	err = json.Unmarshal(bs, &resp)
	if err != nil {
		return resp, err
	}
	//fmt.Printf("resp: %+v\n", resp)
	//fmt.Println()
	if resp.Code != 0 {
		return resp, ErrCodeNot0
	}

	return resp, nil
}

//	func (ws *wps) GetDepartmentChildrenList(ctx context.Context, accessToken string, input GetDepartmentChildrenListRequest) (GetDepartmentChildrenListResponse, error) {
//		var resp GetDepartmentChildrenListResponse
//		return resp, nil
//	}
func (ws *wps) GetUserByUserId(ctx context.Context, accessToken string, input GetUserByUserIdRequest) (GetUserByUserIdResponse, error) {
	log := ws.log.WithContext(ctx)
	log.Infof("GetUserByUserId req %v", input)

	var resp GetUserByUserIdResponse

	ak := ws.cfg.ClientId
	sk := ws.cfg.ClientSecret
	wpsReq := NewWPSRequest(DOMAIN, ak, sk, WithLogger(ws.log))

	uri := strings.Replace(GET_USER_DEPT_BY_USERID, "{user_id}", input.UserID, 1)

	bs, err := wpsReq.GET(context.Background(), uri, accessToken, "")

	log.Infof("GetUserByUserId: %s, err: %+v\n", string(bs), err)
	if err != nil {
		return resp, err
	}

	err = json.Unmarshal(bs, &resp)
	if err != nil {
		return resp, err
	}
	//fmt.Printf("resp: %+v\n", resp)
	//fmt.Println()
	if resp.Code != 0 {
		return resp, ErrCodeNot0
	}

	return resp, nil

}

// TODO 内部网关的签名是不是不对
func (ws *wps) CacheSet(ctx context.Context, accessToken string, key string, value interface{}, expiration time.Duration) error {

	ws.log.WithContext(ctx).Infof("CacheSet req %v", key)
	return status.Error(codes.Unimplemented, "GetDeptByPage")
	// log := ws.log.WithContext(ctx)
	// log.Infof("CacheSet key: %s, value: %v", key, value)

	// ak := ws.serviceConf.Auth.App.ClientId
	// sk := ws.serviceConf.Auth.App.ClientSecret
	// wpsReq := NewWPSRequest(DOMAIN, ak, sk, WithLogger(ws.log))

	// bs, err := wpsReq.PostJSON(context.Background(), POST_CACHE_SET, accessToken, map[string]interface{}{
	// 	"key":       key,
	// 	"value":     value,
	// 	"expire":    expiration,
	// 	"namespace": "nancalacc",
	// })

	// log.Infof("CacheSet: %s, err: %+v\n", string(bs), err)
	// if err != nil {
	// 	return err
	// }

}
func (ws *wps) CacheGet(ctx context.Context, accessToken string, key string) (interface{}, error) {

	ws.log.WithContext(ctx).Infof("CacheGet req %v", key)
	return nil, status.Error(codes.Unimplemented, "GetDeptByPage")
	// log := ws.log.WithContext(ctx)
	// log.Infof("CacheGet key %v", key)

	// ak := ws.serviceConf.Auth.App.ClientId
	// sk := ws.serviceConf.Auth.App.ClientSecret
	// wpsReq := NewWPSRequest(DOMAIN, ak, sk, WithLogger(ws.log))

	// bs, err := wpsReq.PostJSON(context.Background(), POST_CACHE_GET, accessToken, map[string]interface{}{
	// 	"key":       key,
	// 	"namespace": "nancalacc",
	// })

	// log.Infof("CacheGet: %s, err: %+v\n", string(bs), err)
	// if err != nil {
	// 	return nil, err
	// }

	// var resp map[string]interface{}
	// err = json.Unmarshal(bs, &resp)
	// if err != nil {
	// 	return nil, err
	// }
}

func (ws *wps) CacheDel(ctx context.Context, accessToken, input string) error {
	ws.log.WithContext(ctx).Infof("CacheDel req %v", input)
	return status.Error(codes.Unimplemented, "GetDeptByPage")

	// log := ws.log.WithContext(ctx)
	// log.Infof("CacheDel key %v", key)

	// ak := ws.serviceConf.Auth.App.ClientId
	// sk := ws.serviceConf.Auth.App.ClientSecret
	// wpsReq := NewWPSRequest(DOMAIN, ak, sk, WithLogger(ws.log))

	// bs, err := wpsReq.PostJSON(context.Background(), POST_CACHE_DEL, accessToken, map[string]interface{}{
	// 	"key":       key,
	// 	"namespace": "nancalacc",
	// })

	// log.Infof("CacheDel: %s, err: %+v\n", string(bs), err)
	// if err != nil {
	// 	return err
	// }

	// var resp map[string]interface{}
	// err = json.Unmarshal(bs, &resp)
	// if err != nil {
	// 	return err
	// }
}

func (ws *wps) PostUpdateDept(ctx context.Context, accessToken string, input PostUpdateDeptRequest) (*PostUpdateDeptResponse, error) {
	ws.log.WithContext(ctx).Infof("PostUpdateDept req %v", input)
	return nil, status.Error(codes.Unimplemented, "PostUpdateDept")
}
func (ws *wps) PostUpdateUser(ctx context.Context, accessToken string, input PostUpdateUserRequest) (*PostUpdateUserResponse, error) {
	ws.log.WithContext(ctx).Infof("PostUpdateUser req %v", input)
	return nil, status.Error(codes.Unimplemented, "PostUpdateUser")
}

func (ws *wps) PostBatchUserByPage(ctx context.Context, accessToken string, input PostBatchUserByPageRequest) (*PostBatchUserByPageResponse, error) {
	ws.log.WithContext(ctx).Infof("PostCreateUser req %v", input)
	return nil, status.Error(codes.Unimplemented, "PostBatchUserByPage")
}

func (ws *wps) GetDeptByPage(ctx context.Context, accessToken string, input GetDeptByPageRequest) (*GetDeptByPageResponse, error) {
	ws.log.WithContext(ctx).Infof("GetDeptByPage req %v", input)
	return nil, status.Error(codes.Unimplemented, "GetDeptByPage")
}

func (ws *wps) GetUserDeptsByUserId(ctx context.Context, accessToken string, input GetUserDeptsByUserIdRequest) (*GetUserDeptsByUserIdResponse, error) {
	log := ws.log.WithContext(ctx)
	log.Infof("GetUserDeptsByUserId req %v", input)

	var resp GetUserDeptsByUserIdResponse

	ak := ws.cfg.ClientId
	sk := ws.cfg.ClientSecret
	wpsReq := NewWPSRequest(DOMAIN, ak, sk, WithLogger(ws.log))

	uri := strings.Replace(GET_USER_DEPT_BY_USERID, "{user_id}", input.UserID, 1)
	bs, err := wpsReq.GET(context.Background(), uri, accessToken, "")

	log.Infof("GetUserDeptsByUserId: %s, err: %+v\n", string(bs), err)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bs, &resp)
	if err != nil {
		return nil, err
	}
	log.Infof("GetUserDeptsByUserId res %v", resp)
	//fmt.Printf("resp: %+v\n", resp)
	//fmt.Println()
	if resp.Code != 0 {
		return nil, ErrCodeNot0
	}
	return &resp, nil
}

func (ws *wps) GetDeptChildren(ctx context.Context, accessToken string, input GetDeptChildrenRequest) (*GetDeptChildrenResponse, error) {

	log := ws.log.WithContext(ctx)
	log.Infof("GetDeptChildren req %v", input)

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

	log.Infof("GetDeptChildren: %s, err: %+v\n", string(bs), err)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bs, &resp)
	if err != nil {
		return nil, err
	}
	//fmt.Printf("resp: %+v\n", resp)
	//fmt.Println()
	if resp.Code != 0 {
		return nil, ErrCodeNot0
	}

	return &resp, nil

}

func (ws *wps) GetCompAllUsers(ctx context.Context, accessToken string, input GetCompAllUsersRequest) (*GetCompAllUsersResponse, error) {
	ws.log.WithContext(ctx).Infof("GetCompAllUsers req %v", input)

	log := ws.log.WithContext(ctx)
	log.Infof("GetCompAllUsers req %v", input)

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
	log.Infof("GetCompAllUsers req: %s\n", uri)
	bs, err := wpsReq.GET(context.Background(), uri, accessToken, "")

	log.Infof("GetCompAllUsers res: %s, err: %+v\n", string(bs), err)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bs, &resp)
	if err != nil {
		return nil, err
	}
	//fmt.Printf("resp: %+v\n", resp)
	//fmt.Println()
	if resp.Code != 0 {
		return nil, ErrCodeNot0
	}

	return &resp, nil
}

func (ws *wps) PostCreateDept(ctx context.Context, accessToken string, input PostCreateDeptRequest) (*PostCreateDeptResponse, error) {

	log := ws.log.WithContext(ctx)
	log.Infof("PostCreateDept req %v", input)

	var resp *PostCreateDeptResponse

	// input := &EcisaccountsyncIncrementRequest{
	// 	ThirdCompanyId: thirdCompanyId,
	// }
	ak := ws.cfg.ClientId
	sk := ws.cfg.ClientSecret
	wpsReq := NewWPSRequest(DOMAIN, ak, sk, WithLogger(ws.log))

	log.Infof("PostCreateDept uri: %s, input: %+v\n", POST_CREATE_DEPT_PATH, input)
	bs, err := wpsReq.PostJSON(context.Background(), POST_CREATE_DEPT_PATH, accessToken, input)

	log.Infof("PostCreateDept res: %s, err: %+v\n", string(bs), err)
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
	log := ws.log.WithContext(ctx)
	log.Infof("PostCreateUser req %v", input)

	var resp *PostCreateUserResponse

	// input := &EcisaccountsyncIncrementRequest{
	// 	ThirdCompanyId: thirdCompanyId,
	// }
	ak := ws.cfg.ClientId
	sk := ws.cfg.ClientSecret
	wpsReq := NewWPSRequest(DOMAIN, ak, sk, WithLogger(ws.log))

	log.Infof("PostCreateUser uri: %s, input: %+v\n", POST_CREATE_USER_PATH, input)
	bs, err := wpsReq.PostJSON(context.Background(), POST_CREATE_USER_PATH, accessToken, input)

	log.Infof("PostCreateUser res: %s, err: %+v\n", string(bs), err)
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
	log := ws.log.WithContext(ctx)
	log.Infof("GetUsersSearch req %v", input)

	var resp *GetUsersSearchResponse

	ak := ws.cfg.ClientId
	sk := ws.cfg.ClientSecret
	wpsReq := NewWPSRequest(DOMAIN, ak, sk, WithLogger(ws.log))

	log.Infof("GetUsersSearch uri: %s, input: %+v\n", GET_USERS_SEARCH_PATH, input)
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

	log.Infof("GetUsersSearch res: %s, err: %+v\n", string(bs), err)
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
	log := ws.log.WithContext(ctx)
	log.Infof("GetContactPermission req %v", input)

	var resp *GetContactPermissionResponse

	ak := ws.cfg.ClientId
	sk := ws.cfg.ClientSecret
	wpsReq := NewWPSRequest(DOMAIN, ak, sk, WithLogger(ws.log))

	log.Infof("GetUsersSearch uri: %s, input: %+v\n", GET_CONTACT_PERMISSION_PATH, input)
	var uri string
	for _, status := range input.Scopes {
		uri += fmt.Sprintf("&scope=%s", status)
	}
	bs, err := wpsReq.GET(context.Background(), GET_CONTACT_PERMISSION_PATH+"?"+uri, accessToken, "")

	log.Infof("GetUsersSearch res: %s, err: %+v\n", string(bs), err)
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
	log := ws.log.WithContext(ctx)
	log.Infof("GetObjUploadUrl req %v", input)

	var resp *GetObjUploadUrlResponse

	ak := ws.cfg.ClientId
	sk := ws.cfg.ClientSecret
	wpsReq := NewWPSRequest(DOMAIN, ak, sk, WithLogger(ws.log))

	// http: //encs-pri-cams-engine/{c}/asyncacc/v1/task

	uri := "http://119.3.173.229/api/cams/sdk/api/v1/wps3/presigned_upload"
	uri = "http://119.3.173.229/api/cams/sdk/api/v1/wps3/presigned_upload"

	log.Infof("GetObjUploadUrl uri: %s, input: %+v\n", uri, "")
	bs, err := wpsReq.GET(context.Background(), uri, accessToken, "")

	log.Infof("GetObjUploadUrl res: %s, err: %+v\n", string(bs), err)
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
