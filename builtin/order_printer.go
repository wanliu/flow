package builtin

import (
	"log"

	. "github.com/wanliu/flow/builtin/resolves"
	"github.com/wanliu/flow/context"

	config "github.com/wanliu/flow/builtin/config"
	flow "github.com/wanliu/goflow"
)

func NewOrderPrinter() interface{} {
	return new(OrderPrinter)
}

type OrderPrinter struct {
	flow.Component

	Ctx <-chan context.Context
	Out chan<- ReplyData
}

func (s *OrderPrinter) OnCtx(ctx context.Context) {
	if context.GroupChat(ctx) {
		log.Printf("不回应非开单相关的普通群聊")
		return
	}

	currentOrder := ctx.Value(config.CtxKeyOrder)

	if nil == currentOrder {
		s.Out <- ReplyData{"当前没有正在进行中的订单", ctx}
	} else {
		curOrder := currentOrder.(OrderResolve)

		if curOrder.Expired(config.SesssionExpiredMinutes) {
			s.Out <- ReplyData{"当前没有正在进行中的订单", ctx}
		} else {
			orderDetail := "-----------订单详情-------------\n"
			orderDetail = orderDetail + curOrder.AnswerBody()
			s.Out <- ReplyData{orderDetail, ctx}
		}
	}
}
