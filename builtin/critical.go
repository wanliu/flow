package builtin

import (
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
	// entities := ctx.Value("Result").(ResultParams).Entities
	output := "对不起，辜负了您的期望，请给我们时间，我们会改进的"

	replyData := ReplyData{output, ctx}
	order.Out <- replyData
}
