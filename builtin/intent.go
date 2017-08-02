package builtin

import (
	"log"

	"github.com/hysios/apiai-go"
	. "github.com/wanliu/flow/context"
	flow "github.com/wanliu/goflow"
)

type IntentCheck struct {
	flow.Component
	Ctx <-chan Context
	// Query  <-chan ResultParams
	_intent string
	_score  float64
	Intent  <-chan string
	Score   <-chan float64
	Out     chan<- Context
	Next    chan<- Context
}

func NewIntentCheck() interface{} {
	return new(IntentCheck)
}

func (ic *IntentCheck) OnIntent(intent string) {
	ic._intent = intent
}

func (ic *IntentCheck) OnScore(score float64) {
	ic._score = score
}

func (ic *IntentCheck) OnCtx(ctx Context) {
	if res, ok := ctx.Value("Result").(apiai.Result); ok {
		log.Printf("ai: %v", res)
		if res.Metadata.IntentName == ic._intent && res.Score >= ic._score {
			ic.Out <- ctx
		} else {
			ic.Next <- ctx
		}
	} else {
		log.Printf("无效的 Context Value Result")
	}
}
