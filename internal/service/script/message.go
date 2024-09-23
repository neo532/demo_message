package script

import (
	"context"

	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/neo532/gokit/log"

	"demo_message/internal/biz"
)

type MessageScript struct {
	log *log.Helper
	msg *biz.MessageUsecase
}

func NewMessageScript(
	msg *biz.MessageUsecase,
	logger klog.Logger) *MessageScript {
	return &MessageScript{
		msg: msg,
		log: log.NewHelper(logger),
	}
}

func (s *MessageScript) ScanMessage(c context.Context, args string) (err error) {
	err = s.msg.ScanMessage(c)
	return
}
