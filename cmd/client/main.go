package main

import (
	"flag"
	"fmt"

	// "fmt"
	"nancalacc/internal/conf"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
)

var (
	bc conf.Bootstrap
)

func init() {
	var flagconf string
	flag.StringVar(&flagconf, "conf", "../../configs", "config path, eg: -conf config.yaml")
	flag.Parse()

	c := config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
	)
	defer c.Close()

	if err := c.Load(); err != nil {
		panic(err)
	}

	if err := c.Scan(&bc); err != nil {
		panic(err)
	}
}

func main() {

	fmt.Println("start...")

	// CheckReadExcell()
	// ctx := context.Background()
	// fmt.Printf("bc: %+v\n", bc.Service)
	// 初始化 WpsSync
	// token, err := auth.NewAppAuthenticator(bc.Service).GetAccessToken(ctx)

	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Printf("token: %+v\n", token)

	//token := GetToken()
	//fmt.Println(token)
	// token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTM2MTU2MTgsImNvbXBfaWQiOiIxIiwiY2xpZW50X2lkIjoiY29tLmFjYy5hc3luYyIsInRrX3R5cGUiOiJhcHAiLCJzY29wZSI6Imtzby5hY2NvdW50c3luYy5zeW5jLGtzby5jb250YWN0LnJlYWQsa3NvLmNvbnRhY3QucmVhZHdyaXRlIiwiY29tcGFueV9pZCI6MSwiY2xpZW50X3ByaW5jaXBhbF9pZCI6IjczIiwiaXNfd3BzMzY1Ijp0cnVlfQ.ZOkiwnZ6f1uW45_sq5uT_ZW3dmA6yCXuKetMaUI7mCw"
	// 29290326581145992
	//CheckBatchGetDepartment(token)
	//CheckPostBatchUsersByExDepIds(token)
	// CheckPostBatchDepartmentsByExDepIds(token)
	//CheckGetUserByUserId(token)
	//CheckBatchPostUsers(token)

	// CheckPostBatchDepartmentsByExDepIds(token)
	// CheckUserLeaveOrg()
	// CheckPostBatchUsersByExDepIds(token)
	// CheckGetDepartmentRoot(token)
	//CheckBatchGetDepartment(token)
	// 033014104332101118010 test
	// CheckCallEcisaccountsync(token)
	// CheckGetDingtalkUserDetail()
	// "29290326581145992"
	// 03301410433273270
	//users, err := FindWpsUser(context.Background(), []string{"29290326581145992"})

	// if err != nil {
	// 	panic(err)
	// }
	// for _, u := range users {
	// 	fmt.Printf("FindWpsUser user: %v\n", *u)
	// }
	//CheckGetCompAllUsers()
	// CheckGetCompAllDepts()
	// FindAndDeleteUser()
	//FindAndDeleteDept("存在应用授权")
	// ctx := context.Background()
	// appAccessToken, err := auth.NewAppAuthenticator(bc.Service).GetAccessToken(ctx)
	// if err != nil {
	// 	panic(err)
	// }
	//fmt.Printf("appAccessToken: %s\n", appAccessToken.AccessToken)
	// CheckPostCreateUser(appAccessToken.AccessToken)
	//CheckDeleteDept(appAccessToken.AccessToken)
	// authApp := auth.NewAppAuthenticator(bc.Service)

	// authCache := auth.NewAppCacheAuthenticator(authApp)

	// AesEncryptGcmByKey

	// mobile, err := cipherutil.DecryptValueWithEnvSalt("HyyjnqUeVqHoid9cprHMoPgkOAVu8farJigGpvOi+xm0aLO2ZytG")
	// fmt.Printf("mobile: %s, err:%v\n", mobile, err)
	// CheckGetCompAllUsers(appAccessToken.AccessToken)
	// CheckInternalGateWay(appAccessToken.AccessToken)

	// 81
	// CheckGetUsersSearch(appAccessToken.AccessToken)
	// authDingtalk := auth.NewDingTalkAuthenticator(bc.Service)
	// authCache := auth.NewDingtalkCacheAuthenticator(authDingtalk)

	// // authCache := auth.NewDingtalkCacheAuthenticator(authDingtalk, auth.WithKey[*auth.DingtalkCacheConfig]("custom_key"))

	// token, err := authCache.GetAccessToken(ctx)
	// if err != nil {
	// 	panic(err)
	// }
	// dingtalkRepo := dingtalk.NewDingTalkRepo(bc.Service.Auth.Dingtalk, authCache, log.GetLogger())

	// depts, err := dingtalkRepo.FetchDepartments(ctx, token.AccessToken)
	// log.Infof("CreateSyncAccount.FetchDepartments: depts: %+v, err: %v", depts, err)
	// if err != nil {
	// 	panic(err)
	// }
	// for _, dept := range depts {
	// 	log.Infof("biz.CreateSyncAccount: dept: %+v", dept)
	// }
	// var deptIds []int64
	// for _, dept := range depts {
	// 	deptIds = append(deptIds, dept.DeptID)
	// }

	// deptUsers, err := dingtalkRepo.FetchDepartmentUsers(ctx, token.AccessToken, deptIds)

	// log.Infof("CreateSyncAccount.FetchDepartmentUsers deptUsers: %v, err: %v", deptUsers, err)
	// for _, deptUser := range deptUsers {
	// 	log.Infof("biz.CreateSyncAccount: deptUser: %+v", deptUser)
	// }
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("success")
	// for i := 1; i <= 3; i++ {
	// 	token, err := authCache.GetAccessToken(ctx)
	// 	fmt.Printf("GetAccessToken i:%d, token:%s, err:%v", i, token.AccessToken, err)
	// 	// return
	// }

}

