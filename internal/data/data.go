package data

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	"gorm.io/gorm"
)

// ProviderSet is data providers.
// var ProviderSet = wire.NewSet(NewData, NewMysqlDB, NewMysqlSyncDB, NewRedisClient, NewAccounterRepo)
// var ProviderSet = wire.NewSet(NewData, NewMysqlDB, NewMysqlSyncDB, NewRedisClient, NewAccounterRepo)
var ProviderSet = wire.NewSet(
	NewMysqlDB,     // 会绑定到 MainDB
	NewMysqlDBSync, // 会绑定到 SyncDB
	NewRedisClient,
	NewAccounterRepo,
	NewData,
)

type (
	MainDB struct{ *gorm.DB } // 包装结构体
	SyncDB struct{ *gorm.DB }
)

type Data struct {
	db       *gorm.DB // 主数据库
	nancalDB *gorm.DB // 同步数据库
	redis    *redis.Client
}

func NewData(
	// mainDB *gorm.DB,
	// syncDB *gorm.DB,

	syncDB *SyncDB, // ← 使用类型别名
	mainDB *MainDB, // ← 使用类型别名
	redis *redis.Client, logger log.Logger) (*Data, func(), error) {
	return &Data{
			// db:       syncDB,
			// nancalDB: mainDB,
			db:       syncDB.DB, // 类型转换
			nancalDB: mainDB.DB, // 类型转换
			redis:    redis,
		}, func() {
			// 清理函数保持原有逻辑
			// if redis != nil {
			// 	_ = redis.Close()
			// }
			// if sqlDB, err := mainDB.DB(); err == nil {
			// 	_ = sqlDB.Close()
			// }
			// if sqlDB, err := nancalDB.DB(); err == nil {
			// 	_ = sqlDB.Close()
			// }
			cleanup(mainDB.DB, syncDB.DB, redis, logger)
		}, nil
}
func cleanup(mainDB *gorm.DB, syncDB *gorm.DB, redis *redis.Client, logger log.Logger) {
	if redis != nil {
		if err := redis.Close(); err != nil {
			logger.Log(log.LevelError, "msg", "failed to close redis", "error", err)
		}
	}

	if mainDB != nil {
		gormDB := (*gorm.DB)(mainDB)
		if sqlDB, err := gormDB.DB(); err == nil {
			if err := sqlDB.Close(); err != nil {
				logger.Log(log.LevelError, "msg", "failed to close main DB", "error", err)
			}
		}
	}

	if syncDB != nil {
		gormDB := (*gorm.DB)(syncDB)
		if sqlDB, err := gormDB.DB(); err == nil {
			if err := sqlDB.Close(); err != nil {
				logger.Log(log.LevelError, "msg", "failed to close sync DB", "error", err)
			}
		}
	}

	logger.Log(log.LevelInfo, "msg", "all database connections closed")
}

// func NewData1(c *conf.Data, logger log.Logger) (*Data, func(), error) {
// 	log.NewHelper(logger).Info("=====newData.c: %v\n", c)

// 	var db *gorm.DB
// 	var err error

// 	db, err = NewMysqlDB(c)
// 	if err != nil {
// 		log.NewHelper(logger).Error("NewData: init db env failed")
// 		return nil, nil, nil
// 	}
// 	rdb, err := NewRedisClient(c)
// 	if err != nil {
// 		log.NewHelper(logger).Error("NewRedisClient: init db env failed")
// 		return nil, nil, nil
// 	}
// 	cleanup := func() {
// 		log.NewHelper(logger).Info("closing the data resources")
// 		if rdb != nil {
// 			_ = rdb.Close()
// 		}
// 		sqlDB, _ := db.DB()
// 		err = sqlDB.Close()
// 		if err != nil {
// 			log.NewHelper(logger).Error(err)
// 		}
// 	}

// 	return &Data{
// 		db:    db,
// 		redis: rdb,
// 	}, cleanup, nil
// 	// tags := strings.Split(c.Database.Tag, ",")

// 	// if len(tags) > 0 {
// 	// 	for _, tag := range tags {
// 	// 		if tag == "migrate" {
// 	// 			if err := Migrate(db); err != nil {
// 	// 				return nil, cleanup, err
// 	// 			}
// 	// 		}
// 	// 	}
// 	// }

// }

// func initDbEnv(c *conf.Data, logger log.Logger) (*gorm.DB, error) {
// 	encryptedDsn, err := conf.GetEnv("ECIS_ECISACCOUNTSYNC_DB")

// 	log.NewHelper(logger).Info("initDbEnv: %s", encryptedDsn)
// 	if err != nil {
// 		log.NewHelper(logger).Error("initDbEnv: %w", err)
// 		return nil, err
// 	}
// 	appSecret := c.Auth.AppSecret
// 	dsn, err := cipherutil.DecryptByAes(encryptedDsn, appSecret)
// 	if err != nil {
// 		log.NewHelper(logger).Error("initDbEnvDecryptByAes: %w", err)
// 		return nil, err
// 	}
// 	if len(dsn) == 0 {
// 		log.NewHelper(logger).Error("initDbEnvDecryptByAes: dsn is empty")
// 		return nil, err
// 	}

// 	if !strings.Contains(dsn, "parseTime=True") {
// 		dsn = dsn + "&parseTime=True"
// 	}
// 	return gorm.Open(mysql.Open(dsn), &gorm.Config{
// 		Logger: gormlogger.Default.LogMode(gormlogger.Info),
// 	})

// }

// 保持现有MySQL初始化逻辑
// func NewMysqlDB(c *conf.Data) (*gorm.DB, error) {
// 	db, err := gorm.Open(mysql.Open(c.Database.Source), &gorm.Config{})
// 	if err != nil {
// 		return nil, err
// 	}

// 	sqlDB, err := db.DB()
// 	if err != nil {
// 		return nil, err
// 	}
// 	sqlDB.SetMaxOpenConns(10)
// 	sqlDB.SetMaxIdleConns(10)
// 	sqlDB.SetConnMaxLifetime(time.Hour)

// 	return db, nil
// }

// func NewRedisClient(c *conf.Data) (*redis.Client, error) {
// 	rdb := redis.NewClient(&redis.Options{
// 		Addr:     c.Redis.Addr,
// 		Password: c.Redis.Password,
// 		DB:       int(c.Redis.Db),
// 		// PoolSize: int(c.Redis.Pool_size),
// 	})

// 	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
// 		return nil, err
// 	}
// 	return rdb, nil
// }
