package builtin

import (
	"log"

	. "github.com/wanliu/flow/builtin/ai"
	config "github.com/wanliu/flow/builtin/config"
	. "github.com/wanliu/flow/context"
	flow "github.com/wanliu/goflow"
)

func NewApiAi() interface{} {
	return new(ApiAi)
}

type ApiAi struct {
	flow.Component

	MultiField

	token     string
	sessionId string
	proxyUrl  string

	Echo      <-chan bool
	In        <-chan string
	Token     <-chan string
	SessionId <-chan string
	ProxyUrl  <-chan string
	Out       chan<- Context
	Ctx       <-chan Context
}

func (c *ApiAi) Init() {
	c.Fields = []string{config.ValueKeyCtx, config.ValueKeyText}

	c.Process = func() error {
		txt, tok := c.Value(config.ValueKeyText).(string)
		ctx, cok := c.Value(config.ValueKeyCtx).(Context)

		if tok && cok {
			go func() {
				res, _ := ApiAiQuery(txt, c.token, c.sessionId, c.proxyUrl)

				ctx.SetValue(config.CtxkeyResult, res)
				intent := res.Metadata.IntentName
				score := res.Score
				query := res.ResolvedQuery

				log.Printf("意图解析\"%s\" -> %s 准确度: %2.2f%%", query, intent, score*100)
				c.Out <- ctx
			}()
		}

		return nil
	}
}

func (c *ApiAi) OnIn(input string) {
	c.SetValue(config.ValueKeyText, input)
}

func (c *ApiAi) OnToken(token string) {
	c.token = token
}

func (c *ApiAi) OnSessionId(sessionId string) {
	c.sessionId = sessionId
}

func (c *ApiAi) OnProxyUrl(proxy string) {
	c.proxyUrl = proxy
}

func (c *ApiAi) OnCtx(ctx Context) {
	c.SetValue(config.ValueKeyCtx, ctx)
}
