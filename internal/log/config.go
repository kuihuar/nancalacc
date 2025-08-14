package log

import (
	"time"
)

// Config 日志配置
type Config struct {
	Level          string     `json:"level" yaml:"level"`
	Format         string     `json:"format" yaml:"format"`
	Output         string     `json:"output" yaml:"output"`
	FilePath       string     `json:"file_path" yaml:"file_path"`
	MaxSize        int        `json:"max_size" yaml:"max_size"`
	MaxBackups     int        `json:"max_backups" yaml:"max_backups"`
	MaxAge         int        `json:"max_age" yaml:"max_age"`
	Compress       bool       `json:"compress" yaml:"compress"`
	Caller         bool       `json:"caller" yaml:"caller"`
	Stacktrace     bool       `json:"stacktrace" yaml:"stacktrace"`
	EscapeNewlines bool       `json:"escape_newlines" yaml:"escape_newlines"` // 是否转义换行符
	Gorm           GormConfig `json:"gorm" yaml:"gorm"`
	Loki           LokiConfig `json:"loki" yaml:"loki"`
}

// GormConfig GORM日志配置
type GormConfig struct {
	SlowThreshold string `json:"slow_threshold" yaml:"slow_threshold"`
	LogLevel      string `json:"log_level" yaml:"log_level"`
}

// LokiConfig Loki配置
type LokiConfig struct {
	URL       string            `json:"url" yaml:"url"`
	Username  string            `json:"username" yaml:"username"`
	Password  string            `json:"password" yaml:"password"`
	TenantID  string            `json:"tenant_id" yaml:"tenant_id"`
	Labels    map[string]string `json:"labels" yaml:"labels"`
	Enable    bool              `json:"enable" yaml:"enable"`
	BatchSize int               `json:"batch_size" yaml:"batch_size"`
	BatchWait time.Duration     `json:"batch_wait" yaml:"batch_wait"`
	Timeout   time.Duration     `json:"timeout" yaml:"timeout"`
}

// DefaultConfig 默认配置
func DefaultConfig() *Config {
	return &Config{
		Level:          "info",
		Format:         "json",
		Output:         "stdout",
		FilePath:       "logs/app.log",
		MaxSize:        100,
		MaxBackups:     10,
		MaxAge:         30,
		Compress:       true,
		Caller:         true,
		Stacktrace:     false,
		EscapeNewlines: false,
		Gorm: GormConfig{
			SlowThreshold: "200ms",
			LogLevel:      "warn",
		},
		Loki: LokiConfig{
			URL:       "http://localhost:3100",
			Enable:    false,
			BatchSize: 100,
			BatchWait: time.Second,
			Timeout:   10 * time.Second,
			Labels: map[string]string{
				"service": "nancalacc",
			},
		},
	}
}
