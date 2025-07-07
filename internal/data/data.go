package data

import (
	"errors"
	"fmt"
	"nancalacc/internal/conf"
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
}

type ThirdParty struct {
	Endpoint  string
	AppKey    string
	AppSecret string
	Timeout   string
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	fmt.Printf("=====newData.c: %v", c)

	var db *gorm.DB
	var err error

	db, err = initDB(c)
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
	}, cleanup, nil
}

func descryt(key, sk string) string {
	return key
}
func initDbEnv(c *conf.Data) (*gorm.DB, error) {
	ak, err := conf.GetEnv("ECIS_ECISACCOUNTSYNC_DB")
	if err != nil {
		panic(err)
	}
	sk := c.ServiceConf.SecretKey
	fmt.Printf("aaaaaaaaaaaak: %v, sk: %v", ak, sk)
	dsn := descryt(ak, sk)

	// [POST] http://encs-pri-proxy-gateway/ecisaccountsync/api/sync/all
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})

}

func initDB(c *conf.Data) (*gorm.DB, error) {
	//initDbEnv(c)
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

func getGormLogLevel(level string) gormlogger.LogLevel {
	switch level {
	case "silent":
		return gormlogger.Silent
	case "error":
		return gormlogger.Error
	case "warn":
		return gormlogger.Warn
	case "info":
		return gormlogger.Info
	default:
		return gormlogger.Warn // 默认级别
	}
}
