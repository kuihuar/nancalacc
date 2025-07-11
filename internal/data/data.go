package data

import (
	"errors"
	"fmt"
	"nancalacc/internal/conf"
	"nancalacc/pkg/cipherutil"
	"strings"
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
	AppPackage         string
	AppSecret          string
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	fmt.Printf("=====newData.c: %v\n", c)

	var db *gorm.DB
	var err error

	db, err = initDbEnv(c, logger)
	// db, err = initDB(c, logger)
	if err != nil {
		log.NewHelper(logger).Error("NewData: init db env failed")
		return nil, nil, nil
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
		AppPackage:         c.ServiceConf.AppPackage,
		AppSecret:          c.ServiceConf.AppSecret,
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
	// packagename := "com.acc.async"

	fmt.Printf("=====encryptedDsn.ECIS_ECISACCOUNTSYNC_DB: %s, err: %v\n", encryptedDsn, err)
	log.NewHelper(logger).Info("initDbEnv: %s", encryptedDsn)
	if err != nil {
		log.NewHelper(logger).Error("initDbEnv: %w", err)
		return nil, err
	}
	appSecret := c.ServiceConf.AppSecret
	dsn, err := cipherutil.DecryptByAes(encryptedDsn, appSecret)
	fmt.Printf("=====ECIS_ECISACCOUNTSYNC_DB dsn: %s, appSecret:%s, err: %v\n", dsn, appSecret, err)
	if err != nil {
		log.NewHelper(logger).Error("initDbEnvDecryptByAes: %w", err)
		return nil, err
	}
	if len(dsn) == 0 {
		log.NewHelper(logger).Error("initDbEnvDecryptByAes: dsn is empty")
		return nil, err
	}
	fmt.Printf("=====ECIS_ECISACCOUNTSYNC_DB dsn: %s\n", dsn)

	if !strings.Contains(dsn, "parseTime=True") {
		dsn = dsn + "&parseTime=True"
	}
	fmt.Printf("=====ECIS_ECISACCOUNTSYNC_DB dsn: %s\n", dsn)
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
