package biz

import (
	"context"

	"demo_message/internal/biz/entity"
	"demo_message/internal/biz/repo"
	"demo_message/internal/data/db"
	"demo_message/util"

	"github.com/neo532/gofr/tool"
)

type MessageUsecase struct {
	tx        repo.TransactionMessageRepo
	cp        repo.CampaignRepo
	msg       repo.MessageRepo
	rcp       repo.RecipientRepo
	queueSize int
}

func NewMessageUsecase(
	tx repo.TransactionMessageRepo,
	cp repo.CampaignRepo,
	msg repo.MessageRepo,
	rcp repo.RecipientRepo,
) *MessageUsecase {
	return &MessageUsecase{
		tx:        tx,
		cp:        cp,
		msg:       msg,
		rcp:       rcp,
		queueSize: 10,
	}
}

func (uc *MessageUsecase) Create(c context.Context, cp *entity.Campaign, msgs []*entity.Message) (err error) {

	isImmediately := uc.cp.IsImmediately(c, cp)

	tool.PageExec(int64(len(msgs)), uc.queueSize, func(begin, end int64, page int) (er error) {
		ms := msgs[begin:end]

		er = uc.tx.Transaction(c, func(ctx context.Context) (e error) {

			if e = uc.rcp.CreateByMessage(ctx, ms); e != nil {
				return
			}

			if _, e = uc.msg.Create(ctx, ms); e != nil {
				return
			}

			if !isImmediately {
				return
			}
			qms := make([]*entity.Message, 0, len(ms))
			for _, msg := range ms {
				qms = append(qms, msg)
			}
			if len(qms) > 0 {
				if e = uc.msg.PushQueue(ctx, qms); e != nil {
					return
				}
			}
			return
		})
		if er != nil {
			err = util.WrapErr(err, er)
		}
		return
	})

	return
}

func (uc *MessageUsecase) Send(c context.Context, msgs []*entity.Message) (err error) {

	l := len(msgs)
	succIDs := make([]int64, 0, l)
	failIDs := make([]int64, 0, l)

	succIDs, failIDs, err = uc.msg.Send(c, msgs)

	if len(succIDs) > 0 {
		if er := uc.msg.SaveStatus(c, succIDs, db.MessageStatusSended); er != nil {
			err = util.WrapErr(err, er)
		}
	}
	if len(failIDs) > 0 {
		if er := uc.msg.SaveStatus(c, failIDs, db.MessageStatusSendFail); er != nil {
			err = util.WrapErr(err, er)
		}
	}
	return
}

func (uc *MessageUsecase) ScanMessage(c context.Context) (err error) {

	var msgs []*entity.Message
	pageSize := 1000
	for i := 0; ; i += pageSize {
		msgs, err = uc.msg.ScanToSend(c, pageSize, i)
		if len(msgs) == 0 || err != nil {
			return
		}

		tool.PageExec(int64(len(msgs)), uc.queueSize, func(begin, end int64, page int) (er error) {

			er = uc.tx.Transaction(c, func(ctx context.Context) (e error) {

				ms := msgs[begin:end]

				IDs := make([]int64, 0, uc.queueSize)
				for _, v := range ms {
					IDs = append(IDs, v.ID)
				}
				if e = uc.msg.SaveStatus(c, IDs, db.MessageStatusSendMessage); e != nil {
					return
				}

				if e = uc.msg.PushQueue(c, ms); e != nil {
					return
				}
				if e != nil {
					err = util.WrapErr(err, e)
				}
				return
			})
			return
		})
	}

	return
}

func (uc *MessageUsecase) UpdateStatus(c context.Context, ID int64, status int) (err error) {
	return uc.msg.SaveStatus(c, []int64{ID}, status)
}
