package data

import (
	"context"
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
	"gorm.io/gorm/clause"
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
		if user.Name == "" {
			user.Name = user.Userid
		}
		if user.Nickname == "" {
			user.Nickname = user.Name
		}
		var account string
		if user.Name == "" {
			account = user.Userid
		} else {
			account = user.Name
		}

		email, errEmail := cipherutil.AesEncryptGcmByKey(user.Email, r.data.serviceConf.SecretKey)

		phone, errPhone := cipherutil.AesEncryptGcmByKey(user.Mobile, r.data.serviceConf.SecretKey)

		r.log.Infof("AesEncryptGcmByKey email: %v, ncrypt email: %v, err: %v", user.Email, email, errEmail)
		r.log.Infof("AesEncryptGcmByKey phone: %v, ncrypt phone: %v, err: %v", user.Mobile, phone, errPhone)
		entities = append(entities, &models.TbLasUser{
			TaskID:         taskId,
			ThirdCompanyID: thirdCompanyID,
			PlatformID:     platformID,
			Uid:            user.Userid,
			DefDid:         "-1",
			DefDidOrder:    0,
			Account:        account,
			NickName:       user.Nickname,
			Email:          email,
			Phone:          phone,
			Title:          user.Title,
			//Leader:         sql.NullString{String: strconv.FormatBool(account.Leader)},
			Source:           Source,
			Ctime:            time.Now(),
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
		Did:            "0",
		TaskID:         taskId,
		Name:           companyID,
		ThirdCompanyID: thirdCompanyID,
		PlatformID:     platformID,
		Pid:            "-1",
		Order:          0,
		Source:         "sync",
		Ctime:          time.Now(),
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
			Pid:            strconv.FormatInt(dep.ParentID, 10),
			Order:          int(dep.Order),
			Source:         "sync",
			Ctime:          time.Now(),
			Mtime:          time.Now(),
			CheckType:      1,
			//Type:           sql.NullString{String: "dept", Valid: true},
		})
	}
	entities = append(entities, rootDep)
	result := r.data.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "did"}, {Name: "task_id"}, {Name: "third_company_id"}, {Name: "platform_id"}},
		DoNothing: true,
	}).Create(&entities)

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
			Order:          relation.Order,
			CheckType:      1,
		})
	}

	result := r.data.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "did"}, {Name: "uid"}, {Name: "task_id"}, {Name: "third_company_id"}, {Name: "platform_id"}},
		DoNothing: true,
	}).Create(&entities)

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
		Ctime:          time.Now(),
		Mtime:          time.Now(),
	}

	err := r.data.db.WithContext(ctx).Where(models.TbCompanyCfg{
		ThirdCompanyId: thirdCompanyID,
		CompanyId:      companyID,
	}).FirstOrCreate(entity).Error

	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			r.log.Infof("company config already exists")
		} else {
			return err
		}

	}

	return nil
}

func (r *accounterRepo) CallEcisaccountsyncAll(ctx context.Context, taskId string) (biz.EcisaccountsyncAllResponse, error) {
	r.log.Infof("CallEcisaccountsyncAll: %v", taskId)

	path := r.data.serviceConf.EcisaccountsyncUrl

	// path := "http://encs-pri-proxy-gateway/ecisaccountsync/api/sync/all"
	var resp biz.EcisaccountsyncAllResponse
	thirdCompanyID := r.data.serviceConf.ThirdCompanyId
	collectCost := "1100000"
	uri := fmt.Sprintf("%s?taskId=%s&thirdCompanyId=%s&collectCost=%s", path, taskId, thirdCompanyID, collectCost)

	r.log.Infof("CallEcisaccountsyncAll uri: %s", uri)
	bs, err := httputil.PostJSON(uri, nil, time.Second*10)
	r.log.Infof("CallEcisaccountsyncAll.Post output: bs:%s, err:%w", string(bs), err)

	if err != nil {
		return resp, err
	}
	err = json.Unmarshal(bs, &resp)
	if err != nil {
		return resp, fmt.Errorf("Unmarshal err: %w", err)
	}
	if resp.Code != "200" {
		return resp, fmt.Errorf("code not 200: %s", resp.Code)
	}

	return resp, nil

}

func (r *accounterRepo) ClearAll(ctx context.Context) error {
	err := r.data.db.WithContext(ctx).Exec("truncate table tb_company_cfg").Error
	if err != nil {
		return err
	}
	err = r.data.db.WithContext(ctx).Exec("truncate table tb_las_department").Error
	if err != nil {
		return err
	}
	err = r.data.db.WithContext(ctx).Exec("truncate table tb_las_department_user").Error
	if err != nil {
		return err
	}
	err = r.data.db.WithContext(ctx).Exec("truncate table tb_las_account").Error
	if err != nil {
		return err
	}
	return nil
}

