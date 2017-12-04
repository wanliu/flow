package builtin

import (
	// "github.com/wanliu/flow/builtin/ai"
	"github.com/wanliu/flow/builtin/config"
	"github.com/wanliu/flow/builtin/resolves"
	"github.com/wanliu/flow/context"
)

const (
	Per = 5
)

type CustomerOrders struct {
	TryGetEntities
	Type <-chan string

	Ctx  <-chan context.Context
	Page <-chan context.Context

	Out chan<- ReplyData
}

func NewCustomerOrders() interface{} {
	return new(CustomerOrders)
}

func (c *CustomerOrders) OnCtx(ctx context.Context) {
	rsv := resolves.NewCusOrdersResolve(ctx, Per)

	reply := rsv.Answer()

	rsv.Setup(ctx)

	c.Out <- ReplyData{reply, ctx}
}

func (c *CustomerOrders) OnPage(ctx context.Context) {
	in := ctx.Value(config.CtxKeyCusOrders)
	if in == nil {
		c.Out <- ReplyData{"当前没有正在进行的查询", ctx}
	} else {
		rsv := in.(*resolves.CustomerOrdersResolve)
		reply := rsv.Answer()
		c.Out <- ReplyData{reply, ctx}

		if rsv.Done {
			rsv.Clear(ctx)
		}
	}
}