// func CheckInternalGateWay(token string) {
// 	wpsClient := wps.NewWps(bc.Service, log.GetLogger())
// 	// 批量删除用户
// 	res, err := wpsClient.GetObjUploadUrl(context.Background(), token, wps.GetObjUploadUrlRequest{})

// 	fmt.Printf("res:%+v, err:%+v", res, err)

// }
// func CheckPostBatchDeleteUser(token string, userids []string) {
// 	wpsClient := wps.NewWps(bc.Service, log.GetLogger())
// 	// 批量删除用户
// 	delUserRes, err := wpsClient.PostBatchDeleteUser(context.Background(), token, wps.PostBatchDeleteUserRequest{
// 		UserIDs: userids,
// 	})
// 	fmt.Printf("delUserRes:%v, err:%v", delUserRes, err)
// }
// func CheckPostCreateUser(appToken string) {
// 	wpsClient := wps.NewWps(bc.Service, log.GetLogger())
// 	createUserRes, err := wpsClient.PostCreateUser(context.Background(), appToken, wps.PostCreateUserRequest{
// 		ExUserID:  "test01_user",
// 		Email:     "test01@163.com",
// 		UserName:  "test01",
// 		LoginName: "13888888888",
// 		Phone:     "13888888888",
// 		DeptIDs:   []string{"1"},
// 		Source:    "sync",
// 		WorkPlace: "bj",
// 	})

// 	fmt.Printf("createUserRes:%v, err:%v", createUserRes, err)
// }

// // 创建的也存在授权
// func CheckPostCreateDept(appToken string) {
// 	wpsClient := wps.NewWps(bc.Service, log.GetLogger())
// 	parentID := "1"
// 	createDeptRes, err := wpsClient.PostCreateDept(context.Background(), appToken, wps.PostCreateDeptRequest{
// 		ExDeptID: "test01",
// 		Name:     "test01_dep",
// 		ParentID: parentID,
// 		Source:   "sync",
// 		Order:    99,
// 	})
// 	fmt.Printf("createDeptRes:%v, err:%v", createDeptRes, err)
// }
// func CheckGetUsersSearch(appToken string) {

