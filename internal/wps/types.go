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

// dept...
type BatchPostDepartmentsRequest struct {
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

type BatchPostDepartmentsResponse struct {
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
	Detail string `json:"detail"`
	Data   struct {
		Items []WpsDepartmentItem `json:"items"`
	} `json:"data"`
}

type PostBatchDepartmentsByExDepIdsRequest struct {
	ExDeptIDs []string `json:"ex_dept_ids"`
}
type PostBatchDepartmentsByExDepIdsResponse struct {
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
	Detail string `json:"detail"`
	Data   struct {
		Items []WpsDepartmentItem `json:"items"`
	} `json:"data"`
}
type GetDepartmentRootRequest struct {
}

type GetDepartmentRootResponse struct {
	Code   int               `json:"code"`
	Msg    string            `json:"msg"`
	Detail string            `json:"detail"`
	Data   WpsDepartmentItem `json:"data"`
}

type GetDepartmentChildrenListRequest struct {
	Recursive bool   `json:"recursive"`
	PageSize  int    `json:"page_size"`
	PageToken string `json:"page_token"`
	WithTotal bool   `json:"with_total"`
}
type GetDepartmentChildrenListResponse struct {
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
	Detail string `json:"detail"`
	Data   struct {
		Items []WpsDepartmentItem `json:"items"`
	} `json:"data"`
	NextPageToken string `json:"next_page_token"`
	Total         int    `json:"total"`
}

var (
	UserStatusActive   = "active"
	UserStatusNoActive = "notactive"
	UserStatusDisabled = "disabled"
)

// user...
type PostBatchUsersByExDepIdsRequest struct {
	ExUserIDs []string `json:"ex_user_ids"`
	Status    []string `json:"status"`
}
type PostBatchUsersByExDepIdsResponse struct {
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
	Detail string `json:"detail"`
	Data   struct {
		Items []WpsUserItem `json:"items"`
	} `json:"data"`
}
type WpsUserItem struct {
	ID        string `json:"id"`
	UserName  string `json:"user_name"`
	LoginName string `json:"login_name"`
	Avatar    string `json:"avatar"`
	Email     string `json:"email"`
	ExUserId  string `json:"ex_user_id"`
	Phone     string `json:"phone"`
	Role      string `json:"role"`
	Status    string `json:"status"`
	Ctime     int    `json:"ctime"`
}

type BatchPostUsersRequest struct {
	UserIDs  []string `json:"user_ids"`
	Status   []string `json:"status"`
	WithDept bool     `json:"with_dept"`
}
type BatchPostUsersResponse struct {
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
	Detail string `json:"detail"`
	Data   struct {
		Items []WpsUserWithDept `json:"items"`
	} `json:"data"`
}

type WpsUserWithDept struct {
	ID        string `json:"id"`
	ExUserID  string `json:"ex_user_id"`
	Gender    string `json:"gender"`
	Telephone string `json:"telephone"`
	Status    string `json:"status"`
	Depts     []struct {
		AbsPath string `json:"abs_path"`
		ID      string `json:"id"`
		Name    string `json:"name"`
	} `json:"depts"`

	Ctime      int    `json:"ctime"`
	Role       string `json:"role"`
	LoginName  string `json:"login_name"`
	UserName   string `json:"user_name"`
	Avatar     string `json:"avatar"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	Title      string `json:"title"`
	WorkPlace  string `json:"work_place"`
	EmployeeID string `json:"employee_id"`

	//Type string `json:"type"`
}

type GetUserByUserIdRequest struct {
	UserID string `json:"user_id"`
}

type GetUserByUserIdResponse struct {
	Code   int             `json:"code"`
	Msg    string          `json:"msg"`
	Detail string          `json:"detail"`
	Data   WpsUserWithDept `json:"data"`
}
