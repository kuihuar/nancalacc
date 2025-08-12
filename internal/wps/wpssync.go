package wps

import (
	"context"
	"encoding/json"
)

var (
	Source = "sync"
)

func (ws *wps) PostEcisaccountsyncAll(ctx context.Context, accessToken string, input *EcisaccountsyncAllRequest) (EcisaccountsyncAllResponse, error) {

	ws.log.Infof("PostEcisaccountsyncAll input:%+v", input)
	var resp EcisaccountsyncAllResponse

	// input := &EcisaccountsyncAllRequest{
	// 	ThirdCompanyId: thirdCompanyId,
	// }
	input.CollectCost = CollectCost
	ak := ws.cfg.ClientId
	sk := ws.cfg.ClientSecret
	wpsReq := NewWPSRequest(DOMAIN, ak, sk)

	ws.log.Infof("PostEcisaccountsyncAll.req path:%s,input:%+v\n", ECISACCOUNTSYNC_PATH_INCREMENT, input)

	bs, err := wpsReq.PostJSON(context.Background(), ECISACCOUNTSYNC_PATH_INCREMENT, accessToken, input)

	ws.log.Infof("PostEcisaccountsyncAll.res bs:%s,err:%+v\n", string(bs), err)
	if err != nil {
		return resp, err
	}

	err = json.Unmarshal(bs, &resp)
	if err != nil {
		return resp, err
	}
	if resp.Code != "200" {
		return resp, ErrCodeNot200
	}

	return resp, nil
}

func (ws *wps) PostEcisaccountsyncIncrement(ctx context.Context, accessToken string, input *EcisaccountsyncIncrementRequest) (EcisaccountsyncIncrementResponse, error) {

	var resp EcisaccountsyncIncrementResponse

	// input := &EcisaccountsyncIncrementRequest{
	// 	ThirdCompanyId: thirdCompanyId,
	// }
	ak := ws.cfg.ClientId
	sk := ws.cfg.ClientSecret
	wpsReq := NewWPSRequest(DOMAIN, ak, sk)

	ws.log.Infof("PostEcisaccountsyncIncrement.req path:%s,input:%+v\n", ECISACCOUNTSYNC_PATH_INCREMENT, input)

	bs, err := wpsReq.PostJSON(context.Background(), ECISACCOUNTSYNC_PATH_INCREMENT, accessToken, input)

	ws.log.Infof("PostEcisaccountsyncIncrement.res bs:%s,err:%+v\n", string(bs), err)
	if err != nil {
		return resp, err
	}

	err = json.Unmarshal(bs, &resp)
	if err != nil {
		return resp, err
	}
	if resp.Code != "200" {
		return resp, ErrCodeNot200
	}

	return resp, nil
}