// 	wpsClient := wps.NewWps(bc.Service, log.GetLogger())
// 	// 批量删除部门
// 	// := []string{"6", "4"}
// 	// deleDeptRes, err := wpsClient.GetUsersSearch(context.Background(), appToken, wps.GetUsersSearchRequest{
// 	// 	Keyword:  "18910953345",
// 	// 	PageSize: 10,
// 	// 	//Status:                   []string{"active", "notactive", "disabled"},
// 	// 	//SearchSource:             []string{"company_user", "external_contact", "enterprise_partner"},
// 	// 	//SearchFieldConfigEnabled: false,
// 	// })
// 	// fmt.Printf("deleDeptRes:%v, err:%v", deleDeptRes, err)

// 	// 获取通讯录权限
// 	contactPermissionRes, err := wpsClient.GetContactPermission(context.Background(), appToken, wps.GetContactPermissionRequest{
// 		Scopes: []string{"org"},
// 	})
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Printf("contactPermissionRes:%v, err:%v", contactPermissionRes, err)
// }

// func CheckDeleteDept(appToken string) {

// 	wpsClient := wps.NewWps(bc.Service, log.GetLogger())
// 	// 批量删除部门
// 	needDeldept := []string{"6000", "4000"}
// 	deleDeptRes, err := wpsClient.PostBatchDeleteDept(context.Background(), appToken, wps.PostBatchDeleteDeptRequest{
// 		DeptIDs: needDeldept,
// 	})
// 	fmt.Printf("deleDeptRes:%v, err:%v", deleDeptRes, err)
// }
// func FindAndDeleteDept() {
// 	ctx := context.Background()
// 	appAccessToken, err := auth.NewAppAuthenticator(bc.Service).GetAccessToken(ctx)
// 	if err != nil {
// 		panic(err)
// 	}
// 	wpsClient := wps.NewWps(bc.Service, log.GetLogger())
// 	//查询根部门
// 	rootDept, err := wpsClient.GetDepartmentRoot(ctx, appAccessToken.AccessToken, wps.GetDepartmentRootRequest{})
// 	if err != nil {
// 		panic(err)
// 	}
// 	log.Infof("rootDept: %v", rootDept)

// 	// 2. 查询部门下的子部门(要递归)
// 	allDepts, err := wpsClient.GetDeptChildren(ctx, appAccessToken.AccessToken, wps.GetDeptChildrenRequest{
// 		DeptID:    rootDept.Data.ID,
// 		Recursive: true,
// 		PageSize:  50,
// 		WithTotal: true,
// 	})
// 	if err != nil {
// 		panic(err)
// 	}
// 	log.Infof("children: %v", allDepts)

// 	//删除部门除了根部门
// 	var alldeptes []string
// 	for _, dept := range allDepts.Data.Items {
// 		if dept.ID == rootDept.Data.ID {
// 			continue
// 		}
// 		alldeptes = append(alldeptes, dept.ID)
// 	}

// 	needDeldept := alldeptes[:2]
// 	fmt.Printf("deletedept: %v", needDeldept)
// 	// 批量删除部门
// 	deleDeptRes, err := wpsClient.PostBatchDeleteDept(ctx, appAccessToken.AccessToken, wps.PostBatchDeleteDeptRequest{
// 		DeptIDs: needDeldept,
// 	})
// 	fmt.Printf("deleDeptRes:%v, err:%v", deleDeptRes, err)
// }
// func FindAndDeleteUser() {
// 	ctx := context.Background()
// 	appAccessToken, err := auth.NewAppAuthenticator(bc.Service).GetAccessToken(ctx)
// 	if err != nil {
// 		panic(err)
// 	}
// 	wpsClient := wps.NewWps(bc.Service, log.GetLogger())

// 	users, err := wpsClient.GetCompAllUsers(ctx, appAccessToken.AccessToken, wps.GetCompAllUsersRequest{
// 		Recursive: true,
// 		PageSize:  50,
// 		WithTotal: true,
// 		Status:    []string{"active", "notactive", "disabled"},
// 	})
// 	if err != nil {
// 		panic(err)
// 	}

