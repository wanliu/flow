package builtin

import (
	"time"

	. "github.com/wanliu/flow/builtin/resolves"
	. "github.com/wanliu/flow/context"

	config "github.com/wanliu/flow/builtin/config"
)

type OrderTimeout struct {
	mins int

	Ctx  <-chan Context
	Mins <-chan float64
	Out  chan<- ReplyData
}

func NewOrderTimeout() interface{} {
	return new(OrderTimeout)
}

func (c *OrderTimeout) OnMins(t float64) {
	c.mins = int(t)
}

func (c *OrderTimeout) OnCtx(ctx Context) {
	go func() {
		expiredMins := config.SesssionExpiredMinutes
		if c.mins != 0 {
			expiredMins = c.mins
		}

		time.Sleep(time.Duration(expiredMins) * time.Minute)

		order := ctx.Value(config.CtxKeyOrder)

		if order != nil {
			cOrder := order.(OrderResolve)
			if cOrder.Expired(expiredMins) {
				ctx.SetValue(config.CtxKeyOrder, nil)
				c.Out <- ReplyData{"由于长时间未操作完成，当前订单已经失效", ctx}
			}
		}
	}()
}
