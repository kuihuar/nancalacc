package models

import (
	"gorm.io/gorm"
)

// 数据库实体
type Account struct {
	gorm.Model
	ID        int64  `gorm:"column:id;primaryKey"`
	Username  string `gorm:"column:username;size:255;uniqueIndex"`
	Email     string `gorm:"column:email;size:255;uniqueIndex"`
	Phone     string `gorm:"column:phone;size:20"`
	Password  string `gorm:"column:password;size:255"`
	Status    int32  `gorm:"column:status;default:1"`
	CreatedAt int64  `gorm:"column:created_at"`
	UpdatedAt int64  `gorm:"column:updated_at"`
}

func (a *Account) TableName() string {
	return "nancal_account"
}
