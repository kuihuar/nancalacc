package main

import (
	"context"
	"flag"
	"os"

	"nancalacc/internal/conf"
	"nancalacc/internal/otel"
	"nancalacc/internal/service"
	"nancalacc/internal/task"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"

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
func newApp(integration *otel.Integration, gs *grpc.Server, hs *http.Server, cronService *task.CronService, eventService *service.DingTalkEventService) *kratos.App {
	// 从 integration 获取 Kratos 兼容的 logger
	logger := integration.CreateLogger()

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
			logger.Log(log.LevelInfo, "msg", "starting application with OpenTelemetry")
			cronService.Start()
			return nil
		}),
		kratos.AfterStop(func(ctx context.Context) error {
			logger.Log(log.LevelInfo, "msg", "stopping application with OpenTelemetry")
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

	// 加载配置
	bc, err := conf.Load(flagconf)
	if err != nil {
		panic("failed to load config: " + err.Error())
	}

	// 初始化OpenTelemetry
	otelIntegration := initOpenTelemetry(bc)
	defer otelIntegration.Shutdown(context.Background())

	// 创建日志器
	otelIntegration.CreateLogger().Log(log.LevelInfo,
		"msg", "app start",
		"company:", bc.GetApp().GetCompanyId(),
		"third company:", bc.GetApp().GetThirdCompanyId(),
		"platform ids:", bc.GetApp().GetPlatformIds())
	// 创建应用
	app, cleanup, err := wireApp(bc.Server, bc.Data, bc.Otel, otelIntegration)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// 启动应用
	if err := app.Run(); err != nil {
		panic(err)
	}
}

// initOpenTelemetry 初始化OpenTelemetry
func initOpenTelemetry(bc *conf.Bootstrap) *otel.Integration {
	// 创建配置适配器
	adapter := otel.NewConfigAdapter()
	config := adapter.FromBootstrap(bc)

	// 创建集成器
	integration := otel.NewIntegration(config)

	// 初始化
	if err := integration.Init(context.Background()); err != nil {
		panic("failed to init OpenTelemetry: " + err.Error())
	}

	return integration
}
