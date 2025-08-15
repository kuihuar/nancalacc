package otel

import (
	"nancalacc/internal/conf"

	"go.opentelemetry.io/otel/log"
	logNoop "go.opentelemetry.io/otel/log/noop"
	"go.opentelemetry.io/otel/metric"
	metricNoop "go.opentelemetry.io/otel/metric/noop"
	"go.opentelemetry.io/otel/trace"
	traceNoop "go.opentelemetry.io/otel/trace/noop"
)

// Config OpenTelemetry 配置
type Config struct {
	Enabled        bool          `json:"enabled"`
	ServiceName    string        `json:"service_name"`
	ServiceVersion string        `json:"service_version"`
	Environment    string        `json:"environment"`
	Traces         TracesConfig  `json:"traces"`
	Metrics        MetricsConfig `json:"metrics"`
	Logs           LogsConfig    `json:"logs"`
}

// TracesConfig 追踪配置
type TracesConfig struct {
	Enabled bool         `yaml:"enabled" json:"enabled"`
	Jaeger  JaegerConfig `yaml:"jaeger" json:"jaeger"`
}

// JaegerConfig Jaeger 配置
type JaegerConfig struct {
	Enabled  bool   `yaml:"enabled" json:"enabled"`
	Endpoint string `yaml:"endpoint" json:"endpoint"`
}

// MetricsConfig 指标配置
type MetricsConfig struct {
	Enabled    bool             `yaml:"enabled" json:"enabled"`
	Prometheus PrometheusConfig `yaml:"prometheus" json:"prometheus"`
}

// PrometheusConfig Prometheus 配置
type PrometheusConfig struct {
	Enabled  bool   `yaml:"enabled" json:"enabled"`
	Endpoint string `yaml:"endpoint" json:"endpoint"`
	Interval string `yaml:"interval" json:"interval"`
}

// LogsConfig 日志配置
type LogsConfig struct {
	Enabled  bool       `yaml:"enabled" json:"enabled"`
	Level    string     `yaml:"level" json:"level"`
	Format   string     `yaml:"format" json:"format"`
	Output   string     `yaml:"output" json:"output"`
	FilePath string     `yaml:"file_path" json:"file_path"`
	Loki     LokiConfig `yaml:"loki" json:"loki"`
}

// LokiConfig Loki 配置
type LokiConfig struct {
	Enabled  bool   `yaml:"enabled" json:"enabled"`
	Endpoint string `yaml:"endpoint" json:"endpoint"`
}

// Logs 日志配置
type Logs struct {
	Enabled        bool     `json:"enabled"`
	Level          string   `json:"level"`
	Format         string   `json:"format"`
	Output         string   `json:"output"`
	FilePath       string   `json:"file_path"`
	MaxSize        int      `json:"max_size"`
	MaxBackups     int      `json:"max_backups"`
	MaxAge         int      `json:"max_age"`
	Compress       bool     `json:"compress"`
	Caller         bool     `json:"caller"`
	Stacktrace     bool     `json:"stacktrace"`
	EscapeNewlines bool     `json:"escape_newlines"`
	Gorm           GormLogs `json:"gorm"`
	Loki           LokiLogs `json:"loki"`
}

// GormLogs GORM日志配置
type GormLogs struct {
	SlowThreshold string `json:"slow_threshold"`
	LogLevel      string `json:"log_level"`
}

// LokiLogs Loki日志配置
type LokiLogs struct {
	Enabled  bool   `json:"enabled"`
	Endpoint string `json:"endpoint"`
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Enabled:        true,
		ServiceName:    "nancalacc",
		ServiceVersion: "1.0.0",
		Environment:    "development",
		Traces: TracesConfig{
			Enabled: true,
			Jaeger: JaegerConfig{
				Enabled:  true,
				Endpoint: "http://localhost:14268/api/traces",
			},
		},
		Metrics: MetricsConfig{
			Enabled: true,
			Prometheus: PrometheusConfig{
				Enabled:  true,
				Endpoint: "localhost:9090",
				Interval: "15s",
			},
		},
		Logs: LogsConfig{
			Enabled: true,
			Level:   "info",
			Format:  "json",
			Loki: LokiConfig{
				Enabled:  true,
				Endpoint: "http://localhost:3100/loki/api/v1/push",
			},
		},
	}
}

