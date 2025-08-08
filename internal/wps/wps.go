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
	serviceConf *conf.Service
	log         *log.Helper
}

func NewWps(serviceConf *conf.Service, logger log.Logger) Wps {
	return &wps{
		serviceConf: serviceConf,
		log:         log.NewHelper(logger),
	}
}

// BATCH_POST_USERS_PATH        = "/v7/users/batch_read"
func (ws *wps) BatchPostUsers(ctx context.Context, accessToken string, input BatchPostUsersRequest) (BatchPostUsersResponse, error) {
	log := ws.log.WithContext(ctx)
	log.Infof("PostBatchUsersByExDepIds req %v", input)

	var resp BatchPostUsersResponse

	ak := ws.serviceConf.Auth.App.ClientId
	sk := ws.serviceConf.Auth.App.ClientSecret
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

	ak := ws.serviceConf.Auth.App.ClientId
	sk := ws.serviceConf.Auth.App.ClientSecret
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

	ak := ws.serviceConf.Auth.App.ClientId
	sk := ws.serviceConf.Auth.App.ClientSecret
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
	log.Infof("PostBatchDeleteDepartmentsByDepIds req %v", input)
	var resp *PostBatchDeleteDeptResponse

	ak := ws.serviceConf.Auth.App.ClientId
	sk := ws.serviceConf.Auth.App.ClientSecret
	wpsReq := NewWPSRequest(DOMAIN, ak, sk)

	log.Infof("PostBatchDeleteDepartmentsByDepIds uri: %s, input: %+v\n", POST_DELETE_DEPTS_BY_EXDEPTIDS_PATH, input)
	bs, err := wpsReq.PostJSON(context.Background(), POST_DELETE_DEPTS_BY_EXDEPTIDS_PATH, accessToken, input)

	log.Infof("PostBatchDeleteDepartmentsByDepIds: %s, err: %+v\n", string(bs), err)
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

	return resp, nil
}

func (ws *wps) PostRomoveUserIdFromDeptId(ctx context.Context, accessToken string, input PostRomoveUserIdFromDeptIdRequest) (*PostRomoveUserIdFromDeptIdResponse, error) {
	log := ws.log.WithContext(ctx)
	log.Infof("PostRomoveUserIdFromDeptId req %v", input)
	var resp *PostRomoveUserIdFromDeptIdResponse

	ak := ws.serviceConf.Auth.App.ClientId
	sk := ws.serviceConf.Auth.App.ClientSecret
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

	ak := ws.serviceConf.Auth.App.ClientId
	sk := ws.serviceConf.Auth.App.ClientSecret
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
	ak := ws.serviceConf.Auth.App.ClientId
	sk := ws.serviceConf.Auth.App.ClientSecret
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
	ak := ws.serviceConf.Auth.App.ClientId
	sk := ws.serviceConf.Auth.App.ClientSecret
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

	ak := ws.serviceConf.Auth.App.ClientId
	sk := ws.serviceConf.Auth.App.ClientSecret
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

	return nil

}
func (ws *wps) CacheGet(ctx context.Context, accessToken string, key string) (interface{}, error) {

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
	return nil, nil
}

func (ws *wps) CacheDel(ctx context.Context, accessToken, key string) error {

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
	return nil
}

func (ws *wps) PostCreateUser(ctx context.Context, accessToken string, input PostCreateUserRequest) (*PostCreateUserResponse, error) {
	ws.log.WithContext(ctx).Infof("PostCreateUser req %v", input)
	return nil, status.Error(codes.Unimplemented, "PostCreateUser")
}

func (ws *wps) PostCreateDept(ctx context.Context, accessToken string, input PostCreateDeptRequest) (*PostCreateDeptResponse, error) {
	ws.log.WithContext(ctx).Infof("PostCreateUser req %v", input)
	return nil, status.Error(codes.Unimplemented, "PostCreateDept")
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

	ak := ws.serviceConf.Auth.App.ClientId
	sk := ws.serviceConf.Auth.App.ClientSecret
	wpsReq := NewWPSRequest(DOMAIN, ak, sk, WithLogger(ws.log))

	uri := strings.Replace(GET_USER_DEPT_BY_USERID, "{user_id}", input.UserID, 1)
	bs, err := wpsReq.GET(context.Background(), uri, accessToken, "")

	log.Infof("GetUserByUserId: %s, err: %+v\n", string(bs), err)
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

func (ws *wps) GetDeptChildren(ctx context.Context, accessToken string, input GetDeptChildrenRequest) (*GetDeptChildrenResponse, error) {
	ws.log.WithContext(ctx).Infof("GetDeptChildren req %v", input)

	log := ws.log.WithContext(ctx)
	log.Infof("GetUserDeptsByUserId req %v", input)

	var resp GetDeptChildrenResponse

	ak := ws.serviceConf.Auth.App.ClientId
	sk := ws.serviceConf.Auth.App.ClientSecret
	wpsReq := NewWPSRequest(DOMAIN, ak, sk, WithLogger(ws.log))

	uri := strings.Replace(GET_DEPT_CHILDREN, "{dept_id}", input.DeptID, 1)

	uri += fmt.Sprintf(
		"recursive=%t&page_size=%d&with_total=%t",
		input.Recursive,
		input.PageSize,
		// input.PageToken,
		input.WithTotal,
	)
	if len(input.PageToken) > 0 {
		uri += fmt.Sprintf("&page_token=%s", input.PageToken)
	}

	bs, err := wpsReq.GET(context.Background(), uri, accessToken, "")

	log.Infof("GetUserByUserId: %s, err: %+v\n", string(bs), err)
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

	ak := ws.serviceConf.Auth.App.ClientId
	sk := ws.serviceConf.Auth.App.ClientSecret
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
