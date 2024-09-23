package server

import (
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/neo532/gokit/middleware"
	"github.com/neo532/gokit/middleware/log"
	"github.com/neo532/gokit/middleware/server"
	"github.com/neo532/gokit/middleware/tracing"
	"github.com/neo532/gokit/transport/http"

	"demo_message/internal/conf"
)

// NewHTTPServer new a HTTP server.
func NewHTTPServer(bs *conf.Bootstrap, logging klog.Logger) *khttp.Server {

	var opts = []khttp.ServerOption{
		khttp.Middleware(
			server.SetEnv(bs.General.Env),
			server.SetEntry(middleware.EntryApi),
			tracing.Server(),
			recovery.Recovery(),
			log.Server(logging),
			validate.Validator(),
		),
		khttp.ResponseEncoder(http.ResponseEncoder),
		khttp.ErrorEncoder(http.ErrorEncoder),
	}
	if bs.Server.Http.Network != "" {
		opts = append(opts, khttp.Network(bs.Server.Http.Network))
	}
	if bs.Server.Http.Addr != "" {
		opts = append(opts, khttp.Address(bs.Server.Http.Addr))
	}
	if bs.Server.Http.Timeout != nil {
		opts = append(opts, khttp.Timeout(bs.Server.Http.Timeout.AsDuration()))
	}
	srv := khttp.NewServer(opts...)

	return srv
}
