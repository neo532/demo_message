package consumer

import (
	"context"
	"encoding/json"

	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/neo532/gokit/log"
	"github.com/pkg/errors"

	"demo_message/internal/biz"
	"demo_message/internal/biz/entity"
)

type MessageConsumer struct {
	msg *biz.MessageUsecase
	log *log.Helper
	tag string
}

func NewMessageConsumer(
	msg *biz.MessageUsecase,
	logger klog.Logger) *MessageConsumer {
	return &MessageConsumer{
		msg: msg,
		tag: "internal/service/consumer/message",
		log: log.NewHelper(logger),
	}
}

func (s *MessageConsumer) Send(c context.Context, message []byte) (err error) {

	var msg []*entity.Message
	if err = json.Unmarshal(message, &msg); err != nil {
		err = errors.Wrap(err, s.tag+".Send.json.Unmarshal")
		return
	}

	err = s.msg.Send(c, msg)
	return
}
