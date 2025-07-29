package data

import (
	"nancalacc/internal/data/models"
	"time"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {

	if err := db.Migrator().DropTable(&models.TbLasDepartment{}, &models.TbLasUser{}, &models.TbLasDepartmentUser{}, &models.TbCompanyCfg{}, &models.Task{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&models.TbLasDepartment{}, &models.TbLasUser{}, &models.TbLasDepartmentUser{}, &models.TbCompanyCfg{}, &models.Task{}); err != nil {
		return err
	}
	db.Create(&models.TbCompanyCfg{
		ThirdCompanyId: "1",
		PlatformIds:    "1",
		CompanyId:      "1",
		Status:         0,
		Ctime:          time.Now(),
		Mtime:          time.Now(),
	})
	return nil
}
