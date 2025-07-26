package wps

import (
	"context"
)

type Wps interface {
	CallEcisaccountsyncAll(ctx context.Context, accessToken string, taskId string) (EcisaccountsyncAllResponse, error)
	CallEcisaccountsyncIncrement(ctx context.Context, accessToken string, thirdCompanyId string) (EcisaccountsyncIncrementResponse, error)
}
