package api

import (
	"bufio"
	"context"
	"io"
	"os"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/neo532/gokit/log"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/emptypb"

	"demo_message/internal/biz"
	"demo_message/internal/biz/entity"
	"demo_message/internal/biz/repo"
	"demo_message/internal/data/db"
	"demo_message/util"

	pb "demo_message/proto/api/campaign/v1"
)

type CampaignApi struct {
	uc  *biz.CampaignUsecase
	cp  repo.CampaignRepo
	msg *biz.MessageUsecase
	log *log.Helper
	tag string
}

func NewCampaignApi(
	uc *biz.CampaignUsecase,
	cp repo.CampaignRepo,
	msg *biz.MessageUsecase,
	logger klog.Logger,
) *CampaignApi {
	return &CampaignApi{
		uc:  uc,
		cp:  cp,
		msg: msg,
		log: log.NewHelper(logger),
		tag: "api.CampaignApi",
	}
}

func (a *CampaignApi) Post(c context.Context, req *pb.PostRequest) (reply *pb.PostReply, err error) {

	reply = &pb.PostReply{}

	// campaign
	cp := &entity.Campaign{
		OriginType:    int(req.OriginType),
		OriginContent: req.OriginContent,
		MessageType:   int(req.MessageType),
		Message:       req.Message,
		TimeSend:      time.Now(),
	}
	if req.TimeSend != "" {
		if cp.TimeSend, err = time.Parse(time.DateTime, req.TimeSend); err != nil {
			return
		}
	}
	messageStatus := db.MessageStatusToSend
	reply.Id, err = a.uc.Create(c, cp)

	// message
	var fileR *os.File
	if fileR, err = os.OpenFile(cp.OriginContent, os.O_RDONLY, 0666); err != nil {
		return
	}
	defer fileR.Close()
	reader := bufio.NewReader(fileR)

	var line string
	pageSize := 2
	lRunning := 10
	msgs := make([]*entity.Message, 0, pageSize)
	var wg sync.WaitGroup
	var lock sync.Mutex
	running := make(chan int, lRunning)
	for {

		// parse
		if line, err = reader.ReadString('\n'); err == io.EOF {
			err = nil
			break
		}
		attr := strings.SplitN(strings.TrimSpace(line), ",", 3)
		if len(attr) != 3 || attr[0] == "Mobile" {
			continue
		}

		// save
		row := &entity.Message{
			CampaignID: reply.Id,
			Recipient: &entity.Recipient{
				Mobile: strings.TrimSpace(attr[0]),
				Name:   strings.TrimSpace(attr[1]),
			},
			Campaign: &entity.Campaign{
				Message: strings.Trim(attr[1], `"`),
			},
			Status: messageStatus,
		}
		msgs = append(msgs, row)
		if len(msgs) == pageSize {
			wg.Add(1)
			running <- 1

			go func(ms []*entity.Message) {
				if r := recover(); r != nil {
					lock.Lock()
					err = util.WrapErr(err, errors.Errorf(a.tag+"recover[%+v][%s]", r, string(debug.Stack())))
					lock.Unlock()
				}
				defer func() {
					<-running
					wg.Done()
				}()

				if er := a.msg.Create(c, cp, ms); er != nil {
					lock.Lock()
					err = util.WrapErr(err, er)
					lock.Unlock()
				}
			}(msgs)
			msgs = make([]*entity.Message, 0, pageSize)
		}
	}
	wg.Wait()
	close(running)
	if len(msgs) > 0 {
		if er := a.msg.Create(c, cp, msgs); er != nil {
			err = util.WrapErr(err, er)
		}
	}

	return
}

func (a *CampaignApi) PutStatus(c context.Context, req *pb.PutStatusRequest) (reply *emptypb.Empty, err error) {
	err = a.uc.UpdateStatus(c, req.Id, int(req.Status))
	return
}
