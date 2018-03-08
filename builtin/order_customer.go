package builtin

import (
	"log"

	"github.com/hysios/apiai-go"
	"github.com/wanliu/flow/builtin/ai"
	"github.com/wanliu/flow/context"

	"github.com/wanliu/flow/builtin/config"
	"github.com/wanliu/flow/builtin/resolves"
)

type OrderCustomer struct {
	TryGetEntities

	expMins float64

	Ctx <-chan context.Request
	Out chan<- context.Request
}

func NewOrderCustomer() interface{} {
	return new(OrderCustomer)
}

func (c *OrderCustomer) OnCtx(req context.Request) {
	ctx := req.Ctx
	orderRsv := resolves.GetCtxOrder(ctx)

	if nil != orderRsv {
		aiResult := ctx.Value("Result").(apiai.Result)

		params := ai.ApiAiOrder{AiResult: aiResult}
		customer := params.Customer()

		orderRsv.ExtractedCustomer = customer
		orderRsv.CheckExtractedCustomer()

		reply, d := orderRsv.Answer(ctx)
		reply = "收到客户信息：" + customer + "\n" + reply
		data := map[string]interface{}{
			"type":   "info",
			"on":     "order",
			"action": "update",
			"data":   d,
		}

		req.Res = context.Response{reply, ctx, data}
		c.Out <- req

		if orderRsv.Resolved() {
			ctx.SetCtxValue(config.CtxKeyOrder, nil)
			ctx.SetCtxValue(config.CtxKeyLastOrder, orderRsv)
		} else if orderRsv.Failed() {
			ctx.SetCtxValue(config.CtxKeyOrder, nil)
		}
	} else {
		if context.GroupChat(ctx) {
			log.Printf("不回应单独输入客户的的普通群聊")
			return
		}

		req.Res = context.Response{"客户输入无效，当前没有正在进行中的订单", ctx, nil}
		c.Out <- req
	}
}
