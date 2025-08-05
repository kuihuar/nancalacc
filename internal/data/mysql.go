package data

import (
	"errors"
	"nancalacc/internal/conf"
	"nancalacc/pkg/cipherutil"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func NewMysqlDB(c *conf.Data, logger log.Logger) (*MainDB, error) {
	db, err := gorm.Open(mysql.Open(c.Database.Source), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Info),
	})
	if err != nil {
		logger.Log(log.LevelError, "open mysql failed", err)
		return nil, nil
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return &MainDB{DB: db}, nil
}

func NewMysqlSyncDB(c *conf.Data, logger log.Logger) (*SyncDB, error) {
	if c.GetDatabaseSync().Env == "dev" {
		return newMysqlDBSyncTest(c, logger)
	}
	envkey := c.DatabaseSync.Source
	encryptedDsn, err := conf.GetEnv(envkey)

	logger.Log(log.LevelDebug, "envkey:", envkey, "encryptedDsn:", encryptedDsn, "err", err)
	if err != nil {
		logger.Log(log.LevelError, envkey, err)
		return nil, err
	}
	appSecret := c.Auth.AppSecret
	dsn, err := cipherutil.DecryptByAes(encryptedDsn, appSecret)
	if err != nil {
		logger.Log(log.LevelError, envkey, err)
		return nil, err
	}
	if len(dsn) == 0 {
		logger.Log(log.LevelError, "ECIS_ECISACCOUNTSYNC_DB len(dsn) == 0")
		return nil, errors.New("ECIS_ECISACCOUNTSYNC_DB len(dsn) == 0")
	}

	if !strings.Contains(dsn, "parseTime=True") {
		dsn = dsn + "&parseTime=True"
	}

	db, err := gorm.Open(mysql.Open(c.Database.Source), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Info),
	})
	if err != nil {
		logger.Log(log.LevelError, "open mysql failed", err)
		return nil, nil
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return &SyncDB{DB: db}, nil
}

func newMysqlDBSyncTest(c *conf.Data, logger log.Logger) (*SyncDB, error) {
	db, err := gorm.Open(mysql.Open(c.DatabaseSync.Source), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Info),
	})
	if err != nil {
		logger.Log(log.LevelError, "open mysql failed", err)
		return nil, nil
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(int(c.DatabaseSync.GetMaxOpenConns()))
	sqlDB.SetMaxIdleConns(int(c.DatabaseSync.GetMaxOpenConns()))
	sqlDB.SetConnMaxLifetime(time.Hour)

	return &SyncDB{DB: db}, nil
}
