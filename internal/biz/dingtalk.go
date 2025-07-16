package biz

import "context"

type DingTalkRepo interface {
	GetAccessToken(ctx context.Context, code string) (string, error)
	FetchDepartments(ctx context.Context, token string) ([]*DingtalkDept, error)
	FetchDepartmentUsers(ctx context.Context, token string, deptId []int64) ([]*DingtalkDeptUser, error)

	GetUserAccessToken(ctx context.Context, code string) (*AuthResponse, error)
	GetUserInfo(ctx context.Context, token, unionId string) (*DingTalkUserInfo, error)

	GetUseridByUnionid(ctx context.Context, token, unionid string) (string, error)
}
