package data

import (
	"context"

	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/neo532/apitool/transport/http/xhttp/client"
	"github.com/neo532/gokit/database/orm"
	"github.com/neo532/gokit/queue"

	//"github.com/neo532/gokit/database/redis"

	"demo_message/internal/conf"
	message "demo_message/proto/client/message/v1"
)

type (
	DatabaseMessage *orm.Orms
	ProducerMessage *queue.Producers
	//RedisLock       *redis.Rediss
)

// ========== Database ==========
func NewDatabaseMessage(c context.Context, bs *conf.Bootstrap, logger klog.Logger) (DatabaseMessage, func(), error) {
	dbs := newDatabase(c, bs.Data.DatabaseMessage.Conf, logger)
	return dbs, dbs.Cleanup(), dbs.Err
}

// ========== /Database ==========

// ========== Redis ==========

// func NewRedisLock(c context.Context, bs *conf.Bootstrap, logger klog.Logger) (RedisLock, func(), error) {
// 	rdbs := newRedis(c, bs.General, bs.Data.RedisLock.Conf, logger)
// 	return rdbs, rdbs.Cleanup(), rdbs.Err
// }

// func NewToolDistributedLock(rdb RedisLock) *tool.DistributedLock {
// 	return tool.NewDistributedLock(&lredis.GoRedis{Rdb: rdb})
// }

// ========== /Redis ==========

// ========== Producer ==========
func NewProducerMessage(c context.Context, bs *conf.Bootstrap, logger klog.Logger) (ProducerMessage, func(), error) {
	pdcs := newProducer(c, bs.Data.ProducerMessage.Conf, logger)
	return pdcs, pdcs.CleanUp(), pdcs.Err
}

// ========== /Producer ==========

// ========== Client ==========
func NewMessageXHttpClient(clt client.Client, bs *conf.Bootstrap) (xclt *message.MessageXHttpClient) {
	xclt = message.NewMessageXHttpClient(clt)
	xclt.Domain = bs.Third.Message.Domain
	//xclt.WithMiddleware(message.Demo())
	return
}

// ========== /Client ==========
