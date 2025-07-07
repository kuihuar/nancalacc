package models

import (
	"database/sql"
	"time"
)

// 租户关系配置表
type TbCompanyCfg struct {
	Id             uint         `gorm:"primaryKey;autoIncrement;column:id;type:int unsigned;comment:主键id" json:"id"`
	ThirdCompanyId string       `gorm:"column:third_company_id;type:varchar(20);comment:三方租户id;NOT NULL" json:"third_company_id"`
	PlatformIds    string       `gorm:"column:platform_ids;type:varchar(100);comment:平台id, 用来区分多种数据源,多个用逗号分隔;NOT NULL" json:"platform_ids"`
	CompanyId      string       `gorm:"column:company_id;type:varchar(20);comment:云文档租户id;NOT NULL" json:"company_id"`
	Status         int          `gorm:"column:status;type:tinyint(4);default:1;comment:状态,0-禁用,1-启用;NOT NULL" json:"status"`
	Ctime          sql.NullTime `gorm:"column:ctime;type:timestamp;default:CURRENT_TIMESTAMP;comment:创建时间" json:"ctime"`
	Mtime          time.Time    `gorm:"column:mtime;type:timestamp;default:CURRENT_TIMESTAMP;comment:更新时间;NOT NULL" json:"mtime"`
}

func (m *TbCompanyCfg) TableName() string {
	return "tb_company_cfg"
}
