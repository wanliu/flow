package builtin

import (
	// "log"

	"github.com/hysios/apiai-go"
	"github.com/wanliu/flow/builtin/ai"

	. "github.com/wanliu/flow/builtin/config"
	. "github.com/wanliu/flow/builtin/resolves"
	. "github.com/wanliu/flow/context"
)

type OrderCustomer struct {
	TryGetEntities

	expMins float64

	Ctx <-chan Context
	Out chan<- ReplyData
}

func NewOrderCustomer() interface{} {
	return new(OrderCustomer)
}

func (c *OrderCustomer) OnCtx(ctx Context) {
	currentOrder := ctx.Value(CtxKeyOrder)

	if nil != currentOrder {
		aiResult := ctx.Value("Result").(apiai.Result)

		params := ai.ApiAiOrder{AiResult: aiResult}
		customer := params.Customer()

		cOrder := currentOrder.(OrderResolve)
		cOrder.Customer = customer

		reply := "收到客户信息：" + customer + "\n" + cOrder.Answer()
		c.Out <- ReplyData{reply, ctx}

		if cOrder.Fulfiled() {
			ctx.SetValue(CtxKeyOrder, nil)
			ctx.SetValue(CtxKeyLastOrder, cOrder)
		}
	} else {
		c.Out <- ReplyData{"客户输入无效，当前没有正在进行中的订单", ctx}
	}
}
