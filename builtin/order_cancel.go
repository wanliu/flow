package builtin

import (
	"fmt"
	"time"

	"github.com/wanliu/brain_data/database"
	"github.com/wanliu/flow/builtin/config"
	"github.com/wanliu/flow/builtin/resolves"
	"github.com/wanliu/flow/context"

	flow "github.com/wanliu/goflow"
)

func NewOrderCancel() interface{} {
	return new(OrderCancel)
}

type OrderCancel struct {
	flow.Component

	Ctx <-chan context.Context
	Out chan<- ReplyData
}

func (s *OrderCancel) OnCtx(ctx context.Context) {
	currentOrder := ctx.Value(config.CtxKeyOrder)

	if nil == currentOrder {
		preOrderInt := ctx.Value(config.CtxKeyLastOrder)

		if preOrderInt != nil {
			preOrder := preOrderInt.(resolves.OrderResolve)

			eTime := time.Now().Add(-config.PreModifSecs * time.Second)
			if preOrder.UpdatedAt.After(eTime) || preOrder.UpdatedAt.Equal(eTime) {
				// TODO delete the last order
				if preOrder.Id != 0 {
					order, err := database.GetOrder(preOrder.Id)
					if err == nil {
						orderNo := order.No
						err = order.Delete()
						if err == nil {
							ctx.SetValue(config.CtxKeyLastOrder, nil)

							msg := fmt.Sprintf("成功删除订单号为：%v 的订单", orderNo)
							s.Out <- ReplyData{msg, ctx}
							return
						}
					}

					s.Out <- ReplyData{fmt.Sprintf("订单取消失败: %v", err.Error()), ctx}
					return
				}
			}
		}

		s.Out <- ReplyData{"没有可以取消的订单", ctx}
	} else {
		curOrder := currentOrder.(resolves.OrderResolve)

		if curOrder.Cancelable() {
			if curOrder.Cancel() {
				ctx.SetValue(config.CtxKeyOrder, nil)
				s.Out <- ReplyData{"当前订单取消成功", ctx}
			} else {
				s.Out <- ReplyData{"很抱歉，订单取消失败！请联系客服处理", ctx}
			}
		} else {
			s.Out <- ReplyData{"没有可以取消的订单", ctx}
		}
	}
}
