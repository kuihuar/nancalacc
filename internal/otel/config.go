package otel

import (
	"context"
	stdlog "log"
	"time"

	"nancalacc/internal/conf"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/metric"
	traceNoop "go.opentelemetry.io/otel/sdk/trace"
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
	Otlp    OtlpConfig   `yaml:"otlp" json:"otlp"`
}

// JaegerConfig Jaeger 配置
type JaegerConfig struct {
	Enabled  bool   `yaml:"enabled" json:"enabled"`
	Endpoint string `yaml:"endpoint" json:"endpoint"`
}

// OtlpConfig OTLP 配置
type OtlpConfig struct {
	Enabled  bool   `yaml:"enabled" json:"enabled"`
	Endpoint string `yaml:"endpoint" json:"endpoint"`
	Timeout  int    `yaml:"timeout" json:"timeout"` // 超时时间（秒）
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
	Enabled        bool       `yaml:"enabled" json:"enabled"`
	Level          string     `yaml:"level" json:"level"`
	Format         string     `yaml:"format" json:"format"`
	Output         string     `yaml:"output" json:"output"`
	FilePath       string     `yaml:"file_path" json:"file_path"`
	MaxSize        int        `yaml:"max_size" json:"max_size"`
	MaxBackups     int        `yaml:"max_backups" json:"max_backups"`
	MaxAge         int        `yaml:"max_age" json:"max_age"`
	Compress       bool       `yaml:"compress" json:"compress"`
	Caller         bool       `yaml:"caller" json:"caller"`
	Stacktrace     bool       `yaml:"stacktrace" json:"stacktrace"`
	EscapeNewlines bool       `yaml:"escape_newlines" json:"escape_newlines"`
	Gorm           GormLogs   `yaml:"gorm" json:"gorm"`
	Loki           LokiConfig `yaml:"loki" json:"loki"`
	// Zap配置
	UseZap               bool   `yaml:"use_zap" json:"use_zap"`
	ZapDevelopment       bool   `yaml:"zap_development" json:"zap_development"`
	ZapDisableCaller     bool   `yaml:"zap_disable_caller" json:"zap_disable_caller"`
	ZapDisableStacktrace bool   `yaml:"zap_disable_stacktrace" json:"zap_disable_stacktrace"`
	ZapEncoding          string `yaml:"zap_encoding" json:"zap_encoding"`
	ZapTimeKey           string `yaml:"zap_time_key" json:"zap_time_key"`
	ZapLevelKey          string `yaml:"zap_level_key" json:"zap_level_key"`
	ZapNameKey           string `yaml:"zap_name_key" json:"zap_name_key"`
	ZapCallerKey         string `yaml:"zap_caller_key" json:"zap_caller_key"`
	ZapFunctionKey       string `yaml:"zap_function_key" json:"zap_function_key"`
	ZapMessageKey        string `yaml:"zap_message_key" json:"zap_message_key"`
	ZapStacktraceKey     string `yaml:"zap_stacktrace_key" json:"zap_stacktrace_key"`
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
			Otlp: OtlpConfig{
				Enabled:  true,
				Endpoint: "localhost:4317",
				Timeout:  30,
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
			// Zap默认配置
			UseZap:               true,   // 默认启用zap
			ZapDevelopment:       false,  // 默认生产模式
			ZapDisableCaller:     false,  // 默认启用调用者信息
			ZapDisableStacktrace: false,  // 默认启用堆栈跟踪
			ZapEncoding:          "json", // 默认JSON编码
			ZapTimeKey:           "timestamp",
			ZapLevelKey:          "level",
			ZapNameKey:           "logger",
			ZapCallerKey:         "caller",
			ZapFunctionKey:       "func",
			ZapMessageKey:        "message",
			ZapStacktraceKey:     "stacktrace",
		},
	}
}

// GetLogger 获取日志器
func (c *Config) GetLogger() log.Logger {
	if !c.Enabled || !c.Logs.Enabled {
		stdlog.Printf("🔍 [DEBUG] OpenTelemetry logs disabled, using noop logger")
		// 暂时返回 nil，避免导入问题
		return nil
	}

	stdlog.Printf("🔍 [DEBUG] OpenTelemetry logs enabled, attempting to create configured logger")
	// 创建基于配置的日志器
	return c.createConfiguredLogger()
}

// createConfiguredLogger 根据配置创建日志器
func (c *Config) createConfiguredLogger() log.Logger {
	stdlog.Printf("🔍 [DEBUG] Creating configured logger...")

	// 这里应该根据配置创建真正的日志器
	// 由于OpenTelemetry的日志API比较复杂，我们使用Kratos的标准日志器
	// 并在适配器中处理级别过滤

	// 暂时返回 nil，避免导入问题
	stdlog.Printf("🔍 [DEBUG] Returning nil logger (not yet implemented)")
	return nil
}

