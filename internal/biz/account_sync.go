package biz

import (
	"context"
	"fmt"
	v1 "nancalacc/api/account/v1"
	"nancalacc/internal/data/models"
	"nancalacc/internal/wps"

	"github.com/go-kratos/kratos/v2/log"
)

func (uc *AccounterUsecase) CreateSyncTask(ctx context.Context, req *v1.CreateSyncAccountRequest) (*v1.CreateSyncAccountReply, error) {

	wpsAllDept, _ := uc.GetAllWpsDept(ctx)

	fmt.Printf("wpsAllDept: %v", wpsAllDept)
	wpsAlluser, _ := uc.GetAllWpsUser(ctx)

	fmt.Printf("wpsAlluser: %v", wpsAlluser)
	allTableDept, _ := uc.repo.BatchGetDepts(ctx, "taskName")
	allTableUser, _ := uc.repo.BatchGetUsers(ctx, "taskName")
	allTableDeptuser, _ := uc.repo.BatchGetDeptUsers(ctx, "taskName")
	allTableDeptMap := make(map[string]*models.TbLasDepartment, len(allTableDept))
	allTableUserMap := make(map[string]*models.TbLasUser, len(allTableUser))
	allTableDeptuserMap := make(map[string][]*models.TbLasDepartmentUser, 0)
	for _, item := range allTableDept {
		allTableDeptMap[item.Did] = item
	}
	for _, item := range allTableUser {
		allTableUserMap[item.Uid] = item
	}
	for _, item := range allTableDeptuser {
		allTableDeptuserMap[item.Uid] = append(allTableDeptuserMap[item.Uid], item)
	}

	var addp, updp, delp []*models.TbLasDepartment
	for _, item := range wpsAllDept {
		if v, ok := allTableDeptMap[item.ExDeptID]; ok {
			updp = append(updp, v)
			delete(allTableDeptMap, item.ExDeptID)
		} else {
			addp = append(addp, v)
		}
	}
	if len(allTableDeptMap) > 0 {
		for _, v := range allTableDeptMap {
			delp = append(delp, v)
		}
	}

	fmt.Printf("addp: %v, updp: %v, del: %v", addp, updp, delp)

	var addu, updu, delu []wps.UserItem
	for _, item := range wpsAlluser {
		if v, ok := allTableUserMap[item.ExUserID]; ok {
			//这里可以将tabuser里 的属性拷贝到WPS user
			item.UserName = v.NickName
			updu = append(updu, item)
			delete(allTableUserMap, item.ExUserID)
		} else {
			//这个也调用makeWpsUser生成的wps user

			delu = append(addu, item)
		}
	}
	if len(allTableUserMap) > 0 {
		for _, v := range allTableUserMap {
			addu = append(delu, makeWpsUser(*v))
		}
	}

	fmt.Printf("addu: %v, updu: %v, del: %v", addu, updu, delu)

	for _, item := range addu {
		if relations, ok := allTableDeptuserMap[item.ExUserID]; ok {
			var depts []wps.Dept
			for _, relation := range relations {
				depts = append(depts, wps.Dept{DeptID: relation.Did, Name: allTableDeptMap[relation.Did].Name})
			}
			item.Depts = depts
		}

	}

	for _, item := range updu {
		if relations, ok := allTableDeptuserMap[item.ExUserID]; ok {
			var depts []wps.Dept
			for _, relation := range relations {
				depts = append(depts, wps.Dept{DeptID: relation.Did, Name: allTableDeptMap[relation.Did].Name})
			}
			item.Depts = depts
		}

	}
	fmt.Printf("addu: %v, updu: %v, del: %v", addu, updu, delu)

	return nil, nil
}
func makeWpsUser(tbuser models.TbLasUser) wps.UserItem {
	user := wps.UserItem{
		ExUserID: tbuser.Uid,
	}
	return user
}
func (uc *AccounterUsecase) GetAllWpsUser(ctx context.Context) (users []wps.UserItem, err error) {
	appAccessToken, err := uc.appAuth.GetAccessToken(ctx)
	if err != nil {
		return
	}
	for {
		// 查询企业下所有用户
		alluser, err := uc.wps.GetCompAllUsers(ctx, appAccessToken.AccessToken, wps.GetCompAllUsersRequest{
			Recursive: true,
			PageSize:  50,
			WithTotal: true,
			Status:    []string{"active", "notactive", "disabled"},
		})
		if err != nil {
			break
		}
		for _, u := range alluser.Data.Items {
			if u.ID == "1" {
				continue
			}
			if u.Source == "sync" {
				users = append(users, u)
			}

		}

	}
	return

}

func (uc *AccounterUsecase) GetAllWpsDept(ctx context.Context) (depts []wps.DeptItem, err error) {
	appAccessToken, err := uc.appAuth.GetAccessToken(ctx)
	if err != nil {
		return nil, err
	}
	//查询根部门
	rootDept, err := uc.wps.GetDepartmentRoot(ctx, appAccessToken.AccessToken, wps.GetDepartmentRootRequest{})
	if err != nil {
		return nil, err
	}
	log.Infof("rootDept: %v", rootDept)
	queue := []string{rootDept.Data.ID}
	for len(queue) > 0 {
		var pageToken string
		currDepId := queue[0]

		// 2. 查询部门下的子部门(要递归)
		deptInfo, err := uc.wps.GetDeptChildren(ctx, appAccessToken.AccessToken, wps.GetDeptChildrenRequest{
			DeptID:    currDepId,
			Recursive: true,
			PageSize:  50,
			WithTotal: true,
			PageToken: pageToken,
		})
		if err != nil {
			return nil, err
		}

		if deptInfo.Data.Total == 0 {
			queue = queue[1:]
			continue
		} else if deptInfo.Data.Total > 0 && deptInfo.Data.Total < 50 {
			queue = queue[1:]
			for _, item := range deptInfo.Data.Items {
				depts = append(depts, item)
				queue = append(queue, item.ID)
			}
			continue
		} else if deptInfo.Data.NextPageToken != "" {
			pageToken = deptInfo.Data.NextPageToken
		} else {
			queue = queue[1:]
			continue
		}

	}

	return depts, nil

}

