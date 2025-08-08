package biz

import (
	"context"
	"nancalacc/internal/wps"

	"github.com/go-kratos/kratos/v2/log"
)

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
