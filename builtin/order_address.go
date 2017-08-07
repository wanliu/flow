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

		cOrder := currentOrder.(OrderResolve)
		cOrder.Address = address

		if cOrder.Fulfiled() {
			ctx.SetValue(CtxKeyOrder, nil)
			ctx.SetValue(CtxKeyLastOrder, cOrder)
		}

		reply := "收到地址信息：" + address + "\n" + cOrder.Answer()
		c.Out <- ReplyData{reply, ctx}
	} else {
		c.Out <- ReplyData{"地址输入无效，当前没有正在进行中的订单", ctx}
	}
}
