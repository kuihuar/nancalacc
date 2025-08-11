package dingtalk

import (
	"context"
	"errors"
)

func ValidateDingTalkUser(ctx context.Context, user *DingtalkDeptUser) error {
	if user.Name == "" {
		return errors.New("DingtalkDeptUser userid is empty")
	}
	if user.Mobile == "" {
		return errors.New("DingtalkDeptUser mobile is empty")
	}
	if len(user.DeptIDList) == 0 {
		return errors.New("DingtalkDeptUser deptid is empty")
	}
	return nil
}
