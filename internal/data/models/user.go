package models

import (
	"database/sql"
	"time"
)

type TbLasUser struct {
	ID               uint           `gorm:"primaryKey;autoIncrement;column:id;type:int unsigned;comment:主键id"`
	TaskID           string         `gorm:"not null;column:task_id;type:varchar(20);comment:任务id"`
	ThirdCompanyID   string         `gorm:"not null;column:third_company_id;type:varchar(20);comment:租户id"`
	PlatformID       string         `gorm:"not null;column:platform_id;type:varchar(60);comment:平台id"`
	Uid              string         `gorm:"not null;column:uid;type:varchar(255);comment:用户id"`
	DefDid           sql.NullString `gorm:"column:def_did;type:varchar(255);comment:默认部门"`
	DefDidOrder      int            `gorm:"column:def_did_order;type:int;default:0;comment:排序"`
	Account          string         `gorm:"not null;column:account;type:varchar(255);comment:登录名"`
	NickName         string         `gorm:"not null;column:nick_name;type:varchar(255);comment:用户昵称"`
	Password         sql.NullString `gorm:"column:password;type:varchar(255);comment:密码"`
	Avatar           sql.NullString `gorm:"column:avatar;type:varchar(255);comment:头像"`
	Email            sql.NullString `gorm:"column:email;type:varchar(80);comment:邮箱"`
	Gender           sql.NullString `gorm:"column:gender;type:varchar(60);comment:性别"`
	Title            sql.NullString `gorm:"column:title;type:varchar(255);comment:职称"`
	WorkPlace        sql.NullString `gorm:"column:work_place;type:varchar(255);comment:办公地点"`
	Leader           sql.NullString `gorm:"column:leader;type:varchar(255);comment:上级主管ID"`
	Employer         sql.NullString `gorm:"column:employer;type:varchar(255);comment:员工工号"`
	EmploymentStatus string         `gorm:"column:employment_status;type:varchar(60);default:notactive;comment:就职状态"`
	EmploymentType   sql.NullString `gorm:"column:employment_type;type:varchar(60);comment:就职类型"`
	Phone            sql.NullString `gorm:"column:phone;type:varchar(200);comment:手机号"`
	Telephone        sql.NullString `gorm:"column:telephone;type:varchar(200);comment:座机号"`
	Source           string         `gorm:"column:source;type:varchar(20);default:sync;comment:来源"`
	CustomFields     sql.NullString `gorm:"column:custom_fields;type:varchar(5000);comment:自定义字段"`
	Ctime            sql.NullTime   `gorm:"column:ctime;type:timestamp;default:CURRENT_TIMESTAMP;comment:创建时间"`
	Mtime            time.Time      `gorm:"not null;column:mtime;type:timestamp;comment:更新时间"`
	CheckType        int8           `gorm:"not null;column:check_type;type:tinyint;default:0;comment:勾选状态"`
}

func (TbLasUser) TableName() string {
	return "tb_las_user"
}

type TbLasUserIncrement struct {
	ID               uint           `gorm:"primaryKey;autoIncrement;column:id;type:int unsigned;comment:主键id"`
	ThirdCompanyID   string         `gorm:"not null;column:third_company_id;type:varchar(20);comment:租户id"`
	PlatformID       string         `gorm:"not null;column:platform_id;type:varchar(60);comment:平台id"`
	Uid              string         `gorm:"not null;column:uid;type:varchar(255);comment:用户id"`
	DefDid           sql.NullString `gorm:"column:def_did;type:varchar(255);comment:默认部门"`
	DefDidOrder      int            `gorm:"column:def_did_order;type:int;default:0;comment:排序"`
	Account          string         `gorm:"not null;column:account;type:varchar(255);comment:登录名"`
	NickName         string         `gorm:"not null;column:nick_name;type:varchar(255);comment:用户昵称"`
	Password         sql.NullString `gorm:"column:password;type:varchar(255);comment:密码"`
	Avatar           sql.NullString `gorm:"column:avatar;type:varchar(255);comment:头像"`
	Email            sql.NullString `gorm:"column:email;type:varchar(80);comment:邮箱"`
	Gender           sql.NullString `gorm:"column:gender;type:varchar(60);comment:性别"`
	Title            sql.NullString `gorm:"column:title;type:varchar(255);comment:职称"`
	WorkPlace        sql.NullString `gorm:"column:work_place;type:varchar(255);comment:办公地点"`
	Leader           sql.NullString `gorm:"column:leader;type:varchar(255);comment:上级主管ID"`
	Employer         sql.NullString `gorm:"column:employer;type:varchar(255);comment:员工工号"`
	EmploymentStatus string         `gorm:"column:employment_status;type:varchar(60);comment:就职状态"`
	EmploymentType   sql.NullString `gorm:"column:employment_type;type:varchar(60);comment:就职类型"`
	Phone            sql.NullString `gorm:"column:phone;type:varchar(200);comment:手机号"`
	Telephone        sql.NullString `gorm:"column:telephone;type:varchar(200);comment:座机号"`
	Source           string         `gorm:"column:source;type:varchar(20);default:sync;comment:来源"`
	CustomFields     sql.NullString `gorm:"column:custom_fields;type:varchar(5000);comment:自定义字段"`
	SyncType         string         `gorm:"column:sync_type;type:varchar(20);default:auto;comment:同步方式"`
	UpdateType       string         `gorm:"not null;column:update_type;type:varchar(20);comment:修改类型"`
	Status           int            `gorm:"column:status;type:int;default:0;comment:状态"`
	Msg              sql.NullString `gorm:"column:msg;type:varchar(2000);comment:错误详情"`
	Operator         string         `gorm:"column:operator;type:varchar(100);default:系统;comment:operator"`
	SyncTime         sql.NullTime   `gorm:"column:sync_time;type:timestamp;default:CURRENT_TIMESTAMP;comment:增量数据变动时间"`
	Ctime            sql.NullTime   `gorm:"column:ctime;type:timestamp;default:CURRENT_TIMESTAMP;comment:创建时间"`
	Mtime            time.Time      `gorm:"not null;column:mtime;type:timestamp;comment:更新时间"`
}

func (TbLasUserIncrement) TableName() string {
	return "tb_las_user_increment"
}
