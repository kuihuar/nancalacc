package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"nancalacc/internal/conf"
	"nancalacc/internal/service"
	"nancalacc/internal/task"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
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
	//serverJson, _ := json.Marshal(bc.Server)
	//fmt.Printf("key: %s\n, value: %s\n", "/configs/nancalacc/server.json", string(serverJson))
	//cfg.Watch(configSource, bc)

	stdLogger := log.NewStdLogger(os.Stdout)

	//fmt.Println(bc.App.GetLogLevel())
	// 创建级别过滤器
	//_ = log.NewFilter(stdLogger, log.FilterLevel(log.LevelDebug), log.FilterLevel(log.LevelError), log.FilterLevel(log.LevelFatal))

	id := bc.App.GetId()
	Name := bc.App.GetName()
	Version := bc.App.GetVersion()
	Env := bc.App.GetEnv()

	fmt.Println(Env)

	logger := log.With(stdLogger,
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", id,
		"service.name", Name,
		"service.version", Version,
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID(),
	)

	app, cleanup, err := wireApp(bc.Server, bc.Service, bc.Data, logger)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	//app.Use(middleware.LogMiddleware(log.LevelError))

	// ttc := tracer.NewTracerManager()
	// ttc.Init(Env, Name)
	// defer ttc.Shutdown()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}
