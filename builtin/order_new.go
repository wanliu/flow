package builtin

import (
	"log"

	"github.com/wanliu/flow/builtin/config"
	"github.com/wanliu/flow/builtin/resolves"

	"github.com/wanliu/flow/context"
)

type NewOrder struct {
	TryGetEntities
	DefTime    string
	retryCount int

	Ctx     <-chan context.Request
	Deftime <-chan string
	Out     chan<- context.Request
	Notice  chan<- context.Request
	Timeout chan<- context.Request

	RetryOut chan<- context.Request
	RetryIn  <-chan context.Request

	RetryCount <-chan float64
}

func NewNewOrder() interface{} {
	return new(NewOrder)
}

// 默认送货时间
func (c *NewOrder) OnDeftime(t string) {
	c.DefTime = t
}

func (c *NewOrder) OnRetryCount(count float64) {
	c.retryCount = int(count)
}

func (c *NewOrder) OnCtx(req context.Request) {
	ctx := req.Ctx
	orderResolve := resolves.NewOrderResolve(req)

	if c.DefTime != "" {
		orderResolve.SetDefTime(c.DefTime)
	}

	output := ""

	if orderResolve.EmptyProducts() {
		if c.retryCount > 0 {
			log.Printf("重新获取开单产品，第1次，共%v次", c.retryCount)
			c.RetryOut <- req
		} else {
			if context.GroupChat(ctx) {
				c.GroupAnswer(req, orderResolve)
				return
			}

			output = "没有相关的产品"
			req.Res = context.Response{
				Reply: output,
				Ctx:   ctx,
			}
			c.Out <- req
		}
	} else {
		reply, d := orderResolve.Answer(ctx)

		if orderResolve.Resolved() {
			ctx.SetCtxValue(config.CtxKeyLastOrder, &orderResolve)
			ctx.SetCtxValue(config.CtxKeyOrder, nil)
		} else if orderResolve.Failed() {
			ctx.SetCtxValue(config.CtxKeyOrder, nil)
		} else if orderResolve.MismatchQuantity() {
			ctx.SetCtxValue(config.CtxKeyOrder, nil)
		} else {
			ctx.SetCtxValue(config.CtxKeyOrder, &orderResolve)
		}

		c.Timeout <- req

		data := map[string]interface{}{
			"type":   "info",
			"on":     "order",
			"action": "create",
			"data":   d,
		}

		req.Res = context.Response{
			Reply: reply,
			Ctx:   ctx,
			Data:  data,
		}
		c.Out <- req
	}
}

func (c *NewOrder) OnRetryIn(req context.Request) {
	ctx := req.Ctx
	orderResolve := resolves.NewOrderResolve(req)

	if c.DefTime != "" {
		orderResolve.SetDefTime(c.DefTime)
	}

	if orderResolve.EmptyProducts() {
		retriedCount := 1
		retriedCountInt := ctx.Value(config.CtxKeyRetriedCount)

		if retriedCountInt != nil {
			retriedCount = retriedCountInt.(int)
		}

		if retriedCount >= c.retryCount {
			if context.GroupChat(ctx) {
				c.GroupAnswer(req, orderResolve)
				return
			}

			output := "没有相关的产品"

			req.Res = context.Response{
				Reply: output,
				Ctx:   ctx,
			}
			c.Out <- req
		} else {
			retriedCount++
			log.Printf("重新获取开单产品，第%v次，共%v次", retriedCount, c.retryCount)

			ctx.SetValue(config.CtxKeyRetriedCount, retriedCount)
			c.RetryOut <- req
		}
	} else {
		reply, d := orderResolve.Answer(ctx)

		if orderResolve.Resolved() {
			ctx.SetCtxValue(config.CtxKeyLastOrder, &orderResolve)
			ctx.SetCtxValue(config.CtxKeyOrder, nil)
		} else if orderResolve.Failed() {
			ctx.SetCtxValue(config.CtxKeyOrder, nil)
		} else if orderResolve.MismatchQuantity() {
			ctx.SetCtxValue(config.CtxKeyOrder, nil)
		} else {
			ctx.SetCtxValue(config.CtxKeyOrder, &orderResolve)
		}

		// c.Notice <- ctx
		c.Timeout <- req

		data := map[string]interface{}{
			"type":   "info",
			"on":     "order",
			"action": "create",
			"data":   d,
		}

		req.Res = context.Response{
			Reply: reply,
			Ctx:   ctx,
			Data:  data,
		}
		c.Out <- req
	}

}

func (c *NewOrder) GroupAnswer(req context.Request, orderResolve *resolves.OrderResolve) {
	ctx := req.Ctx
	reply, d := orderResolve.Answer(ctx)

	data := map[string]interface{}{
		"type":   "info",
		"on":     "order",
		"action": "create",
		"data":   d,
	}

	if orderResolve.Fulfiled() {
		req.Res = context.Response{
			Reply: reply,
			Ctx:   ctx,
			Data:  data,
		}
		c.Out <- req
	} else {
		log.Printf("群聊开单失败, 取消回复。失败原因：%v", reply)
	}
}
