package data

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"nancalacc/internal/biz"
	"nancalacc/internal/data/models"
	"nancalacc/pkg/cipherutil"
	"nancalacc/pkg/httputil"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type accounterRepo struct {
	data *Data
	log  *log.Helper
}

var (
	Source = "sync"
)

// NewAccounterRepo .
func NewAccounterRepo(data *Data, logger log.Logger) biz.AccounterRepo {
	return &accounterRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *accounterRepo) SaveUsers(ctx context.Context, users []*biz.DingtalkDeptUser, taskId string) (int, error) {
	r.log.Infof("SaveUsers: %v", users)
	if len(users) == 0 {
		r.log.Infof("users is empty")
		return 0, nil
	}
	entities := make([]*models.TbLasUser, 0, len(users))

	thirdCompanyID := r.data.serviceConf.ThirdCompanyId
	platformID := r.data.serviceConf.PlatformIds
	for _, user := range users {
		if user.Nickname == "" {
			user.Nickname = user.Name
		}
		email, _ := cipherutil.AesEncryptGcmByKey(user.Email, r.data.serviceConf.SecretKey)
		phone, _ := cipherutil.AesEncryptGcmByKey(user.Mobile, r.data.serviceConf.SecretKey)

		entities = append(entities, &models.TbLasUser{
			TaskID:         taskId,
			ThirdCompanyID: thirdCompanyID,
			PlatformID:     platformID,
			Uid:            user.Unionid,
			DefDid:         sql.NullString{String: "-1", Valid: true},
			DefDidOrder:    0,
			Account:        user.Userid,
			NickName:       user.Nickname,
			Email:          sql.NullString{String: email, Valid: true},
			Phone:          sql.NullString{String: phone, Valid: true},
			Title:          sql.NullString{String: user.Title, Valid: true},
			//Leader:         sql.NullString{String: strconv.FormatBool(account.Leader)},
			Source:           Source,
			Ctime:            sql.NullTime{Time: time.Now(), Valid: true},
			Mtime:            time.Now(),
			CheckType:        1,
			EmploymentStatus: "active",
			//Type:           sql.NullString{String: "dept", Valid: true},
		})
	}

	result := r.data.db.WithContext(ctx).Create(&entities)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			r.log.Infof("user already exists")
		} else {
			return 0, result.Error
		}

	}

	return int(result.RowsAffected), nil
}

func (r *accounterRepo) SaveDepartments(ctx context.Context, depts []*biz.DingtalkDept, taskId string) (int, error) {
	r.log.Infof("SaveDepartments: %v", depts)
	entities := make([]*models.TbLasDepartment, 0, len(depts))

	thirdCompanyID := r.data.serviceConf.ThirdCompanyId
	platformID := r.data.serviceConf.PlatformIds
	companyID := r.data.serviceConf.CompanyId
	rootDep := &models.TbLasDepartment{
		Did:            companyID,
		TaskID:         taskId,
		Name:           companyID,
		ThirdCompanyID: thirdCompanyID,
		PlatformID:     platformID,
		Pid:            sql.NullString{String: "-1", Valid: true},
		Order:          0,
		Source:         "sync",
		Ctime:          sql.NullTime{Time: time.Now(), Valid: true},
		Mtime:          time.Now(),
		CheckType:      1,
		//Type:           sql.NullString{String: "dept", Valid: true},
	}
	for _, dep := range depts {
		entities = append(entities, &models.TbLasDepartment{
			Did:            strconv.FormatInt(dep.DeptID, 10),
			TaskID:         taskId,
			Name:           dep.Name,
			ThirdCompanyID: thirdCompanyID,
			PlatformID:     platformID,
			Pid:            sql.NullString{String: strconv.FormatInt(dep.ParentID, 10), Valid: true},
			Order:          int(dep.Order),
			Source:         "sync",
			Ctime:          sql.NullTime{Time: time.Now(), Valid: true},
			Mtime:          time.Now(),
			CheckType:      1,
			//Type:           sql.NullString{String: "dept", Valid: true},
		})
	}
	entities = append(entities, rootDep)
	result := r.data.db.WithContext(ctx).Create(&entities)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			r.log.Infof("department already exists")
		} else {
			return 0, result.Error
		}

	}

	return int(result.RowsAffected), nil
}

func (r *accounterRepo) SaveDepartmentUserRelations(ctx context.Context, relations []*biz.DingtalkDeptUserRelation, taskId string) (int, error) {
	r.log.Infof("SaveDepartmentUserRelations: %v", relations)

	entities := make([]*models.TbLasDepartmentUser, 0, len(relations))

	thirdCompanyID := r.data.serviceConf.ThirdCompanyId
	platformID := r.data.serviceConf.PlatformIds
	for _, relation := range relations {
		entities = append(entities, &models.TbLasDepartmentUser{
			Did:            relation.Did,
			TaskID:         taskId,
			ThirdCompanyID: thirdCompanyID,
			PlatformID:     platformID,
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

	thirdCompanyID := r.data.serviceConf.ThirdCompanyId
	platformID := r.data.serviceConf.PlatformIds
	companyID := r.data.serviceConf.CompanyId
	entity := &models.TbCompanyCfg{
		ThirdCompanyId: thirdCompanyID,
		PlatformIds:    platformID,
		CompanyId:      companyID,
		Status:         1,
		Ctime:          sql.NullTime{Time: time.Now(), Valid: true},
		Mtime:          time.Now(),
	}

	err := r.data.db.WithContext(ctx).Create(entity).Error

	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			r.log.Infof("company config already exists")
		} else {
			return err
		}

	}

	return nil
}

func (r *accounterRepo) CallEcisaccountsyncAll(ctx context.Context, taskId string) (biz.EcisaccountsyncResponse, error) {
	r.log.Infof("CallEcisaccountsyncAll: %v", taskId)

	path := r.data.serviceConf.EcisaccountsyncUrl
	thirdCompanyID := r.data.serviceConf.ThirdCompanyId
	collectCost := "1100000"
	uri := fmt.Sprintf("%s?taskId=%s&thirdCompanyId=%s&collectCost=%s", path, taskId, thirdCompanyID, collectCost)
	var resp biz.EcisaccountsyncResponse
	r.log.Infof("CallEcisaccountsyncAll: %s", uri)
	bs, err := httputil.PostJSON(uri, nil, time.Second*10)
	r.log.Infof("CallEcisaccountsyncAll: %s", string(bs))
	if err != nil {
		return biz.EcisaccountsyncResponse{}, fmt.Errorf("CallEcisaccountsyncAll: %w", err)
	}
	err = json.Unmarshal(bs, &resp)
	if err != nil {
		return biz.EcisaccountsyncResponse{}, fmt.Errorf("CallEcisaccountsyncAll: %w", err)
	}
	if resp.Code != "200" {
		return biz.EcisaccountsyncResponse{}, fmt.Errorf("CallEcisaccountsyncAll: %s", resp.Msg)
	}
	return resp, nil

}
