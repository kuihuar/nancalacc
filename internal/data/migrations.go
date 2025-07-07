package data

import (
	"nancalacc/internal/data/models"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	// return db.AutoMigrate(&models.Account{}, &models.TbLasDepartment{}, &models.TbLasUser{}, &models.TbCompanyCfg{})
	return db.AutoMigrate(&models.TbLasDepartment{}, &models.TbLasUser{}, &models.TbLasDepartmentUser{}, &models.TbCompanyCfg{})

}
