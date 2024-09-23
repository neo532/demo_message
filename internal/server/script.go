package server

import (
	"context"
	"time"

	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/neo532/gokit/log"
	"github.com/neo532/gokit/middleware"
	mdwServer "github.com/neo532/gokit/middleware/server"
	"github.com/neo532/gokit/middleware/tracing"
	"github.com/pkg/errors"

	"demo_message/internal/service/script"
)

type Script struct {
	arguments      []string
	log            *log.Helper
	err            error
	scriptFinishCh chan struct{}
	startTime      time.Time

	msg *script.MessageScript
	srv chan bool
}

func (s *Script) Start(c context.Context) (err error) {
	defer func() {
		s.scriptFinishCh <- struct{}{}
	}()

	if len(s.arguments) < 1 {
		err = errors.Errorf("%+v", s.arguments)
		return
	}
	c = tracing.Script(c)
	c = mdwServer.SetEntryForCtx(c, middleware.EntryScript)
	s.startTime = time.Now()

	switch s.arguments[0] {
	case "Get":
		err = s.msg.Push(c, s.arguments[1])
	default:
		err = errors.Errorf("Invaild command![args:%+v]", s.arguments)
	}
	return
}
func (s *Script) Stop(c context.Context) (err error) {
	if s.err != nil {
		s.log.WithContext(c).Errorf("Has err[%+v]", s.err)
	}
	cost := time.Since(s.startTime)
	s.log.WithContext(c).Infof("Script stop:%v startTime:%v costTime:%v", s.arguments[0], s.startTime, cost)
	return
}

// RunScript new a script.
func NewScript(args []string, finish chan struct{}, logging klog.Logger, dm *script.MessageScript) (srv *Script) {
	srv = &Script{
		arguments:      args,
		log:            log.NewHelper(logging),
		msg:            dm,
		scriptFinishCh: finish,
	}
	return
}
