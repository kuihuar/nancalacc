# 配置管理类 (Configuration) - Configuration Management

## 概述
配置管理类用于管理系统配置，包括应用配置、数据库配置、第三方服务配置等，支持从环境变量、配置文件等多种方式加载配置。

## 分类

### 3.1 应用配置
```go
// 应用基础配置
type AppConfig struct {
    Name        string `mapstructure:"name"`
    Version     string `mapstructure:"version"`
    Environment string `mapstructure:"environment"`
    Port        int    `mapstructure:"port"`
    Host        string `mapstructure:"host"`
    Debug       bool   `mapstructure:"debug"`
}

// 日志配置
type LogConfig struct {
    Level      string `mapstructure:"level"`
    Format     string `mapstructure:"format"`
    Output     string `mapstructure:"output"`
    FilePath   string `mapstructure:"file_path"`
    MaxSize    int    `mapstructure:"max_size"`
    MaxBackups int    `mapstructure:"max_backups"`
    MaxAge     int    `mapstructure:"max_age"`
    Compress   bool   `mapstructure:"compress"`
}

// 数据库配置
type DatabaseConfig struct {
    Driver   string `mapstructure:"driver"`
    Host     string `mapstructure:"host"`
    Port     int    `mapstructure:"port"`
    Username string `mapstructure:"username"`
    Password string `mapstructure:"password"`
    Database string `mapstructure:"database"`
    Charset  string `mapstructure:"charset"`
    Timeout  int    `mapstructure:"timeout"`
    MaxIdle  int    `mapstructure:"max_idle"`
    MaxOpen  int    `mapstructure:"max_open"`
}
```

### 3.2 第三方服务配置
```go
// 钉钉配置
type DingtalkConfig struct {
    AppKey    string `mapstructure:"app_key"`
    AppSecret string `mapstructure:"app_secret"`
    AgentId   string `mapstructure:"agent_id"`
    BaseURL   string `mapstructure:"base_url"`
    Timeout   int    `mapstructure:"timeout"`
    Retry     int    `mapstructure:"retry"`
}

// Redis配置
type RedisConfig struct {
    Host     string `mapstructure:"host"`
    Port     int    `mapstructure:"port"`
    Password string `mapstructure:"password"`
    Database int    `mapstructure:"database"`
    PoolSize int    `mapstructure:"pool_size"`
    Timeout  int    `mapstructure:"timeout"`
}

// 消息队列配置
type MessageQueueConfig struct {
    Type     string `mapstructure:"type"`
    Host     string `mapstructure:"host"`
    Port     int    `mapstructure:"port"`
    Username string `mapstructure:"username"`
    Password string `mapstructure:"password"`
    VHost    string `mapstructure:"vhost"`
    Queue    string `mapstructure:"queue"`
}
```

### 3.3 业务配置
```go
// 同步任务配置
type SyncConfig struct {
    BatchSize     int           `mapstructure:"batch_size"`
    MaxConcurrent int           `mapstructure:"max_concurrent"`
    Timeout       time.Duration `mapstructure:"timeout"`
    RetryCount    int           `mapstructure:"retry_count"`
    RetryInterval time.Duration `mapstructure:"retry_interval"`
}

// 缓存配置
type CacheConfig struct {
    Type     string        `mapstructure:"type"`
    TTL      time.Duration `mapstructure:"ttl"`
    MaxSize  int           `mapstructure:"max_size"`
    Strategy string        `mapstructure:"strategy"`
}

// 监控配置
type MonitorConfig struct {
    Enabled    bool   `mapstructure:"enabled"`
    MetricsURL string `mapstructure:"metrics_url"`
    TraceURL   string `mapstructure:"trace_url"`
    LogURL     string `mapstructure:"log_url"`
}
```

### 3.4 配置加载器
```go
// 配置管理器
type ConfigManager struct {
    config *Config
    viper  *viper.Viper
}

// 主配置结构
type Config struct {
    App      AppConfig           `mapstructure:"app"`
    Log      LogConfig           `mapstructure:"log"`
    Database DatabaseConfig      `mapstructure:"database"`
    Dingtalk DingtalkConfig      `mapstructure:"dingtalk"`
    Redis    RedisConfig         `mapstructure:"redis"`
    MQ       MessageQueueConfig  `mapstructure:"mq"`
    Sync     SyncConfig          `mapstructure:"sync"`
    Cache    CacheConfig         `mapstructure:"cache"`
    Monitor  MonitorConfig       `mapstructure:"monitor"`
}

// 配置加载方法
func (cm *ConfigManager) Load(configPath string) error {
    cm.viper.SetConfigFile(configPath)
    cm.viper.AutomaticEnv()
    cm.viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
    
    if err := cm.viper.ReadInConfig(); err != nil {
        return fmt.Errorf("failed to read config file: %w", err)
    }
    
    if err := cm.viper.Unmarshal(&cm.config); err != nil {
        return fmt.Errorf("failed to unmarshal config: %w", err)
    }
    
    return nil
}

// 获取配置
func (cm *ConfigManager) Get() *Config {
    return cm.config
}
```

## 特点

1. **多源支持**: 支持文件、环境变量、命令行参数等多种配置源
2. **类型安全**: 使用结构体定义配置，提供类型安全
3. **热重载**: 支持配置热重载，无需重启应用
4. **验证支持**: 支持配置验证和默认值设置
5. **环境隔离**: 支持不同环境的配置隔离

## 使用场景

- 应用启动配置
- 运行时配置管理
- 环境特定配置
- 功能开关控制
- 性能调优参数

## 最佳实践

1. **分层配置**: 按功能模块分层组织配置
2. **环境变量**: 敏感信息使用环境变量
3. **默认值**: 为所有配置项提供合理的默认值
4. **验证**: 在启动时验证配置的有效性
5. **文档化**: 为配置项提供清晰的说明文档 