func (r *accounterRepo) SaveIncrementUsers(ctx context.Context, usersAdd, usersDel []*biz.DingtalkDeptUser) error {
	entities := make([]*models.TbLasUserIncrement, 0, len(usersAdd)+len(usersDel))

	thirdCompanyID := r.data.serviceConf.ThirdCompanyId
	platformID := r.data.serviceConf.PlatformIds
	for _, user := range usersAdd {
		if user.Name == "" {
			user.Name = user.Userid
		}
		if user.Nickname == "" {
			user.Nickname = user.Name
		}
		var account string
		if user.Name == "" {
			account = user.Userid
		} else {
			account = user.Name
		}
		email, errEmail := cipherutil.AesEncryptGcmByKey(user.Email, r.data.serviceConf.SecretKey)

		phone, errPhone := cipherutil.AesEncryptGcmByKey(user.Mobile, r.data.serviceConf.SecretKey)

		r.log.Infof("AesEncryptGcmByKey email: %v, ncrypt email: %v, err: %v", user.Email, email, errEmail)
		r.log.Infof("AesEncryptGcmByKey phone: %v, ncrypt phone: %v, err: %v", user.Mobile, phone, errPhone)

		entities = append(entities, &models.TbLasUserIncrement{
			ThirdCompanyID: thirdCompanyID,
			PlatformID:     platformID,
			Uid:            user.Userid,
			DefDid:         "-1",
			DefDidOrder:    0,
			Account:        account,
			NickName:       user.Nickname,
			Email:          email,
			Phone:          phone,
			Title:          user.Title,
			//Leader:         sql.NullString{String: strconv.FormatBool(account.Leader)},
			Source:           Source,
			Ctime:            time.Now(),
			Mtime:            time.Now(),
			EmploymentStatus: "active",
			UpdateType:       "user_add",
			SyncType:         "auto",
			SyncTime:         time.Now(),
			Status:           0,
			//Type:           sql.NullString{String: "dept", Valid: true},
		})
	}
	for _, user := range usersDel {
		entities = append(entities, &models.TbLasUserIncrement{
			ThirdCompanyID: thirdCompanyID,
			PlatformID:     platformID,
			Uid:            user.Userid,
			DefDid:         "-1",
			DefDidOrder:    0,
			Account:        user.Userid,
			NickName:       user.Nickname,
			Email:          "email",
			Phone:          "phone",
			Title:          user.Title,
			//Leader:         sql.NullString{String: strconv.FormatBool(account.Leader)},
			Source:           Source,
			Ctime:            time.Now(),
			Mtime:            time.Now(),
			EmploymentStatus: "active",
			UpdateType:       "user_del",
			SyncType:         "auto",
			SyncTime:         time.Now(),
			Status:           0,
			//Type:           sql.NullString{String: "dept", Valid: true},
		})
	}
	result := r.data.db.WithContext(ctx).Create(&entities)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			r.log.Infof("user already exists")
		} else {
			return result.Error
		}

	}

	return nil
}
func (r *accounterRepo) SaveIncrementDepartments(ctx context.Context, deptsAdd, deptsDel []*biz.DingtalkDept) error {
	r.log.Infof("SaveIncrementDepartments deptsAdd: %v, deptsDel: %v", deptsAdd, deptsDel)
	entities := make([]*models.TbLasDepartmentIncrement, 0, len(deptsAdd)+len(deptsDel))

	thirdCompanyID := r.data.serviceConf.ThirdCompanyId
	platformID := r.data.serviceConf.PlatformIds

	for _, dep := range deptsDel {
		entities = append(entities, &models.TbLasDepartmentIncrement{
			Did:            strconv.FormatInt(dep.DeptID, 10),
			Name:           dep.Name,
			ThirdCompanyID: thirdCompanyID,
			PlatformID:     platformID,
			// 删除的时候这个父id, 怎么传
			Pid:        "",
			Order:      int32(dep.Order),
			Source:     "sync",
			Ctime:      time.Now(),
			Mtime:      time.Now(),
			UpdateType: "dept_move",
			SyncTime:   time.Now(),
			SyncType:   "auto",
			Status:     0,
		})
	}

	for _, dep := range deptsAdd {
		entities = append(entities, &models.TbLasDepartmentIncrement{
			Did:            strconv.FormatInt(dep.DeptID, 10),
			Name:           dep.Name,
			ThirdCompanyID: thirdCompanyID,
			PlatformID:     platformID,
			Pid:            strconv.FormatInt(dep.ParentID, 10),
			Order:          int32(dep.Order),
			Source:         "sync",
			Ctime:          time.Now(),
			Mtime:          time.Now(),
			UpdateType:     "dept_add",
			SyncTime:       time.Now(),
			SyncType:       "auto",
			Status:         0,
		})
	}
	result := r.data.db.WithContext(ctx).Create(&entities)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			r.log.Infof("user already exists")
		} else {
			return result.Error
		}

	}

	return nil
}

