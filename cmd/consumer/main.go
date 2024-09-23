package main

import (
	"context"
	"flag"
	"os"
	"time"

	"github.com/go-kratos/kratos/v2"
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/neo532/apitool/transport"
	"github.com/neo532/apitool/transport/http/xhttp/client"
	"github.com/neo532/gokit/log"
	"github.com/neo532/gokit/middleware"
	"github.com/neo532/gokit/middleware/tracing"
	"github.com/neo532/gokit/queue"
	"github.com/neo532/gokit/server"

	"demo_message/cmd"
	"demo_message/internal/conf"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	//Name string
	// Version is the version of the compiled software.
	Version string
	// flagconf is the config flag.
	flagConf string
)

func init() {
	flag.StringVar(&flagConf, "conf", "../../configs", "config path, eg: -conf config.yaml")
}

func newApp(ctx context.Context, bs *conf.Bootstrap, logger klog.Logger, csm queue.Consumer) *kratos.App {
	return kratos.New(
		kratos.ID(bs.General.Ip),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			csm,
		),
		kratos.Context(ctx),
	)
}

func main() {
	flag.Parse()

	// config
	bs, err := cmd.InitConfig(flagConf)
	if err != nil {
		panic(err)
	}

	// log
	zapLogger := cmd.InitLogger(bs, middleware.EntryConsumer, bs.General.Logger.FilenameConsumer)
	defer func() {
		zapLogger.Sync()
		time.Sleep(3 * time.Second)
	}()
	logger := log.AddGlobalVariable(zapLogger)

	envEnv, er := transport.String2Env(bs.General.Env)
	if er != nil {
		panic(er)
	}

	// traceID
	tracing.SetGroupForTracing(bs.General.Name)
	tracing.SetNameForTracing(bs.General.Name)

	// app
	app, cleanup, err := initApp(
		cmd.BootContext(),
		bs,
		client.New(
			client.WithLogger(log.NewXHttpLogger(logger)),
			client.WithEnv(envEnv),
		),
		logger,
	)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// pid
	// prestop: ["cat","/home/www/insurance/pid","|","xargs","kill","-2"]
	if err := server.WritePID(os.Getpid(), ""); err != nil {
		panic(err)
	}

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}
