package models

import (
	"fmt"
	"nancalacc/internal/dingtalk"
	"nancalacc/internal/pkg/cipherutil"
	"strconv"
	"time"
)

func MakeLasUser(user *dingtalk.DingtalkDeptUser, thirdCompanyID, platformID, source, taskId string) *TbLasUser {
	//var account string
	// if user.Mobile != "" {
	// 	account = user.Mobile
	// } else {
	// 	account = user.Userid
	// }
	account := user.Mobile
	now := time.Now()
	mobile, err := cipherutil.EncryptUserInfo(user.Mobile, user.Userid)
	fmt.Printf("EncryptUserInfo userid:%s, mobile: %s mobile: %s, err: %v\n", user.Userid, user.Mobile, mobile, err)
	entity := &TbLasUser{
		TaskID:         taskId,
		ThirdCompanyID: thirdCompanyID,
		PlatformID:     platformID,
		Uid:            user.Userid,
		DefDid:         "-1",
		DefDidOrder:    0,
		Account:        account,
		NickName:       user.Name,
		Email:          user.Email,
		// Phone:         user.Mobile,
		Phone:            mobile,
		Title:            user.Title,
		Source:           source,
		Ctime:            now,
		Mtime:            now,
		CheckType:        1,
		EmploymentStatus: "active",
		EmploymentType:   "permanent",
	}
	if len(user.LeaderInDept) > 0 {
		//entity.Leader = user.LeaderInDept[0].Leader
	}
	return entity
}

func MakeTbLasDepartment(dep *dingtalk.DingtalkDept, thirdCompanyID, platformID, companyID, source, taskId string) *TbLasDepartment {

	now := time.Now()
	return &TbLasDepartment{
		Did:            strconv.FormatInt(dep.DeptID, 10),
		TaskID:         taskId,
		Name:           dep.Name,
		ThirdCompanyID: thirdCompanyID,
		PlatformID:     platformID,
		Pid:            strconv.FormatInt(dep.ParentID, 10),
		Order:          int(dep.Order),
		Source:         "sync",
		Ctime:          now,
		Mtime:          now,
		CheckType:      1,
	}
}

func MakeTbLasRootDepartment(thirdCompanyID, platformID, companyID, source, taskId string) *TbLasDepartment {

	now := time.Now()
	return &TbLasDepartment{
		Did:            "0",
		TaskID:         taskId,
		Name:           companyID,
		ThirdCompanyID: thirdCompanyID,
		PlatformID:     platformID,
		Pid:            "-1",
		Order:          0,
		Source:         source,
		Ctime:          now,
		Mtime:          now,
		CheckType:      1,
	}
}

func MakeTbLasDepartmentUser(relation *dingtalk.DingtalkDeptUserRelation, thirdCompanyID, platformID, companyID, source, taskId string) *TbLasDepartmentUser {

	return &TbLasDepartmentUser{
		Did:            relation.Did,
		TaskID:         taskId,
		ThirdCompanyID: thirdCompanyID,
		PlatformID:     platformID,
		Uid:            relation.Uid,
		Ctime:          time.Now(),
		Order:          int(relation.Order),
		CheckType:      1,
	}
}

func MakeLasUserIncrement(user *dingtalk.DingtalkDeptUser, thirdCompanyID, platformID, companyID, source, updateType string) *TbLasUserIncrement {

	// var account string
	// if user.Mobile != "" {
	// 	account = user.Mobile
	// } else {
	// 	account = user.Userid
	// }
	account := user.Mobile
	now := time.Now()

	mobile, err := cipherutil.EncryptUserInfo(user.Mobile, user.Userid)
	fmt.Printf("EncryptUserInfo userid:%s, mobile: %s mobile: %s, err: %v\n", user.Userid, user.Mobile, mobile, err)

	entity := &TbLasUserIncrement{
		ThirdCompanyID: thirdCompanyID,
		PlatformID:     platformID,
		Uid:            user.Userid,
		DefDid:         "-1",
		DefDidOrder:    0,
		Account:        account,
		NickName:       user.Name,
		Email:          user.Email,
		// Phone:            user.Mobile,
		Phone:            mobile,
		Title:            user.Title,
		Source:           source,
		Ctime:            now,
		Mtime:            now,
		EmploymentStatus: "active",
		EmploymentType:   "permanent",
		UpdateType:       updateType,
		SyncType:         "auto",
		SyncTime:         now,
		Status:           0,
	}
	if len(user.LeaderInDept) > 0 {
		//entity.Leader = user.LeaderInDept[0].Leader
	}
	return entity
}

func MakeDepartmentIncrement(dept *dingtalk.DingtalkDept, thirdCompanyID, platformID, companyID, source, updateType string) *TbLasDepartmentIncrement {

	now := time.Now()
	entity := &TbLasDepartmentIncrement{
		Did:            strconv.FormatInt(dept.DeptID, 10),
		Name:           dept.Name,
		ThirdCompanyID: thirdCompanyID,
		PlatformID:     platformID,
		//Pid:            strconv.FormatInt(dept.ParentID, 10),
		Order:      int32(dept.Order),
		Source:     "sync",
		Ctime:      now,
		Mtime:      now,
		UpdateType: updateType,
		SyncTime:   now,
		SyncType:   "auto",
		Status:     0,
	}
	if dept.ParentID != 0 {
		entity.Pid = strconv.FormatInt(dept.ParentID, 10)
	} else {
		entity.Pid = "-1"
	}
	return entity

}

func MmakeTbLasDepartmentUserIncrement(relation *dingtalk.DingtalkDeptUserRelation, thirdCompanyID, platformID, companyID, source, updateType string) *TbLasDepartmentUserIncrement {
	now := time.Now()
	entity := &TbLasDepartmentUserIncrement{
		Did:            relation.Did,
		ThirdCompanyID: thirdCompanyID,
		PlatformID:     platformID,
		Uid:            relation.Uid,
		Ctime:          now,
		Order:          1,
		UpdateType:     updateType,
		SyncType:       "auto",
		SyncTime:       now,
		Status:         0,
	}

	return entity

}
