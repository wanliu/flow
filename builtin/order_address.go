package builtin

import (
	"log"

	"github.com/hysios/apiai-go"
	"github.com/wanliu/flow/builtin/ai"
	"github.com/wanliu/flow/builtin/resolves"
	"github.com/wanliu/flow/context"

	config "github.com/wanliu/flow/builtin/config"
)

type OrderAddress struct {
	TryGetEntities

	expMins      float64
	confirmScore float64

	Ctx          <-chan context.Context
	ConfirmScore <-chan float64
	Out          chan<- ReplyData
}

func NewOrderAddress() interface{} {
	return new(OrderAddress)
}

func (c *OrderAddress) OnConfirmScore(score float64) {
	c.confirmScore = score
}

func (c *OrderAddress) OnCtx(ctx context.Context) {
	// if context.GroupChat(ctx) {
	// 	log.Printf("不回应非开单相关的普通群聊")
	// 	return
	// }

	currentOrder := ctx.CtxValue(config.CtxKeyOrder)

	if nil != currentOrder {
		aiResult := ctx.Value(config.ValueKeyResult).(apiai.Result)
		cOrder := currentOrder.(resolves.OrderResolve)

		if cOrder.Expired(config.SesssionExpiredMinutes) {
			c.Out <- ReplyData{"会话已经过时，当前没有正在进行中的订单", ctx, nil}
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
				cOrder.Customer = customer
			}

			reply := "收到客户/地址信息：" + address + customer + "\n" + cOrder.Answer(ctx)
			c.Out <- ReplyData{reply, ctx, nil}

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
			c.Out <- ReplyData{reply, ctx, nil}
		}
	} else {
		if context.GroupChat(ctx) {
			log.Printf("不回应群聊无效的客户输入")
			return
		}

		c.Out <- ReplyData{"客户输入无效，当前没有正在进行中的订单", ctx, nil}
	}
}
