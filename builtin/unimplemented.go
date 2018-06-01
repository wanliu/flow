package builtin

import (
	"log"

	// "github.com/hysios/apiai-go"
	"github.com/wanliu/flow/context"
)

type Unimplemented struct {
	TryGetEntities
	Ctx  <-chan context.Request
	Type <-chan string
	Out  chan<- context.Request
}

func NewUnimplemented() interface{} {
	return new(Unimplemented)
}

func (c *Unimplemented) OnCtx(req context.Request) {
	ctx := req.Ctx

	if req.IsCommand() {
		output := "指令无法识别"

		req.Res = context.Response{
			Reply: output,
			Ctx:   ctx,
		}

		c.Out <- req
		return
	}

	if context.GroupChat(ctx) {
		log.Printf("不回应非开单相关的普通群聊")
		return
	}

	output := "你好，请问有什么可以帮您？"
	aiResult := req.ApiAiResult

	if r := aiResult.Fulfillment.Speech; r != "" {
		output = r
	}

	req.Res = context.Response{
		Reply: output,
		Ctx:   ctx,
	}
	c.Out <- req
}
