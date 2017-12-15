package builtin

import (
	. "github.com/wanliu/flow/builtin/config"
	. "github.com/wanliu/flow/builtin/resolves"
	. "github.com/wanliu/flow/context"
	flow "github.com/wanliu/goflow"
)

func NewOrderCancel() interface{} {
	return new(OrderCancel)
}

type OrderCancel struct {
	flow.Component

	Ctx <-chan Context
	Out chan<- ReplyData
}

func (s *OrderCancel) OnCtx(ctx Context) {
	currentOrder := ctx.Value(CtxKeyOrder)

	if nil == currentOrder {
		s.Out <- ReplyData{"没有可以取消的订单", ctx}
	} else {
		curOrder := currentOrder.(OrderResolve)

		if curOrder.Cancelable() {
			if curOrder.Cancel() {
				ctx.SetValue(CtxKeyOrder, nil)
				s.Out <- ReplyData{"当前订单取消成功", ctx}
			} else {
				s.Out <- ReplyData{"很抱歉，订单取消失败！请联系客服处理", ctx}
			}
		} else {
			s.Out <- ReplyData{"没有可以取消的订单", ctx}
		}
	}
}
