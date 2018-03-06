package builtin

import (
	// "github.com/oleiade/lane"
	// "log"
	"math/rand"
	// "strconv"
	// "sync"
	"time"

	"github.com/wanliu/flow/context"

	flow "github.com/wanliu/goflow"
)

func NewFinal() interface{} {
	return new(Final)
}

type Final struct {
	// sync.RWMutex

	flow.Component

	delayMin int
	delayMax int

	In       <-chan context.Request
	DelayMin <-chan float64
	DelayMax <-chan float64
}

func (s *Final) OnIn(req context.Request) {
	newReq := context.Request{
		Id:   req.Id,
		Text: req.Text,
	}
	// ctx := req.Ctx
	// req.Ctx = nil
	req.ctx.Post(req.Res.Reply, req.Res.Data, newReq)
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
		return 3 + rand.Intn(2)
	} else {
		if s.delayMax > s.delayMin {
			return s.delayMin + rand.Intn(s.delayMax-s.delayMin)
		} else {
			return s.delayMin + rand.Intn(3)
		}
	}
}
