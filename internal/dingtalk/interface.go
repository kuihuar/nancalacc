package dingtalk

import (
	"context"
	"nancalacc/internal/auth"
)

type Dingtalk interface {
	GetAccessToken(ctx context.Context) (*auth.AccessTokenResp, error)
	FetchDepartments(ctx context.Context, token string) ([]*DingtalkDept, error)
	FetchDepartmentUsers(ctx context.Context, token string, deptIds []int64) ([]*DingtalkDeptUser, error)

	GetUserAccessToken(ctx context.Context, code string) (*AuthResponse, error)
	GetUserInfo(ctx context.Context, token, unionId string) (*DingTalkUserInfo, error)

	GetUseridByUnionid(ctx context.Context, token, unionid string) (string, error)

	FetchDeptDetails(ctx context.Context, token string, deptIds []int64) ([]*DingtalkDept, error)

	FetchUserDetail(ctx context.Context, token string, userIds []string) ([]*DingtalkDeptUser, error)
}
