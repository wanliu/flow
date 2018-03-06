package builtin

import (
	"log"

	"github.com/wanliu/flow/context"
)

type Abuse struct {
	TryGetEntities
	Ctx  <-chan context.Request
	Type <-chan string
	Out  chan<- context.Request
}

func NewAbuse() interface{} {
	return new(Abuse)
}

// entity: 贬低
func (order *Abuse) OnCtx(req context.Request) {
	ctx := req.Ctx

	if context.GroupChat(ctx) {
		log.Printf("不回应非开单相关的普通群聊")
		return
	}
	// entities := ctx.Value("Result").(ResultParams).Entities
	output := "请不要脏话哦"

	req.Res = context.Response{output, ctx, nil}
	order.Out <- req
}
