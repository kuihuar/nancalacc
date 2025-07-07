package data

import (
	"context"
	"database/sql"
	"strconv"
	"time"

	"nancalacc/internal/biz"
	"nancalacc/internal/data/models"

	"github.com/go-kratos/kratos/v2/log"
)

type accounterRepo struct {
	data *Data
	log  *log.Helper
}

var (
	ThirdCompanyID = "nancal"
	PlatformID     = "dingtalk"
	Source         = "sync"
)

// NewAccounterRepo .
func NewAccounterRepo(data *Data, logger log.Logger) biz.AccounterRepo {
	return &accounterRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *accounterRepo) SaveUsers(ctx context.Context, users []*biz.DingtalkDeptUser) (int, error) {
	r.log.Infof("SaveUsers: %v", users)
	entities := make([]*models.TbLasUser, 0, len(users))
	var taskIds []string
	for i := 1; i <= len(users); i++ {
		taskId := time.Now().Add(time.Duration(i) * time.Second).Format("20060102150405")
		taskIds = append(taskIds, taskId)
	}

	for index, user := range users {
		entities = append(entities, &models.TbLasUser{
			TaskID:         taskIds[index],
			ThirdCompanyID: ThirdCompanyID,
			PlatformID:     PlatformID,
			Uid:            user.Userid,
			DefDid:         sql.NullString{String: "1", Valid: true},
			DefDidOrder:    0,
			Account:        user.Userid,
			NickName:       user.Nickname,
			Email:          sql.NullString{String: user.Email, Valid: true},
			Phone:          sql.NullString{String: user.Mobile, Valid: true},
			Title:          sql.NullString{String: user.Title, Valid: true},
			//Leader:         sql.NullString{String: strconv.FormatBool(account.Leader)},
			Source:    Source,
			Ctime:     sql.NullTime{Time: time.Now(), Valid: true},
			Mtime:     time.Now(),
			CheckType: 1,
			//Type:           sql.NullString{String: "dept", Valid: true},
		})
	}

	result := r.data.db.WithContext(ctx).Create(&entities)

	if result.Error != nil {
		return 0, result.Error
	}

	return int(result.RowsAffected), nil
}

func (r *accounterRepo) SaveDepartments(ctx context.Context, depts []*biz.DingtalkDept) (int, error) {
	r.log.Infof("SaveDepartments: %v", depts)
	entities := make([]*models.TbLasDepartment, 0, len(depts))
	var taskIds []string
	for i := 1; i <= len(depts); i++ {
		taskId := time.Now().Add(time.Duration(i) * time.Second).Format("20060102150405")
		taskIds = append(taskIds, taskId)
	}
	for index, dep := range depts {
		entities = append(entities, &models.TbLasDepartment{
			Did:            strconv.FormatInt(dep.DeptID, 10),
			TaskID:         taskIds[index],
			Name:           dep.Name,
			ThirdCompanyID: ThirdCompanyID,
			PlatformID:     PlatformID,
			Pid:            sql.NullString{String: strconv.FormatInt(dep.ParentID, 10), Valid: true},
			Order:          int(dep.Order),
			Source:         "sync",
			Ctime:          sql.NullTime{Time: time.Now(), Valid: true},
			Mtime:          time.Now(),
			CheckType:      1,
			//Type:           sql.NullString{String: "dept", Valid: true},
		})
	}
	result := r.data.db.WithContext(ctx).Create(&entities)

	if result.Error != nil {
		return 0, result.Error
	}

	return int(result.RowsAffected), nil
}

func (r *accounterRepo) SaveDepartmentUserRelations(ctx context.Context, relations []*biz.DingtalkDeptUserRelation) (int, error) {
	r.log.Infof("SaveDepartmentUserRelations: %v", relations)

	entities := make([]*models.TbLasDepartmentUser, 0, len(relations))
	var taskIds []string
	for i := 1; i <= len(relations); i++ {
		taskId := time.Now().Add(time.Duration(i) * time.Second).Format("20060102150405")
		taskIds = append(taskIds, taskId)
	}
	for index, relation := range relations {
		entities = append(entities, &models.TbLasDepartmentUser{
			Did:            relation.Did,
			TaskID:         taskIds[index],
			ThirdCompanyID: ThirdCompanyID,
			PlatformID:     PlatformID,
			Uid:            relation.Uid,
			Ctime:          time.Now(),
			Order:          sql.NullInt32{Int32: int32(relation.Order), Valid: true},
			CheckType:      1,
		})
	}
	result := r.data.db.WithContext(ctx).Create(&entities)

	if result.Error != nil {
		return 0, result.Error
	}

	return int(result.RowsAffected), nil
}

func (r *accounterRepo) SaveCompanyCfg(ctx context.Context, cfg *biz.DingtalkCompanyCfg) error {
	r.log.Infof("SaveCompanyCfg: %v", cfg)
	entity := &models.TbCompanyCfg{
		ThirdCompanyId: ThirdCompanyID,
		PlatformIds:    PlatformID,
		CompanyId:      ThirdCompanyID,
		Status:         1,
		Ctime:          sql.NullTime{Time: time.Now(), Valid: true},
		Mtime:          time.Now(),
	}
	return r.data.db.WithContext(ctx).Create(entity).Error
}
