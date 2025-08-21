package mysql

import (
	"context"
	"nancalacc/internal/data"
	"nancalacc/internal/data/models"
	"nancalacc/internal/repository/contracts"

	"github.com/go-kratos/kratos/v2/log"
)

type taskRepository struct {
	data *data.Data
	log  log.Logger
}

// NewTaskRepository 创建任务Repository
func NewTaskRepository(data *data.Data, logger log.Logger) contracts.TaskRepository {
	return &taskRepository{
		data: data,
		log:  logger,
	}
}

func (r *taskRepository) CreateTask(ctx context.Context, taskName string) (int, error) {
	db, err := r.data.GetSyncDB()
	if err != nil {
		return 0, err
	}

	task := &models.Task{
		Title:  taskName,
		Status: "running",
	}

	result := db.WithContext(ctx).Create(task)
	if result.Error != nil {
		return 0, result.Error
	}

	return int(task.ID), nil
}

func (r *taskRepository) UpdateTask(ctx context.Context, taskName, status string) error {
	db, err := r.data.GetSyncDB()
	if err != nil {
		return err
	}

	result := db.WithContext(ctx).Model(&models.Task{}).
		Where("task_name = ?", taskName).
		Update("status", status)

	return result.Error
}

func (r *taskRepository) GetTask(ctx context.Context, taskName string) (*models.Task, error) {
	db, err := r.data.GetSyncDB()
	if err != nil {
		return nil, err
	}

	var task models.Task
	result := db.WithContext(ctx).Where("task_name = ?", taskName).First(&task)
	if result.Error != nil {
		return nil, result.Error
	}

	return &task, nil
}
