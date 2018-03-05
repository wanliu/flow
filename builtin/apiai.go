package builtin

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/hysios/apiai-go"
	"github.com/oleiade/lane"
	// "github.com/wanliu/flow/builtin/config"
	"github.com/wanliu/flow/context"

	. "github.com/wanliu/flow/builtin/ai"

	flow "github.com/wanliu/goflow"
)

func NewApiAi() interface{} {
	return new(ApiAi)
}

type ApiAi struct {
	flow.Component

	MultiField

	sync.Mutex

	token      string
	sessionId  string
	proxyUrl   string
	retryCount int

	CtxQueue *lane.Queue
	TxtQueue *lane.Queue

	Echo      <-chan bool
	In        <-chan context.Request
	Token     <-chan string
	SessionId <-chan string
	ProxyUrl  <-chan string
	Out       chan<- context.Request

	RetryCount <-chan float64

	RetryIn  <-chan context.Request
	RetryOut chan<- context.Request
}

func (c *ApiAi) Init() {
	c.CtxQueue = lane.NewQueue()
	c.TxtQueue = lane.NewQueue()
}

func (c *ApiAi) OnIn(req context.Request) {
	text := req.Text
	res := c.SendQuery(text)

	intent := res.Metadata.IntentName
	score := res.Score
	data, _ := json.Marshal(res)

	log.Printf("意图解析\"%s\" -> %s 准确度: %2.2f%%\n结果:%v", text, intent, score*100, string(data))

	req.ApiAiResult = res

	c.Out <- req
}

func (c *ApiAi) OnToken(token string) {
	c.token = token
}

func (c *ApiAi) OnRetryCount(count float64) {
	c.retryCount = int(count)
}

func (c *ApiAi) OnSessionId(sessionId string) {
	c.sessionId = sessionId
}

func (c *ApiAi) OnProxyUrl(proxy string) {
	c.proxyUrl = proxy
}

func (c *ApiAi) OnRetryIn(req context.Request) {
	text := req.Text
	res := c.SendQuery(text)

	intent := res.Metadata.IntentName
	score := res.Score
	data, _ := json.Marshal(res)

	log.Printf("意图解析\"%s\" -> %s 准确度: %2.2f%%\n结果:%v", text, intent, score*100, string(data))

	req.ApiAiResult = res

	c.RetryOut <- req
}

func (c *ApiAi) SendQuery(txt string) apiai.Result {
	count := 0

	tBegin := time.Now()

	res, err := ApiAiQuery(txt, c.token, c.sessionId, c.proxyUrl)

	tEnd := time.Now()
	log.Printf("ApiAi request cost time %v", tEnd.Sub(tBegin))

	for err != nil && count < c.retryCount {
		count++

		log.Printf("意图\"%s\"重新解析失败:%s", txt, err.Error())
		log.Printf("尝试第%v/%v次重新解析", count, c.retryCount)

		res, err = ApiAiQuery(txt, c.token, c.sessionId, c.proxyUrl)
		if err != nil {
			log.Printf("意图\"%s\"再第%v/%v次重新解析失败:%s", txt, count, c.retryCount, err.Error())
		}
	}

	return res
}
