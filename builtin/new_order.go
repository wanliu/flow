package builtin

import (
	. "github.com/wanliu/flow/builtin/resolves"
	. "github.com/wanliu/flow/context"
)

type NewOrder struct {
	TryGetEntities
	DefTime string
	Ctx     <-chan Context
	Deftime <-chan string
	Out     chan<- ReplyData
	Notice  chan<- Context
}

func NewNewOrder() interface{} {
	return new(NewOrder)
}

// 默认送货时间
func (c *NewOrder) OnDeftime(t string) {
	c.DefTime = t
}

func (c *NewOrder) OnCtx(ctx Context) {
	orderResolve := NewOrderResolve(ctx)

	if c.DefTime != "" {
		orderResolve.SetDefTime(c.DefTime)
	}

	ctx.SetValue("Order", *orderResolve)

	output := ""

	if orderResolve.EmptyProducts() {
		output = "没有相关的产品"
	} else {
		output = orderResolve.Answer()
		c.Notice <- ctx
	}

	replyData := ReplyData{output, ctx}
	c.Out <- replyData
}
