package builtin

import (
	flow "github.com/wanliu/goflow"
)

func NewFinal() interface{} {
	return new(Final)
}

type Final struct {
	flow.Component
	In <-chan ReplyData
}

func (s *Final) OnIn(data ReplyData) {
	data.Ctx.Post(data.Reply)
}
