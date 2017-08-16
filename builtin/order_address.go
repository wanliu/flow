package builtin

import (
	// "log"

	"github.com/hysios/apiai-go"
	"github.com/wanliu/flow/builtin/ai"

	. "github.com/wanliu/flow/builtin/resolves"
	. "github.com/wanliu/flow/context"

	config "github.com/wanliu/flow/builtin/config"
)

type OrderAddress struct {
	TryGetEntities

	expMins      float64
	confirmScore float64

	Ctx          <-chan Context
	ConfirmScore <-chan float64
	Out          chan<- ReplyData
}

func NewOrderAddress() interface{} {
	return new(OrderAddress)
}

func (c *OrderAddress) OnConfirmScore(score float64) {
	c.confirmScore = score
}

func (c *OrderAddress) OnCtx(ctx Context) {
	currentOrder := ctx.Value(config.CtxKeyOrder)

	if nil != currentOrder {
		aiResult := ctx.Value(config.ValueKeyResult).(apiai.Result)
		cOrder := currentOrder.(OrderResolve)

		if cOrder.Expired(config.SesssionExpiredMinutes) {
			c.Out <- ReplyData{"当前没有正在进行中的订单", ctx}
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

			if cOrder.Fulfiled() {
				ctx.SetValue(config.CtxKeyOrder, nil)
				ctx.SetValue(config.CtxKeyLastOrder, cOrder)
			}

			reply := "收到客户/地址信息：" + address + customer + "\n" + cOrder.Answer()
			c.Out <- ReplyData{reply, ctx}
		} else {
			query := params.Query()
			address := params.Address() + params.Customer()

			if address == "" {
				address = query
			}

			addressConfirm := AddressConfirm{Order: &cOrder, Value: address}

			ctx.SetValue(config.CtxKeyConfirm, addressConfirm)

			reply := "收到您的回复:" + query + "\n"
			reply = reply + "是否将 “" + address + "” 做为收货地址?"
			c.Out <- ReplyData{reply, ctx}
		}
	} else {
		c.Out <- ReplyData{"地址输入无效，当前没有正在进行中的订单", ctx}
	}
}
