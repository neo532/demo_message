package api

import (
	"context"

	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/neo532/gokit/log"
	"google.golang.org/protobuf/types/known/emptypb"

	"demo_message/internal/biz"

	pb "demo_message/proto/api/message/v1"
)

type MessageApi struct {
	uc  *biz.MessageUsecase
	log *log.Helper
	tag string
}

func NewMessageApi(
	uc *biz.MessageUsecase,
	logger klog.Logger,
) *MessageApi {
	return &MessageApi{
		uc:  uc,
		log: log.NewHelper(logger),
		tag: "api.messageApi",
	}
}

func (a *MessageApi) Post(c context.Context, req *pb.PostRequest) (reply *emptypb.Empty, err error) {
	return
}

func (a *MessageApi) PutStatus(c context.Context, req *pb.PutStatusRequest) (reply *emptypb.Empty, err error) {
	err = a.uc.UpdateStatus(c, req.Id, int(req.Status))
	return
}
