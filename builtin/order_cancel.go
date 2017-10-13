package builtin

import (
	// "fmt"
	"time"

	"github.com/wanliu/brain_data/database"
	"github.com/wanliu/flow/builtin/config"
	"github.com/wanliu/flow/builtin/confirm"
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

func (c *OrderCancel) OnCtx(ctx context.Context) {
	currentOrder := ctx.Value(config.CtxKeyOrder)

	if nil == currentOrder {
		preOrderInt := ctx.Value(config.CtxKeyLastOrder)

		if preOrderInt != nil {
			preOrder := preOrderInt.(resolves.OrderResolve)

			eTime := time.Now().Add(-config.PreModifSecs * time.Second)
			if preOrder.UpdatedAt.After(eTime) || preOrder.UpdatedAt.Equal(eTime) {

				if preOrder.Id != 0 {
					order, err := database.GetOrder(preOrder.Id)
					if err == nil {
						orderNo := order.No
						deleteComfirm := confirm.OrderDelete{OrderNo: orderNo}

						deleteComfirm.SetUp(ctx)

						notice := deleteComfirm.Notice(ctx)
						c.Out <- ReplyData{notice, ctx}
						return
					}
				}
			}
		}

		deleteResolve := resolves.OrderDeleteResolve{}
		deleteResolve.SetUp(ctx)

		c.Out <- ReplyData{deleteResolve.Hint(), ctx}
	} else {
		curOrder := currentOrder.(resolves.OrderResolve)

		if curOrder.Cancelable() {
			if curOrder.Cancel() {
				ctx.SetValue(config.CtxKeyOrder, nil)
				c.Out <- ReplyData{"当前订单取消成功", ctx}
			} else {
				c.Out <- ReplyData{"很抱歉，订单取消失败！请联系客服处理", ctx}
			}
		} else {
			c.Out <- ReplyData{"没有可以取消的订单", ctx}
		}
	}
}
