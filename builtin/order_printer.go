package builtin

import (
	"log"

	"github.com/wanliu/flow/builtin/resolves"
	"github.com/wanliu/flow/context"

	config "github.com/wanliu/flow/builtin/config"
	flow "github.com/wanliu/goflow"
)

func NewOrderPrinter() interface{} {
	return new(OrderPrinter)
}

type OrderPrinter struct {
	flow.Component

	Ctx <-chan context.Request
	Out chan<- context.Request
}

func (s *OrderPrinter) OnCtx(req context.Request) {
	ctx := req.Ctx

	if context.GroupChat(ctx) {
		log.Printf("不回应非开单相关的普通群聊")
		return
	}

	orderRsv := resolves.GetCtxOrder(ctx)

	if nil == orderRsv {
		req.Res = context.Response{"当前没有正在进行中的订单", ctx, nil}
		s.Out <- req
	} else {
		if orderRsv.Expired(config.SesssionExpiredMinutes) {
			req.Res = context.Response{"当前没有正在进行中的订单", ctx, nil}
			s.Out <- req
		} else {
			// orderDetail := "-----------订单详情-------------\n"
			// orderDetail = orderDetail + orderRsv.AnswerBody()
			req.Res = context.Response{"订单详情", ctx, orderRsv.ToDescStruct()}
			s.Out <- req
		}
	}
}
