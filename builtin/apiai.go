package builtin

import (
	"github.com/hysios/apiai-go"

	. "github.com/wanliu/flow/builtin/ai"
	flow "github.com/wanliu/goflow"
)

func NewApiAi() interface{} {
	return new(ApiAi)
}

type ApiAi struct {
	flow.Component
	token     string
	sessionId string

	Echo      <-chan bool
	In        <-chan string
	Token     <-chan string
	SessionId <-chan string
	Out       chan<- apiai.Result
}

func (l *ApiAi) OnIn(input string) {
	result, _ := ApiAiQuery(input, l.token, l.sessionId)

	// if err != nil {
	l.Out <- result
	// }
}

func (l *ApiAi) OnToken(token string) {
	l.token = token
}

func (l *ApiAi) OnSessionId(sessionId string) {
	l.sessionId = sessionId
}
