package data

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/neo532/gofr/gofunc"
	"github.com/neo532/gofr/tool"
	"github.com/neo532/gokit/database/orm"
	"github.com/neo532/gokit/database/redis"
	"github.com/neo532/gokit/middleware/tracing"
	"github.com/neo532/gokit/queue"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"demo_message/internal/biz/entity"
	"demo_message/internal/biz/repo"
	"demo_message/internal/data/db"
	message "demo_message/proto/client/message/v1"
	"demo_message/util"
)

type MessageRepo struct {
	db   *orm.Orms
	pdc  *queue.Producers
	rdb  *redis.Rediss
	freq *tool.Freq
	xclt *message.MessageXHttpClient
	tag  string
}

func NewMessageRepo(
	db DatabaseMessage,
	pdc ProducerMessage,
	rdb RedisFreq,
	xclt *message.MessageXHttpClient,
) repo.MessageRepo {
	freq := tool.NewFreq(&rDb{rdb})
	freq.Timezone("Local")

	return &MessageRepo{
		tag:  "internal/data/message",
		db:   db,
		pdc:  pdc,
		freq: freq,
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
			Status:      d.Status,
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

		preKey := strings.NewReplacer("{duration}", time.Now().Format("2006-01-02_15:04")).Replace("sendmessge_{duration}")
		rule := []tool.FreqRule{
			tool.FreqRule{Duri: "60", Times: 100},
		}

		for {
			if ok, err := r.freq.IncrCheck(c, preKey, rule...); err == nil && ok {
				break
			}
			time.Sleep(time.Second)
		}

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
	//

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

func (r *MessageRepo) ScanToSend(c context.Context, limit, offset int) (rs []*entity.Message, err error) {
	rs = make([]*entity.Message, 0, limit)

	var ds []*db.Message
	err = r.db.Read(c).
		WithContext(c).
		Select("id", "campaign_id", "recipient_id").
		Order("id desc").
		Limit(limit).
		Offset(offset).
		Where("status=?", db.MessageStatusToSend).
		Preload("Campaign", func(odb *gorm.DB) *gorm.DB {
			return odb.
				Select("id", "message").
				Where("status=? and time_send>?", db.CampaignStatusOn, time.Now().Format(time.DateTime))
		}).
		Preload("Recipient", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "mobile", "name")
		}).
		Find(&ds).
		Error
	if err != nil {
		return
	}
	badIDs := make([]int64, 0, len(ds))
	for _, v := range ds {
		if v.Campaign == nil || v.Recipient == nil {
			badIDs = append(badIDs, v.ID)
			continue
		}
		row := &entity.Message{
			ID: v.ID,
			Campaign: &entity.Campaign{
				Message: v.Campaign.Message,
			},
			Recipient: &entity.Recipient{
				Mobile: v.Recipient.Mobile,
				Name:   v.Recipient.Name,
			},
		}
		rs = append(rs, row)
	}

	if len(badIDs) > 0 {
		if er := r.SaveStatus(c, badIDs, db.MessageStatusBadMessage); er != nil {
			err = util.WrapErr(err, er)
		}
	}
	return
}
