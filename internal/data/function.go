package data

import (
	"context"

	"github.com/IBM/sarama"
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/neo532/gokit/database/orm"
	"github.com/neo532/gokit/database/redis"
	"github.com/neo532/gokit/queue"
	"github.com/neo532/gokit/queue/kafka/producer"
	"gorm.io/driver/mysql"

	"demo_message/internal/conf"
)

// ========== Database ==========
func newDatabase(c context.Context, cfg *conf.Data_Database, logger klog.Logger) (dbs *orm.Orms) {
	connet := func(c context.Context, d *conf.Data_Database, db *conf.Data_DatabaseDb, logger klog.Logger) *orm.Orm {
		return orm.New(
			db.Name,
			mysql.Open(db.Dsn),
			orm.WithTablePrefix(d.TablePrefix),
			orm.WithConnMaxLifetime(d.ConnMaxLifetime.AsDuration()),
			orm.WithMaxIdleConns(int(d.MaxIdleConns)),
			orm.WithMaxOpenConns(int(d.MaxOpenConns)),
			orm.WithLogger(logger),
			orm.WithSingularTable(),
			orm.WithContext(c),
			orm.WithSlowLog(d.MaxSlowtime.AsDuration()),
		)
	}
	dbs = orm.News(
		connet(c, cfg, cfg.Read, logger),
		connet(c, cfg, cfg.Write, logger),
	)
	if cfg.ShadowRead != nil && cfg.ShadowWrite != nil {
		dbs.SetShadow(
			connet(c, cfg, cfg.ShadowRead, logger),
			connet(c, cfg, cfg.ShadowWrite, logger),
		)
	}
	return
}

// ========== /Database ==========

// ========== Redis ==========
func newRedis(c context.Context, cg *conf.General, cfg *conf.Data_Redis, logger klog.Logger) (rdbs *redis.Rediss) {
	connnet := func(c context.Context, d *conf.Data_Redis, r *conf.Data_RedisRdb, logger klog.Logger) *redis.Redis {
		return redis.New(
			r.Name,
			r.Addr,
			redis.WithPassword(r.Password),
			redis.WithSlowTime(d.MaxSlowtime.AsDuration()),
			redis.WithDb(int(r.Db)),
			redis.WithLogger(logger),
			redis.WithContext(c),
		)
	}

	rdbs = redis.News(connnet(c, cfg, cfg.Default, logger))
	if cfg.Shadow != nil {
		rdbs.SetShadow(connnet(c, cfg, cfg.Shadow, logger))
	}
	if cfg.Gray != nil {
		rdbs.SetGray(connnet(c, cfg, cfg.Gray, logger))
	}
	return
}

// ========== /Redis ==========

// ========== Producer ==========
func newProducer(c context.Context, cfg *conf.Data_Producer, logger klog.Logger) (pds *queue.Producers) {
	connect := func(c context.Context, p *conf.Data_ProducerPdc, logger klog.Logger) queue.Producer {
		return producer.New(
			p.Name,
			p.Addrs,
			producer.WithReturnSucesses(true),
			producer.WithPartitioner(sarama.NewHashPartitioner),
			producer.WithLogger(logger, c),
			producer.WithRequiredAcks(sarama.WaitForAll),
			producer.WithSync(true),
			producer.WithTopic(p.Topic),
			producer.WithIdempotent(p.Idempotent),
			producer.WithNetMaxOpenRequest(int(p.NetMaxOpenRequest)),
		)
	}
	pds = queue.NewProducers(connect(c, cfg.Default, logger))
	if cfg.Shadow != nil {
		pds.SetShadow(connect(c, cfg.Shadow, logger))
	}
	if cfg.Gray != nil {
		pds.SetGray(connect(c, cfg.Gray, logger))
	}
	return pds
}

// ========== /Producer ==========
