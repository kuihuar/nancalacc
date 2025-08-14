package main

import (
	"context"
	"flag"
	"os"

	"nancalacc/internal/conf"
	"nancalacc/internal/service"
	"nancalacc/internal/task"
	"nancalacc/internal/tracer"

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

// newApp 创建应用实例
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
			logger.Log(log.LevelInfo, "msg", "starting application with database factory")
			cronService.Start()
			return nil
		}),
		kratos.AfterStop(func(ctx context.Context) error {
			logger.Log(log.LevelInfo, "msg", "stopping application with database factory")
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

	// 初始化OpenTelemetry追踪系统
	tracerManager := tracer.NewTracerManager()
	if err := tracerManager.Init(bc.GetApp().GetEnv(), bc.GetApp().GetName()); err != nil {
		panic("failed to init tracer: " + err.Error())
	}
	defer tracerManager.Shutdown()

	// 4. 创建Kratos日志适配器
	kratosLogger, err := nancalaccLog.NewLoggerFromBootstrap(bc)
	if err != nil {
		panic(err)
	}

	kratosLogger.Log(log.LevelInfo,
		"company:", bc.GetApp().GetCompanyId(),
		"third company:", bc.GetApp().GetThirdCompanyId(),
		"platform ids:", bc.GetApp().GetPlatformIds())

	app, cleanup, err := wireApp(bc.Server, bc.Data, bc.Tracing, kratosLogger)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}
