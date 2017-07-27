package builtin

import (
	. "github.com/wanliu/flow/context"
)

type PrintVoice struct {
	TryGetEntities
	Ctx <-chan Context
	Out chan<- ReplyData
}

func NewPrintVoice() interface{} {
	return new(PrintVoice)
}

func (order *PrintVoice) OnCtx(ctx Context) {
	text := ctx.Value("Text").(string)

	replyData := ReplyData{text, ctx}
	order.Out <- replyData
}
