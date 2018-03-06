package builtin

import (
	"time"

	"github.com/wanliu/flow/builtin/resolves"
	"github.com/wanliu/flow/context"

	config "github.com/wanliu/flow/builtin/config"
)

type OrderTimeout struct {
	TryGetEntities

	mins int

	Ctx <-chan context.Request
	// Mins <-chan float64
	Out chan<- context.Request
}

func NewOrderTimeout() interface{} {
	return new(OrderTimeout)
}

// func (c *OrderTimeout) OnMins(t float64) {
// 	c.mins = int(t)
// }

func (c *OrderTimeout) OnCtx(req context.Request) {
	go func() {
		ctx := req.Ctx

		expiredMins := config.SesssionExpiredMinutes
		settedMins := ctx.CtxValue(config.CtxKeyExpiredMinutes)

		if settedMins != nil {
			expiredMins = settedMins.(int)
		}

		time.Sleep(time.Duration(expiredMins) * time.Minute)

		order := ctx.CtxValue(config.CtxKeyOrder)

		if order != nil {
			cOrder := order.(resolves.OrderResolve)
			if cOrder.Expired(expiredMins) {
				ctx.SetCtxValue(config.CtxKeyOrder, nil)
				req.Res = context.Response{"由于长时间未操作完成，当前订单已经失效", ctx, nil}
				c.Out <- req
			}
		}
	}()
}
