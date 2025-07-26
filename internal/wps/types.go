package wps

type EcisaccountsyncAllRequest struct {
	TaskId         string `json:"task_id"`
	ThirdCompanyId string `json:"third_company_id"`
	CollectCost    int    `json:"collect_cost"`
}

type EcisaccountsyncAllResponse struct {
	Code   string `json:"code"`
	Msg    string `json:"msg"`
	Data   any    `json:"data"`
	Detail string `json:"detail"`
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

type BatchGetDepartmentRequest struct {
	DeptIDs []string `json:"dept_ids"`
}

type WpsDepartmentItem struct {
	ID       string                `json:"id"`
	Name     string                `json:"name"`
	ParentID string                `json:"parent_id"`
	Order    int                   `json:"order"`
	Leaders  []WpsDepartmentLeader `json:"leaders"`
	AbsPath  string                `json:"abs_path"`
	Ctime    int                   `json:"ctime"`
	ExDeptID string                `json:"ex_dept_id"`
}

type WpsDepartmentLeader struct {
	Order  int    `json:"order"`
	UserID string `json:"user_id"`
}
type BatchGetDepartmentResponse struct {
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
	Detail string `json:"detail"`
	Data   struct {
		DepartmentList []struct {
			DeptID   string `json:"dept_id"`
			DeptName string `json:"dept_name"`
		} `json:"items"`
	} `json:"data"`
}
