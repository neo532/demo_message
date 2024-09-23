//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"context"

	kratos "github.com/go-kratos/kratos/v2"
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/neo532/apitool/transport/http/xhttp/client"

	"demo_message/cmd"
	"demo_message/internal/biz"
	"demo_message/internal/conf"
	"demo_message/internal/data"
	"demo_message/internal/server"
	"demo_message/internal/service/api"
)

// initApp init kratos application.
func initApp(
	context.Context,
	*conf.Bootstrap,
	client.Client,
	klog.Logger,
) (*kratos.App, func(), error) {
	panic(wire.Build(
		server.NewHTTPServer,
		server.InitHTTPRouter,

		newApp,
		api.ProviderSet,

		biz.ProviderSet,
		data.ProviderSet,
	))
}

func InitDemo(client.Client) (*biz.MessageUsecase, func(), error) {
	panic(wire.Build(
		cmd.InitUnitTestSet,

		biz.ProviderSet,
		data.ProviderSet,
	))
}
