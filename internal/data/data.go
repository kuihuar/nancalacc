package data

import (
	"errors"
	"fmt"
	"nancalacc/internal/conf"
	"nancalacc/pkg/cipherutil"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewAccounterRepo)

type Data struct {
	db *gorm.DB
}

func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	fmt.Printf("=====newData.c: %v\n", c)

	var db *gorm.DB
	var err error

	// db, err = initDbEnv(c, logger)
	db, err = initDB(c, logger)
	if err != nil {
		log.NewHelper(logger).Error("NewData: init db env failed")
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

	return &Data{
		db: db,
	}, cleanup, nil
}

func initDbEnv(c *conf.Data, logger log.Logger) (*gorm.DB, error) {
	encryptedDsn, err := conf.GetEnv("ECIS_ECISACCOUNTSYNC_DB")

	log.NewHelper(logger).Info("initDbEnv: %s", encryptedDsn)
	if err != nil {
		log.NewHelper(logger).Error("initDbEnv: %w", err)
		return nil, err
	}
	appSecret := c.Auth.AppSecret
	dsn, err := cipherutil.DecryptByAes(encryptedDsn, appSecret)
	if err != nil {
		log.NewHelper(logger).Error("initDbEnvDecryptByAes: %w", err)
		return nil, err
	}
	if len(dsn) == 0 {
		log.NewHelper(logger).Error("initDbEnvDecryptByAes: dsn is empty")
		return nil, err
	}

	if !strings.Contains(dsn, "parseTime=True") {
		dsn = dsn + "&parseTime=True"
	}
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
