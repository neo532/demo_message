// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"context"
	"demo_message/cmd"
	"demo_message/internal/biz"
	"demo_message/internal/conf"
	"demo_message/internal/data"
	"demo_message/internal/server"
	"demo_message/internal/service/api"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/neo532/apitool/transport/http/xhttp/client"
)

// Injectors from wire.go:

// initApp init kratos application.
func initApp(contextContext context.Context, bootstrap *conf.Bootstrap, clientClient client.Client, logger log.Logger) (*kratos.App, func(), error) {
	httpServer := server.NewHTTPServer(bootstrap, logger)
	databaseMessage, cleanup, err := data.NewDatabaseMessage(contextContext, bootstrap, logger)
	if err != nil {
		return nil, nil, err
	}
	transactionMessageRepo := data.NewTransactionMessageRepo(databaseMessage, logger)
	campaignRepo := data.NewCampaignRepo(databaseMessage)
	campaignUsecase := biz.NewCampaignUsecase(transactionMessageRepo, campaignRepo)
	producerMessage, cleanup2, err := data.NewProducerMessage(contextContext, bootstrap, logger)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	messageXHttpClient := data.NewMessageXHttpClient(clientClient, bootstrap)
	messageRepo := data.NewMessageRepo(databaseMessage, producerMessage, messageXHttpClient)
	recipientRepo := data.NewRecipientRepo(databaseMessage)
	messageUsecase := biz.NewMessageUsecase(transactionMessageRepo, messageRepo, recipientRepo)
	campaignApi := api.NewCampaignApi(campaignUsecase, messageUsecase, logger)
	messageApi := api.NewMessageApi(messageUsecase, logger)
	systemApi := api.NewSystemApi(logger)
	router := server.InitHTTPRouter(httpServer, campaignApi, messageApi, systemApi)
	app := newApp(contextContext, bootstrap, httpServer, router, logger)
	return app, func() {
		cleanup2()
		cleanup()
	}, nil
}

func InitDemo(clientClient client.Client) (*biz.MessageUsecase, func(), error) {
	contextContext := cmd.BootContext()
	bootstrap := cmd.ConfBootstap()
	logger := cmd.ConfLogger(bootstrap)
	databaseMessage, cleanup, err := data.NewDatabaseMessage(contextContext, bootstrap, logger)
	if err != nil {
		return nil, nil, err
	}
	transactionMessageRepo := data.NewTransactionMessageRepo(databaseMessage, logger)
	producerMessage, cleanup2, err := data.NewProducerMessage(contextContext, bootstrap, logger)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	messageXHttpClient := data.NewMessageXHttpClient(clientClient, bootstrap)
	messageRepo := data.NewMessageRepo(databaseMessage, producerMessage, messageXHttpClient)
	recipientRepo := data.NewRecipientRepo(databaseMessage)
	messageUsecase := biz.NewMessageUsecase(transactionMessageRepo, messageRepo, recipientRepo)
	return messageUsecase, func() {
		cleanup2()
		cleanup()
	}, nil
}
