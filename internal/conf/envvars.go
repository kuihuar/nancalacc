package conf

import (
	sync "sync"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/env"
)

var (
	envConfig config.Config
	onceEnv   sync.Once
)

func GetEnv(key string) (string, error) {
	onceEnv.Do(func() {
		envSource := env.NewSource()
		envConfig = config.New(
			config.WithSource(envSource),
		)
		if err := envConfig.Load(); err != nil {
			panic(err) // 或者记录日志后退出
		}
		var data map[string]interface{}
		if err := envConfig.Scan(&data); err == nil {
			// fmt.Println("All Environment Variables:")
			// for k, v := range data {
			// 	fmt.Printf("%s=%v\n", k, v)
			// }
		}
	})
	return envConfig.Value(key).String()
}
