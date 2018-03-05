package builtin

import (
	"log"

	"github.com/wanliu/flow/builtin/config"
	"github.com/wanliu/flow/builtin/resolves"
	"github.com/wanliu/flow/context"
)

type OrderTouch struct {
	TryGetEntities

	mins int

	Ctx     <-chan context.Request
	Next    chan<- context.Request
	Timeout chan<- context.Request
}

func NewOrderTouch() interface{} {
	return new(OrderTouch)
}

func (c *OrderTouch) OnCtx(req context.Request) {
	ctx := req.Ctx

	order := ctx.CtxValue(config.CtxKeyOrder)

	if order != nil {
		cOrder := order.(resolves.OrderResolve)
		log.Printf("[Update] Current order touched.")
		cOrder.Touch()
		ctx.SetCtxValue(config.CtxKeyOrder, cOrder)

		c.Timeout <- req
	}

	c.Next <- req
}
