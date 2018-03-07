package builtin

import (
	"fmt"
	"log"

	// "github.com/hysios/apiai-go"
	"github.com/wanliu/flow/builtin/ai"
	"github.com/wanliu/flow/builtin/resolves"
	"github.com/wanliu/flow/context"

	"github.com/wanliu/flow/builtin/config"
)

type OrderAddress struct {
	TryGetEntities

	expMins      float64
	confirmScore float64

	Ctx          <-chan context.Request
	ConfirmScore <-chan float64
	Out          chan<- context.Request
}

func NewOrderAddress() interface{} {
	return new(OrderAddress)
}

func (c *OrderAddress) OnConfirmScore(score float64) {
	c.confirmScore = score
}

func (c *OrderAddress) OnCtx(req context.Request) {
	// if context.GroupChat(ctx) {
	// 	log.Printf("不回应非开单相关的普通群聊")
	// 	return
	// }

	ctx := req.Ctx
	currentOrder := ctx.CtxValue(config.CtxKeyOrder)

	if nil != currentOrder {
		aiResult := req.ApiAiResult

		cOrder := currentOrder.(resolves.OrderResolve)

		if cOrder.Expired(config.SesssionExpiredMinutes) {
			req.Res = context.Response{"会话已经过时，当前没有正在进行中的订单", ctx, nil}
			c.Out <- req
			return
		}

		params := ai.ApiAiOrder{AiResult: aiResult}

		if c.confirmScore != 0 && params.Score() >= c.confirmScore {
			address := params.Address()
			customer := params.Customer()

			if address != "" {
				cOrder.Address = address
			}

			if customer != "" {
				cOrder.ExtractedCustomer = customer
				cOrder.CheckExtractedCustomer()
			}

			reply, d := cOrder.Answer(ctx)
			reply = fmt.Sprintf("收到客户/地址信息：%v%v\n%v", address, customer, reply)
			data := map[string]interface{}{
				"type":   "info",
				"on":     "order",
				"action": "update",
				"data":   d,
			}
			req.Res = context.Response{reply, ctx, data}
			c.Out <- req

			if cOrder.Resolved() {
				ctx.SetCtxValue(config.CtxKeyOrder, nil)
				ctx.SetCtxValue(config.CtxKeyLastOrder, cOrder)
			} else if cOrder.Failed() {
				ctx.SetCtxValue(config.CtxKeyOrder, nil)
			}
		} else {
			var values []string

			query := params.Query()
			customer := params.Customer()

			if customer == "" {
				values = []string{query}
			} else {
				values = []string{customer, query}
			}

			addressConfirm := resolves.AddressConfirm{Values: values}

			ctx.SetCtxValue(config.CtxKeyConfirm, addressConfirm)

			reply := "收到您的回复:" + query + "\n"
			reply = reply + addressConfirm.Notice(ctx)
			req.Res = context.Response{reply, ctx, nil}
			c.Out <- req
		}
	} else {
		if context.GroupChat(ctx) {
			log.Printf("不回应群聊无效的客户输入")
			return
		}

		req.Res = context.Response{"客户输入无效，当前没有正在进行中的订单", ctx, nil}
		c.Out <- req
	}
}