// 	var deleteUser *dingtalk.DingtalkDeptUser
// 	for _, user := range users.Data.Items {
// 		log.Infof("user: %+v", user)
// 		if user.Phone == "18910953345" {
// 			deleteUser = &dingtalk.DingtalkDeptUser{
// 				Userid: user.ExUserID,
// 				Mobile: user.Phone,
// 				Name:   user.UserName,
// 				Email:  user.Email,
// 			}
// 		}
// 		for _, dep := range user.Depts {
// 			depId, _ := strconv.ParseInt(dep.DeptID, 10, 64)
// 			deleteUser.DeptIDList = append(deleteUser.DeptIDList, depId)
// 		}
// 	}
// 	syncDB, err := data.NewMysqlSyncDB(bc.Data, log.GetLogger())
// 	if err != nil {
// 		panic(err)
// 		//return nil, nil, err
// 	}
// 	mainDB, err := data.NewMysqlDB(bc.Data, log.GetLogger())
// 	if err != nil {
// 		panic(err)
// 		//return nil, nil, err
// 	}
// 	client, err := data.NewRedisClient(bc.Data, log.GetLogger())
// 	if err != nil {
// 		panic(err)
// 		//return nil, nil, err
// 	}
// 	dataData, _, err := data.NewData(syncDB, mainDB, client, log.GetLogger())
// 	if err != nil {
// 		panic(err)
// 		// return nil, nil, err
// 	}

// 	accounterRepo := data.NewAccounterRepo(bc.Service, dataData, log.GetLogger())
// 	err = accounterRepo.SaveIncrementUsers(ctx, nil, []*dingtalk.DingtalkDeptUser{deleteUser}, nil)
// 	if err != nil {
// 		panic(err)
// 		//return err
// 	}

// 	wpsSync := wps.NewWpsSync(bc.Service, log.GetLogger())
// 	res, err := wpsSync.PostEcisaccountsyncIncrement(ctx, appAccessToken.AccessToken, &wps.EcisaccountsyncIncrementRequest{
// 		ThirdCompanyId: "1",
// 	})
// 	fmt.Printf("res:%v, err:%v", res, err)
// 	return
// 	// var alluserids []string
// 	// for _, u := range users.Data.Items {
// 	// 	fmt.Printf("FindWpsUser user: %v\n", u)
// 	// 	if u.ID == "1" {
// 	// 		continue
// 	// 	}
// 	// 	alluserids = append(alluserids, u.ID)
// 	// }
// 	// deleteuser := alluserids[:2]
// 	// fmt.Printf("deleteuser: %v", deleteuser)
// 	// // 存在授权问题
// 	// delRes, err := wpsClient.PostBatchDeleteUser(ctx, appAccessToken.AccessToken, wps.PostBatchDeleteUserRequest{
// 	// 	UserIDs: deleteuser,
// 	// })
// 	// fmt.Printf("delRes:%v, err:%v", delRes, err)

// }
// func CheckGetCompAllDepts() {
// 	appAccessToken, err := auth.NewAppAuthenticator(bc.Service).GetAccessToken(context.Background())
// 	if err != nil {
// 		panic(err)
// 	}
// 	wpsClient := wps.NewWps(bc.Service, log.GetLogger())
// 	rootDept, err := wpsClient.GetDepartmentRoot(context.Background(), appAccessToken.AccessToken, wps.GetDepartmentRootRequest{})
// 	if err != nil {
// 		panic(err)
// 	}
// 	log.Infof("rootDept: %v", rootDept)

// 	allDepts, err := wpsClient.GetDeptChildren(context.Background(), appAccessToken.AccessToken, wps.GetDeptChildrenRequest{
// 		DeptID:    rootDept.Data.ID,
// 		Recursive: true,
// 		PageSize:  50,
// 		WithTotal: true,
// 	})
// 	if err != nil {
// 		panic(err)
// 	}

