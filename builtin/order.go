package builtin

import (
	// "log"

	. "github.com/wanliu/flow/builtin/resolves"
	. "github.com/wanliu/flow/context"
)

type Order struct {
	TryGetEntities
	Ctx   <-chan Context
	New   chan<- ReplyData
	Patch chan<- ReplyData
}

func GetOrder() interface{} {
	return new(Order)
}

func (c *Order) OnCtx(ctx Context) {
	currentOrder := ctx.Value("Order")
	if nil != currentOrder {
		cOrder = currentOrder.(OrderResolve)

		if cOrder.Modifable() {
			c.Patch <- ctx
		} else {
			c.New <- ctx
		}
	} else {
		c.New <- ctx
	}
}
