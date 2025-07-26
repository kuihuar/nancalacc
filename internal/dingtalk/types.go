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
	DeptID   int64  `json:"dept_id"`
	Name     string `json:"name"`
	ParentID int64  `json:"parent_id"`

	SourceIdentifier  string `json:"source_identifier"`
	CreateDeptGroup   bool   `json:"create_dept_group"`
	AutoAddUser       bool   `json:"auto_add_user"`
	Order             int64  `json:"order"`
	MemberCount       int64  `json:"member_count"`
	OwningMemberCount int64  `json:"owning_member_count"`
	//Tags             string `json:"tags"`
	//FromUnionOrg     bool   `json:"from_union_org"`

	//DeptGroupChatID     string   `json:"dept_group_chat_id"`
	//DeptPermits          []int64         `json:"dept_permits"`
	// GroupContainSubDept bool `json:"group_contain_sub_dept"`
	// OrgDeptOwner       string   `json:"org_dept_owner"`
	//OuterPermitUsers []string `json:"outer_permit_users"`
	//DeptManagerUserIDs []string `json:"dept_manager_userid_list"`

	//OuterDept bool `json:"outer_dept"`

	//HideDept bool `json:"hide_dept"`

	//OuterPermitDepts []int64  `json:"outer_permit_depts"`
	//UserPermits      []string `json:"user_permits"`

	Code         string `json:"code"`
	UnionDeptExt struct {
		CorpID string `json:"corp_id"`
		DeptID int64  `json:"dept_id"`
	} `json:"union_dept_ext"`
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

type DingtalkDeptUser struct {
	Userid     string  `json:"userid"`  //zhangsan 用户的userId
	Unionid    string  `json:"unionid"` //用户姓名
	Name       string  `json:"name"`
	Avatar     string  `json:"avatar"`
	Mobile     string  `json:"mobile"`
	Email      string  `json:"email"`
	Remark     string  `json:"remark"`
	DeptIDList []int64 `json:"dept_id_list"`
	Extension  string  `json:"extension"`
	Active     bool    `json:"active"`
	Boss       bool    `json:"boss"`
	Admin      bool    `json:"admin"`
	Title      string  `json:"title"`
	Leader     bool    `json:"leader"`
	Nickname   string  `json:"nickname"`

	//ExclusiveAccountType string `json:"exclusive_account_type"`

	//ExclusiveAccount string `json:"exclusive_account"`
	//HiredDate        string `json:"hired_date"`

	//WorkPlace string `json:"work_place"`

	//JobNumber string `json:"job_number"`
	DeptOrder int `json:"dept_order"`
	//LoginID                  string `json:"login_id"`
	//ExclusiveAccountCorpName string `json:"exclusive_account_corp_name"`
	//Telephone                string `json:"telephone"`

	//HideMobile             string `json:"hide_mobile"`
	//ExclusiveAccountCorpID string `json:"exclusive_account_corp_id"`
	//OrgEmail               string `json:"org_email"`
	StateCode string `json:"state_code"`
}

type DingtalkDeptUserRelation struct {
	Uid            string `json:"uid"`
	Did            string `json:"did"`
	ThirdCompanyID string `json:"third_company_id"`

	PlatformID string `json:"platform_id"`
	Order      int    `json:"order"`
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