// 	for _, d := range allDepts.Data.Items {
// 		fmt.Printf("CheckGetCompAllDepts dept: %v\n", d)
// 	}
// }
// func CheckGetCompAllUsers(token string) {

// 	wpsClient := wps.NewWps(bc.Service, log.GetLogger())

// 	users, err := wpsClient.GetCompAllUsers(context.Background(), token, wps.GetCompAllUsersRequest{
// 		Recursive: true,
// 		PageSize:  50,
// 		WithTotal: true,
// 		Status:    []string{"active", "notactive", "disabled"},
// 	})
// 	if err != nil {
// 		panic(err)
// 	}
// 	for _, u := range users.Data.Items {
// 		fmt.Printf("FindWpsUser user: %+v\n", u)
// 	}

// }
// func FindWpsUser(ctx context.Context, userids []string) ([]*dingtalk.DingtalkDeptUser, error) {
// 	fmt.Printf("FindWpsUser userids: %v\n", userids)
// 	var users []*dingtalk.DingtalkDeptUser

// 	appAccessToken, err := auth.NewAppAuthenticator(bc.Service).GetAccessToken(context.Background())
// 	if err != nil {
// 		return nil, err
// 	}
// 	wpsClient := wps.NewWps(bc.Service, log.GetLogger())

// 	for _, userId := range userids {
// 		wpsUserInfo, err := wpsClient.PostBatchUsersByExDepIds(ctx, appAccessToken.AccessToken, wps.PostBatchUsersByExDepIdsRequest{
// 			ExUserIDs: []string{userId},
// 			Status:    []string{wps.UserStatusActive, wps.UserStatusNoActive, wps.UserStatusDisabled},
// 		})
// 		fmt.Printf("FindWpsUser wpsUserInfo: %v, err: %v\n", wpsUserInfo, err)
// 		if err != nil {
// 			return nil, err
// 		}

// 		if len(wpsUserInfo.Data.Items) == 1 {
// 			fmt.Printf(">>>>>>>>>>>>>wpsUserInfo: %v\n", wpsUserInfo)
// 			wpsUserid := wpsUserInfo.Data.Items[0].ID

// 			wpsDeptInfo, err := wpsClient.GetUserDeptsByUserId(ctx, appAccessToken.AccessToken, wps.GetUserDeptsByUserIdRequest{
// 				UserID: wpsUserid,
// 			})
// 			if err != nil {
// 				return nil, err
// 			}
// 			if len(wpsDeptInfo.Data.Items) > 0 {
// 				for _, item := range wpsDeptInfo.Data.Items {

// 					//if _, ok := relationsMap[wpsUserid+item.ID]; !ok {
// 					user := &dingtalk.DingtalkDeptUser{
// 						Userid: userId,
// 					}
// 					deptId, err := strconv.ParseInt(item.ExDeptID, 10, 64)
// 					if err != nil {
// 						return nil, err
// 					}
// 					user.DeptIDList = append(user.DeptIDList, deptId)
// 					users = append(users, user)
// 					//}
// 				}
// 			}

// 		}
// 	}
// 	return users, nil
// }

// func CheckGetDingtalkUserDetail() {
// 	confService := bc.GetService()
// 	auth := auth.NewDingtalkCacheAuthenticator(auth.NewDingTalkAuthenticator(confService))
// 	dingtalkRepo := dingtalk.NewDingTalkRepo(confService.Auth.Dingtalk, auth, log.GetLogger())
// 	ctx := context.Background()
// 	// accessToken, err := dingtalkRepo.GetAccessToken(ctx, "code")
// 	// log.Infof("UserAddOrg.GetAccessToken accessToken: %v, err: %v", accessToken, err)
// 	// if err != nil {
// 	// 	panic(err)
// 	// }
// 	accessToken := "2a99752c4fd9317f81d4a20c0f1d7c7e"
// 	//user, err := dingtalkRepo.FetchUserDetail(ctx, accessToken, []string{"03301410433273270"})

