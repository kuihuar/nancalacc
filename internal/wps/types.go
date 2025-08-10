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

type PostBatchDeleteDeptRequest struct {
	DeptIDs []string `json:"dept_ids"`
}
type PostBatchDeleteDeptResponse struct {
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
	Detail string `json:"detail"`
	Data   struct {
	} `json:"data"`
}

type PostBatchDeleteUserRequest struct {
	UserIDs []string `json:"user_ids"`
}
type PostBatchDeleteUserResponse struct {
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
	Detail string `json:"detail"`
	Data   struct {
	} `json:"data"`
}
type PostRomoveUserIdFromDeptIdRequest struct {
	UserID string `json:"user_id"`
	DeptID string `json:"dept_id"`
}
type PostRomoveUserIdFromDeptIdResponse struct {
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
	Detail string `json:"detail"`
	Data   struct {
	} `json:"data"`
}

type PostAddUserIdToDeptIdRequest struct {
	UserID string `json:"user_id"`
	DeptID string `json:"dept_id"`
}
type PostAddUserIdToDeptIdResponse struct {
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
	Detail string `json:"detail"`
	Data   struct {
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
		DeptID  string `json:"dept_id"`
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

type PostCreateUserRequest struct {
	Avatar           string   `json:"avatar"`
	DeptIDs          []string `json:"dept_ids"`
	Email            string   `json:"email"`
	EmployeeID       string   `json:"employee_id"`
	EmploymentStatus string   `json:"employment_status"`
	EmploymentType   string   `json:"employment_type"`
	ExUserID         string   `json:"ex_user_id"`
	Gender           string   `json:"gender"`
	LeaderID         string   `json:"leader_id"`
	LoginName        string   `json:"login_name"`
	Password         string   `json:"password"`
	Phone            string   `json:"phone"`
	Source           string   `json:"source"`
	Telephone        string   `json:"telephone"`
	Title            string   `json:"title"`
	Titles           struct {
		Mode      string `json:"mode"`
		TitleList []struct {
			DeptID  string `json:"dept_id"`
			TitleID string `json:"title_id"`
		} `json:"title_list"`
	} `json:"titles"`
	UserName  string `json:"user_name"`
	WorkPlace string `json:"work_place"`
}

type PostCreateUserResponse struct {
	Data struct {
		Avatar string `json:"avatar"`
		Ctime  int64  `json:"ctime"`
		Depts  []struct {
			AbsPath string `json:"abs_path"`
			DeptID  string `json:"dept_id"`
			Name    string `json:"name"`
		} `json:"depts"`
		Email            string `json:"email"`
		EmployeeID       string `json:"employee_id"`
		EmploymentStatus string `json:"employment_status"`
		EmploymentType   string `json:"employment_type"`
		ExUserID         string `json:"ex_user_id"`
		Gender           string `json:"gender"`
		ID               string `json:"id"`
		LeaderID         string `json:"leader_id"`
		LoginName        string `json:"login_name"`
		Phone            string `json:"phone"`
		Role             string `json:"role"`
		Source           string `json:"source"`
		Status           string `json:"status"`
		Telephone        string `json:"telephone"`
		Title            string `json:"title"`
		Type             string `json:"type"`
		UserName         string `json:"user_name"`
		WorkPlace        string `json:"work_place"`
	} `json:"data"`
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type PostCreateDeptRequest struct {
	ExDeptID string `json:"ex_dept_id"`
	Leaders  []struct {
		Order  int    `json:"order"`
		UserID string `json:"user_id"`
	} `json:"leaders"`
	Name     string `json:"name"`
	Order    int    `json:"order"`
	ParentID string `json:"parent_id"`
	Source   string `json:"source"` // Default: "inner"
}

type PostCreateDeptResponse struct {
	Data struct {
		AbsPath  string `json:"abs_path"`
		Ctime    int64  `json:"ctime"`
		ExDeptID string `json:"ex_dept_id"`
		ID       string `json:"id"`
		Leaders  []struct {
			Order  int    `json:"order"`
			UserID string `json:"user_id"`
		} `json:"leaders"`
		Name     string `json:"name"`
		Order    int    `json:"order"`
		ParentID string `json:"parent_id"`
		Source   string `json:"source"` // Default: "inner"
	} `json:"data"`
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type PostUpdateDeptRequest struct {
	ExDeptID string   `json:"ex_dept_id"`
	Leaders  []Leader `json:"leaders"`
	Name     string   `json:"name"`
	Order    int      `json:"order"`
	ParentID string   `json:"parent_id"`
	Source   string   `json:"source"`
}

type Leader struct {
	Order  int    `json:"order"`
	UserID string `json:"user_id"`
}

type PostUpdateDeptResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type PostUpdateUserRequest struct {
	Avatar           string    `json:"avatar"`
	Email            string    `json:"email"`
	EmployeeID       string    `json:"employee_id"`
	EmploymentStatus string    `json:"employment_status"`
	EmploymentType   string    `json:"employment_type"`
	ExUserID         string    `json:"ex_user_id"`
	Gender           string    `json:"gender"`
	LeaderID         string    `json:"leader_id"`
	LoginName        string    `json:"login_name"`
	Phone            string    `json:"phone"`
	Source           string    `json:"source"` // Default: "inner"
	Telephone        string    `json:"telephone"`
	Title            string    `json:"title"`
	Titles           TitleInfo `json:"titles"`
	UserName         string    `json:"user_name"`
	WorkPlace        string    `json:"work_place"`
}

type TitleInfo struct {
	Mode      string      `json:"mode"` // Default: "general"
	TitleList []UserTitle `json:"title_list"`
}

type UserTitle struct {
	DeptID  string `json:"dept_id"`
	TitleID string `json:"title_id"`
}

type PostUpdateUserResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type PostBatchUserByPageRequest struct {
	PageSize  int      `json:"page_size"`
	PageNum   int      `json:"page_num"`
	PageToken string   `json:"page_token"`
	WithTotal bool     `json:"with_total"`
	WithDept  bool     `json:"with_dept"`
	Status    []string `json:"status"`
}

type PostBatchUserByPageResponse struct {
	Data struct {
		Items         []UserItem `json:"items"`
		NextPageToken string     `json:"next_page_token"`
		Total         int        `json:"total"`
	} `json:"data"`
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type UserItem struct {
	Avatar           string `json:"avatar"`
	Ctime            int64  `json:"ctime"`
	Depts            []Dept `json:"depts"`
	Email            string `json:"email"`
	EmployeeID       string `json:"employee_id"`
	EmploymentStatus string `json:"employment_status"`
	EmploymentType   string `json:"employment_type"`
	ExUserID         string `json:"ex_user_id"`
	Gender           string `json:"gender"`
	ID               string `json:"id"`
	LeaderID         string `json:"leader_id"`
	LoginName        string `json:"login_name"`
	Phone            string `json:"phone"`
	Role             string `json:"role"`
	Source           string `json:"source"` // Default: "inner"
	Status           string `json:"status"`
	Telephone        string `json:"telephone"`
	Title            string `json:"title"`
	Type             string `json:"type"` // Default: "company_user"
	UserName         string `json:"user_name"`
	WorkPlace        string `json:"work_place"`
}

type Dept struct {
	AbsPath string `json:"abs_path"`
	DeptID  string `json:"dept_id"`
	Name    string `json:"name"`
}

type GetDeptByPageRequest struct {
	PageSize  int    `json:"page_size"`
	PageToken string `json:"page_token"`
	WithTotal bool   `json:"with_total"`
	Recursive bool   `json:"recursive"`
}

type GetDeptByPageResponse struct {
	Data struct {
		Items         []Department `json:"items"`
		NextPageToken string       `json:"next_page_token"`
		Total         int          `json:"total"`
	} `json:"data"`
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type Department struct {
	AbsPath  string   `json:"abs_path"`
	Ctime    int64    `json:"ctime"`
	ExDeptID string   `json:"ex_dept_id"`
	ID       string   `json:"id"`
	Leaders  []Leader `json:"leaders"`
	Name     string   `json:"name"`
	Order    int      `json:"order"`
	ParentID string   `json:"parent_id"`
	Source   string   `json:"source"` // Default: "inner"
}

type GetUserDeptsByUserIdRequest struct {
	UserID string `json:"user_id"`
}
type GetUserDeptsByUserIdResponse struct {
	Code   int    `json:"code"`
	Detail string `json:"detail"`
	Msg    string `json:"msg"`
	Data   struct {
		Items []struct {
			AbsPath  string `json:"abs_path"`
			Ctime    int64  `json:"ctime"`
			ExDeptID string `json:"ex_dept_id"`
			ID       string `json:"id"`
			Name     string `json:"name"`
			Order    int64  `json:"order"`
			ParentID string `json:"parent_id"`
		} `json:"items"`
	} `json:"data"`
}

type GetDeptChildrenRequest struct {
	DeptID    string `json:"dept_id"`
	Recursive bool   `json:"recursive"`
	PageSize  int    `json:"page_size"`
	PageToken string `json:"page_token"`
	WithTotal bool   `json:"with_total"`
}

//	type GetDeptChildrenResponse struct {
//		Code int    `json:"code"`
//		Msg  string `json:"msg"`
//		Data struct {
//			Items []Department `json:"items"`
//		} `json:"data"`
//	}
type DeptItem struct {
	AbsPath  string `json:"abs_path"`
	Ctime    int64  `json:"ctime"`
	ExDeptID string `json:"ex_dept_id"`
	ID       string `json:"id"`
	Leaders  []struct {
		Order  int    `json:"order"`
		UserID string `json:"user_id"`
	} `json:"leaders"`
	Name     string `json:"name"`
	Order    int    `json:"order"`
	ParentID string `json:"parent_id"`
	Source   string `json:"source"` // 枚举值如 "inner" 可以用 string 类型
}
type GetDeptChildrenResponse struct {
	Data struct {
		Items         []DeptItem `json:"items"`
		NextPageToken string     `json:"next_page_token"`
		Total         int        `json:"total"`
	} `json:"data"`
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type GetCompAllUsersRequest struct {
	PageSize  int      `json:"page_size"`
	PageToken string   `json:"page_token"`
	WithTotal bool     `json:"with_total"`
	WithDept  bool     `json:"with_dept"`
	Status    []string `json:"status"`
	Recursive bool     `json:"recursive"`
}

type GetCompAllUsersResponse struct {
	Data struct {
		Items         []UserItem `json:"items"`
		NextPageToken string     `json:"next_page_token"`
		Total         int        `json:"total"`
	} `json:"data"`
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}
