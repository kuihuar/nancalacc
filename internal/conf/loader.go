package conf

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"time"

	etcdConfig "github.com/go-kratos/kratos/contrib/config/etcd/v2"
	"gopkg.in/yaml.v3"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/env"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/go-cmp/cmp"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"
)

var (
	globalConf  *Bootstrap
	confLock    sync.RWMutex
	once        sync.Once
	etcdClient  *clientv3.Client // 全局ETCD客户端
	etcdSource  config.Source    // 全局ETCD配置源
	cancelWatch context.CancelFunc
)

// Load loads the configuration from file and etcd.
func Load(configPath string) (*Bootstrap, error) {
	var loadErr error
	once.Do(func() {
		// 1. 首先加载文件配置
		fileConf, err := loadFileConfig(configPath)
		if err != nil {
			loadErr = err
			return
		}

		// 2. 如果配置了ETCD，尝试加载ETCD配置
		if fileConf.Data != nil && fileConf.Data.Etcd != nil && fileConf.Data.Etcd.Enable && len(fileConf.Data.Etcd.Endpoints) > 0 {
			etcdConf, cli, src, err := loadEtcdConfig(fileConf)
			if err != nil {
				log.Warnf("Failed to load etcd config: %v, using file config only", err)
				setGlobalConfig(fileConf)
				return
			}

			// 保存ETCD客户端和source
			etcdClient = cli
			etcdSource = src

			// 3. 合并配置(ETCD配置优先)
			merged := mergeConfigs(fileConf, etcdConf)
			setGlobalConfig(merged)

			// 4. 如果启用配置监听，启动监听协程
			if fileConf.Data.Etcd.EnableConfigWatch {
				ctx, cancel := context.WithCancel(context.Background())
				cancelWatch = cancel
				go watchEtcdConfigChanges(ctx, src, merged)
			}
		} else {
			setGlobalConfig(fileConf)
		}
	})

	if loadErr != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", loadErr)
	}

	return Get(), nil
}

// Close 释放资源
func Close() {
	if cancelWatch != nil {
		cancelWatch()
	}
	if etcdClient != nil {
		etcdClient.Close()
	}
}

// loadFileConfig 从文件加载配置
func loadFileConfig(configPath string) (*Bootstrap, error) {
	fmt.Printf("load file config: %s\n", configPath)
	fileSource := file.NewSource(configPath)
	fileConf := config.New(
		config.WithSource(
			fileSource,
			env.NewSource("KRATOS_"), // 支持环境变量覆盖
		),
	)

	if err := fileConf.Load(); err != nil {
		return nil, fmt.Errorf("failed to load file config: %w", err)
	}

	var bc Bootstrap
	if err := fileConf.Scan(&bc); err != nil {
		return nil, fmt.Errorf("failed to scan file config: %w", err)
	}

	return &bc, nil
}

// loadEtcdConfig 从ETCD加载配置，返回配置、客户端和source
func loadEtcdConfig(baseConf *Bootstrap) (*Bootstrap, *clientv3.Client, config.Source, error) {
	fmt.Printf("load etcd config: %+vv\n", baseConf.Data.Etcd)
	// 1. 创建ETCD客户端
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   baseConf.Data.Etcd.Endpoints,
		DialTimeout: parseDurationToTime(baseConf.Data.Etcd.DialTimeout),
	})
	if err != nil {
		return nil, nil, nil, fmt.Errorf("create etcd client failed: %w", err)
	}
	// 2. 创建ETCD配置源
	source, err := etcdConfig.New(cli,
		etcdConfig.WithPath(baseConf.Data.Etcd.ConfigPrefix),
		etcdConfig.WithPrefix(true),
	)
	if err != nil {
		cli.Close() // 创建失败时关闭客户端
		return nil, nil, nil, fmt.Errorf("create etcd config source failed: %w", err)
	}

	// 3. 加载配置
	conf := config.New(config.WithSource(source))
	if err := conf.Load(); err != nil {
		cli.Close() // 加载失败时关闭客户端
		return nil, nil, nil, fmt.Errorf("load etcd config failed: %w", err)
	}

	var etcdBc Bootstrap
	if err := conf.Scan(&etcdBc); err != nil {
		cli.Close() // 解析失败时关闭客户端
		return nil, nil, nil, fmt.Errorf("scan etcd config failed: %w", err)
	}

	return &etcdBc, cli, source, nil
}

