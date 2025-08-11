package wps

import (
	"context"
	"errors"
)

func ValidateWpsUser(ctx context.Context, user *UserItem) error {
	if user.ExUserID == "" {
		return errors.New("ex_user_id is empty")
	}
	return nil
}
