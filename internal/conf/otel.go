package conf

// OpenTelemetry 配置
type OpenTelemetry struct {
	// 是否启用OpenTelemetry
	Enabled bool `json:"enabled"`
	// 服务名称
	ServiceName string `json:"service_name"`
	// 服务版本
	ServiceVersion string `json:"service_version"`
	// 环境
	Environment string `json:"environment"`
	// Jaeger配置
	Jaeger *JaegerConfig `json:"jaeger"`
	// OTLP配置
	OTLP *OTLPConfig `json:"otlp"`
}

// Jaeger配置
type JaegerConfig struct {
	// Jaeger端点
	Endpoint string `json:"endpoint"`
	// 是否启用
	Enabled bool `json:"enabled"`
}

// OTLP配置
type OTLPConfig struct {
	// OTLP端点
	Endpoint string `json:"endpoint"`
	// 是否启用
	Enabled bool `json:"enabled"`
	// 超时时间（秒）
	Timeout int `json:"timeout"`
}
