package builtin

import (
	// "log"
	"strings"

	// "github.com/hysios/apiai-go"
	"github.com/wanliu/flow/context"

	flow "github.com/wanliu/goflow"
)

type IntentCheck struct {
	flow.Component
	Ctx <-chan context.Request

	_intent string
	_score  float64
	// _flow   bool

	Intent <-chan string
	Score  <-chan float64
	// Flow   <-chan bool

	Out  chan<- context.Request
	Next chan<- context.Request

	// 即使意图和得分不满足，也向这个端口发送ｃｔｘ，列如使确认信息失效这种情况
	FlowOut chan<- context.Request
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

func (ic *IntentCheck) OnCtx(req context.Request) {
	cmd := req.Command

	if cmd != nil {
		if cmd.Action == ic._intent {
			ic.Out <- req
		} else {
			ic.Next <- req
		}

		return
	}

	res := req.ApiAiResult

	if strings.HasPrefix(res.Metadata.IntentName, ic._intent) && res.Score >= ic._score {
		ic.Out <- req
	} else {
		ic.Next <- req

		// log.Printf("...detect flow: %v", ic._flow)
		// if ic._flow {
		// log.Printf("...sending flow")

		ic.FlowOut <- req
		// }
	}
}
