package wps

import (
	"context"
)

var (
	Source = "sync"
)

func (ws *wps) PostEcisaccountsyncAll(ctx context.Context, accessToken string, input *EcisaccountsyncAllRequest) (EcisaccountsyncAllResponse, error) {
	var resp EcisaccountsyncAllResponse

	// 设置收集成本标志
	input.CollectCost = CollectCost

	// 记录请求
	logAPIRequest(ctx, ws.log, "PostEcisaccountsyncAll", "POST", ECISACCOUNTSYNC_PATH_INCREMENT, input)

	// 创建请求
	wpsReq := NewWPSRequest(DOMAIN, ws.cfg.ClientId, ws.cfg.ClientSecret, WithLogger(ws.log))

	// 执行请求
	bs, err := wpsReq.PostJSON(ctx, ECISACCOUNTSYNC_PATH_INCREMENT, accessToken, input)
	if err != nil {
		ws.log.WithContext(ctx).Errorf("PostEcisaccountsyncAll request failed: %v", err)
		return resp, err
	}

	// 使用泛型函数处理响应
	if err := handleAPIResponse(ctx, ws.log, "PostEcisaccountsyncAll", bs, &resp, "200"); err != nil {
		return resp, err
	}

	return resp, nil
}

func (ws *wps) PostEcisaccountsyncIncrement(ctx context.Context, accessToken string, input *EcisaccountsyncIncrementRequest) (EcisaccountsyncIncrementResponse, error) {
	var resp EcisaccountsyncIncrementResponse

	// 记录请求
	logAPIRequest(ctx, ws.log, "PostEcisaccountsyncIncrement", "POST", ECISACCOUNTSYNC_PATH_INCREMENT, input)

	// 创建请求
	wpsReq := NewWPSRequest(DOMAIN, ws.cfg.ClientId, ws.cfg.ClientSecret, WithLogger(ws.log))

	// 执行请求
	bs, err := wpsReq.PostJSON(ctx, ECISACCOUNTSYNC_PATH_INCREMENT, accessToken, input)
	if err != nil {
		ws.log.WithContext(ctx).Errorf("PostEcisaccountsyncIncrement request failed: %v", err)
		return resp, err
	}

	// 使用泛型函数处理响应
	if err := handleAPIResponse(ctx, ws.log, "PostEcisaccountsyncIncrement", bs, &resp, "200"); err != nil {
		return resp, err
	}

	return resp, nil
}
