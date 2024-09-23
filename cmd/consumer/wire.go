//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"context"

	"github.com/google/wire"

	kratos "github.com/go-kratos/kratos/v2"
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/neo532/apitool/transport/http/xhttp/client"

	"demo_message/internal/biz"
	"demo_message/internal/conf"
	"demo_message/internal/data"
	"demo_message/internal/server"
	"demo_message/internal/service/consumer"
)

// initApp init kratos application.
func initApp(
	context.Context,
	*conf.Bootstrap,
	client.Client,
	klog.Logger,
) (*kratos.App, func(), error) {
	panic(wire.Build(
		server.NewConsumerMessage,
		newApp,
		consumer.ProviderSet,

		biz.ProviderSet,
		data.ProviderSet,
	))
}
