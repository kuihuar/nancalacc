//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"nancalacc/internal/biz"
	"nancalacc/internal/conf"
	"nancalacc/internal/data"
	"nancalacc/internal/dingtalk"
	"nancalacc/internal/server"
	"nancalacc/internal/service"
	"nancalacc/internal/task"
	"nancalacc/internal/wps"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(*conf.Server, *conf.App, *conf.Data, *conf.Auth, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(
		server.ProviderSet,
		data.ProviderSet,
		wps.WpsProviderSet,
		dingtalk.DingtalkProviderSet,
		biz.ProviderSet,
		service.ProviderSet,
		task.ProviderSet,
		newApp))
}
