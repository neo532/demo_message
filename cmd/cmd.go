package cmd

import (
	"context"
	"os"
	"reflect"
	"strings"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/neo532/gokit/log/zap"
	"github.com/neo532/gokit/middleware"
	"github.com/neo532/gokit/middleware/tracing"

	"demo_message/internal/conf"
)

type pkg struct{}

var (
	InitUnitTestSet = wire.NewSet(BootContext, ConfBootstap, ConfLogger)
	rootPath        string
)

func ConfBootstap() (bs *conf.Bootstrap) {
	rootPath := "."

	if pwd, err := os.Getwd(); err == nil {

		pn := strings.Split(reflect.TypeOf(pkg{}).PkgPath(), "/")
		pkgName := "/" + pn[len(pn)-2]

		tmp := strings.SplitN(pwd, pkgName, 2)
		if len(tmp) > 0 {
			rootPath = tmp[0] + pkgName
		}
	}

	var err error
	if bs, err = InitConfig(rootPath + "/configs/config.yaml"); err != nil {
		panic(err)
	}
	return bs
}

func ConfLogger(bs *conf.Bootstrap) klog.Logger {
	return InitLogger(bs, middleware.EntryTest, bs.General.Logger.FilenameTest)
}

func InitConfig(path string) (bs *conf.Bootstrap, err error) {
	bs = &conf.Bootstrap{}
	c := config.New(
		config.WithSource(
			file.NewSource(path),
		),
	)
	defer c.Close()
	if err = c.Load(); err != nil {
		return
	}
	if err = c.Scan(bs); err != nil {
		return
	}

	bs.General.Ip, _ = os.Hostname()
	return
}

func InitLogger(bc *conf.Bootstrap, entry string, fileName string) (zapLogger *zap.ZapLogger) {
	zapLogger = zap.NewLogger(
		zap.WithIP(bc.General.Ip),
		zap.WithName(bc.General.Name),
		zap.WithDepartment(bc.General.Department),
		zap.WithLevel(bc.General.Logger.Level),
		zap.WithEnv(bc.General.Env),
		zap.WithCompress(bc.General.Logger.Compress),
		zap.WithMaxAge(int(bc.General.Logger.MaxAge)),
		zap.WithMaxBackups(int(bc.General.Logger.MaxBackup)),
		zap.WithMaxSize(int(bc.General.Logger.MaxSize)),
		zap.WithVersion(bc.General.Version),
		zap.WithFilename(fileName),
		zap.WithEntry(entry),
	)
	return
}

func BootContext() (ctx context.Context) {
	ctx = context.Background()
	ctx = tracing.SetTraceIDForServer(ctx, "")
	ctx = tracing.SetFromForServer(ctx, "bootstrap")
	return
}
