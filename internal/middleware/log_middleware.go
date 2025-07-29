package middleware

// import (
// 	"context"
// 	"fmt"

// 	"github.com/go-kratos/kratos/v2/log"
// 	"github.com/go-kratos/kratos/v2/middleware"
// 	"github.com/go-kratos/kratos/v2/transport"
// )

// // TraceLogMiddleware 记录请求开始和结束的日志，包括返回结果
// func TraceLogMiddleware(logger log.Logger) middleware.Middleware {
// 	return func(handler middleware.Handler) middleware.Handler {
// 		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
// 			// 获取当前请求的元信息（如方法名、服务名）
// 			var operation string
// 			if tr, ok := transport.FromServerContext(ctx); ok {
// 				operation = tr.Operation() // 例如："/package.Service/Method"
// 			}

// 			// 记录请求开始（包含 trace.id）
// 			log.FromContext(ctx).Infow(
// 				"request started",
// 				"operation", operation,
// 				"req", fmt.Sprintf("%+v", req), // 记录请求参数
// 			)

// 			// 执行实际业务逻辑（调用结束后才会继续执行）
// 			reply, err = handler(ctx, req)

// 			// 调用结束后，记录返回结果（包含 trace.id）
// 			if err != nil {
// 				log.FromContext(ctx).Errorw(
// 					"request finished with error",
// 					"operation", operation,
// 					"err", err,
// 				)
// 			} else {
// 				log.FromContext(ctx).Infow(
// 					"request finished successfully",
// 					"operation", operation,
// 					"reply", fmt.Sprintf("%+v", reply), // 记录返回结果
// 				)
// 			}
// 			return reply, err
// 		}
// 	}
// }
