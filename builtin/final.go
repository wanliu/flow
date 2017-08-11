package builtin

import (
	"github.com/oleiade/lane"
	"log"
	"math/rand"
	"strconv"
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

	delayMin int
	delayMax int

	In       <-chan ReplyData
	DelayMin <-chan float64
	DelayMax <-chan float64
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

func (s *Final) OnDelayMin(min float64) {
	s.delayMin = int(min)
}

func (s *Final) OnDelayMax(max float64) {
	s.delayMax = int(max)
}

func (s Final) DelayRange() int {
	rand.Seed(time.Now().UnixNano())

	if s.delayMin == 0 {
		return 5 + rand.Intn(3)
	} else {
		if s.delayMax > s.delayMin {
			return s.delayMin + rand.Intn(s.delayMax-s.delayMin)
		} else {
			return s.delayMin + rand.Intn(3)
		}
	}
}

func (s *Final) SendReply() {
	secs := s.DelayRange()

	s.RLock()
	for s.ReplyQueue.Head() != nil {
		data := s.ReplyQueue.Dequeue().(ReplyData)

		log.Printf("[Delay]Delay reply for " + strconv.Itoa(secs) + " seconds.")
		time.Sleep(time.Second * time.Duration(secs))

		data.Ctx.Post(data.Reply)
	}
	s.RUnlock()
}
