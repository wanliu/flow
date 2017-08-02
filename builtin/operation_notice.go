package builtin

import (
	. "github.com/wanliu/flow/context"
	flow "github.com/wanliu/goflow"
)

type OperationNotice struct {
	flow.Component
	Ctx <-chan Context
	Out chan<- ReplyData
}

func NewOperationNotice() interface{} {
	return new(OperationNotice)
}

func (s OperationNotice) OnCtx(ctx Context) {
	s.Out <- ReplyData{"你可以继续提交产品到订单，也可以立刻取消当前任务（3分钟以内）", ctx}
}
