package data

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/neo532/gofr/gofunc"
	"github.com/neo532/gokit/database/orm"
	"github.com/neo532/gokit/middleware/tracing"
	"github.com/neo532/gokit/queue"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"demo_message/internal/biz/entity"
	"demo_message/internal/biz/repo"
	"demo_message/internal/data/db"
	message "demo_message/proto/client/message/v1"
)

type MessageRepo struct {
	db   *orm.Orms
	pdc  *queue.Producers
	xclt *message.MessageXHttpClient
	tag  string
}

func NewMessageRepo(
	db DatabaseMessage,
	pdc ProducerMessage,
	xclt *message.MessageXHttpClient,
) repo.MessageRepo {
	return &MessageRepo{
		tag:  "internal/data/message",
		db:   db,
		pdc:  pdc,
		xclt: xclt,
	}
}

func (r *MessageRepo) Create(c context.Context, ds []*entity.Message) (insIDs []int64, err error) {

	data := make([]*db.Message, 0, len(ds))
	def := time.Unix(1, 0).UTC()
	for _, d := range ds {
		data = append(data, &db.Message{
			CampaignID:  d.CampaignID,
			RecipientID: d.RecipientID,
			Status:      db.MessageStatusToSend,
			// Recipient: &db.Recipient{
			// 	Mobile: d.Recipient.Mobile,
			// 	Name:   d.Recipient.Name,
			// },
			TimeSend: def,
		})
	}

	err = r.db.Write(c).
		WithContext(c).
		Save(data).
		Error

	insIDs = make([]int64, 0, len(data))
	for i, d := range data {
		insIDs = append(insIDs, d.ID)
		ds[i].ID = d.ID
	}
	return
}

func (r *MessageRepo) PushQueue(c context.Context, msg []*entity.Message) (err error) {
	var b []byte
	if b, err = json.Marshal(msg); err != nil {
		err = errors.Wrap(err, r.tag+".Send")
		return
	}
	err = r.pdc.Producer(c).Send(c, b, tracing.GetTraceIDByCtx(c))
	return
}

func (r *MessageRepo) Send(c context.Context, msgs []*entity.Message) (succIDs, failIDs []int64, err error) {

	l := len(msgs)
	succIDs = make([]int64, l, l)
	failIDs = make([]int64, l, l)

	fn := func(i int) (err error) {
		msg := msgs[i]
		if msg.Campaign == nil || msg.Recipient == nil {
			err = errors.New(r.tag + ".Send has empty param")
			failIDs[i] = msg.ID
			return
		}
		req := &message.SendRequest{
			Mobile:  msg.Recipient.Mobile,
			Message: msg.Campaign.Message,
		}

		var resp *message.SendReply
		resp, err = r.xclt.Send(c, req)
		if err != nil {
			failIDs[i] = msg.ID
			return
		}
		if resp != nil && resp.Code != http.StatusOK {
			err = errors.New(resp.Msg)
			return
		}
		succIDs[i] = msg.ID
		return
	}

	log := &gofunc.DefaultLogger{}
	gofn := gofunc.NewGoFunc(gofunc.WithLogger(log), gofunc.WithMaxGoroutine(20))
	fns := make([]func(i int) error, 0, l)
	for i := 0; i < l; i++ {
		fns = append(fns, fn)
	}
	gofn.WithTimeout(c, 3*time.Second, fns...)
	err = log.Err()
	return
}

func (r *MessageRepo) SaveStatus(c context.Context, oIDs []int64, status int) (err error) {
	IDs := make([]int64, 0, len(oIDs))
	for _, ID := range oIDs {
		if ID != 0 {
			IDs = append(IDs, ID)
		}
	}

	if len(IDs) == 0 {
		return
	}
	d := map[string]interface{}{
		"status": status,
	}
	switch status {
	case db.MessageStatusSended:
		d["time_send"] = time.Now()
	case db.MessageStatusSendFail:
		d["log_id"] = tracing.GetTraceIDByCtx(c)
	}
	err = r.db.Write(c).
		WithContext(c).
		Model(&db.Message{}).
		Where("id in ?", IDs).
		UpdateColumns(d).
		Error
	return
}

func (r *MessageRepo) ScanToSend(c context.Context) (rs []*entity.Message, err error) {
	rs = make([]*entity.Message, 0, 5)

	var ds []*db.Message
	err = r.db.Read(c).
		WithContext(c).
		Select("id", "campaign_id", "recipient_id").
		Order("id desc").
		Limit(1000).
		Joins("left join campaign as c on c.id=message.campaign_id").
		Where("c.status=? and c.time_send>? and message.status=?", db.CampaignStatusOn, time.Now().Format(time.DateTime), db.MessageStatusToSend).
		Preload("Campaign", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "message")
		}).
		Preload("Recipient", func(db *gorm.DB) *gorm.DB {
			return db.Select("mobile", "name")
		}).
		Find(&ds).
		Error
	if err != nil {
		return
	}
	for _, v := range ds {
		row := &entity.Message{
			ID: v.ID,
		}
		if v.Campaign != nil {
			row.Campaign = &entity.Campaign{
				Message: v.Campaign.Message,
			}
		}
		if v.Recipient != nil {
			row.Recipient = &entity.Recipient{
				Mobile: v.Recipient.Mobile,
				Name:   v.Recipient.Name,
			}
		}
		rs = append(rs, row)
	}
	return
}
