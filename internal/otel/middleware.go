package otel

import (
	"context"
	"net/http"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// HTTPMiddleware HTTP 中间件
type HTTPMiddleware struct {
	tracer trace.Tracer
	meter  metric.Meter
	logger log.Logger
}

// NewHTTPMiddleware 创建 HTTP 中间件
func NewHTTPMiddleware(tracer trace.Tracer, meter metric.Meter, logger log.Logger) *HTTPMiddleware {
	return &HTTPMiddleware{
		tracer: tracer,
		meter:  meter,
		logger: logger,
	}
}

// Wrap 包装 HTTP 处理器
func (m *HTTPMiddleware) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// 创建 span
		ctx, span := m.tracer.Start(r.Context(), "HTTP "+r.Method+" "+r.URL.Path)
		defer span.End()

		// 设置 span 属性
		span.SetAttributes(
			attribute.String("http.method", r.Method),
			attribute.String("http.url", r.URL.String()),
			attribute.String("http.user_agent", r.UserAgent()),
			attribute.String("http.remote_addr", r.RemoteAddr),
		)

		// 包装响应写入器以获取状态码
		wrappedWriter := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// 执行处理器
		next.ServeHTTP(wrappedWriter, r.WithContext(ctx))

		// 记录响应信息
		duration := time.Since(start)
		span.SetAttributes(
			attribute.Int("http.status_code", wrappedWriter.statusCode),
			attribute.String("http.duration", duration.String()),
		)

		// 记录指标
		if m.meter != nil {
			if counter, _ := m.meter.Int64Counter("http_requests_total"); counter != nil {
				counter.Add(ctx, 1,
					metric.WithAttributes(attribute.String("method", r.Method)),
					metric.WithAttributes(attribute.String("path", r.URL.Path)),
					metric.WithAttributes(attribute.Int("status_code", wrappedWriter.statusCode)),
				)
			}

			if histogram, _ := m.meter.Float64Histogram("http_request_duration_seconds"); histogram != nil {
				histogram.Record(ctx, duration.Seconds(),
					metric.WithAttributes(attribute.String("method", r.Method)),
					metric.WithAttributes(attribute.String("path", r.URL.Path)),
				)
			}
		}

		// 记录请求日志
		if m.logger != nil {
			record := log.Record{}
			record.SetTimestamp(time.Now())
			record.SetSeverity(log.SeverityInfo)
			record.SetBody(log.StringValue("HTTP Request"))
			record.AddAttributes(
				log.String("method", r.Method),
				log.String("url", r.URL.String()),
				log.String("user_agent", r.UserAgent()),
				log.String("remote_addr", r.RemoteAddr),
			)
			m.logger.Emit(ctx, record)
		}
	})
}

// responseWriter 包装的响应写入器
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

// GRPCUnaryInterceptor gRPC 一元拦截器
func GRPCUnaryInterceptor(tracer trace.Tracer, meter metric.Meter, logger log.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()

		// 创建 span
		ctx, span := tracer.Start(ctx, "gRPC "+info.FullMethod)
		defer span.End()

		// 设置 span 属性
		span.SetAttributes(
			attribute.String("grpc.method", info.FullMethod),
		)

		// 执行处理器
		resp, err := handler(ctx, req)

		// 记录响应信息
		duration := time.Since(start)
		span.SetAttributes(attribute.String("grpc.duration", duration.String()))

		// 处理错误
		if err != nil {
			st, _ := status.FromError(err)
			span.SetAttributes(
				attribute.String("grpc.status_code", st.Code().String()),
				attribute.String("grpc.error_message", err.Error()),
			)
			span.RecordError(err)
		} else {
			span.SetAttributes(attribute.String("grpc.status_code", codes.OK.String()))
		}

		// 记录指标
		if meter != nil {
			if counter, _ := meter.Int64Counter("grpc_requests_total"); counter != nil {
				counter.Add(ctx, 1,
					metric.WithAttributes(attribute.String("method", info.FullMethod)),
					metric.WithAttributes(attribute.String("status_code", getGRPCStatusCode(err))),
				)
			}

			if histogram, _ := meter.Float64Histogram("grpc_request_duration_seconds"); histogram != nil {
				histogram.Record(ctx, duration.Seconds(),
					metric.WithAttributes(attribute.String("method", info.FullMethod)),
				)
			}
		}

		// 记录日志
		if logger != nil {
			level := log.SeverityInfo
			if err != nil {
				level = log.SeverityError
			}

			record := log.Record{}
			record.SetTimestamp(time.Now())
			record.SetSeverity(level)
			record.SetBody(log.StringValue("gRPC " + info.FullMethod + " " + getGRPCStatusCode(err)))
			record.AddAttributes(
				log.String("grpc.method", info.FullMethod),
				log.String("grpc.status_code", getGRPCStatusCode(err)),
				log.String("grpc.duration", duration.String()),
			)
			logger.Emit(ctx, record)
		}

		return resp, err
	}
}

// getGRPCStatusCode 获取 gRPC 状态码
func getGRPCStatusCode(err error) string {
	if err == nil {
		return codes.OK.String()
	}

	if st, ok := status.FromError(err); ok {
		return st.Code().String()
	}

	return codes.Unknown.String()
}
