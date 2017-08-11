package builtin

import (
	"github.com/oleiade/lane"
	"math/rand"
	"sync"
	"time"

	flow "github.com/wanliu/goflow"
)

func NewFinal() interface{} {
	return new(Final)
}

type Final struct {
	sync.RWMutex

	flow.Component

	ReplyQueue *lane.Queue

	In <-chan ReplyData
}

func (s *Final) Init() {
	s.ReplyQueue = lane.NewQueue()
}

func (s *Final) OnIn(data ReplyData) {
	s.Lock()
	s.ReplyQueue.Enqueue(data)
	s.Unlock()

	s.SendReply()
}

func (s *Final) SendReply() {
	rand.Seed(time.Now().UnixNano())

	s.RLock()
	for s.ReplyQueue.Head() != nil {
		data := s.ReplyQueue.Dequeue().(ReplyData)

		secs := 3 + rand.Intn(3)
		time.Sleep(time.Second * time.Duration(secs))

		data.Ctx.Post(data.Reply)
	}
	s.RUnlock()
}
