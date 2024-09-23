package repo

import (
	"context"

	"demo_message/internal/biz/entity"
)

type TransactionMessageRepo interface {
	Transaction(c context.Context, fn func(c context.Context) (err error)) (err error)
}

type CampaignRepo interface {
	Create(c context.Context, cps *entity.Campaign) (insID int64, err error)
	UpdateStatus(c context.Context, ID int64, status int) (err error)
	IsImmediately(c context.Context, d *entity.Campaign) (b bool)
}

type RecipientRepo interface {
	CreateByMessage(c context.Context, ds []*entity.Message) (err error)
}

type MessageRepo interface {
	Create(c context.Context, msg []*entity.Message) (insIDs []int64, err error)
	PushQueue(c context.Context, msg []*entity.Message) (err error)
	Send(c context.Context, msg []*entity.Message) (succIDs, failIDs []int64, err error)
	SaveStatus(c context.Context, IDs []int64, status int) (err error)
	ScanToSend(c context.Context, limit, offset int) (rs []*entity.Message, err error)
}
