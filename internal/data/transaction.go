package data

import (
	"context"

	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/neo532/gokit/database/orm"
	"github.com/neo532/gokit/log"

	"demo_message/internal/biz/repo"
)

type TransactionMessageRepo struct {
	db  *orm.Orms
	log *log.Helper
}

func NewTransactionMessageRepo(messageDB DatabaseMessage, logger klog.Logger) repo.TransactionMessageRepo {
	return &TransactionMessageRepo{
		db:  messageDB,
		log: log.NewHelper(logger),
	}
}

func (r *TransactionMessageRepo) Transaction(c context.Context, fn func(ctx context.Context) error) (err error) {
	err = r.db.Transaction(c, fn)
	return
}
