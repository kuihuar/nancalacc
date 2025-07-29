package models

import "time"

const (
	TaskStatusPending    = "pending"
	TaskStatusInProgress = "in_progress"
	TaskStatusCompleted  = "completed"
	TaskStatusCancelled  = "cancelled"
)

type Task struct {
	ID            uint      `gorm:"primaryKey;autoIncrement;column:id;comment:任务ID"`
	Title         string    `gorm:"type:varchar(255);not null;column:title;comment:任务标题"`
	Description   string    `gorm:"type:text;column:description;comment:任务描述"`
	Status        string    `gorm:"type:varchar(20);not null;default:pending;column:status;comment:任务状态(pending/in_progress/completed/cancelled)"`
	CreatedAt     time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP;column:created_at;comment:创建时间"`
	UpdatedAt     time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;column:updated_at;comment:更新时间"`
	DueDate       time.Time `gorm:"type:timestamp;column:due_date;comment:截止时间"`
	StartDate     time.Time `gorm:"type:timestamp;column:start_date;comment:开始时间"`
	CompletedAt   time.Time `gorm:"type:timestamp;column:completed_at;comment:完成时间"`
	CreatorID     uint      `gorm:"type:bigint;not null;column:creator_id;comment:创建人ID"`
	Progress      int8      `gorm:"type:tinyint;default:0;column:progress;comment:进度(0-100)"`
	EstimatedTime int       `gorm:"type:int;column:estimated_time;comment:预估耗时(分钟)"`
	ActualTime    int       `gorm:"type:int;column:actual_time;comment:实际耗时(分钟)"`
}

func (Task) TableName() string {
	return "task"
}

func (Task) Indexes() []string {
	return []string{"idx_status"}
}