func (r *accounterRepo) SaveIncrementDepartmentUserRelations(ctx context.Context, relationsAdd, relationsDel []*biz.DingtalkDeptUserRelation) error {

	entities := make([]*models.TbLasDepartmentUserIncrement, 0, len(relationsAdd)+len(relationsDel))

	thirdCompanyID := r.data.serviceConf.ThirdCompanyId
	platformID := r.data.serviceConf.PlatformIds

	for _, relation := range relationsAdd {
		entities = append(entities, &models.TbLasDepartmentUserIncrement{
			Did:            relation.Did,
			ThirdCompanyID: thirdCompanyID,
			PlatformID:     platformID,
			Uid:            relation.Uid,
			Ctime:          time.Now(),
			Order:          1,
			UpdateType:     "user_dept_add",
			SyncType:       "auto",
			SyncTime:       time.Now(),
			Status:         0,
		})
	}

	for _, relation := range relationsDel {
		entities = append(entities, &models.TbLasDepartmentUserIncrement{
			Did:            relation.Did,
			ThirdCompanyID: thirdCompanyID,
			PlatformID:     platformID,
			Uid:            relation.Uid,
			Ctime:          time.Now(),
			Order:          1,
			UpdateType:     "user_dept_del",
			SyncType:       "auto",
			SyncTime:       time.Now(),
			Status:         0,
		})
	}
	result := r.data.db.WithContext(ctx).Create(&entities)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			r.log.Infof("user already exists")
		} else {
			return result.Error
		}

	}

	return nil
}

func (r *accounterRepo) CallEcisaccountsyncIncrement(ctx context.Context, thirdCompanyID string) (biz.EcisaccountsyncIncrementResponse, error) {

	//待开发
	path := r.data.serviceConf.EcisaccountsyncUrlIncrement
	r.log.Infof("CallEcisaccountsyncIncrement path : %v", path)

	uri := path
	var resp biz.EcisaccountsyncIncrementResponse
	thirdCompanyID = r.data.serviceConf.ThirdCompanyId
	input := &biz.EcisaccountsyncIncrementRequest{
		ThirdCompanyId: thirdCompanyID,
	}
	jsonData, err := json.Marshal(input)
	if err != nil {
		return resp, err
	}

	r.log.Infof("CallEcisaccountsyncIncrement uri: %s, input: %s", uri, string(jsonData))
	bs, err := httputil.PostJSON(uri, jsonData, time.Second*10)
	r.log.Infof("CallEcisaccountsyncIncrement.Post output: bs:%s, err:%w", string(bs), err)

	if err != nil {
		return resp, err
	}
	err = json.Unmarshal(bs, &resp)
	if err != nil {
		return resp, fmt.Errorf("Unmarshal err: %w", err)
	}
	if resp.Code != "200" {
		return resp, fmt.Errorf("code not 200: %s", resp.Code)
	}

	r.callEcisaccountsyncIncrementTest(ctx, thirdCompanyID)

	return resp, nil

}

func (r *accounterRepo) callEcisaccountsyncIncrementTest(ctx context.Context, thirdCompanyID string) (biz.EcisaccountsyncIncrementResponse, error) {

	path := r.data.serviceConf.EcisaccountsyncUrlIncrement
	thirdCompanyID = r.data.serviceConf.ThirdCompanyId
	r.log.Infof("callEcisaccountsyncIncrementTest path : %v", path)

	uri := fmt.Sprintf("%s?&thirdCompanyId=%s", path, thirdCompanyID)

	var resp biz.EcisaccountsyncIncrementResponse

	r.log.Infof("callEcisaccountsyncIncrementTest uri: %s", uri)
	bs, err := httputil.PostJSON(uri, nil, time.Second*10)
	r.log.Infof("callEcisaccountsyncIncrementTest output: bs:%s, err:%w", string(bs), err)

	if err != nil {
		return resp, err
	}
	err = json.Unmarshal(bs, &resp)
	if err != nil {
		return resp, fmt.Errorf("Unmarshal err: %w", err)
	}
	if resp.Code != "200" {
		return resp, fmt.Errorf("code not 200: %s", resp.Code)
	}

	return resp, nil

}
