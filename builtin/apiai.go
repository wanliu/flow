package builtin

import (
	"log"
	"sync"

	"github.com/oleiade/lane"

	. "github.com/wanliu/flow/builtin/ai"
	. "github.com/wanliu/flow/context"

	config "github.com/wanliu/flow/builtin/config"
	flow "github.com/wanliu/goflow"
)

func NewApiAi() interface{} {
	return new(ApiAi)
}

type ApiAi struct {
	flow.Component

	MultiField

	sync.Mutex

	token     string
	sessionId string
	proxyUrl  string

	CtxQueue *lane.Queue
	TxtQueue *lane.Queue

	Echo      <-chan bool
	In        <-chan string
	Token     <-chan string
	SessionId <-chan string
	ProxyUrl  <-chan string
	Out       chan<- Context
	Ctx       <-chan Context
}

func (c *ApiAi) Init() {
	c.CtxQueue = lane.NewQueue()
	c.TxtQueue = lane.NewQueue()
}

func (c *ApiAi) OnIn(input string) {
	// c.SetValue(config.ValueKeyText, input)
	c.Lock()
	c.TxtQueue.Enqueue(input)
	c.Unlock()

	c.SendQuery()
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
	// c.SetValue(config.ValueKeyCtx, ctx)
	c.Lock()
	c.CtxQueue.Enqueue(ctx)
	c.Unlock()

	c.SendQuery()
}

func (c *ApiAi) SendQuery() {
	c.Lock()

	for c.CtxQueue.Head() != nil && c.TxtQueue.Head() != nil {
		txt := c.TxtQueue.Dequeue().(string)
		ctx := c.CtxQueue.Dequeue().(Context)

		res, err := ApiAiQuery(txt, c.token, c.sessionId, c.proxyUrl)
		if err != nil {
			log.Printf("意图\"%s\"第一次解析失败:%s", txt, err.Error())
			log.Printf("尝试再一次解析")

			res, err = ApiAiQuery(txt, c.token, c.sessionId, c.proxyUrl)
			if err != nil {
				log.Printf("意图\"%s\"再一次解析失败:%s", txt, err.Error())
			}
		}

		ctx.SetValue(config.CtxkeyResult, res)

		intent := res.Metadata.IntentName
		score := res.Score
		// query := res.ResolvedQuery

		log.Printf("意图解析\"%s\" -> %s 准确度: %2.2f%%", txt, intent, score*100)

		c.Out <- ctx
	}

	c.Unlock()
}
