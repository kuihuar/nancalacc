package main

import (
	"context"
	"flag"
	"os"

	"nancalacc/internal/conf"
	"nancalacc/internal/service"
	"nancalacc/internal/task"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"

	nancalaccLog "nancalacc/internal/log"

	_ "go.uber.org/automaxprocs"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string
	// Version is the version of the compiled software.
	Version string

	Env string
	// flagconf is the config flag.
	flagconf string

	id, _ = os.Hostname()
)

func init() {
	flag.StringVar(&flagconf, "conf", "../../configs", "config path, eg: -conf config.yaml")
}

func newApp(logger log.Logger, gs *grpc.Server, hs *http.Server, cronService *task.CronService, eventService *service.DingTalkEventService) *kratos.App {

	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			gs,
			hs,
		),
		kratos.BeforeStart(func(ctx context.Context) error {
			cronService.Start()
			return nil
		}),
		kratos.AfterStop(func(ctx context.Context) error {
			cronService.Stop()
			return nil
		}),
		kratos.BeforeStart(func(ctx context.Context) error {
			eventService.Start()
			return nil
		}),
		kratos.AfterStop(func(ctx context.Context) error {
			eventService.Stop()
			return nil
		}),
	)
}

func main() {
	flag.Parse()
	var bc *conf.Bootstrap
	bc, err := conf.Load(flagconf)

	if err != nil {
		panic("failed to load config: " + err.Error())
	}

	// 4. 创建Kratos日志适配器
	kratosLogger, err := nancalaccLog.NewLoggerFromBootstrap(bc)
	if err != nil {
		panic(err)
	}

	kratosLogger.Log(log.LevelInfo,
		"company:", bc.GetApp().GetCompanyId(),
		"third company:", bc.GetApp().GetThirdCompanyId(),
		"platform ids:", bc.GetApp().GetPlatformIds())

	app, cleanup, err := wireApp(bc.Server, bc.App, bc.Data, bc.Auth, kratosLogger)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}

// func initLogger(bc *conf.Bootstrap) (log.Logger, error) {
// 	logConfig := &nancalaccLog.Config{
// 		Level:      bc.GetLogging().GetLevel(),
// 		Format:     bc.GetLogging().GetFormat(),
// 		Output:     bc.GetLogging().GetOutput(),
// 		FilePath:   bc.GetLogging().GetFilePath(),
// 		MaxSize:    int(bc.GetLogging().GetMaxSize()),
// 		MaxBackups: int(bc.GetLogging().GetMaxBackups()),
// 		MaxAge:     int(bc.GetLogging().GetMaxAge()),
// 		Compress:   bc.GetLogging().GetCompress(),
// 		Stacktrace: bc.GetLogging().GetStacktrace(),
// 		Loki: nancalaccLog.LokiConfig{
// 			Enable: bc.GetLogging().GetLoki().GetEnable(),
// 		},
// 	}
// 	customLogger, err := nancalaccLog.NewLogger(logConfig)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to create logger: %w", err)
// 	}

// 	// 添加基础字段
// 	//id := bc.App.GetId()
// 	Name := bc.App.GetName()
// 	//Version := bc.App.GetVersion()

// 	loggerWithFields := customLogger.WithFields(map[string]interface{}{
// 		//"caller": log.DefaultCaller,
// 		//"service.id":      id,
// 		"service.name": Name,
// 		//"service.version": Version,
// 		// trace.id 和 span.id 需要在有追踪上下文时动态添加
// 	})

// 	// 创建Kratos日志适配器
// 	kratosLogger := nancalaccLog.NewKratosLoggerAdapter(loggerWithFields)

// 	return kratosLogger, nil
// }

// func TracingMiddleware(logger log.Logger) middleware.Middleware {
//     return func(handler middleware.Handler) middleware.Handler {
//         return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
//             // 从上下文中获取追踪信息
//             if span := trace.SpanFromContext(ctx); span != nil {
//                 logger = logger.With(
//                     "trace.id", span.SpanContext().TraceID().String(),
//                     "span.id", span.SpanContext().SpanID().String(),
//                 )
//             }

//             return handler(ctx, req)
//         }
//     }
// }