// GetTracer 获取追踪器
func (c *Config) GetTracer() trace.Tracer {
	if !c.Enabled || !c.Traces.Enabled {
		stdlog.Printf("🔍 [DEBUG] OpenTelemetry traces disabled, using noop tracer")
		return traceNoop.NewTracerProvider().Tracer("nancalacc")
	}

	stdlog.Printf("🔍 [DEBUG] OpenTelemetry traces enabled, attempting to create real tracer")

	// 尝试创建真正的追踪器
	if c.Traces.Jaeger.Enabled {
		stdlog.Printf("🔍 [DEBUG] Jaeger tracing enabled, endpoint: %s", c.Traces.Jaeger.Endpoint)
		// 创建真正的 Jaeger 导出器
		return c.createJaegerTracer()
	}

	// 如果没有配置 Jaeger，尝试使用 OTLP
	if c.Traces.Otlp.Enabled {
		stdlog.Printf("🔍 [DEBUG] OTLP tracing enabled, endpoint: %s", c.Traces.Otlp.Endpoint)
		// 创建真正的 OTLP 导出器
		return c.createOTLPTracer()
	}

	stdlog.Printf("🔍 [DEBUG] No tracing backend configured, using noop tracer")
	// 使用 noop 追踪器作为默认实现
	return traceNoop.NewTracerProvider().Tracer("nancalacc")
}

// GetMeter 获取指标器
func (c *Config) GetMeter() metric.Meter {
	if !c.Enabled || !c.Metrics.Enabled {
		// 暂时返回 nil，避免导入问题
		return nil
	}
	// 暂时返回 nil，避免导入问题
	return nil
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
				Otlp: OtlpConfig{
					Enabled:  bc.Otel.Traces.Otlp.Enabled,
					Endpoint: bc.Otel.Traces.Otlp.Endpoint,
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
				// Zap配置
				UseZap:               bc.Otel.Logs.UseZap,
				ZapDevelopment:       bc.Otel.Logs.ZapDevelopment,
				ZapDisableCaller:     bc.Otel.Logs.ZapDisableCaller,
				ZapDisableStacktrace: bc.Otel.Logs.ZapDisableStacktrace,
				ZapEncoding:          bc.Otel.Logs.ZapEncoding,
				ZapTimeKey:           bc.Otel.Logs.ZapTimeKey,
				ZapLevelKey:          bc.Otel.Logs.ZapLevelKey,
				ZapNameKey:           bc.Otel.Logs.ZapNameKey,
				ZapCallerKey:         bc.Otel.Logs.ZapCallerKey,
				ZapFunctionKey:       bc.Otel.Logs.ZapFunctionKey,
				ZapMessageKey:        bc.Otel.Logs.ZapMessageKey,
				ZapStacktraceKey:     bc.Otel.Logs.ZapStacktraceKey,
			},
		}
		return config
	}

	return DefaultConfig()
}

// NewConfigFromConf 从 conf.OpenTelemetry 创建 Config
func NewConfigFromConf(otelConf *conf.OpenTelemetry) *Config {
	if otelConf == nil {
		return DefaultConfig()
	}

	config := &Config{
		Enabled:        otelConf.Enabled,
		ServiceName:    otelConf.ServiceName,
		ServiceVersion: otelConf.ServiceVersion,
		Environment:    otelConf.Environment,
	}

	// 设置追踪配置
	if otelConf.Traces != nil {
		config.Traces.Enabled = otelConf.Traces.Enabled
		if otelConf.Traces.Jaeger != nil {
			config.Traces.Jaeger.Enabled = otelConf.Traces.Jaeger.Enabled
			config.Traces.Jaeger.Endpoint = otelConf.Traces.Jaeger.Endpoint
		}
		if otelConf.Traces.Otlp != nil {
			config.Traces.Otlp.Enabled = otelConf.Traces.Otlp.Enabled
			config.Traces.Otlp.Endpoint = otelConf.Traces.Otlp.Endpoint
		}
	}

	// 设置指标配置
	if otelConf.Metrics != nil {
		config.Metrics.Enabled = otelConf.Metrics.Enabled
		if otelConf.Metrics.Prometheus != nil {
			config.Metrics.Prometheus.Enabled = otelConf.Metrics.Prometheus.Enabled
			config.Metrics.Prometheus.Endpoint = otelConf.Metrics.Prometheus.Endpoint
			config.Metrics.Prometheus.Interval = otelConf.Metrics.Prometheus.Interval
		}
	}

	// 设置日志配置
	if otelConf.Logs != nil {
		config.Logs.Enabled = otelConf.Logs.Enabled
		config.Logs.Level = otelConf.Logs.Level
		config.Logs.Format = otelConf.Logs.Format
		config.Logs.Output = otelConf.Logs.Output
		config.Logs.FilePath = otelConf.Logs.FilePath
		if otelConf.Logs.Loki != nil {
			config.Logs.Loki.Enabled = otelConf.Logs.Loki.Enabled
			config.Logs.Loki.Endpoint = otelConf.Logs.Loki.Endpoint
		}
		// Zap配置
		config.Logs.UseZap = otelConf.Logs.UseZap
		config.Logs.ZapDevelopment = otelConf.Logs.ZapDevelopment
		config.Logs.ZapDisableCaller = otelConf.Logs.ZapDisableCaller
		config.Logs.ZapDisableStacktrace = otelConf.Logs.ZapDisableStacktrace
		config.Logs.ZapEncoding = otelConf.Logs.ZapEncoding
		config.Logs.ZapTimeKey = otelConf.Logs.ZapTimeKey
		config.Logs.ZapLevelKey = otelConf.Logs.ZapLevelKey
		config.Logs.ZapNameKey = otelConf.Logs.ZapNameKey
		config.Logs.ZapCallerKey = otelConf.Logs.ZapCallerKey
		config.Logs.ZapFunctionKey = otelConf.Logs.ZapFunctionKey
		config.Logs.ZapMessageKey = otelConf.Logs.ZapMessageKey
		config.Logs.ZapStacktraceKey = otelConf.Logs.ZapStacktraceKey
	}

	return config
}

