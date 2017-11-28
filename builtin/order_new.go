package builtin

import (
	"github.com/wanliu/flow/builtin/config"
	"github.com/wanliu/flow/builtin/resolves"
	"github.com/wanliu/flow/context"
)

type NewOrder struct {
	TryGetEntities
	DefTime string
	Ctx     <-chan context.Context
	Deftime <-chan string
	Out     chan<- ReplyData
	Notice  chan<- context.Context
	Timeout chan<- context.Context
}

func NewNewOrder() interface{} {
	return new(NewOrder)
}

// 默认送货时间
func (c *NewOrder) OnDeftime(t string) {
	c.DefTime = t
}

func (c *NewOrder) OnCtx(ctx context.Context) {
	orderResolve := resolves.NewOrderResolve(ctx)

	if c.DefTime != "" {
		orderResolve.SetDefTime(c.DefTime)
	}

	output := ""

	if orderResolve.EmptyProducts() {
		output = "没有相关的产品"
	} else {
		output = orderResolve.Answer(ctx)

		if orderResolve.Resolved() {
			ctx.SetValue(config.CtxKeyLastOrder, *orderResolve)
			ctx.SetValue(config.CtxKeyOrder, nil)
		} else if orderResolve.Failed() {
			ctx.SetValue(config.CtxKeyOrder, nil)
		} else {
			ctx.SetValue(config.CtxKeyOrder, *orderResolve)
		}

		// c.Notice <- ctx
		c.Timeout <- ctx
	}

	replyData := ReplyData{output, ctx}
	c.Out <- replyData
}
