package dingtalk

type ListDeptIDRequest struct {
	DeptID int64 `json:"dept_id"`
}
type ListDeptIDResponse struct {
	Errcode int `json:"errcode"`
	Result  struct {
		DeptIDList []int64 `json:"dept_id_list"`
	} `json:"result"`
	Errmsg    string `json:"errmsg"`
	RequestID string `json:"request_id"`
}
type DingtalkDept struct {
	AutoAddUser           bool     `json:"auto_add_user"`
	AutoApproveApply      bool     `json:"auto_approve_apply"`
	Brief                 string   `json:"brief"`
	CreateDeptGroup       bool     `json:"create_dept_group"`
	DeptGroupChatID       string   `json:"dept_group_chat_id"`
	DeptID                int64    `json:"dept_id"`
	DeptManagerUseridList []string `json:"dept_manager_userid_list"`
	DeptPermits           []string `json:"dept_permits"`
	EmpApplyJoinDept      bool     `json:"emp_apply_join_dept"`
	GroupContainSubDept   bool     `json:"group_contain_sub_dept"`
	HideDept              bool     `json:"hide_dept"`
	MemberCount           int      `json:"member_count"`
	Name                  string   `json:"name"`
	Order                 int64    `json:"order"`
	OrgDeptOwner          string   `json:"org_dept_owner"`
	OuterDept             bool     `json:"outer_dept"`
	OuterPermitDepts      []string `json:"outer_permit_depts"`
	OuterPermitUsers      []string `json:"outer_permit_users"`
	OwningMemberCount     int      `json:"owning_member_count"`
	ParentID              int64    `json:"parent_id"`
	UserPermits           []string `json:"user_permits"`
}

type DingtalkDeptRequest struct {
	DeptID   int64  `json:"dept_id"`
	Language string `json:"language"`
}
type DingtalkDeptResponse struct {
	Errcode int          `json:"errcode"`
	Result  DingtalkDept `json:"result"`
	Errmsg  string       `json:"errmsg"`
}

// type ListDeptResponse struct {
// 	Errcode int            `json:"errcode"`
// 	Result  []DingtalkDept `json:"result"`
// 	Errmsg  string         `json:"errmsg"`
// }

// type ListDeptRequest struct {
// 	DeptID int64 `json:"dept_id"`
// }

type FetchUserDetailRequest struct {
	UserIDs []string `json:"userid"`
}
type FetchUserDetailResponse struct {
	UserIDs []string `json:"userid"`
}

type DingtalkDeptUser struct {
	Active        bool    `json:"active"`
	Admin         bool    `json:"admin"`
	Avatar        string  `json:"avatar"`
	Boss          bool    `json:"boss"`
	CreateTime    string  `json:"create_time"`
	DeptIDList    []int64 `json:"dept_id_list"`
	DeptOrderList []struct {
		DeptID int64 `json:"dept_id"`
		Order  int64 `json:"order"`
	} `json:"dept_order_list"`
	Email            string `json:"email"`
	ExclusiveAccount bool   `json:"exclusive_account"`
	HideMobile       bool   `json:"hide_mobile"`
	HiredDate        int64  `json:"hired_date"`
	JobNumber        string `json:"job_number"`
	LeaderInDept     []struct {
		DeptID int64 `json:"dept_id"`
		Leader bool  `json:"leader"`
	} `json:"leader_in_dept"`
	Mobile     string `json:"mobile"`
	Name       string `json:"name"`
	RealAuthed bool   `json:"real_authed"`
	Remark     string `json:"remark"`
	Senior     bool   `json:"senior"`
	StateCode  string `json:"state_code"`
	Telephone  string `json:"telephone"`
	Title      string `json:"title"`
	Unionid    string `json:"unionid"`
	Userid     string `json:"userid"`
	WorkPlace  string `json:"work_place"`
}

type DingtalkDeptUserRelation struct {
	Uid            string `json:"uid"`
	Did            string `json:"did"`
	ThirdCompanyID string `json:"third_company_id"`

	PlatformID string `json:"platform_id"`
	Order      int64  `json:"order"`
}

type ListDeptUserRequest struct {
	Cursor             int64  `json:"cursor"`
	ContainAccessLimit bool   `json:"contain_access_limit"`
	Size               int64  `json:"size"`
	OrderField         string `json:"order_field"`
	Language           string `json:"language"`
	DeptID             int64  `json:"dept_id"`
}

type ListDeptUserResponse struct {
	Errcode int `json:"errcode"`
	Result  struct {
		NextCursor int64              `json:"next_cursor"`
		HasMore    bool               `json:"has_more"`
		List       []DingtalkDeptUser `json:"list"`
	} `json:"result"`
	Errmsg string `json:"errmsg"`
}

//	//{
//		"expireIn":7200,
//		"accessToken":"6f874309a6c031f9a2033a54dcafadae",
//		"refreshToken":"d5a84a019bf23e2e9d7803666175c337"
//	}
type AuthResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpireIn     int    `json:"expireIn"`
}

type DingTalkUserInfo struct {
	AvatarUrl string `json:"avatarUrl,omitempty"`
	Email     string `json:"email,omitempty"`
	Mobile    string `json:"mobile,omitempty"`
	Nick      string `json:"nick,omitempty"`
	OpenId    string `json:"openId,omitempty"`
	StateCode string `json:"stateCode,omitempty"`
	UnionId   string `json:"unionId,omitempty"`
	Visitor   bool   `json:"visitor,omitempty"`
}

type DingtalkCompanyCfg struct {
}

type DingTalkUseridByUnionidRequest struct {
	Unionid string `json:"unionid"`
}

type DingTalkUseridByUnionidResponse struct {
	Errcode int `json:"errcode"`
	Result  struct {
		ContactType int    `json:"contact_type"`
		Userid      string `json:"userid"`
	} `json:"result"`
	Errmsg string `json:"errmsg"`
}

type DingTalkUserDetailRequest struct {
	Userid string `json:"userid"`
}

type DingTalkUserDetailResponse struct {
	Errcode int              `json:"errcode"`
	Result  DingtalkDeptUser `json:"result"`
	Errmsg  string           `json:"errmsg"`
}

type UserModifyOrgEventData struct {
	TimeStamp string   `json:"timeStamp"`
	UserId    []string `json:"userId"`
	DiffInfo  struct {
		Prev struct {
			ManagerUserid string `json:"managerUserid"`
			HiredDate     string `json:"hiredDate"`
			Name          string `json:"name"`
			Telephone     string `json:"telephone"`
			Email         string `json:"email"`
			JobNumber     string `json:"jobNumber"`
			WorkPlace     string `json:"workPlace"`
		} `json:"prev"`
		Curr struct {
			ManagerUserid string `json:"managerUserid"`
			HiredDate     string `json:"hiredDate"`
			Name          string `json:"name"`
			Email         string `json:"email"`
			JobNumber     string `json:"jobNumber"`
			WorkPlace     string `json:"workPlace"`
		} `json:"curr"`
		Userid string `json:"userid"`
	} `json:"diffInfo"`
}
type UserModifyOrgEvent struct {
	EventUnifiedAppId string                 `json:"eventUnifiedAppId"`
	EventCorpId       string                 `json:"eventCorpId"`
	EventType         string                 `json:"eventType"`
	EventId           string                 `json:"eventId"`
	EventBornTime     int64                  `json:"eventBornTime"`
	Data              UserModifyOrgEventData `json:"data"`
}
