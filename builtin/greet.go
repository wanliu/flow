package builtin

import (
	. "github.com/wanliu/flow/context"
	"log"
)

type Greet struct {
	TryGetEntities
	Ctx  <-chan Context
	Type <-chan string
	Out  chan<- ReplyData
}

func NewGreet() interface{} {
	return new(Greet)
}

func (order *Greet) OnCtx(ctx Context) {
	// entities := ctx.Value("Result").(ResultParams).Entities
	params := ctx.Value("Result").(ResultParams)
	output := ""

	replies := map[string]string{
		"Hi":    "Hello, 请问你需要什么",
		"你好":    "你好，很高兴为您服务",
		"你好吗":   "我很好，你什么什么服务吗",
		"Hello": "Hello, 请问你需要什么",
		"天气不错":  "是啊，天气不错，可以外出喔",
	}

	text := params.Query
	r, hasKey := replies[text]

	if hasKey {
		output = r
	}
	// "default":

	if output == "" {
		output = "你好，很高兴为您服务"
	}

	log.Printf("OUTPUT: %v", output)

	replyData := ReplyData{output, ctx}
	order.Out <- replyData
}
