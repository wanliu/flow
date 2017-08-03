package builtin

import (
	// "log"

	. "github.com/wanliu/flow/builtin/config"
	. "github.com/wanliu/flow/builtin/resolves"
	. "github.com/wanliu/flow/context"
)

type Order struct {
	TryGetEntities

	expMins float64

	Ctx           <-chan Context
	ExpireMinutes <-chan float64
	New           chan<- Context
	Patch         chan<- Context
	Out           chan<- ReplyData
}

func GetOrder() interface{} {
	return new(Order)
}

func (c *Order) OnCtx(ctx Context) {
	currentOrder := ctx.Value(CtxKeyOrder)

	if c.expMins != 0 {
		ctx.SetValue(CtxKeyExpiredMinutes, int(c.expMins))
	}

	if nil != currentOrder {
		cOrder := currentOrder.(OrderResolve)

		exMin := SesssionExpiredMinutes

		if c.expMins != 0 {
			exMin = int(c.expMins)
		}

		if cOrder.Modifable(exMin) {
			c.Patch <- ctx
		} else {
			c.New <- ctx
		}
	} else {
		c.New <- ctx
	}
}

func (c *Order) OnExpireMinutes(minutes float64) {
	c.expMins = minutes
}
