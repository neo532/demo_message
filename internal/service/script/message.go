package script

import (
	"context"
	"encoding/csv"
	"os"
	"strconv"

	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/neo532/gokit/log"

	"demo_message/internal/biz"
)

type MessageScript struct {
	log *log.Helper
	msg *biz.MessageUsecase
}

func NewMessageScript(
	msg *biz.MessageUsecase,
	logger klog.Logger) *MessageScript {
	return &MessageScript{
		msg: msg,
		log: log.NewHelper(logger),
	}
}

func (s *MessageScript) ScanMessage(c context.Context, args string) (err error) {
	s.GenerateCsv(c, args)
	//err = s.msg.ScanMessage(c)
	return
}

func (s *MessageScript) GenerateCsv(c context.Context, args string) (err error) {

	var file *os.File
	file, err = os.Create("/Users/neo/devspace/go/iv/demo_message/internal/conf/data2.csv")

	w := csv.NewWriter(file)

	w.Write([]string{"Mobile", "Name", "Message"})

	for i := 0; i < 100000; i++ {
		w.Write([]string{strconv.Itoa(i), "aaaa", "message"})
	}
	w.Flush()
	return
}
