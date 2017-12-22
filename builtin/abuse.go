package builtin

import (
	"log"

	"github.com/wanliu/flow/context"
)

type Abuse struct {
	TryGetEntities
	Ctx  <-chan context.Context
	Type <-chan string
	Out  chan<- ReplyData
}

func NewAbuse() interface{} {
	return new(Abuse)
}

// entity: 贬低
func (order *Abuse) OnCtx(ctx context.Context) {
	if context.GroupChat(ctx) {
		log.Printf("不回应非开单相关的普通群聊")
		return
	}
	// entities := ctx.Value("Result").(ResultParams).Entities
	output := "请不要脏话哦"

	replyData := ReplyData{output, ctx}
	order.Out <- replyData
}
