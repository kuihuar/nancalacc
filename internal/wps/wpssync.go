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

func (ws *wpsSync) PostEcisaccountsyncAll(ctx context.Context, accessToken string, input *EcisaccountsyncAllRequest) (EcisaccountsyncAllResponse, error) {

	ws.log.Infof("PostEcisaccountsyncAll input:%+v", input)
	var resp EcisaccountsyncAllResponse

	// input := &EcisaccountsyncAllRequest{
	// 	ThirdCompanyId: thirdCompanyId,
	// }
	input.CollectCost = CollectCost
	ak := ws.serviceConf.Auth.App.ClientId
	sk := ws.serviceConf.Auth.App.ClientSecret
	wpsReq := httpwps.NewWPSRequest(DOMAIN, ak, sk)

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

func (ws *wpsSync) PostEcisaccountsyncIncrement(ctx context.Context, accessToken string, input *EcisaccountsyncIncrementRequest) (EcisaccountsyncIncrementResponse, error) {

	var resp EcisaccountsyncIncrementResponse

	// input := &EcisaccountsyncIncrementRequest{
	// 	ThirdCompanyId: thirdCompanyId,
	// }
	ak := ws.serviceConf.Auth.App.ClientId
	sk := ws.serviceConf.Auth.App.ClientSecret
	wpsReq := httpwps.NewWPSRequest(DOMAIN, ak, sk)

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
