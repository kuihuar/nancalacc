package wps

type EcisaccountsyncAllRequest struct {
	TaskId         string `json:"taskId"`
	ThirdCompanyId string `json:"thirdCompanyId"`
	CollectCost    string `json:"collectCost"`
}

type EcisaccountsyncAllResponse struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
}

type EcisaccountsyncIncrementRequest struct {
	ThirdCompanyId string `json:"third_company_id"`
}

type EcisaccountsyncIncrementResponse struct {
	Code   string `json:"code"`
	Msg    string `json:"msg"`
	Data   any    `json:"data"`
	Detail string `json:"detail"`
}
