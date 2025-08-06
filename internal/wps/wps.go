package wps

import (
	"context"
	"encoding/json"
	"fmt"
	"nancalacc/internal/conf"
	"time"

	//httpwps "nancalacc/pkg/httputil/wps"

	"github.com/go-kratos/kratos/v2/log"
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
	log.Infof("PostBatchUsersByExDepIds req %v", input)

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
	if resp.Code != 200 {
		return resp, ErrCodeNot200
	}

	return resp, nil
}

func (ws *wps) PostBatchDepartmentsByExDepIds(ctx context.Context, accessToken string, input PostBatchDepartmentsByExDepIdsRequest) (*PostBatchDepartmentsByExDepIdsResponse, error) {
	log := ws.log.WithContext(ctx)
	log.Infof("PostBatchUsersByExDepIds req %v", input)
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
func (ws *wps) GetUserByUserId(ctx context.Context, accessToken string, req GetUserByUserIdRequest) (GetUserByUserIdResponse, error) {
	log := ws.log.WithContext(ctx)
	log.Infof("GetUserByUserId req %v", req)

	var resp GetUserByUserIdResponse

	ak := ws.serviceConf.Auth.App.ClientId
	sk := ws.serviceConf.Auth.App.ClientSecret
	wpsReq := NewWPSRequest(DOMAIN, ak, sk, WithLogger(ws.log))

	path := fmt.Sprintf("%s/%s", GET_USER_BY_USERID, req.UserID)
	log.Infof("GetUserByUserId path: %s\n", path)
	bs, err := wpsReq.GET(context.Background(), path, accessToken, "")

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

	log := ws.log.WithContext(ctx)
	log.Infof("CacheSet key: %s, value: %v", key, value)

	ak := ws.serviceConf.Auth.App.ClientId
	sk := ws.serviceConf.Auth.App.ClientSecret
	wpsReq := NewWPSRequest(DOMAIN, ak, sk, WithLogger(ws.log))

	bs, err := wpsReq.PostJSON(context.Background(), POST_CACHE_SET, accessToken, map[string]interface{}{
		"key":       key,
		"value":     value,
		"expire":    expiration,
		"namespace": "nancalacc",
	})

	log.Infof("CacheSet: %s, err: %+v\n", string(bs), err)
	if err != nil {
		return err
	}

	return nil

}
func (ws *wps) CacheGet(ctx context.Context, accessToken string, key string) (interface{}, error) {

	log := ws.log.WithContext(ctx)
	log.Infof("CacheGet key %v", key)

	ak := ws.serviceConf.Auth.App.ClientId
	sk := ws.serviceConf.Auth.App.ClientSecret
	wpsReq := NewWPSRequest(DOMAIN, ak, sk, WithLogger(ws.log))

	bs, err := wpsReq.PostJSON(context.Background(), POST_CACHE_GET, accessToken, map[string]interface{}{
		"key":       key,
		"namespace": "nancalacc",
	})

	log.Infof("CacheGet: %s, err: %+v\n", string(bs), err)
	if err != nil {
		return nil, err
	}

	var resp map[string]interface{}
	err = json.Unmarshal(bs, &resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (ws *wps) CacheDel(ctx context.Context, accessToken, key string) error {

	log := ws.log.WithContext(ctx)
	log.Infof("CacheDel key %v", key)

	ak := ws.serviceConf.Auth.App.ClientId
	sk := ws.serviceConf.Auth.App.ClientSecret
	wpsReq := NewWPSRequest(DOMAIN, ak, sk, WithLogger(ws.log))

	bs, err := wpsReq.PostJSON(context.Background(), POST_CACHE_DEL, accessToken, map[string]interface{}{
		"key":       key,
		"namespace": "nancalacc",
	})

	log.Infof("CacheDel: %s, err: %+v\n", string(bs), err)
	if err != nil {
		return err
	}

	var resp map[string]interface{}
	err = json.Unmarshal(bs, &resp)
	if err != nil {
		return err
	}
	return nil
}
