package conf

import (
	"fmt"
	"os"
	"sync"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/env"
	"github.com/go-kratos/kratos/v2/log"
)

var (
	envConfig config.Config
	onceEnv   sync.Once
	logger    log.Logger
)

// InitEnvConfig 初始化环境变量配置
func InitEnvConfig(l log.Logger) {
	logger = l
}

// GetEnv 获取环境变量值，支持默认值
func GetEnv(key string) (string, error) {
	// 首先尝试从系统环境变量获取
	if value := os.Getenv(key); value != "" {
		return value, nil
	}

	// 如果系统环境变量不存在，尝试从配置的环境变量源获取
	onceEnv.Do(func() {
		envSource := env.NewSource()
		envConfig = config.New(
			config.WithSource(envSource),
		)
		if err := envConfig.Load(); err != nil {
			if logger != nil {
				logger.Log(log.LevelError, "failed to load env config", "error", err)
			}
			return
		}
	})

	if envConfig == nil {
		return "", fmt.Errorf("environment variable %s not found", key)
	}

	value := envConfig.Value(key)
	if value == nil {
		return "", fmt.Errorf("environment variable %s not found", key)
	}

	str, err := value.String()
	if err != nil {
		return "", fmt.Errorf("failed to get string value for %s: %w", key, err)
	}
	return str, nil
}

// GetEnvWithDefault 获取环境变量值，如果不存在则返回默认值
func GetEnvWithDefault(key, defaultValue string) string {
	value, err := GetEnv(key)
	if err != nil {
		if logger != nil {
			logger.Log(log.LevelWarn, "environment variable not found, using default",
				"key", key, "default", defaultValue, "error", err)
		}
		return defaultValue
	}
	return value
}

// MustGetEnv 获取环境变量值，如果不存在则panic
func MustGetEnv(key string) string {
	value, err := GetEnv(key)
	if err != nil {
		panic(fmt.Sprintf("required environment variable %s not found: %v", key, err))
	}
	return value
}

// ValidateRequiredEnvVars 验证必需的环境变量是否存在
func ValidateRequiredEnvVars(requiredKeys []string) error {
	var missingKeys []string

	for _, key := range requiredKeys {
		if _, err := GetEnv(key); err != nil {
			missingKeys = append(missingKeys, key)
		}
	}

	if len(missingKeys) > 0 {
		return fmt.Errorf("missing required environment variables: %v", missingKeys)
	}

	return nil
}
