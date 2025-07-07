package models

import (
	"database/sql"
	"time"
)

// 全量部门用户关系表
type TbLasDepartmentUser struct {
	ID             uint          `gorm:"primaryKey;autoIncrement;column:id;type:int unsigned;comment:主键id"`
	TaskID         string        `gorm:"not null;column:task_id;type:varchar(20);comment:任务id"`
	ThirdCompanyID string        `gorm:"not null;column:third_company_id;type:varchar(20);comment:租户id"`
	PlatformID     string        `gorm:"not null;column:platform_id;type:varchar(60);comment:平台id"`
	Uid            string        `gorm:"not null;column:uid;type:varchar(255);comment:用户id"`
	Did            string        `gorm:"not null;column:did;type:varchar(255);comment:部门id"`
	Order          sql.NullInt32 `gorm:"column:order;type:int;comment:排序"`
	Main           int           `gorm:"column:main;type:int;default:0;comment:是否是主部门"`
	Ctime          time.Time     `gorm:"not null;column:ctime;type:timestamp;default:CURRENT_TIMESTAMP;comment:创建时间"`
	CheckType      int8          `gorm:"not null;column:check_type;type:tinyint;default:0;comment:勾选状态"`
}

func (TbLasDepartmentUser) TableName() string {
	return "tb_las_department_user"
}

// 增量部门用户关系表
type TbLasDepartmentUserIncrement struct {
	ID             uint           `gorm:"primaryKey;autoIncrement;column:id;type:int unsigned;comment:主键id"`
	ThirdCompanyID string         `gorm:"not null;column:third_company_id;type:varchar(20);comment:租户id"`
	PlatformID     string         `gorm:"not null;column:platform_id;type:varchar(60);comment:平台id"`
	Uid            string         `gorm:"not null;column:uid;type:varchar(255);comment:用户id"`
	Did            string         `gorm:"not null;column:did;type:varchar(255);comment:默认部门id"`
	Order          sql.NullInt32  `gorm:"column:order;type:int;comment:排序"`
	Main           int            `gorm:"column:main;type:int;default:0;comment:是否是主部门"`
	SyncType       string         `gorm:"column:sync_type;type:varchar(20);default:auto;comment:同步方式"`
	UpdateType     string         `gorm:"not null;column:update_type;type:varchar(20);comment:修改类型"`
	Status         int            `gorm:"column:status;type:int;default:0;comment:状态"`
	Msg            sql.NullString `gorm:"column:msg;type:varchar(2000);comment:错误详情"`
	Operator       string         `gorm:"column:operator;type:varchar(100);default:系统;comment:operator"`
	SyncTime       sql.NullTime   `gorm:"column:sync_time;type:timestamp;default:CURRENT_TIMESTAMP;comment:增量数据变动时间"`
	Ctime          time.Time      `gorm:"not null;column:ctime;type:timestamp;default:CURRENT_TIMESTAMP;comment:创建时间"`
	Mtime          time.Time      `gorm:"not null;column:mtime;type:timestamp;comment:修改时间"`
	Dids           sql.NullString `gorm:"column:dids;type:varchar(5000);comment:dids"`
}

func (TbLasDepartmentUserIncrement) TableName() string {
	return "tb_las_department_user_increment"
}
