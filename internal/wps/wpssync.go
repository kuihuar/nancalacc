package wps

import (
	"context"
	"encoding/json"
	"nancalacc/internal/conf"
	httpwps "nancalacc/pkg/httputil/wps"

	"github.com/go-kratos/kratos/v2/log"
)

type wpsSync struct {
	serviceConf *conf.Service
	log         *log.Helper
}

var (
	Source = "sync"
)

func NewWpsSync(serviceConf *conf.Service, logger log.Logger) WpsSync {
	return &wpsSync{
		serviceConf: serviceConf,
		log:         log.NewHelper(logger),
	}
}

func (ws *wpsSync) CallEcisaccountsyncAll(ctx context.Context, accessToken string, input *EcisaccountsyncAllRequest) (EcisaccountsyncAllResponse, error) {

	var resp EcisaccountsyncAllResponse

	// input := &EcisaccountsyncAllRequest{
	// 	ThirdCompanyId: thirdCompanyId,
	// }
	input.CollectCost = CollectCost
	ak := ws.serviceConf.Auth.App.ClientId
	sk := ws.serviceConf.Auth.App.ClientSecret
	wpsReq := httpwps.NewWPSRequest(DOMAIN, ak, sk)

	bs, err := wpsReq.PostJSON(context.Background(), ECISACCOUNTSYNC_PATH_INCREMENT, accessToken, input)
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

	// w.log.Infof("CallEcisaccountsyncAll: %v", acctaskId)

	// path := w.serviceConf.Business.EcisaccountsyncUrl

	// // path := "http://encs-pri-proxy-gateway/ecisaccountsync/api/sync/all"
	// var resp EcisaccountsyncAllResponse

	// thirdCompanyID := w.serviceConf.Business.ThirdCompanyId

	// collectCost := "1100000"
	// uri := fmt.Sprintf("%s?taskId=%s&thirdCompanyId=%s&collectCost=%s", path, acctaskId, thirdCompanyID, collectCost)

	// w.log.Infof("CallEcisaccountsyncAll uri: %s", uri)
	// bs, err := httputil.PostJSON(uri, nil, time.Second*10)
	// w.log.Infof("CallEcisaccountsyncAll.Post output: bs:%s, err:%w", string(bs), err)

	// if err != nil {
	// 	return resp, err
	// }
	// err = json.Unmarshal(bs, &resp)
	// if err != nil {
	// 	return resp, fmt.Errorf("Unmarshal err: %w", err)
	// }
	// if resp.Code != "200" {
	// 	return resp, fmt.Errorf("code not 200: %s", resp.Code)
	// }

	// return resp, nil
}

func (ws *wpsSync) CallEcisaccountsyncIncrement(ctx context.Context, accessToken string, input *EcisaccountsyncIncrementRequest) (EcisaccountsyncIncrementResponse, error) {

	var resp EcisaccountsyncIncrementResponse

	// input := &EcisaccountsyncIncrementRequest{
	// 	ThirdCompanyId: thirdCompanyId,
	// }
	ak := ws.serviceConf.Auth.App.ClientId
	sk := ws.serviceConf.Auth.App.ClientSecret
	wpsReq := httpwps.NewWPSRequest(DOMAIN, ak, sk)

	bs, err := wpsReq.PostJSON(context.Background(), ECISACCOUNTSYNC_PATH_INCREMENT, accessToken, input)
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
