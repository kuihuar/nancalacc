package conf

import (
	"fmt"
	"net"
	"net/url"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

// ConfigValidator 配置验证器
type ConfigValidator struct {
	logger log.Logger
}

// NewConfigValidator 创建新的配置验证器
func NewConfigValidator(l log.Logger) *ConfigValidator {
	return &ConfigValidator{
		logger: l,
	}
}

// ValidateBootstrap 验证Bootstrap配置
func (v *ConfigValidator) ValidateBootstrap(b *Bootstrap) error {
	if err := v.validateServer(b.Server); err != nil {
		return fmt.Errorf("server config validation failed: %w", err)
	}

	if err := v.validateData(b.Data); err != nil {
		return fmt.Errorf("data config validation failed: %w", err)
	}

	if err := v.validateApp(b.App); err != nil {
		return fmt.Errorf("app config validation failed: %w", err)
	}

	if err := v.validateAuth(b.Auth); err != nil {
		return fmt.Errorf("auth config validation failed: %w", err)
	}

	return nil
}

// validateServer 验证服务器配置
func (v *ConfigValidator) validateServer(s *Server) error {
	if s == nil {
		return fmt.Errorf("server config is nil")
	}

	// 验证HTTP配置
	if s.Http != nil {
		if err := v.validateHTTP(s.Http); err != nil {
			return fmt.Errorf("http config validation failed: %w", err)
		}
	}

	// 验证GRPC配置
	if s.Grpc != nil {
		if err := v.validateGRPC(s.Grpc); err != nil {
			return fmt.Errorf("grpc config validation failed: %w", err)
		}
	}

	return nil
}

// validateHTTP 验证HTTP配置
func (v *ConfigValidator) validateHTTP(h *Server_HTTP) error {
	if h.Addr == "" {
		return fmt.Errorf("http addr is empty")
	}

	if _, err := net.ResolveTCPAddr("tcp", h.Addr); err != nil {
		return fmt.Errorf("invalid http addr %s: %w", h.Addr, err)
	}

	if h.Timeout != nil {
		if h.Timeout.AsDuration() <= 0 {
			return fmt.Errorf("http timeout must be positive")
		}
	}

	return nil
}

// validateGRPC 验证GRPC配置
func (v *ConfigValidator) validateGRPC(g *Server_GRPC) error {
	if g.Addr == "" {
		return fmt.Errorf("grpc addr is empty")
	}

	if _, err := net.ResolveTCPAddr("tcp", g.Addr); err != nil {
		return fmt.Errorf("invalid grpc addr %s: %w", g.Addr, err)
	}

	if g.Timeout != nil {
		if g.Timeout.AsDuration() <= 0 {
			return fmt.Errorf("grpc timeout must be positive")
		}
	}

	return nil
}

// validateData 验证数据配置
func (v *ConfigValidator) validateData(d *Data) error {
	if d == nil {
		return fmt.Errorf("data config is nil")
	}

	// 验证数据库配置
	if d.Database != nil {
		if err := v.validateDatabase(d.Database); err != nil {
			return fmt.Errorf("database config validation failed: %w", err)
		}
	}

	// 验证同步数据库配置
	if d.DatabaseSync != nil {
		if err := v.validateDatabaseSync(d.DatabaseSync); err != nil {
			return fmt.Errorf("database sync config validation failed: %w", err)
		}
	}

	// 验证Redis配置
	if d.Redis != nil {
		if err := v.validateRedis(d.Redis); err != nil {
			return fmt.Errorf("redis config validation failed: %w", err)
		}
	}

	// 验证Etcd配置
	if d.Etcd != nil {
		if err := v.validateEtcd(d.Etcd); err != nil {
			return fmt.Errorf("etcd config validation failed: %w", err)
		}
	}

	return nil
}

// validateDatabase 验证数据库配置
func (v *ConfigValidator) validateDatabase(db *Data_Database) error {
	if db.Driver == "" {
		return fmt.Errorf("database driver is empty")
	}

	if db.Source == "" {
		return fmt.Errorf("database source is empty")
	}

	if db.MaxOpenConns <= 0 {
		return fmt.Errorf("database max_open_conns must be positive")
	}

	if db.MaxIdleConns <= 0 {
		return fmt.Errorf("database max_idle_conns must be positive")
	}

	if db.ConnMaxLifetime != "" {
		if _, err := time.ParseDuration(db.ConnMaxLifetime); err != nil {
			return fmt.Errorf("invalid database conn_max_lifetime %s: %w", db.ConnMaxLifetime, err)
		}
	}

	return nil
}

// validateDatabaseSync 验证同步数据库配置
func (v *ConfigValidator) validateDatabaseSync(db *Data_DatabaseSync) error {
	if db.Driver == "" {
		return fmt.Errorf("database sync driver is empty")
	}

	if db.Source == "" && db.SourceKey == "" {
		return fmt.Errorf("database sync source or source_key must be provided")
	}

	if db.MaxOpenConns <= 0 {
		return fmt.Errorf("database sync max_open_conns must be positive")
	}

	if db.MaxIdleConns <= 0 {
		return fmt.Errorf("database sync max_idle_conns must be positive")
	}

	if db.ConnMaxLifetime != "" {
		if _, err := time.ParseDuration(db.ConnMaxLifetime); err != nil {
			return fmt.Errorf("invalid database sync conn_max_lifetime %s: %w", db.ConnMaxLifetime, err)
		}
	}

	return nil
}

// validateRedis 验证Redis配置
func (v *ConfigValidator) validateRedis(r *Data_Redis) error {
	if r.Addr == "" {
		return fmt.Errorf("redis addr is empty")
	}

	if _, err := net.ResolveTCPAddr("tcp", r.Addr); err != nil {
		return fmt.Errorf("invalid redis addr %s: %w", r.Addr, err)
	}

	if r.Db < 0 {
		return fmt.Errorf("redis db must be non-negative")
	}

	if r.ReadTimeout != nil {
		if r.ReadTimeout.AsDuration() <= 0 {
			return fmt.Errorf("redis read_timeout must be positive")
		}
	}

	if r.WriteTimeout != nil {
		if r.WriteTimeout.AsDuration() <= 0 {
			return fmt.Errorf("redis write_timeout must be positive")
		}
	}

	return nil
}

// validateEtcd 验证Etcd配置
func (v *ConfigValidator) validateEtcd(e *Data_Etcd) error {
	if len(e.Endpoints) == 0 {
		return fmt.Errorf("etcd endpoints is empty")
	}

	for i, endpoint := range e.Endpoints {
		if endpoint == "" {
			return fmt.Errorf("etcd endpoint[%d] is empty", i)
		}

		if _, err := url.Parse(endpoint); err != nil {
			return fmt.Errorf("invalid etcd endpoint[%d] %s: %w", i, endpoint, err)
		}
	}

	if e.DialTimeout != "" {
		if _, err := time.ParseDuration(e.DialTimeout); err != nil {
			return fmt.Errorf("invalid etcd dial_timeout %s: %w", e.DialTimeout, err)
		}
	}

	if e.ConfigPrefix == "" {
		return fmt.Errorf("etcd config_prefix is empty")
	}

	return nil
}

// validateApp 验证应用配置
func (v *ConfigValidator) validateApp(a *App) error {
	if a == nil {
		return fmt.Errorf("app config is nil")
	}

	if a.Id == "" {
		return fmt.Errorf("app id is empty")
	}

	if a.Name == "" {
		return fmt.Errorf("app name is empty")
	}

	if a.Version == "" {
		return fmt.Errorf("app version is empty")
	}

	if a.Env == "" {
		return fmt.Errorf("app env is empty")
	}

	// LogLevel and LogOut are now in OpenTelemetry.Logs

	if a.AppPackage == "" {
		return fmt.Errorf("app app_package is empty")
	}

	if a.AppId == "" {
		return fmt.Errorf("app app_id is empty")
	}

	if a.AppSecret == "" {
		return fmt.Errorf("app app_secret is empty")
	}

	if a.ThirdCompanyId == "" {
		return fmt.Errorf("app third_company_id is empty")
	}

	if a.PlatformIds == "" {
		return fmt.Errorf("app platform_ids is empty")
	}

	if a.CompanyId == "" {
		return fmt.Errorf("app company_id is empty")
	}

	if a.AccessKey == "" {
		return fmt.Errorf("app access_key is empty")
	}

	if a.SecretKey == "" {
		return fmt.Errorf("app secret_key is empty")
	}

	return nil
}

// validateAuth 验证认证配置
func (v *ConfigValidator) validateAuth(a *Auth) error {
	if a == nil {
		return fmt.Errorf("auth config is nil")
	}

	// 验证WPS应用配置
	if a.Wpsapp != nil {
		if err := v.validateWpsapp(a.Wpsapp); err != nil {
			return fmt.Errorf("wpsapp config validation failed: %w", err)
		}
	}

	// 验证钉钉配置
	if a.Dingtalk != nil {
		if err := v.validateDingtalk(a.Dingtalk); err != nil {
			return fmt.Errorf("dingtalk config validation failed: %w", err)
		}
	}

	return nil
}

// validateWpsapp 验证WPS应用配置
func (v *ConfigValidator) validateWpsapp(w *Auth_Wpsapp) error {
	if w.ClientId == "" {
		return fmt.Errorf("wpsapp client_id is empty")
	}

	if w.ClientSecret == "" {
		return fmt.Errorf("wpsapp client_secret is empty")
	}

	if w.AuthUrl == "" {
		return fmt.Errorf("wpsapp auth_url is empty")
	}

	if _, err := url.Parse(w.AuthUrl); err != nil {
		return fmt.Errorf("invalid wpsapp auth_url %s: %w", w.AuthUrl, err)
	}

	if w.AuthPath == "" {
		return fmt.Errorf("wpsapp auth_path is empty")
	}

	if w.GrantType == "" {
		return fmt.Errorf("wpsapp grant_type is empty")
	}

	return nil
}

// validateDingtalk 验证钉钉配置
func (v *ConfigValidator) validateDingtalk(d *Auth_Dingtalk) error {
	if d.Endpoint == "" {
		return fmt.Errorf("dingtalk endpoint is empty")
	}

	if _, err := url.Parse(d.Endpoint); err != nil {
		return fmt.Errorf("invalid dingtalk endpoint %s: %w", d.Endpoint, err)
	}

	if d.AppKey == "" {
		return fmt.Errorf("dingtalk app_key is empty")
	}

	if d.AppSecret == "" {
		return fmt.Errorf("dingtalk app_secret is empty")
	}

	if d.Timeout != "" {
		if _, err := time.ParseDuration(d.Timeout); err != nil {
			return fmt.Errorf("invalid dingtalk timeout %s: %w", d.Timeout, err)
		}
	}

	if d.MaxConcurrent <= 0 {
		return fmt.Errorf("dingtalk maxConcurrent must be positive")
	}

	return nil
}

// ValidateEnvironment 验证环境变量配置
func (v *ConfigValidator) ValidateEnvironment() error {
	requiredEnvVars := []string{
		"ENCRYPTION_SALT",
	}

	optionalEnvVars := []string{
		"APP_UID",
		"ECIS_ECISACCOUNTSYNC_DB",
	}

	// 验证必需的环境变量
	for _, key := range requiredEnvVars {
		if _, err := GetEnv(key); err != nil {
			return fmt.Errorf("required environment variable %s not found: %w", key, err)
		}
	}

	// 记录可选的环境变量状态
	for _, key := range optionalEnvVars {
		if value, err := GetEnv(key); err != nil {
			v.logger.Log(log.LevelWarn, "optional environment variable not found", "key", key)
		} else {
			v.logger.Log(log.LevelInfo, "environment variable found", "key", key, "has_value", value != "")
		}
	}

	return nil
}
