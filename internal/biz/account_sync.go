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
	//查询根部门
	rootDept, err := uc.wps.GetDepartmentRoot(ctx, appAccessToken.AccessToken, wps.GetDepartmentRootRequest{})
	if err != nil {
		return err
	}
	log.Infof("rootDept: %v", rootDept)

	// 2. 查询部门下的子部门(要递归)
	children, err := uc.wps.GetDeptChildren(ctx, appAccessToken.AccessToken, wps.GetDeptChildrenRequest{
		DeptID: rootDept.Data.ID,
	})
	if err != nil {
		return err
	}
	log.Infof("children: %v", children)

	// 查询企业下所有用户

	// 1. 查询全量部门
	// uc.wps.PostBatchDeptByPage(ctx, taskId, wps.PostBatchDeptByPageRequest{
	// 	PageSize: 100,
	// 	PageNum:  1,
	// })

	// 1.批量删除用户
	// 2. 批量删除部门

	// 1. 创建部门
	uc.wps.PostCreateDept(ctx, taskId, wps.PostCreateDeptRequest{
		Name: "测试部门",
	})

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
