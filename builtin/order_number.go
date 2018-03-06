package builtin

import (
	"log"
	// "time"

	"github.com/hysios/apiai-go"
	"github.com/wanliu/flow/builtin/config"
	"github.com/wanliu/flow/builtin/resolves"

	"github.com/wanliu/flow/context"

	flow "github.com/wanliu/goflow"
)

func NewOrderNumber() interface{} {
	return new(OrderNumber)
}

type OrderNumber struct {
	flow.Component

	Ctx <-chan context.Request
	Out chan<- context.Request
}

func (c *OrderNumber) OnCtx(req context.Request) {
	ctx := req.Ctx
	if context.GroupChat(ctx) {
		log.Printf("不回应非开单相关的普通群聊")
		return
	}

	aiResult := ctx.Value("Result").(apiai.Result)

	if numInt, exist := aiResult.Params["order-numder"]; exist {
		orderNo := numInt.(string)

		resolveInt := ctx.CtxValue(config.CtxKeyOrderNum)
		if resolveInt != nil {
			resolve := resolveInt.(resolves.OrderNumberResolve)
			reply := resolve.Resolve(orderNo, ctx)
			req.Res = context.Response{reply, ctx, nil}
			c.Out <- req
		} else {
			// if context.GroupChat(ctx) {
			// 	log.Printf("不回应非开单相关的普通群聊")
			// 	return
			// }

			req.Res = context.Response{"接收到订单号输入，但是没有对应的操作哦", ctx, nil}
			c.Out <- req
		}

	} else {
		req.Res = context.Response{"无效的订单号输入", ctx, nil}
		c.Out <- req
	}
}
