package data

import (
	"context"

	"github.com/neo532/gokit/database/orm"

	"demo_message/internal/biz/entity"
	"demo_message/internal/biz/repo"
	"demo_message/internal/data/db"
)

type CampaignRepo struct {
	db *orm.Orms
}

func NewCampaignRepo(
	db DatabaseMessage,
) repo.CampaignRepo {
	return &CampaignRepo{
		db: db,
	}
}

func (r *CampaignRepo) Create(c context.Context, d *entity.Campaign) (insID int64, err error) {

	data := &db.Campaign{
		OriginType:    d.OriginType,
		OriginContent: d.OriginContent,
		MessageType:   d.MessageType,
		Message:       d.Message,
		Status:        db.CampaignStatusOn,
		TimeSend:      d.TimeSend,
	}

	err = r.db.Write(c).
		WithContext(c).
		Create(data).
		Error

	insID = data.ID
	return
}

func (r *CampaignRepo) UpdateStatus(c context.Context, ID int64, status int) (err error) {

	err = r.db.Read(c).
		WithContext(c).
		Where("id=?", ID).
		Update("status", status).
		Error
	if err != nil {
		return
	}
	return
}
