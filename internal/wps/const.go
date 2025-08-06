package wps

import "errors"

const (
	DOMAIN                         = "http://119.3.173.229/openapi"
	ECISACCOUNTSYNC_PATH_FULL      = "/v7/ecisaccountsync/full/sync"
	ECISACCOUNTSYNC_PATH_INCREMENT = "/v7/ecisaccountsync/increment/sync"

	GET_DEPARTMENT_ROOT          = "/v7/depts/root"
	BATCH_POST_USERS_PATH        = "/v7/users/batch_read"
	BATCH_POST_DEPTS_PATH        = "/v7/depts/batch_read"
	POST_DEPTS_BY_EXDEPTIDS_PATH = "/v7/depts/by_ex_dept_ids"
	POST_USERS_BY_EXDEPTIDS_PATH = "/v7/users/by_ex_user_ids"

	GET_USER_BY_USERID = "/v7/users"

	POST_CACHE_SET = "http://encs-pri-cams-engine/i/cams/sdk/api/v1/cache/set"
	POST_CACHE_GET = "http://encs-pri-cams-engine/i/cams/sdk/api/v1/cache/get"
	POST_CACHE_DEL = "http://encs-pri-cams-engine/i/cams/sdk/api/v1/cache/del"
)

var (
	ErrCodeNot200     = errors.New("HTTP status code is not 200")
	ErrCodeNot0       = errors.New("HTTP status code is not 0")
	CollectCost   int = 1100000
)
