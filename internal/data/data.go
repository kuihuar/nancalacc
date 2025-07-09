package data

import (
	"errors"
	"fmt"
	"nancalacc/internal/conf"
	"nancalacc/pkg/cipherutil"
	"time"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dingtalkcontact_1_0 "github.com/alibabacloud-go/dingtalk/contact_1_0"
	dingtalkoauth2_1_0 "github.com/alibabacloud-go/dingtalk/oauth2_1_0"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewAccounterRepo, NewDingTalkRepo)

// Data .
type Data struct {
	// TODO wrapped database client
	db *gorm.DB
	// 钉钉配置
	thirdParty  *ThirdParty
	dingtalkCli *dingtalkoauth2_1_0.Client

	dingtalkCliContact *dingtalkcontact_1_0.Client

	// 服务配置
	serviceConf *ServiceConf
}

type ThirdParty struct {
	Endpoint  string
	AppKey    string
	AppSecret string
	Timeout   string
}

type ServiceConf struct {
	CompanyId          string
	ThirdCompanyId     string
	PlatformIds        string
	SecretKey          string
	AccessKey          string
	EcisaccountsyncUrl string
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	fmt.Printf("=====newData.c: %v", c)

	var db *gorm.DB
	var err error

	if c.ServiceConf.Env == "prod" {
		db, err = initDbEnv(c, logger)
	} else {
		db, err = initDB(c, logger)
	}
	if err != nil {
		return nil, nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, nil, err
	}
	sqlDB.SetMaxIdleConns(10)           // 最大空闲连接数
	sqlDB.SetMaxOpenConns(100)          // 最大打开连接数
	sqlDB.SetConnMaxLifetime(time.Hour) // 连接最大存活时间

	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
		if err := sqlDB.Close(); err != nil {
			log.NewHelper(logger).Error(err)
		}
	}
	// if err := Migrate(db); err != nil {
	// 	return nil, cleanup, err
	// }
	// if err = Seed(db); err != nil {
	// 	log.NewHelper(logger).Errorf("seed data failed: %v", err)
	// }
	dingtalk := &ThirdParty{
		Endpoint:  c.Dingtalk.Endpoint,
		AppKey:    c.Dingtalk.AppKey,
		AppSecret: c.Dingtalk.AppSecret,
		Timeout:   c.Dingtalk.Timeout,
	}
	serviceConf := &ServiceConf{
		CompanyId:          c.ServiceConf.CompanyId,
		ThirdCompanyId:     c.ServiceConf.ThirdCompanyId,
		PlatformIds:        c.ServiceConf.PlatformIds,
		SecretKey:          c.ServiceConf.SecretKey,
		AccessKey:          c.ServiceConf.AccessKey,
		EcisaccountsyncUrl: c.ServiceConf.EcisaccountsyncUrl,
	}
	config := &openapi.Config{
		Protocol: tea.String("https"),
		RegionId: tea.String("central"),
	}

	client, err := dingtalkoauth2_1_0.NewClient(config)
	if err != nil {
		return nil, cleanup, err
	}

	clientContact, err := dingtalkcontact_1_0.NewClient(config)
	if err != nil {
		return nil, cleanup, err
	}
	return &Data{
		db:                 db,
		thirdParty:         dingtalk,
		dingtalkCli:        client,
		dingtalkCliContact: clientContact,
		serviceConf:        serviceConf,
	}, cleanup, nil
}

func initDbEnv(c *conf.Data, logger log.Logger) (*gorm.DB, error) {
	encryptedDsn, err := conf.GetEnv("ECIS_ECISACCOUNTSYNC_DB")
	log.NewHelper(logger).Info("initDbEnv: %s", encryptedDsn)
	if err != nil {
		log.NewHelper(logger).Error("initDbEnv: %w", err)
		return nil, err
	}
	sk := c.ServiceConf.SecretKey
	dsn, err := cipherutil.DecryptByAes(encryptedDsn, sk)
	if err != nil {
		log.NewHelper(logger).Error("initDbEnvDecryptByAes: %w", err)
		return nil, err
	}
	if len(dsn) == 0 {
		log.NewHelper(logger).Error("initDbEnvDecryptByAes: dsn is empty")
		return nil, err
	}
	dsn = "mysql://" + dsn
	log.NewHelper(logger).Error("initDbEnvDecrypt.dsn: %s", dsn)
	// [POST] http://encs-pri-proxy-gateway/ecisaccountsync/api/sync/all
	return gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Info),
	})

}

func initDB(c *conf.Data, logger log.Logger) (*gorm.DB, error) {
	var dialector gorm.Dialector
	switch c.Database.Driver {
	case "mysql":
		dialector = mysql.Open(c.Database.Source)
	case "sqlite":
		dialector = sqlite.Open(c.Database.Source)
	default:
		return nil, errors.New("unsupported database driver")
	}

	return gorm.Open(dialector, &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "nancal_", // 表前缀
			SingularTable: true,      // 使用单数表名
		},
	})
}
