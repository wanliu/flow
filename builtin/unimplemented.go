package builtin

import (
	"log"

	"github.com/hysios/apiai-go"
	. "github.com/wanliu/flow/context"
)

type Unimplemented struct {
	TryGetEntities
	Ctx  <-chan Context
	Type <-chan string
	Out  chan<- ReplyData
}

func NewUnimplemented() interface{} {
	return new(Unimplemented)
}

func (order *Unimplemented) OnCtx(ctx Context) {
	output := "你好，请问有什么可以帮您？"

	aiResult := ctx.Value("Result").(apiai.Result)

	if r := aiResult.Fulfillment.Speech; r != "" {
		output = r
	}

	log.Printf("OUTPUT: %v", output)

	replyData := ReplyData{output, ctx}
	order.Out <- replyData
}
