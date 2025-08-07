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

	// 获取根部门(4)
	GetDepartmentRoot(ctx context.Context, accessToken string, req GetDepartmentRootRequest) (GetDepartmentRootResponse, error)
	//GetDepartmentChildrenList(ctx context.Context, accessToken string, req GetDepartmentChildrenListRequest) (GetDepartmentChildrenListResponse, error)
	// user..
	BatchPostUsers(ctx context.Context, accessToken string, input BatchPostUsersRequest) (BatchPostUsersResponse, error)
	PostBatchUsersByExDepIds(ctx context.Context, accessToken string, input PostBatchUsersByExDepIdsRequest) (*PostBatchUsersByExDepIdsResponse, error)

	// dept...
	BatchPostDepartments(ctx context.Context, accessToken string, req BatchPostDepartmentsRequest) (BatchPostDepartmentsResponse, error)
	PostBatchDepartmentsByExDepIds(ctx context.Context, accessToken string, input PostBatchDepartmentsByExDepIdsRequest) (*PostBatchDepartmentsByExDepIdsResponse, error)

	// 批量删除部门
	PostBatchDeleteDept(ctx context.Context, accessToken string, input PostBatchDeleteDeptRequest) (*PostBatchDeleteDeptResponse, error)
	// 批量删除用户
	PostBatchDeleteUser(ctx context.Context, accessToken string, input PostBatchDeleteUserRequest) (*PostBatchDeleteUserResponse, error)

	// 用户离开部门
	PostRomoveUserIdFromDeptId(ctx context.Context, accessToken string, input PostRomoveUserIdFromDeptIdRequest) (*PostRomoveUserIdFromDeptIdResponse, error)
	// 用户加入部门(1)
	PostAddUserIdToDeptId(ctx context.Context, accessToken string, input PostAddUserIdToDeptIdRequest) (*PostAddUserIdToDeptIdResponse, error)
	// 用户创建(2)
	PostCreateUser(ctx context.Context, accessToken string, input PostCreateUserRequest) (*PostCreateUserResponse, error)
	// 部门创建(3)
	PostCreateDept(ctx context.Context, accessToken string, input PostCreateDeptRequest) (*PostCreateDeptResponse, error)

	// 部门更新
	PostUpdateDept(ctx context.Context, accessToken string, input PostUpdateDeptRequest) (*PostUpdateDeptResponse, error)
	// 用户更新
	PostUpdateUser(ctx context.Context, accessToken string, input PostUpdateUserRequest) (*PostUpdateUserResponse, error)

	// 批量获取所有用户（分页）
	PostBatchUserByPage(ctx context.Context, accessToken string, input PostBatchUserByPageRequest) (*PostBatchUserByPageResponse, error)

	// 批量获取所有部门(查询子部门列表分页)
	GetDeptByPage(ctx context.Context, accessToken string, input GetDeptByPageRequest) (*GetDeptByPageResponse, error)

	GetUserByUserId(ctx context.Context, accessToken string, req GetUserByUserIdRequest) (GetUserByUserIdResponse, error)

	CacheSet(ctx context.Context, accessToken string, key string, value interface{}, expiration time.Duration) error
	CacheGet(ctx context.Context, accessToken string, key string) (interface{}, error)
	CacheDel(ctx context.Context, accessToken string, key string) error
}