// watchEtcdConfigChanges 监听ETCD配置变更
func watchEtcdConfigChanges(ctx context.Context, source config.Source, baseConf *Bootstrap) {
	log.Info("watchEtcdConfigChanges watcher start")

	defer log.Info("[TEST] ETCD watcher exit")
	watcher, err := source.Watch()
	log.Infof("watcher: %+v, err: %+v", watcher, err)
	if err != nil {
		log.Errorf("Failed to create config watcher: %v", err)
		return
	}
	defer watcher.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Info("Config watcher stopped by context")
			return
		default:
			log.Debug("Waiting for next config change...")
			values, err := watcher.Next()
			if err != nil {
				log.Errorf("watchEtcdConfigChanges watcher err: %+v", err)
				if errors.Is(err, context.Canceled) {
					log.Info("Config watcher stopped normally")
					return
				}
				log.Errorf("Failed to watch next config: %v", err)
				// select {
				// case <-time.After(5 * time.Second): // 带退避的重试
				// case <-ctx.Done():
				// 	return
				// }
				continue
			}
			log.Debugf("[TEST] 收到变更事件:\nKey: %s\nValue: %s\nVersion: %d",
				string(values[0].Key),
				string(values[0].Value), 100)
			var etcdBc Bootstrap
			if err := unmarshalKeyValues(values, &etcdBc); err != nil {
				log.Errorf("Failed to scan changed config: %v", err)
				continue
			}

			// 合并配置并原子更新
			merged := mergeConfigs(baseConf, &etcdBc)
			setGlobalConfig(merged)
			log.Info("Configuration updated from etcd")
			log.Infof("Global config updated",
				"changes: %s\n", diffConfigs(baseConf, &etcdBc)) // 变更差异日志
		}
	}
}

func diffConfigs(old, new *Bootstrap) string {
	return cmp.Diff(old, new, protocmp.Transform())
}

