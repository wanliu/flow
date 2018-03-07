package builtin

import (
	"fmt"

	"github.com/wanliu/flow/builtin/config"
	"github.com/wanliu/flow/builtin/resolves"
	"github.com/wanliu/flow/context"

	flow "github.com/wanliu/goflow"
)

type OrderItemDelete struct {
	flow.Component

	Ctx <-chan context.Request
	Out chan<- context.Request
}

func NewOrderItemDelete() interface{} {
	return new(OrderItemDelete)
}

func (c *OrderItemDelete) OnCtx(req context.Request) {
	ctx := req.Ctx
	currentOrder := ctx.CtxValue(config.CtxKeyOrder)

	if nil != currentOrder {

		cOrder := currentOrder.(resolves.OrderResolve)

		if cOrder.Expired(config.SesssionExpiredMinutes) {
			req.Res = context.Response{"会话已经过时，当前没有正在进行中的订单", ctx, nil}
			c.Out <- req
			return
		}

		cmd := req.Command
		if cmd != nil {
			data := cmd.Data
			if itemName, ok := data["itemName"].(string); ok {
				removed := cOrder.Products.Remove(itemName)
				if removed {
					_, d := cOrder.Answer(ctx)

					data := map[string]interface{}{
						"type":   "info",
						"on":     "order",
						"action": "update",
						"data":   d,
					}

					reply := fmt.Sprintf("已经删除%v", itemName)
					req.Res = context.Response{reply, ctx, data}
					c.Out <- req
				} else {
					reply := fmt.Sprintf("无效的操作，%v不存在", itemName)
					req.Res = context.Response{reply, ctx, nil}
					c.Out <- req
				}
			} else {
				req.Res = context.Response{"无效的删除操作", ctx, nil}
				c.Out <- req
			}
		}
	} else {
		req.Res = context.Response{"无效的操作，当前没有正在进行中的订单", ctx, nil}
		c.Out <- req
	}
}
