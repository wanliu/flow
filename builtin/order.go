package builtin

import (
	// "log"

	"github.com/wanliu/flow/builtin/config"
	"github.com/wanliu/flow/builtin/resolves"
	"github.com/wanliu/flow/context"
)

type Order struct {
	TryGetEntities

	expMins        float64
	OrderSyncQueue string

	Ctx           <-chan context.Context
	ExpireMinutes <-chan float64
	New           chan<- context.Context
	Patch         chan<- context.Context
	Out           chan<- ReplyData
	SyncQueue     <-chan string
}

func GetOrder() interface{} {
	return new(Order)
}

func (c *Order) OnSyncQueue(queue string) {
	c.OrderSyncQueue = queue
}

func (c *Order) OnCtx(ctx context.Context) {
	currentOrder := ctx.Value(config.CtxKeyOrder)

	if c.expMins != 0 {
		ctx.SetValue(config.CtxKeyExpiredMinutes, int(c.expMins))
	}

	if c.OrderSyncQueue != "" {
		ctx.SetValue(config.CtxKeySyncQueue, c.OrderSyncQueue)
	}

	if nil != currentOrder {
		cOrder := currentOrder.(resolves.OrderResolve)

		exMin := config.SesssionExpiredMinutes

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