// etcdctl put /configs/app.json '{"server":{"port":8080}}'
func unmarshalKeyValues(kvs []*config.KeyValue, target interface{}) error {
	if len(kvs) == 0 {
		return fmt.Errorf("empty config values")
	}
	ext := filepath.Ext(kvs[0].Key)        // 获取后缀如 ".json"
	format := strings.TrimPrefix(ext, ".") // 去掉点

	if len(format) == 0 {
		return fmt.Errorf("unsupported format: %s", format)
	}
	// format := strings.ToLower(kvs[0].Format)
	data := kvs[0].Value
	switch format {
	case "json":
		if err := json.Unmarshal(data, target); err != nil {
			return fmt.Errorf("json unmarshal failed: %w", err)
		}
	case "yaml":
		if err := yaml.Unmarshal(data, target); err != nil {
			return fmt.Errorf("yaml unmarshal failed: %w", err)
		}
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
	return nil
}

// setGlobalConfig 线程安全地设置全局配置
func setGlobalConfig(conf *Bootstrap) {
	confLock.Lock()
	defer confLock.Unlock()
	globalConf = conf
}

// Get 线程安全地获取全局配置
func Get() *Bootstrap {
	confLock.RLock()
	defer confLock.RUnlock()
	return globalConf
}

// parseDurationToTime 将字符串持续时间转换为time.Duration
func parseDurationToTime(durStr string) time.Duration {
	if durStr == "" {
		return 5 * time.Second // 默认超时时间
	}
	d, err := time.ParseDuration(durStr)
	if err != nil {
		return 5 * time.Second // 默认超时时间
	}
	return d
}

// mergeConfigs merges two configurations with etcdConfig taking precedence
func mergeConfigs(fileConfig, etcdConfig *Bootstrap) *Bootstrap {
	// Create a deep copy of the file config
	merged := proto.Clone(fileConfig).(*Bootstrap)

	// Merge logic for each section
	if etcdConfig.Server != nil {
		if merged.Server == nil {
			merged.Server = &Server{}
		}
		mergeServer(merged.Server, etcdConfig.Server)
	}

	if etcdConfig.Data != nil {
		if merged.Data == nil {
			merged.Data = &Data{}
		}
		mergeData(merged.Data, etcdConfig.Data)
	}

	if etcdConfig.App != nil {
		if merged.App == nil {
			merged.App = &App{}
		}
		mergeApp(merged.App, etcdConfig.App)
	}

	if etcdConfig.Service != nil {
		if merged.Service == nil {
			merged.Service = &Service{}
		}
		mergeService(merged.Service, etcdConfig.Service)
	}

	return merged
}

// Helper functions to merge each configuration section
func mergeServer(dst, src *Server) {
	if src.Http != nil {
		if dst.Http == nil {
			dst.Http = &Server_HTTP{}
		}
		if src.Http.Network != "" {
			dst.Http.Network = src.Http.Network
		}
		if src.Http.Addr != "" {
			dst.Http.Addr = src.Http.Addr
		}
		if src.Http.Timeout != nil {
			dst.Http.Timeout = src.Http.Timeout
		}
	}
	if src.Grpc != nil {
		if dst.Grpc == nil {
			dst.Grpc = &Server_GRPC{}
		}
		if src.Grpc.Network != "" {
			dst.Grpc.Network = src.Grpc.Network
		}
		if src.Grpc.Addr != "" {
			dst.Grpc.Addr = src.Grpc.Addr
		}
		if src.Grpc.Timeout != nil {
			dst.Grpc.Timeout = src.Grpc.Timeout
		}
	}
}

func mergeData(dst, src *Data) {
	if src.Database != nil {
		if dst.Database == nil {
			dst.Database = &Data_Database{}
		}
		if src.Database.Driver != "" {
			dst.Database.Driver = src.Database.Driver
		}
		if src.Database.Source != "" {
			dst.Database.Source = src.Database.Source
		}
		if src.Database.Tag != "" {
			dst.Database.Tag = src.Database.Tag
		}
	}
	if src.DatabaseSync != nil {
		if dst.DatabaseSync == nil {
			dst.DatabaseSync = &Data_DatabaseSync{}
		}
		if src.DatabaseSync.Driver != "" {
			dst.DatabaseSync.Driver = src.DatabaseSync.Driver
		}
		if src.DatabaseSync.Source != "" {
			dst.DatabaseSync.Source = src.DatabaseSync.Source
		}
		if src.DatabaseSync.Tag != "" {
			dst.DatabaseSync.Tag = src.DatabaseSync.Tag
		}
		if src.DatabaseSync.MaxOpenConns != 0 {
			dst.DatabaseSync.MaxOpenConns = src.DatabaseSync.MaxOpenConns
		}
		if src.DatabaseSync.MaxIdleConns != 0 {
			dst.DatabaseSync.MaxIdleConns = src.DatabaseSync.MaxIdleConns
		}
		if src.DatabaseSync.SourceKey != "" {
			dst.DatabaseSync.SourceKey = src.DatabaseSync.SourceKey
		}
	}
	if src.Redis != nil {
		if dst.Redis == nil {
			dst.Redis = &Data_Redis{}
		}
		if src.Redis.Network != "" {
			dst.Redis.Network = src.Redis.Network
		}
		if src.Redis.Addr != "" {
			dst.Redis.Addr = src.Redis.Addr
		}
		if src.Redis.Password != "" {
			dst.Redis.Password = src.Redis.Password
		}
		if src.Redis.Db != 0 {
			dst.Redis.Db = src.Redis.Db
		}
		if src.Redis.ReadTimeout != nil {
			dst.Redis.ReadTimeout = src.Redis.ReadTimeout
		}
		if src.Redis.WriteTimeout != nil {
			dst.Redis.WriteTimeout = src.Redis.WriteTimeout
		}
	}
	if src.Etcd != nil {
		if dst.Etcd == nil {
			dst.Etcd = &Data_Etcd{}
		}
		if len(src.Etcd.Endpoints) > 0 {
			dst.Etcd.Endpoints = src.Etcd.Endpoints
		}
		if src.Etcd.DialTimeout != "" {
			dst.Etcd.DialTimeout = src.Etcd.DialTimeout
		}
		if src.Etcd.ConfigPrefix != "" {
			dst.Etcd.ConfigPrefix = src.Etcd.ConfigPrefix
		}
		dst.Etcd.EnableConfigWatch = src.Etcd.EnableConfigWatch
	}
	if src.Auth != nil {
		if dst.Auth == nil {
			dst.Auth = &Data_Auth{}
		}
		if src.Auth.AppId != "" {
			dst.Auth.AppId = src.Auth.AppId
		}
		if src.Auth.AppSecret != "" {
			dst.Auth.AppSecret = src.Auth.AppSecret
		}
	}
}

func mergeApp(dst, src *App) {
	if src.Id != "" {
		dst.Id = src.Id
	}
	if src.Name != "" {
		dst.Name = src.Name
	}
	if src.Version != "" {
		dst.Version = src.Version
	}
	if src.Env != "" {
		dst.Env = src.Env
	}
	if src.LogLevel != "" {
		dst.LogLevel = src.LogLevel
	}
	if src.LogOut != "" {
		dst.LogOut = src.LogOut
	}
}

func mergeService(dst, src *Service) {
	if src.Business != nil {
		if dst.Business == nil {
			dst.Business = &Service_Business{}
		}
		if src.Business.ThirdCompanyId != "" {
			dst.Business.ThirdCompanyId = src.Business.ThirdCompanyId
		}
		if src.Business.PlatformIds != "" {
			dst.Business.PlatformIds = src.Business.PlatformIds
		}
		if src.Business.CompanyId != "" {
			dst.Business.CompanyId = src.Business.CompanyId
		}
		if src.Business.EcisaccountsyncUrl != "" {
			dst.Business.EcisaccountsyncUrl = src.Business.EcisaccountsyncUrl
		}
		if src.Business.EcisaccountsyncUrlIncrement != "" {
			dst.Business.EcisaccountsyncUrlIncrement = src.Business.EcisaccountsyncUrlIncrement
		}
	}
	if src.Auth != nil {
		if dst.Auth == nil {
			dst.Auth = &Service_Auth{}
		}
		mergeServiceAuth(dst.Auth, src.Auth)
	}
}

func mergeServiceAuth(dst, src *Service_Auth) {
	if src.Self != nil {
		if dst.Self == nil {
			dst.Self = &Service_Auth_Self{}
		}
		if src.Self.AppPackage != "" {
			dst.Self.AppPackage = src.Self.AppPackage
		}
		if src.Self.AppSecret != "" {
			dst.Self.AppSecret = src.Self.AppSecret
		}
		if src.Self.AccessKey != "" {
			dst.Self.AccessKey = src.Self.AccessKey
		}
		if src.Self.SecretKey != "" {
			dst.Self.SecretKey = src.Self.SecretKey
		}
	}
	if src.App != nil {
		if dst.App == nil {
			dst.App = &Service_Auth_App{}
		}
		if src.App.ClientId != "" {
			dst.App.ClientId = src.App.ClientId
		}
		if src.App.ClientSecret != "" {
			dst.App.ClientSecret = src.App.ClientSecret
		}
		if src.App.AuthUrl != "" {
			dst.App.AuthUrl = src.App.AuthUrl
		}
		if src.App.AuthPath != "" {
			dst.App.AuthPath = src.App.AuthPath
		}
		if src.App.GrantType != "" {
			dst.App.GrantType = src.App.GrantType
		}
	}
	if src.Third != nil {
		if dst.Third == nil {
			dst.Third = &Service_Auth_Third{}
		}
		if src.Third.ClientId != "" {
			dst.Third.ClientId = src.Third.ClientId
		}
		if src.Third.ClientSecret != "" {
			dst.Third.ClientSecret = src.Third.ClientSecret
		}
		if src.Third.AuthUrl != "" {
			dst.Third.AuthUrl = src.Third.AuthUrl
		}
		if src.Third.AuthPath != "" {
			dst.Third.AuthPath = src.Third.AuthPath
		}
		if src.Third.GrantType != "" {
			dst.Third.GrantType = src.Third.GrantType
		}
		if src.Third.CompanyId != "" {
			dst.Third.CompanyId = src.Third.CompanyId
		}
	}
	if src.User != nil {
		if dst.User == nil {
			dst.User = &Service_Auth_User{}
		}
		if src.User.ClientId != "" {
			dst.User.ClientId = src.User.ClientId
		}
		if src.User.ClientSecret != "" {
			dst.User.ClientSecret = src.User.ClientSecret
		}
		if src.User.AuthUrl != "" {
			dst.User.AuthUrl = src.User.AuthUrl
		}
		if src.User.AuthPath != "" {
			dst.User.AuthPath = src.User.AuthPath
		}
		if src.User.GrantType != "" {
			dst.User.GrantType = src.User.GrantType
		}
		if src.User.RedirectUri != "" {
			dst.User.RedirectUri = src.User.RedirectUri
		}
	}
	if src.Dingtalk != nil {
		if dst.Dingtalk == nil {
			dst.Dingtalk = &Service_Auth_Dingtalk{}
		}
		if src.Dingtalk.Endpoint != "" {
			dst.Dingtalk.Endpoint = src.Dingtalk.Endpoint
		}
		if src.Dingtalk.AppKey != "" {
			dst.Dingtalk.AppKey = src.Dingtalk.AppKey
		}
		if src.Dingtalk.AppSecret != "" {
			dst.Dingtalk.AppSecret = src.Dingtalk.AppSecret
		}
		if src.Dingtalk.Timeout != "" {
			dst.Dingtalk.Timeout = src.Dingtalk.Timeout
		}
		if src.Dingtalk.MaxConcurrent != 0 {
			dst.Dingtalk.MaxConcurrent = src.Dingtalk.MaxConcurrent
		}
	}
}
