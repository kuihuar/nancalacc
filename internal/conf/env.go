package conf

import (
	sync "sync"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/env"
)

var (
	envConfig config.Config
	once      sync.Once
)

func GetEnv(key string) (string, error) {
	once.Do(func() {
		envSource := env.NewSource()
		envConfig = config.New(
			config.WithSource(envSource),
		)
		if err := envConfig.Load(); err != nil {
			panic(err) // 或者记录日志后退出
		}
	})
	return envConfig.Value(key).String()
}
