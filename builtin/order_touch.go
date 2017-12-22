package builtin

import (
	"log"

	. "github.com/wanliu/flow/builtin/resolves"
	. "github.com/wanliu/flow/context"

	config "github.com/wanliu/flow/builtin/config"
)

type OrderTouch struct {
	TryGetEntities

	mins int

	Ctx     <-chan Context
	Next    chan<- Context
	Timeout chan<- Context
}

func NewOrderTouch() interface{} {
	return new(OrderTouch)
}

func (c *OrderTouch) OnCtx(ctx Context) {
	order := ctx.CtxValue(config.CtxKeyOrder)

	if order != nil {
		cOrder := order.(OrderResolve)
		log.Printf("[Update] Current order touched.")
		cOrder.Touch()
		ctx.SetCtxValue(config.CtxKeyOrder, cOrder)

		c.Timeout <- ctx
	}

	c.Next <- ctx
}
