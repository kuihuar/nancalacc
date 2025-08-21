//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"nancalacc/internal/biz"
	"nancalacc/internal/conf"
	"nancalacc/internal/data"
	"nancalacc/internal/dingtalk"
	"nancalacc/internal/otel"
	"nancalacc/internal/repository"
	"nancalacc/internal/server"
	"nancalacc/internal/service"
	"nancalacc/internal/task"
	"nancalacc/internal/wps"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// provideLogger creates a logger from the integration
func provideLogger(integration *otel.Integration) log.Logger {
	return integration.CreateLogger()
}

// wireApp init kratos application with OpenTelemetry integration.
func wireApp(*conf.Server, *conf.Data, *conf.OpenTelemetry, *otel.Integration) (*kratos.App, func(), error) {
	panic(wire.Build(
		provideLogger,
		server.ProviderSet,
		data.ProviderSet,
		repository.ProviderSet,
		wps.WpsProviderSet,
		dingtalk.DingtalkProviderSet,
		biz.ProviderSet,
		service.ProviderSet,
		task.ProviderSet,
		newApp))
}
