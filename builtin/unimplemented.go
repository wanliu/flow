package builtin

import (
	"log"

	"github.com/hysios/apiai-go"
	"github.com/wanliu/flow/context"
)

type Unimplemented struct {
	TryGetEntities
	Ctx  <-chan context.Context
	Type <-chan string
	Out  chan<- ReplyData
}

func NewUnimplemented() interface{} {
	return new(Unimplemented)
}

func (order *Unimplemented) OnCtx(ctx context.Context) {
	if context.GroupChat(ctx) {
		log.Printf("不回应非开单相关的普通群聊")
		return
	}

	output := "你好，请问有什么可以帮您？"
	aiResult := ctx.Value("Result").(apiai.Result)

	if r := aiResult.Fulfillment.Speech; r != "" {
		output = r
	}

	replyData := ReplyData{output, ctx}
	order.Out <- replyData
}
