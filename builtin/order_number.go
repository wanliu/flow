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

	Ctx <-chan context.Context
	Out chan<- ReplyData
}

func (c *OrderNumber) OnCtx(ctx context.Context) {
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
			c.Out <- ReplyData{reply, ctx}
		} else {
			c.Out <- ReplyData{"接收到订单号输入，但是没有对应的操作哦", ctx}
		}

	} else {
		c.Out <- ReplyData{"无效的订单号输入", ctx}
	}
}