// 这个方法是把全量数据执插入表后，可以自已调用原生API去同步
func (uc *AccounterUsecase) StartSync(ctx context.Context, taskId, filename string) (err error) {

	appAccessToken, err := uc.appAuth.GetAccessToken(ctx)
	if err != nil {
		return err
	}
	// 查询企业下所有用户
	alluser, err := uc.wps.GetCompAllUsers(ctx, appAccessToken.AccessToken, wps.GetCompAllUsersRequest{
		Recursive: true,
		PageSize:  50,
		WithTotal: true,
		Status:    []string{"active", "notactive", "disabled"},
	})
	if err != nil {
		return err
	}
	log.Infof("alluser: %v", alluser)
	// 批量删除用户(除了admin)
	var allUsers []string
	for _, user := range alluser.Data.Items {
		if user.ID == "1" {
			continue
		}
		allUsers = append(allUsers, user.ID)
	}
	////////////////////// 存在授权问题
	deluserRes, err := uc.wps.PostBatchDeleteUser(ctx, appAccessToken.AccessToken, wps.PostBatchDeleteUserRequest{
		UserIDs: allUsers,
	})
	if err != nil {
		return err
	}
	log.Infof("deluserRes: %v", deluserRes)

	//查询根部门
	rootDept, err := uc.wps.GetDepartmentRoot(ctx, appAccessToken.AccessToken, wps.GetDepartmentRootRequest{})
	if err != nil {
		return err
	}
	log.Infof("rootDept: %v", rootDept)

	// 2. 查询部门下的子部门(要递归)
	allDepts, err := uc.wps.GetDeptChildren(ctx, appAccessToken.AccessToken, wps.GetDeptChildrenRequest{
		DeptID: rootDept.Data.ID,
	})
	if err != nil {
		return err
	}
	log.Infof("children: %v", allDepts)

	//删除部门除了根部门
	var alldeptes []string
	for _, dept := range allDepts.Data.Items {
		if dept.ID == rootDept.Data.ID {
			continue
		}
		alldeptes = append(alldeptes, dept.ID)
	}

	// 批量删除部门(接口有授权问题)
	uc.wps.PostBatchDeleteDept(ctx, appAccessToken.AccessToken, wps.PostBatchDeleteDeptRequest{
		DeptIDs: alldeptes,
	})

	// 1. 创建部门(接口有授权问题)
	uc.wps.PostCreateDept(ctx, appAccessToken.AccessToken, wps.PostCreateDeptRequest{
		ExDeptID: "test01",
		Name:     "test01_dep",
		ParentID: "1",
		Source:   "sync",
		Order:    99,
	})

	uc.wps.PostCreateUser(ctx, appAccessToken.AccessToken, wps.PostCreateUserRequest{
		ExUserID:  "test01_user",
		Email:     "test01@163.com",
		UserName:  "test01",
		LoginName: "13888888888",
		Phone:     "13888888888",
		DeptIDs:   []string{"1"},
		Source:    "sync",
		WorkPlace: "bj",
	})
	////////////////////////////////////////////////////////////////
	// 2. 创建部门存在的 > 更新部门
	uc.wps.PostUpdateDept(ctx, taskId, wps.PostUpdateDeptRequest{
		ExDeptID: "10000000000000000000000000000000",
		Name:     "测试部门",
	})

	// 3. 全量的的 减去 创建的， 再减去更新的， 就是删除的
	uc.wps.PostBatchDeleteDept(ctx, taskId, wps.PostBatchDeleteDeptRequest{
		DeptIDs: []string{"10000000000000000000000000000000"},
	})

	// 这里的用户包含所在部门
	uc.wps.PostCreateUser(ctx, taskId, wps.PostCreateUserRequest{
		UserName: "测试用户",
	})

	// 2. 更新用户
	uc.wps.PostUpdateUser(ctx, taskId, wps.PostUpdateUserRequest{
		ExUserID: "10000000000000000000000000000000",
		UserName: "测试部门",
	})

	// 3. 更新用户
	uc.wps.PostUpdateUser(ctx, taskId, wps.PostUpdateUserRequest{
		ExUserID: "10000000000000000000000000000000",
		UserName: "测试用户",
	})

	// 全量用户

	uc.wps.PostBatchUserByPage(ctx, taskId, wps.PostBatchUserByPageRequest{
		PageSize: 100,
		PageNum:  1,
	})

	// 5. 全量的的 减去 创建的， 再减去更新的， 就是删除的
	uc.wps.PostBatchDeleteUser(ctx, taskId, wps.PostBatchDeleteUserRequest{
		UserIDs: []string{"10000000000000000000000000000000"},
	})
	return nil
}