// GetLogger 获取日志器
func (c *Config) GetLogger() log.Logger {
	if !c.Enabled || !c.Logs.Enabled {
		return logNoop.NewLoggerProvider().Logger("nancalacc")
	}

	// 创建基于配置的日志器
	return c.createConfiguredLogger()
}

// createConfiguredLogger 根据配置创建日志器
func (c *Config) createConfiguredLogger() log.Logger {
	// 这里应该根据配置创建真正的日志器
	// 由于OpenTelemetry的日志API比较复杂，我们使用Kratos的标准日志器
	// 并在适配器中处理级别过滤

	// 暂时返回noop日志器，实际的日志处理在KratosLoggerAdapter中
	return logNoop.NewLoggerProvider().Logger("nancalacc")
}

// GetTracer 获取追踪器
func (c *Config) GetTracer() trace.Tracer {
	if !c.Enabled || !c.Traces.Enabled {
		return traceNoop.NewTracerProvider().Tracer("nancalacc")
	}
	// 使用 noop 追踪器作为默认实现
	return traceNoop.NewTracerProvider().Tracer("nancalacc")
}

// GetMeter 获取指标器
func (c *Config) GetMeter() metric.Meter {
	if !c.Enabled || !c.Metrics.Enabled {
		return metricNoop.NewMeterProvider().Meter("nancalacc")
	}
	// 使用 noop 指标器作为默认实现
	return metricNoop.NewMeterProvider().Meter("nancalacc")
}

// ConfigAdapter 配置适配器
type ConfigAdapter struct{}

// NewConfigAdapter 创建配置适配器
func NewConfigAdapter() *ConfigAdapter {
	return &ConfigAdapter{}
}

// FromBootstrap 从Bootstrap配置转换为OpenTelemetry配置
func (a *ConfigAdapter) FromBootstrap(bootstrap interface{}) *Config {
	// 使用反射获取Bootstrap配置
	// 这里简化处理，直接返回默认配置
	// 在实际使用中，应该根据Bootstrap配置进行转换

	// 尝试从Bootstrap中获取OpenTelemetry配置
	if bc, ok := bootstrap.(*conf.Bootstrap); ok && bc.Otel != nil {
		config := &Config{
			Enabled:        bc.Otel.Enabled,
			ServiceName:    bc.Otel.ServiceName,
			ServiceVersion: bc.Otel.ServiceVersion,
			Environment:    bc.Otel.Environment,
			Traces: TracesConfig{
				Enabled: bc.Otel.Traces.Enabled,
				Jaeger: JaegerConfig{
					Enabled:  bc.Otel.Traces.Jaeger.Enabled,
					Endpoint: bc.Otel.Traces.Jaeger.Endpoint,
				},
			},
			Metrics: MetricsConfig{
				Enabled: bc.Otel.Metrics.Enabled,
				Prometheus: PrometheusConfig{
					Enabled:  bc.Otel.Metrics.Prometheus.Enabled,
					Endpoint: bc.Otel.Metrics.Prometheus.Endpoint,
					Interval: bc.Otel.Metrics.Prometheus.Interval,
				},
			},
			Logs: LogsConfig{
				Enabled:  bc.Otel.Logs.Enabled,
				Level:    bc.Otel.Logs.Level,
				Format:   bc.Otel.Logs.Format,
				Output:   bc.Otel.Logs.Output,
				FilePath: bc.Otel.Logs.FilePath,
				Loki: LokiConfig{
					Enabled:  bc.Otel.Logs.Loki.Enabled,
					Endpoint: bc.Otel.Logs.Loki.Endpoint,
				},
			},
		}
		return config
	}

	return DefaultConfig()
}
