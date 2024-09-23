package data

import (
	"context"

	"github.com/neo532/gokit/database/orm"

	"demo_message/internal/biz/entity"
	"demo_message/internal/biz/repo"
	"demo_message/internal/data/db"
)

type RecipientRepo struct {
	db *orm.Orms
}

func NewRecipientRepo(
	db DatabaseMessage,
) repo.RecipientRepo {
	return &RecipientRepo{
		db: db,
	}
}

func (r *RecipientRepo) CreateByMessage(c context.Context, ds []*entity.Message) (err error) {
	l := len(ds)

	data := make([]*db.Recipient, 0, l)
	for _, d := range ds {
		if d.Recipient == nil {
			continue
		}
		data = append(data, &db.Recipient{
			Mobile: d.Recipient.Mobile,
			Name:   d.Recipient.Name,
		})
	}

	key := func(mobile, name string) string {
		return mobile + name
	}

	err = r.db.Write(c).
		WithContext(c).
		Create(data).
		Error
	rs := make(map[string]int64, l)
	emptyIDs := make([]string, 0, l)
	for _, d := range data {
		if d.ID == 0 {
			emptyIDs = append(emptyIDs, d.Mobile)
			continue
		}
		rs[key(d.Mobile, d.Name)] = d.ID
	}
	var emptys []*db.Recipient
	if err = r.db.Read(c).
		WithContext(c).
		Select("id", "mobile", "name").
		Where("mobile in ?", emptyIDs).
		Find(&emptys).
		Error; err == nil && emptys != nil {
		for _, e := range emptys {
			rs[key(e.Mobile, e.Name)] = e.ID
		}
	}

	for i, d := range ds {
		if d.Recipient == nil {
			continue
		}
		ds[i].RecipientID, _ = rs[key(d.Recipient.Mobile, d.Recipient.Name)]
	}

	return
}
