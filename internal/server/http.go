package server

import (
	"io"
	v1 "nancalacc/api/account/v1"
	"nancalacc/internal/conf"
	"nancalacc/internal/otel"
	"nancalacc/internal/service"

	nethttp "net/http"

	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new an HTTP server with OpenTelemetry integration.
func NewHTTPServer(c *conf.Server, accountService *service.AccountService, integration *otel.Integration, otelConfig *conf.OpenTelemetry) *http.Server {
	// 创建中间件选项
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
		),
	}

	// 添加OpenTelemetry中间件
	if integration != nil {
		otelOpts := integration.CreateHTTPMiddleware()
		opts = append(opts, otelOpts...)
	}

	// 添加请求解码器
	opts = append(opts, http.RequestDecoder(func(r *http.Request, v interface{}) error {
		if r.URL.Path == "/v1/upload" {
			data, err := io.ReadAll(r.Body)
			if err != nil {
				return err
			}
			if req, ok := v.(*v1.UploadRequest); ok {
				req.File = data
				if filename := r.Header.Get("X-Filename"); filename != "" {
					req.Filename = filename
				}
				return nil
			}
		}
		return http.DefaultRequestDecoder(r, v)
	}))

	// 添加服务器配置
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}

	srv := http.NewServer(opts...)
	v1.RegisterAccountHTTPServer(srv, accountService)

	// 添加健康检查路由
	srv.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(nethttp.StatusOK)
		w.Write([]byte("OK"))
	})

	// 添加指标端点（如果启用）
	if otelConfig != nil && otelConfig.Metrics != nil && otelConfig.Metrics.Enabled {
		srv.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
			// 这里可以添加Prometheus指标端点
			w.WriteHeader(nethttp.StatusOK)
			w.Write([]byte("# OpenTelemetry metrics endpoint"))
		})
	}

	return srv
}
