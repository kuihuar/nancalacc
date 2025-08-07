package main

import (
	"context"
	"flag"
	"fmt"
	"nancalacc/internal/auth"
	"nancalacc/internal/biz"
	"nancalacc/internal/conf"
	"nancalacc/internal/data"
	"nancalacc/internal/dingtalk"
	"nancalacc/internal/wps"
	"time"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/open-dingtalk/dingtalk-stream-sdk-go/clientV2"
	"github.com/xuri/excelize/v2"
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

	token := GetToken()
	fmt.Println(token)
	// token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTM2MTU2MTgsImNvbXBfaWQiOiIxIiwiY2xpZW50X2lkIjoiY29tLmFjYy5hc3luYyIsInRrX3R5cGUiOiJhcHAiLCJzY29wZSI6Imtzby5hY2NvdW50c3luYy5zeW5jLGtzby5jb250YWN0LnJlYWQsa3NvLmNvbnRhY3QucmVhZHdyaXRlIiwiY29tcGFueV9pZCI6MSwiY2xpZW50X3ByaW5jaXBhbF9pZCI6IjczIiwiaXNfd3BzMzY1Ijp0cnVlfQ.ZOkiwnZ6f1uW45_sq5uT_ZW3dmA6yCXuKetMaUI7mCw"

	//CheckPostBatchUsersByExDepIds(token)
	CheckPostBatchDepartmentsByExDepIds(token)
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

}
func CheckGetDingtalkUserDetail() {
	confService := bc.GetService()
	dingtalkRepo := dingtalk.NewDingTalkRepo(confService.Auth.Dingtalk, log.GetLogger())
	ctx := context.Background()
	// accessToken, err := dingtalkRepo.GetAccessToken(ctx, "code")
	// log.Infof("UserAddOrg.GetAccessToken accessToken: %v, err: %v", accessToken, err)
	// if err != nil {
	// 	panic(err)
	// }
	accessToken := "2a99752c4fd9317f81d4a20c0f1d7c7e"
	//user, err := dingtalkRepo.FetchUserDetail(ctx, accessToken, []string{"03301410433273270"})

	depts, _ := dingtalkRepo.FetchDeptDetails(ctx, accessToken, []int64{1002216804})
	fmt.Println()
	for i, dept := range depts {
		fmt.Printf("部门 %d: %+v\n", i, *dept)
	}
}
func CheckReadExcell() {
	f, err := excelize.OpenFile("Book1.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	// Get value from cell by given worksheet name and cell reference.
	cell, err := f.GetCellValue("Sheet1", "B2")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(cell)
	// Get all the rows in the Sheet1.
	rows, err := f.GetRows("Sheet1")
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, row := range rows {
		for _, colCell := range row {
			fmt.Print(colCell, "\t")
		}
		fmt.Println()
	}
}
func CheckGetUserByUserId(token string) {
	ctx := context.Background()
	wpsClient := wps.NewWps(bc.Service, log.GetLogger())

	// 调用 Wps 接口
	res, err := wpsClient.GetUserByUserId(ctx, token, wps.GetUserByUserIdRequest{
		UserID: "81",
	})

	fmt.Printf("CheckGetUserByUserId res: %+v, err:%+v\n", res, err)
}
func CheckUserLeaveOrg() {

	ctx := context.Background()
	nancalDB, _ := data.NewMysqlDB(bc.Data, log.GetLogger())

	syncDB, _ := data.NewMysqlSyncDB(bc.Data, log.GetLogger())

	client, _ := data.NewRedisClient(bc.Data, log.GetLogger())

	dataData, _, err := data.NewData(syncDB, nancalDB, client, log.GetLogger())
	if err != nil {
		panic(err)
	}
	confService := bc.GetService()
	accounterRepo := data.NewAccounterRepo(confService, dataData, log.GetLogger())
	service_Auth_Dingtalk := conf.ProvideDingtalkConfig(confService)
	dingtalkDingtalk := dingtalk.NewDingTalkRepo(service_Auth_Dingtalk, log.GetLogger())
	authenticator := auth.NewAppAuthenticator(confService)
	wpsSync := wps.NewWpsSync(confService, log.GetLogger())
	wpsWps := wps.NewWps(confService, log.GetLogger())
	service_Business := conf.ProvideBusinessConfig(confService)
	accounterIncreUsecase := biz.NewAccounterIncreUsecase(accounterRepo, dingtalkDingtalk, authenticator, wpsSync, wpsWps, service_Business, log.GetLogger())

	// deptId
	// org_dept_create
	// org_dept_modify
	// org_dept_remove
	// userId
	// user_add_org
	// user_modify_org
	// user_leave_org
	// map[string]interface{"userId": []string{"033014104332101118010"}
	event := &clientV2.GenericOpenDingTalkEvent{
		EventId:           "111",
		EventBornTime:     "111",
		EventCorpId:       "111",
		EventType:         "user_leave_org",
		EventUnifiedAppId: "111",
		Data:              make(map[string]interface{}, 0),
	}
	event.Data["userId"] = []string{"033014104332101118010"}

	err = accounterIncreUsecase.UserLeaveOrg(ctx, event)

	fmt.Printf("err:%+v", err)
}
func CheckUserAddOrg(ctx context.Context, event *clientV2.GenericOpenDingTalkEvent) {

}

func GetToken() string {
	// fmt.Printf("bc.Service: %v", bc.Service)
	token, err := auth.NewAppAuthenticator(bc.Service).GetAccessToken(context.Background())

	if err != nil {
		panic(err)
	}

	//fmt.Printf("token: %+v\n", token)
	return token.AccessToken
}
func CheckBatchPostUsers(token string) {
	ctx := context.Background()
	wpsClient := wps.NewWps(bc.Service, log.GetLogger())

	// 调用 Wps 接口
	res, err := wpsClient.BatchPostUsers(ctx, token, wps.BatchPostUsersRequest{
		UserIDs:  []string{"81", "2"},
		Status:   []string{wps.UserStatusActive, wps.UserStatusNoActive, wps.UserStatusDisabled},
		WithDept: true,
	})

	fmt.Printf("CheckBatchPostUsers res: %+v, err:%+v\n", res, err)
}
func CheckPostBatchUsersByExDepIds(token string) {
	ctx := context.Background()
	wpsClient := wps.NewWps(bc.Service, log.GetLogger())

	// 调用 Wps 接口
	res, err := wpsClient.PostBatchUsersByExDepIds(ctx, token, wps.PostBatchUsersByExDepIdsRequest{
		ExUserIDs: []string{"033014104332101118010", "18910953345"},
		Status:    []string{wps.UserStatusActive, wps.UserStatusNoActive, wps.UserStatusDisabled},
	})

	fmt.Printf("CheckPostBatchUsersByExDepIds res: %+v, err:%+v\n", res, err)
}

func CheckPostBatchDepartmentsByExDepIds(token string) {
	ctx := context.Background()
	wpsClient := wps.NewWps(bc.Service, log.GetLogger())

	// 调用 Wps 接口
	res, err := wpsClient.PostBatchDepartmentsByExDepIds(ctx, token, wps.PostBatchDepartmentsByExDepIdsRequest{
		ExDeptIDs: []string{"1002216804"},
	})

	fmt.Printf("CheckPostBatchDepartmentsByExDepIds res: %+v, err:%+v\n", res, err)
}
func CheckGetDepartmentRoot(token string) {
	ctx := context.Background()
	wpsClient := wps.NewWps(bc.Service, log.GetLogger())

	// 调用 Wps 接口
	res, err := wpsClient.GetDepartmentRoot(ctx, token, wps.GetDepartmentRootRequest{})

	fmt.Printf("CheckGetDepartmentRoot res: %+v, err:%+v\n", res, err)

}
func CheckBatchGetDepartment(token string) {
	ctx := context.Background()
	wpsClient := wps.NewWps(bc.Service, log.GetLogger())

	// 调用 Wps 接口
	res, err := wpsClient.BatchPostDepartments(ctx, token, wps.BatchPostDepartmentsRequest{
		DeptIDs: []string{"1005617108", "1006047132", "1"},
	})
	fmt.Printf("CheckBatchGetDepartment res: %+v, err:%+v\n", res, err)

}
func CheckPostEcisaccountsync(token string) {
	ctx := context.Background()
	wpsSync := wps.NewWpsSync(bc.Service, log.GetLogger())
	res, err := wpsSync.PostEcisaccountsyncIncrement(ctx, token, &wps.EcisaccountsyncIncrementRequest{
		ThirdCompanyId: "1",
	})

	fmt.Printf("PostEcisaccountsyncIncrement res: %+v, err:%+v\n", res, err)
	res1, err := wpsSync.PostEcisaccountsyncAll(ctx, token, &wps.EcisaccountsyncAllRequest{
		ThirdCompanyId: "1",
		TaskId:         time.Now().Add(time.Duration(1) * time.Second).Format("20060102150405"),
	})

	fmt.Printf("PostEcisaccountsyncAll res: %+v, err:%+v\n", res1, err)
}
