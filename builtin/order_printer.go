package builtin

import (
	. "github.com/wanliu/flow/builtin/config"
	. "github.com/wanliu/flow/builtin/resolves"
	. "github.com/wanliu/flow/context"
	flow "github.com/wanliu/goflow"
)

func NewOrderPrinter() interface{} {
	return new(OrderPrinter)
}

type OrderPrinter struct {
	flow.Component

	Ctx <-chan Context
	Out chan<- ReplyData
}

func (s *OrderPrinter) OnCtx(ctx Context) {
	currentOrder := ctx.Value(CtxKeyOrder)

	if nil == currentOrder {
		s.Out <- ReplyData{"当前没有正在进行中的订单", ctx}
	} else {
		curOrder := currentOrder.(OrderResolve)

		if curOrder.Expired(SesssionExpiredMinutes) {
			s.Out <- ReplyData{"当前没有正在进行中的订单", ctx}
		} else {
			orderDetail := "-----------订单详情-------------\n"
			orderDetail = orderDetail + curOrder.AnswerBody()
			s.Out <- ReplyData{orderDetail, ctx}
		}
	}
}
