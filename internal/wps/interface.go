package wps

import (
	"context"
)

type WpsSync interface {
	CallEcisaccountsyncAll(ctx context.Context, accessToken string, input *EcisaccountsyncAllRequest) (EcisaccountsyncAllResponse, error)
	CallEcisaccountsyncIncrement(ctx context.Context, accessToken string, input *EcisaccountsyncIncrementRequest) (EcisaccountsyncIncrementResponse, error)
}

type Wps interface {
	BatchGetDepartment(ctx context.Context, accessToken string, req *BatchGetDepartmentRequest) (BatchGetDepartmentResponse, error)
}
