package data

import (
	"context"
	"fmt"
	"nancalacc/internal/data/models"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

const (
	pageSize = 1000 // 每页数据量
	// maxResults = 10000 // 最大结果数限制
)

func (r *accounterRepo) BatchGetDeptUsers(ctx context.Context, taskName, thirdCompanyId, platformId string) ([]*models.TbLasDepartmentUser, error) {
	r.log.Log(log.LevelInfo, "msg", "BatchGetDeptUsers", "taskName", taskName, "platformId", platformId)

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

		// 优化后的查询条件，按照索引顺序排列
		query := db.WithContext(ctx).
			Where("task_id = ?", taskName).
			Where("third_company_id = ?", thirdCompanyId).
			Where("platform_id = ?", platformId).
			Where("check_type = ?", 1). // 只查询勾选的记录
			Order("id ASC")

		if lastID > 0 {
			query = query.Where("id > ?", lastID)
		}

		result := query.Limit(pageSize).Find(&pageUsers)
		if result.Error != nil {
			r.log.Log(log.LevelError, "msg", "BatchGetDeptUsers", "Query failed at lastID", lastID, "err", result.Error)
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

	r.log.Log(log.LevelInfo, "msg", "BatchGetDeptUsers", "Fetched", totalCount, "records")
	return allUsers, nil
}

func (r *accounterRepo) BatchGetUsers(ctx context.Context, taskName, thirdCompanyId, platformId string) ([]*models.TbLasUser, error) {
	r.log.Log(log.LevelInfo, "msg", "BatchGetUsers", "taskName", taskName, "platformId", platformId)

	var (
		allUsers   = make([]*models.TbLasUser, 0, pageSize)
		lastID     uint64
		totalCount int
	)

	// 只选择必要字段，减少数据传输量
	selectedFields := []string{
		"id", "task_id", "third_company_id", "platform_id", "uid",
		"def_did", "def_did_order", "account", "nick_name", "email",
		"phone", "employment_status", "check_type",
	}

	db, err := r.data.GetSyncDB()
	if err != nil {
		return nil, err
	}

	for {
		var pageUsers []*models.TbLasUser

		// 优化后的查询条件，按照索引顺序排列
		query := db.WithContext(ctx).
			Select(selectedFields).
			Where("task_id = ?", taskName).
			Where("third_company_id = ?", thirdCompanyId).
			Where("platform_id = ?", platformId).
			//Where("check_type = ?", 1).               // 只查询勾选的记录
			//Where("employment_status = ?", "active"). // 只查询在职用户
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

func (r *accounterRepo) BatchGetDepts(ctx context.Context, taskName, thirdCompanyId, platformId string) ([]*models.TbLasDepartment, error) {
	r.log.Log(log.LevelInfo, "msg", "BatchGetDepts", "taskName", taskName, "platformId", platformId)

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

		// 优化后的查询条件，按照索引顺序排列
		query := db.WithContext(ctx).
			Where("task_id = ?", taskName).
			Where("third_company_id = ?", thirdCompanyId).
			Where("platform_id = ?", platformId).
			Where("check_type = ?", 1). // 只查询勾选的记录
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
