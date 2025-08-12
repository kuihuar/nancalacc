package conf

import (
	"fmt"
	"log"
)

// ExampleUsage 展示如何正确使用配置系统
func ExampleUsage() {
	// 1. 加载应用配置（使用conf.Load()）
	// 这会加载配置文件和环境变量（KRATOS_前缀）
	bootstrap, err := Load("configs/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 2. 获取应用配置
	fmt.Printf("Server HTTP Address: %s\n", bootstrap.Server.Http.Addr)
	fmt.Printf("Database Source: %s\n", bootstrap.Data.Database.Source)

	// 3. 获取运行时环境变量（使用envvars.go）
	// 这些变量不需要映射到配置结构，只需要在运行时获取

	// 获取加密盐值
	salt, err := GetEnv("ENCRYPTION_SALT")
	if err != nil {
		log.Printf("Warning: ENCRYPTION_SALT not found: %v", err)
	} else {
		fmt.Printf("Encryption Salt: %s\n", salt)
	}

	// 获取带默认值的环境变量
	appUID := GetEnvWithDefault("APP_UID", "nancalacc-default")
	fmt.Printf("App UID: %s\n", appUID)

	// 获取普通环境变量
	dbSource, err := GetEnv("DATABASE_SOURCE")
	if err != nil {
		log.Printf("Warning: DATABASE_SOURCE not found: %v", err)
	} else {
		fmt.Printf("Database Source: %s\n", dbSource)
	}

	// 4. 验证必需的环境变量
	requiredVars := []string{"ENCRYPTION_SALT", "APP_SECRET"}
	if err := ValidateRequiredEnvVars(requiredVars); err != nil {
		log.Fatalf("Missing required environment variables: %v", err)
	}
}

// 环境变量使用指南：
//
// 1. 应用配置（使用conf.Load()）：
//    - 服务器配置、数据库配置等
//    - 只从配置文件加载，不支持环境变量覆盖
//    - 例如：configs/config.yaml
//
// 2. 运行时环境变量（使用envvars.go）：
//    - 加密密钥、API密钥等敏感信息
//    - 不需要映射到配置结构
//    - 例如：ENCRYPTION_SALT=my-secret-salt
//
// 3. 环境变量优先级：
//    - 系统环境变量 > 配置文件
//
// 4. 最佳实践：
//    - 敏感信息使用环境变量
//    - 应用配置使用配置文件
//    - 开发环境使用默认值
//    - 生产环境使用必需变量验证