// 	depts, _ := dingtalkRepo.FetchDeptDetails(ctx, accessToken, []int64{1002216804})
// 	fmt.Println()
// 	for i, dept := range depts {
// 		fmt.Printf("部门 %d: %+v\n", i, *dept)
// 	}
// }
// func CheckReadExcell() {
// 	f, err := excelize.OpenFile("Book1.xlsx")
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	defer func() {
// 		// Close the spreadsheet.
// 		if err := f.Close(); err != nil {
// 			fmt.Println(err)
// 		}
// 	}()
// 	// Get value from cell by given worksheet name and cell reference.
// 	cell, err := f.GetCellValue("Sheet1", "B2")
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	fmt.Println(cell)
// 	// Get all the rows in the Sheet1.
// 	rows, err := f.GetRows("Sheet1")
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	for _, row := range rows {
// 		for _, colCell := range row {
// 			fmt.Print(colCell, "\t")
// 		}
// 		fmt.Println()
// 	}
// }
// func CheckGetUserByUserId(token string) {
// 	ctx := context.Background()
// 	wpsClient := wps.NewWps(bc.Service, log.GetLogger())

// 	// 调用 Wps 接口
// 	res, err := wpsClient.GetUserByUserId(ctx, token, wps.GetUserByUserIdRequest{
// 		UserID: "81",
// 	})

// 	fmt.Printf("CheckGetUserByUserId res: %+v, err:%+v\n", res, err)
// }
// func CheckUserLeaveOrg() {

// 	ctx := context.Background()
// 	nancalDB, _ := data.NewMysqlDB(bc.Data, log.GetLogger())

// 	syncDB, _ := data.NewMysqlSyncDB(bc.Data, log.GetLogger())

// 	client, _ := data.NewRedisClient(bc.Data, log.GetLogger())

// 	dataData, _, err := data.NewData(syncDB, nancalDB, client, log.GetLogger())
// 	if err != nil {
// 		panic(err)
// 	}
// 	confService := bc.GetService()
// 	accounterRepo := data.NewAccounterRepo(confService, dataData, log.GetLogger())
// 	authDingtalk := auth.NewDingtalkCacheAuthenticator(auth.NewDingTalkAuthenticator(confService))
// 	dingtalkDingtalk := dingtalk.NewDingTalkRepo(service_Auth_Dingtalk, authDingtalk, log.GetLogger())
// 	authenticator := auth.NewAppAuthenticator(confService)
// 	wpsSync := wps.NewWpsSync(confService, log.GetLogger())
// 	wpsWps := wps.NewWps(confService, log.GetLogger())
// 	service_Business := conf.ProvideBusinessConfig(confService)
// 	accounterIncreUsecase := biz.NewAccounterIncreUsecase(accounterRepo, dingtalkDingtalk, authenticator, wpsSync, wpsWps, service_Business, log.GetLogger())

// 	// deptId
// 	// org_dept_create
// 	// org_dept_modify
// 	// org_dept_remove
// 	// userId
// 	// user_add_org
// 	// user_modify_org
// 	// user_leave_org
// 	// map[string]interface{"userId": []string{"033014104332101118010"}
// 	event := &clientV2.GenericOpenDingTalkEvent{
// 		EventId:           "111",
// 		EventBornTime:     "111",
// 		EventCorpId:       "111",
// 		EventType:         "user_leave_org",
// 		EventUnifiedAppId: "111",
// 		Data:              make(map[string]interface{}, 0),
// 	}
// 	event.Data["userId"] = []string{"033014104332101118010"}

// 	err = accounterIncreUsecase.UserLeaveOrg(ctx, event)

// 	fmt.Printf("err:%+v", err)
// }
// func CheckUserAddOrg(ctx context.Context, event *clientV2.GenericOpenDingTalkEvent) {

// }

// func GetToken() string {
// 	// fmt.Printf("bc.Service: %v", bc.Service)
// 	token, err := auth.NewAppAuthenticator(bc.Service).GetAccessToken(context.Background())

// 	if err != nil {
// 		panic(err)
// 	}

// 	//fmt.Printf("token: %+v\n", token)
// 	return token.AccessToken
// }
// func CheckBatchPostUsers(token string) {
// 	ctx := context.Background()
// 	wpsClient := wps.NewWps(bc.Service, log.GetLogger())

// 	// 调用 Wps 接口
// 	res, err := wpsClient.BatchPostUsers(ctx, token, wps.BatchPostUsersRequest{
// 		UserIDs:  []string{"81", "2"},
// 		Status:   []string{wps.UserStatusActive, wps.UserStatusNoActive, wps.UserStatusDisabled},
// 		WithDept: true,
// 	})

// 	fmt.Printf("CheckBatchPostUsers res: %+v, err:%+v\n", res, err)
// }
// func CheckPostBatchUsersByExDepIds(token string) {
// 	ctx := context.Background()
// 	wpsClient := wps.NewWps(bc.Service, log.GetLogger())

// 	// 调用 Wps 接口
// 	res, err := wpsClient.PostBatchUsersByExDepIds(ctx, token, wps.PostBatchUsersByExDepIdsRequest{
// 		ExUserIDs: []string{"033014104332101118010", "18910953345"},
// 		Status:    []string{wps.UserStatusActive, wps.UserStatusNoActive, wps.UserStatusDisabled},
// 	})

// 	fmt.Printf("CheckPostBatchUsersByExDepIds res: %+v, err:%+v\n", res, err)
// }

// func CheckPostBatchDepartmentsByExDepIds(token string) {
// 	ctx := context.Background()
// 	wpsClient := wps.NewWps(bc.Service, log.GetLogger())

// 	// 调用 Wps 接口
// 	res, err := wpsClient.PostBatchDepartmentsByExDepIds(ctx, token, wps.PostBatchDepartmentsByExDepIdsRequest{
// 		ExDeptIDs: []string{"1002216804"},
// 	})

// 	fmt.Printf("CheckPostBatchDepartmentsByExDepIds res: %+v, err:%+v\n", res, err)
// }
// func CheckGetDepartmentRoot(token string) {
// 	ctx := context.Background()
// 	wpsClient := wps.NewWps(bc.Service, log.GetLogger())

// 	// 调用 Wps 接口
// 	res, err := wpsClient.GetDepartmentRoot(ctx, token, wps.GetDepartmentRootRequest{})

// 	fmt.Printf("CheckGetDepartmentRoot res: %+v, err:%+v\n", res, err)

// }
// func CheckBatchGetDepartment(token string) {
// 	ctx := context.Background()
// 	wpsClient := wps.NewWps(bc.Service, log.GetLogger())

// 	// 调用 Wps 接口
// 	res, err := wpsClient.BatchPostDepartments(ctx, token, wps.BatchPostDepartmentsRequest{
// 		DeptIDs: []string{"201", "1", "33"},
// 	})
// 	fmt.Printf("CheckBatchGetDepartment res: %+v, err:%+v\n", res, err)

// }
// func CheckPostEcisaccountsync(token string) {
// 	ctx := context.Background()
// 	wpsSync := wps.NewWpsSync(bc.Service, log.GetLogger())
// 	res, err := wpsSync.PostEcisaccountsyncIncrement(ctx, token, &wps.EcisaccountsyncIncrementRequest{
// 		ThirdCompanyId: "1",
// 	})

// 	fmt.Printf("PostEcisaccountsyncIncrement res: %+v, err:%+v\n", res, err)
// 	res1, err := wpsSync.PostEcisaccountsyncAll(ctx, token, &wps.EcisaccountsyncAllRequest{
// 		ThirdCompanyId: "1",
// 		TaskId:         time.Now().Add(time.Duration(1) * time.Second).Format("20060102150405"),
// 	})

// 	fmt.Printf("PostEcisaccountsyncAll res: %+v, err:%+v\n", res1, err)
// }
