package wps

import "errors"

const (
	DOMAIN                         = "http://119.3.173.229/openapi"
	ECISACCOUNTSYNC_PATH_FULL      = "/v7/ecisaccountsync/full/sync"
	ECISACCOUNTSYNC_PATH_INCREMENT = "/v7/ecisaccountsync/increment/sync"
)

var (
	ErrCodeNot200 = errors.New("HTTP status code is not 200")
)
