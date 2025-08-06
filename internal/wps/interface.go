package wps

import (
	"context"
	"time"
)

type WpsSync interface {
	PostEcisaccountsyncAll(ctx context.Context, accessToken string, input *EcisaccountsyncAllRequest) (EcisaccountsyncAllResponse, error)
	PostEcisaccountsyncIncrement(ctx context.Context, accessToken string, input *EcisaccountsyncIncrementRequest) (EcisaccountsyncIncrementResponse, error)
}

type Wps interface {

	// nouse
	GetDepartmentRoot(ctx context.Context, accessToken string, req GetDepartmentRootRequest) (GetDepartmentRootResponse, error)
	//GetDepartmentChildrenList(ctx context.Context, accessToken string, req GetDepartmentChildrenListRequest) (GetDepartmentChildrenListResponse, error)
	// user..
	BatchPostUsers(ctx context.Context, accessToken string, input BatchPostUsersRequest) (BatchPostUsersResponse, error)
	PostBatchUsersByExDepIds(ctx context.Context, accessToken string, input PostBatchUsersByExDepIdsRequest) (*PostBatchUsersByExDepIdsResponse, error)

	// dept...
	BatchPostDepartments(ctx context.Context, accessToken string, req BatchPostDepartmentsRequest) (BatchPostDepartmentsResponse, error)
	PostBatchDepartmentsByExDepIds(ctx context.Context, accessToken string, input PostBatchDepartmentsByExDepIdsRequest) (*PostBatchDepartmentsByExDepIdsResponse, error)

	GetUserByUserId(ctx context.Context, accessToken string, req GetUserByUserIdRequest) (GetUserByUserIdResponse, error)

	CacheSet(ctx context.Context, accessToken string, key string, value interface{}, expiration time.Duration) error
	CacheGet(ctx context.Context, accessToken string, key string) (interface{}, error)
	CacheDel(ctx context.Context, accessToken string, key string) error
}
