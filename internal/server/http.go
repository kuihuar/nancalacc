package server

import (
	"io"
	v1 "nancalacc/api/account/v1"
	"nancalacc/internal/conf"
	"nancalacc/internal/middleware"
	"nancalacc/internal/service"

	nethttp "net/http"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, accountService *service.AccountService, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			tracing.Server(),
			middleware.LoggingMiddleware(logger),
		),

		http.RequestDecoder(func(r *http.Request, v interface{}) error {
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
		}),
	}
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
	return srv
}
