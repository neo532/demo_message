package api

import (
	"context"

	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/neo532/gokit/log"
	"google.golang.org/protobuf/types/known/emptypb"
)

type SystemApi struct {
	log *log.Helper
}

func NewSystemApi(
	logger klog.Logger,
) *SystemApi {
	return &SystemApi{
		log: log.NewHelper(logger),
	}
}

func (a *SystemApi) GetPing(c context.Context, req *emptypb.Empty) (reply *emptypb.Empty, err error) {
	return
}
