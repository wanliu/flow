package builtin

import (
	"log"
	"strings"

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
	// _flow   bool

	Intent <-chan string
	Score  <-chan float64
	// Flow   <-chan bool

	Out  chan<- Context
	Next chan<- Context

	// 即使意图和得分不满足，也向这个端口发送ｃｔｘ，列如使确认信息失效这种情况
	FlowOut chan<- Context
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

func (ic *IntentCheck) OnFlow(flow bool) {
	// ic._flow = flow
}

func (ic *IntentCheck) OnCtx(ctx Context) {
	if res, ok := ctx.Value("Result").(apiai.Result); ok {
		// if res.Metadata.IntentName == ic._intent && res.Score >= ic._score {
		if strings.HasPrefix(res.Metadata.IntentName, ic._intent) && res.Score >= ic._score {
			ic.Out <- ctx
		} else {
			ic.Next <- ctx

			// log.Printf("...detect flow: %v", ic._flow)
			// if ic._flow {
			// log.Printf("...sending flow")

			ic.FlowOut <- ctx
			// }
		}
	} else {
		log.Printf("无效的 Context Value Result")
	}
}
