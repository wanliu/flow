package builtin

import (
	"log"

	"github.com/wanliu/flow/builtin/wechat_type"
	. "github.com/wanliu/flow/context"
)

type Critical struct {
	TryGetEntities
	Ctx  <-chan Context
	Type <-chan string
	Out  chan<- ReplyData
}

func NewCritical() interface{} {
	return new(Critical)
}

// entity: 贬低
func (order *Critical) OnCtx(ctx Context) {
	if wechat_type.GroupChat(ctx) {
		log.Printf("不回应非开单相关的普通群聊")
		return
	}

	// entities := ctx.Value("Result").(ResultParams).Entities
	output := "对不起，辜负了您的期望，请给我们时间，我们会改进的"

	replyData := ReplyData{output, ctx}
	order.Out <- replyData
}
