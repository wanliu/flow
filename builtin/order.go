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

	Ctx           <-chan context.Request
	ExpireMinutes <-chan float64
	New           chan<- context.Request
	Patch         chan<- context.Request
	Out           chan<- context.Request
	SyncQueue     <-chan string
}

func GetOrder() interface{} {
	return new(Order)
}

func (c *Order) OnSyncQueue(queue string) {
	c.OrderSyncQueue = queue
}

func (c *Order) OnCtx(req context.Request) {
	ctx := req.Ctx

	if c.OrderSyncQueue != "" {
		ctx.SetValue(config.CtxKeySyncQueue, c.OrderSyncQueue)
	}

	if c.expMins != 0 {
		ctx.SetCtxValue(config.CtxKeyExpiredMinutes, int(c.expMins))
	}

	if context.GroupChat(ctx) {
		c.New <- req
		return
	}

	// currentOrder := ctx.CtxValue(config.CtxKeyOrder)
	orderRsv := resolves.GetCtxOrder(ctx)

	if nil != orderRsv {
		exMin := config.SesssionExpiredMinutes

		if c.expMins != 0 {
			exMin = int(c.expMins)
		}

		if orderRsv.Modifable(exMin) {
			c.Patch <- req
		} else {
			c.New <- req
		}
	} else {
		c.New <- req
	}
}

func (c *Order) OnExpireMinutes(minutes float64) {
	c.expMins = minutes
}
