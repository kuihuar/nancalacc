package data

import (
	"context"
	"fmt"
	"nancalacc/internal/data/models"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

const (
	pageSize   = 1000  // 每页数据量
	maxResults = 10000 // 最大结果数限制
)

func (r *accounterRepo) BatchGetDeptUsers(ctx context.Context, taskName string) ([]*models.TbLasDepartmentUser, error) {
	r.log.Log(log.LevelInfo, "msg", "BatchGetDeptUsers", "taskName", taskName)

	var (
		allUsers   = make([]*models.TbLasDepartmentUser, 0, pageSize)
		lastID     uint64
		totalCount int
	)
	db, err := r.data.GetSyncDB()
	if err != nil {
		return nil, err
	}

	for {
		var pageUsers []*models.TbLasDepartmentUser

		query := db.WithContext(ctx).
			Where("task_id = ?", taskName).
			Order("id ASC") // 必须按ID排序

		if lastID > 0 {
			query = query.Where("id > ?", lastID)
		}

		result := query.Limit(pageSize).Find(&pageUsers)
		if result.Error != nil {
			r.log.Log(log.LevelError, "msg", "BatchGetDeptUsers", "Query failed at lastID", lastID, "err", result.Error)
			return nil, fmt.Errorf("database error: %w", result.Error)
		}

		if len(pageUsers) == 0 {
			break // 没有更多数据
		}

		// 更新lastID为当前页最后一条记录的ID
		lastID = uint64(pageUsers[len(pageUsers)-1].ID)
		allUsers = append(allUsers, pageUsers...)
		totalCount += len(pageUsers)

		// 如果获取数量小于pageSize，说明是最后一页
		if len(pageUsers) < pageSize {
			break
		}
	}

	if len(allUsers) == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	r.log.Log(log.LevelInfo, "msg", "BatchGetDeptUsers", "Fetched", totalCount, "records")
	return allUsers, nil
}

func (r *accounterRepo) BatchGetUsers(ctx context.Context, taskName string) ([]*models.TbLasUser, error) {
	r.log.Log(log.LevelInfo, "msg", "BatchGetUsers", "taskName", taskName)

	var (
		allUsers   = make([]*models.TbLasUser, 0, pageSize)
		lastID     uint64
		totalCount int
	)

	selectedFields := []string{"id", "name", "dept_id"} // 只选择必要字段

	db, err := r.data.GetSyncDB()
	if err != nil {
		return nil, err
	}
	for {
		var pageUsers []*models.TbLasUser

		query := db.WithContext(ctx).
			Select(selectedFields).
			Where("task_id = ?", taskName).
			Order("id ASC")

		if lastID > 0 {
			query = query.Where("id > ?", lastID)
		}

		result := query.Limit(pageSize).Find(&pageUsers)
		if result.Error != nil {
			r.log.Log(log.LevelError, "msg", "BatchGetUsers", "Query failed at lastID", lastID, "err", result.Error)
			return nil, fmt.Errorf("database error: %w", result.Error)
		}

		if len(pageUsers) == 0 {
			break
		}

		lastID = uint64(pageUsers[len(pageUsers)-1].ID)
		allUsers = append(allUsers, pageUsers...)
		totalCount += len(pageUsers)

		if len(pageUsers) < pageSize {
			break
		}
	}

	if len(allUsers) == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	r.log.Log(log.LevelInfo, "msg", "BatchGetUsers", "Fetched", totalCount, "users")
	return allUsers, nil
}

func (r *accounterRepo) BatchGetDepts(ctx context.Context, taskName string) ([]*models.TbLasDepartment, error) {
	r.log.Log(log.LevelInfo, "msg", "BatchGetDepts", "taskName", taskName)

	var (
		allDepts   = make([]*models.TbLasDepartment, 0, pageSize)
		lastID     uint64
		totalCount int
	)

	db, err := r.data.GetSyncDB()
	if err != nil {
		return nil, err
	}
	for {
		var pageDepts []*models.TbLasDepartment

		query := db.WithContext(ctx).
			Where("task_id = ?", taskName).
			Order("id ASC")

		if lastID > 0 {
			query = query.Where("id > ?", lastID)
		}

		result := query.Limit(pageSize).Find(&pageDepts)
		if result.Error != nil {
			r.log.Log(log.LevelError, "msg", "BatchGetDepts", "Query failed at lastID", lastID, "err", result.Error)
			return nil, fmt.Errorf("database error: %w", result.Error)
		}

		if len(pageDepts) == 0 {
			break
		}

		lastID = uint64(pageDepts[len(pageDepts)-1].ID)
		allDepts = append(allDepts, pageDepts...)
		totalCount += len(pageDepts)

		if len(pageDepts) < pageSize {
			break
		}
	}

	if len(allDepts) == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	r.log.Log(log.LevelInfo, "msg", "BatchGetDepts", "Fetched", totalCount, "departments")
	return allDepts, nil
}
