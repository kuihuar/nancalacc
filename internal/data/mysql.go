package data

import (
	"errors"
	"fmt"
	"nancalacc/internal/conf"
	"nancalacc/internal/pkg/cipherutil"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func NewMysqlDB(c *conf.Data, logger log.Logger) (*MainDB, error) {

	if !c.Database.Enable {
		logger.Log(log.LevelWarn, "database Enable", c.Database.Enable)
		return &MainDB{}, nil
	}
	db, err := gorm.Open(mysql.Open(c.Database.Source), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Info),
	})
	if err != nil {
		logger.Log(log.LevelError, "open mysql failed", err)
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(int(c.Database.MaxOpenConns))
	sqlDB.SetMaxIdleConns(int(c.Database.MaxIdleConns))
	duration, err := time.ParseDuration(c.GetDatabaseSync().GetConnMaxLifetime())
	if err != nil {
		return nil, err
	}
	sqlDB.SetConnMaxLifetime(duration)

	return &MainDB{DB: db}, nil
}

// deleteUser
// ECIS_ECISACCOUNTSYNC_DB=rESQpZX7v1v4YZletn9rCPJXehB9GFT/dzVFk3R99aMEjCKAG6w+vQKYdwFjEil8Lz4JLaZu8ziT1U3oHwv02MwjBfKed1/xNwhlGrt9jQ+zcQod+0W8QS5SNDr3InBlJzFYqPpAu9UkDpHsYsheVgZGouFz6qVKetVn17ZpdFZgnS2Ct8mJgYLFr3Sry9m8
func NewMysqlSyncDB(c *conf.Data, logger log.Logger) (*SyncDB, error) {
	var dsn string
	syncdb := c.GetDatabaseSync()
	if syncdb.GetEnv() == "dev" {
		dsn = syncdb.GetSource()
	} else {
		envkey := syncdb.GetSourceKey()
		encryptedDsn, err := conf.GetEnv(envkey)

		//logger.Log(log.LevelDebug, "envkey:", envkey, "encryptedDsn:", encryptedDsn, "err", err)
		if err != nil {
			logger.Log(log.LevelError, envkey, err)
			return nil, err
		}
		appSecret := conf.Get().GetApp().GetAppSecret()
		fmt.Printf("appSecret: %s\n", appSecret)
		dsn, err = cipherutil.DecryptByAes(encryptedDsn, appSecret)
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
	}
	// logger.Log(log.LevelInfo, "sync dsn", dsn)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Info),
	})
	if err != nil {
		logger.Log(log.LevelError, "open mysql failed", err)
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(int(syncdb.GetMaxOpenConns()))
	sqlDB.SetMaxIdleConns(int(syncdb.GetMaxIdleConns()))
	duration, err := time.ParseDuration(syncdb.GetConnMaxLifetime())
	if err != nil {
		return nil, err
	}

	sqlDB.SetConnMaxLifetime(duration)

	return &SyncDB{DB: db}, nil
}
