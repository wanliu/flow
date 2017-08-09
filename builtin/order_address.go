package builtin

import (
	// "log"

	"github.com/hysios/apiai-go"
	"github.com/wanliu/flow/builtin/ai"

	. "github.com/wanliu/flow/builtin/config"
	. "github.com/wanliu/flow/builtin/resolves"
	. "github.com/wanliu/flow/context"
)

type OrderAddress struct {
	TryGetEntities

	expMins float64

	Ctx <-chan Context
	Out chan<- ReplyData
}

func NewOrderAddress() interface{} {
	return new(OrderAddress)
}

func (c *OrderAddress) OnCtx(ctx Context) {
	currentOrder := ctx.Value(CtxKeyOrder)

	if nil != currentOrder {
		aiResult := ctx.Value("Result").(apiai.Result)

		params := ai.ApiAiOrder{AiResult: aiResult}
		address := params.Address()
		customer := params.Customer()

		cOrder := currentOrder.(OrderResolve)

		if address != "" {
			cOrder.Address = address
		}

		if customer != "" {
			cOrder.Customer = customer
		}

		if cOrder.Fulfiled() {
			ctx.SetValue(CtxKeyOrder, nil)
			ctx.SetValue(CtxKeyLastOrder, cOrder)
		}

		reply := "收到客户/地址信息：" + address + customer + "\n" + cOrder.Answer()
		c.Out <- ReplyData{reply, ctx}
	} else {
		c.Out <- ReplyData{"地址输入无效，当前没有正在进行中的订单", ctx}
	}
}
