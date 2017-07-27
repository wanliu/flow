package builtin

import (
	. "github.com/wanliu/flow/context"
	flow "github.com/wanliu/goflow"
)

func NewPrintVoice() interface{} {
	return new(PrintVoice)
}

type PrintVoice struct {
	flow.Component
	Ctx <-chan Context
	Out chan<- ReplyData
}

func (s *PrintVoice) OnCtx(ctx Context) {
	txt := ctx.Value("Text").(string)
	reply := ReplyData{txt, ctx}
	s.Out <- reply
}
