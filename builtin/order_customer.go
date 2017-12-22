package builtin

import (
	"log"

	"github.com/hysios/apiai-go"
	"github.com/wanliu/flow/builtin/ai"
	"github.com/wanliu/flow/context"

	. "github.com/wanliu/flow/builtin/config"
	. "github.com/wanliu/flow/builtin/resolves"
)

type OrderCustomer struct {
	TryGetEntities

	expMins float64

	Ctx <-chan context.Context
	Out chan<- ReplyData
}

func NewOrderCustomer() interface{} {
	return new(OrderCustomer)
}

func (c *OrderCustomer) OnCtx(ctx context.Context) {
	if context.GroupChat(ctx) {
		log.Printf("不回应非开单相关的普通群聊")
		return
	}

	currentOrder := ctx.Value(CtxKeyOrder)

	if nil != currentOrder {
		aiResult := ctx.Value("Result").(apiai.Result)

		params := ai.ApiAiOrder{AiResult: aiResult}
		customer := params.Customer()

		cOrder := currentOrder.(OrderResolve)
		cOrder.ExtractedCustomer = customer
		cOrder.CheckExtractedCustomer()

		reply := "收到客户信息：" + customer + "\n" + cOrder.Answer(ctx)
		c.Out <- ReplyData{reply, ctx}

		if cOrder.Resolved() {
			ctx.SetValue(CtxKeyOrder, nil)
			ctx.SetValue(CtxKeyLastOrder, cOrder)
		} else if cOrder.Failed() {
			ctx.SetValue(CtxKeyOrder, nil)
		}
	} else {
		c.Out <- ReplyData{"客户输入无效，当前没有正在进行中的订单", ctx}
	}
}