// createJaegerTracer 创建 Jaeger 追踪器
func (c *Config) createJaegerTracer() trace.Tracer {
	stdlog.Printf("🔍 [DEBUG] Creating real Jaeger tracer with endpoint: %s", c.Traces.Jaeger.Endpoint)

	// 创建 Jaeger 导出器
	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(c.Traces.Jaeger.Endpoint)))
	if err != nil {
		stdlog.Printf("🔍 [ERROR] Failed to create Jaeger exporter: %v", err)
		stdlog.Printf("🔍 [DEBUG] Falling back to noop tracer")
		return traceNoop.NewTracerProvider().Tracer("nancalacc")
	}

	// 创建资源
	ctx := context.Background()
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(c.ServiceName),
			semconv.ServiceVersion(c.ServiceVersion),
			semconv.DeploymentEnvironment(c.Environment),
		),
	)
	if err != nil {
		stdlog.Printf("🔍 [ERROR] Failed to create resource: %v", err)
		stdlog.Printf("🔍 [DEBUG] Falling back to noop tracer")
		return traceNoop.NewTracerProvider().Tracer("nancalacc")
	}

	// 创建 TracerProvider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)

	// 设置全局 TracerProvider
	otel.SetTracerProvider(tp)

	// 创建 tracer
	tracer := tp.Tracer(c.ServiceName)

	stdlog.Printf("🔍 [DEBUG] Real Jaeger tracer created successfully")
	return tracer
}

// createOTLPTracer 创建 OTLP 追踪器
func (c *Config) createOTLPTracer() trace.Tracer {
	stdlog.Printf("🔍 [DEBUG] Creating real OTLP tracer with endpoint: %s", c.Traces.Otlp.Endpoint)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.Traces.Otlp.Timeout)*time.Second)
	defer cancel()

	// 创建 gRPC 连接
	conn, err := grpc.DialContext(ctx, c.Traces.Otlp.Endpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		stdlog.Printf("🔍 [ERROR] Failed to connect to OTLP endpoint: %v", err)
		stdlog.Printf("🔍 [DEBUG] Falling back to noop tracer")
		return trace.NewNoopTracerProvider().Tracer("nancalacc")
	}

	// 创建 OTLP 导出器
	exporter, err := otlptrace.New(ctx, otlptracegrpc.NewClient(otlptracegrpc.WithGRPCConn(conn)))
	if err != nil {
		stdlog.Printf("🔍 [ERROR] Failed to create OTLP exporter: %v", err)
		stdlog.Printf("🔍 [DEBUG] Falling back to noop tracer")
		return trace.NewNoopTracerProvider().Tracer("nancalacc")
	}

	// 创建资源
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(c.ServiceName),
			semconv.ServiceVersion(c.ServiceVersion),
			semconv.DeploymentEnvironment(c.Environment),
		),
	)
	if err != nil {
		stdlog.Printf("🔍 [ERROR] Failed to create resource: %v", err)
		stdlog.Printf("🔍 [DEBUG] Falling back to noop tracer")
		return trace.NewNoopTracerProvider().Tracer("nancalacc")
	}

	// 创建 TracerProvider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)

	// 设置全局 TracerProvider
	otel.SetTracerProvider(tp)

	// 创建 tracer
	tracer := tp.Tracer(c.ServiceName)

	stdlog.Printf("🔍 [DEBUG] Real OTLP tracer created successfully")
	return tracer
}
