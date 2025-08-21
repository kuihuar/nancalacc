package contracts

import (
	"context"
	"nancalacc/internal/data/models"
)

// TaskRepository 任务数据访问接口
type TaskRepository interface {
	// 任务管理
	CreateTask(ctx context.Context, taskName string) (int, error)
	UpdateTask(ctx context.Context, taskName, status string) error
	GetTask(ctx context.Context, taskName string) (*models.Task, error)
}
