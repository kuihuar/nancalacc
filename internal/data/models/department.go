package models

import (
	"database/sql"
	"time"
)

// 全量部门采集表
type TbLasDepartment struct {
	ID             uint           `gorm:"primaryKey;autoIncrement;column:id;type:int unsigned;comment:主键id" json:"id"`
	Did            string         `gorm:"not null;column:did;type:varchar(255);comment:部门id" json:"did"`
	TaskID         string         `gorm:"not null;column:task_id;type:varchar(20);comment:任务id" json:"task_id"`
	ThirdCompanyID string         `gorm:"not null;column:third_company_id;type:varchar(20);comment:租户id" json:"third_company_id"`
	PlatformID     string         `gorm:"not null;column:platform_id;type:varchar(60);comment:平台id" json:"platform_id"`
	Pid            sql.NullString `gorm:"column:pid;type:varchar(255);comment:父部门id" json:"pid"`
	Name           string         `gorm:"not null;column:name;type:varchar(255);comment:部门名称" json:"name"`
	Order          int            `gorm:"column:order;type:int;default:0;comment:排序" json:"order"`
	Source         string         `gorm:"column:source;type:varchar(20);default:sync;comment:来源" json:"source"`
	Ctime          sql.NullTime   `gorm:"column:ctime;type:timestamp;default:CURRENT_TIMESTAMP;comment:创建时间" json:"ctime"`
	Mtime          time.Time      `gorm:"not null;column:mtime;type:timestamp;comment:修改时间" json:"mtime"`
	CheckType      int8           `gorm:"not null;column:check_type;type:tinyint;default:0;comment:1-勾选 0-未勾选" json:"check_type"`
	Type           sql.NullString `gorm:"column:type;type:varchar(255);comment:类型" json:"type"`
}

// TableName 设置表名
func (TbLasDepartment) TableName() string {
	return "tb_las_department"
}

// 增量部门采集表
type TbLasDepartmentIncrement struct {
	ID             uint           `gorm:"primaryKey;autoIncrement;column:id;comment:主键id"`
	Did            string         `gorm:"not null;column:did;size:255;comment:部门id"`
	ThirdCompanyID string         `gorm:"not null;column:third_company_id;size:20;comment:租户id"`
	PlatformID     string         `gorm:"not null;column:platform_id;size:60;comment:平台id"`
	Pid            sql.NullString `gorm:"column:pid;size:255;comment:父部门id"`
	Name           string         `gorm:"not null;column:name;size:255;comment:部门名称"`
	Order          sql.NullInt32  `gorm:"column:order;comment:排序"`
	Source         string         `gorm:"column:source;size:20;default:sync;comment:来源"`
	SyncType       string         `gorm:"column:sync_type;size:20;default:auto;comment:同步方式"`
	UpdateType     string         `gorm:"not null;column:update_type;size:20;comment:修改类型"`
	Status         int            `gorm:"column:status;default:0;comment:状态"`
	Msg            sql.NullString `gorm:"column:msg;size:2000;comment:错误详情"`
	Operator       string         `gorm:"column:operator;size:100;default:系统;comment:operator"`
	SyncTime       sql.NullTime   `gorm:"column:sync_time;default:CURRENT_TIMESTAMP;comment:增量数据变动时间"`
	Ctime          sql.NullTime   `gorm:"column:ctime;default:CURRENT_TIMESTAMP;comment:创建时间"`
	Mtime          time.Time      `gorm:"not null;column:mtime;comment:修改时间"`
	Type           sql.NullString `gorm:"column:type;size:255;comment:类型"`
}

func (TbLasDepartmentIncrement) TableName() string {
	return "tb_las_department_increment"
}
