package builtin

import (
	// "log"
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

	Ctx <-chan context.Request
	Out chan<- context.Request
}

func (c *OrderCancel) OnCtx(req context.Request) {
	// if context.GroupChat(ctx) {
	// 	log.Printf("不回应非开单相关的普通群聊")
	// 	return
	// }

	ctx := req.Ctx
	orderRsv := resolves.GetCtxOrder(ctx)

	if nil == orderRsv {
		lastOrderRsv := resolves.GetCtxLastOrder(ctx)

		if lastOrderRsv != nil {
			eTime := time.Now().Add(-config.PreModifSecs * time.Second)
			if lastOrderRsv.UpdatedAt.After(eTime) || lastOrderRsv.UpdatedAt.Equal(eTime) {

				if lastOrderRsv.Id != 0 {
					order, err := database.GetOrder(lastOrderRsv.Id)
					if err == nil {
						orderNo := order.No
						deleteComfirm := resolves.OrderDeleteConfirm{OrderNo: orderNo}

						deleteComfirm.SetUp(ctx)

						notice := deleteComfirm.Notice(ctx)
						req.Res = context.Response{notice, ctx, nil}
						c.Out <- req
						return
					}
				}
			}
		}

		// 暂时不在群聊中提供根据订单号删除订单的功能,只能删除最近订单
		if context.GroupChat(ctx) {
			// log.Printf("不回应非开单相关的普通群聊")
			req.Res = context.Response{"当前没有可以取消的订单", ctx, nil}
			c.Out <- req
			return
		}

		deleteResolve := resolves.OrderDeleteResolve{}
		deleteResolve.SetUp(ctx)

		req.Res = context.Response{deleteResolve.Hint(), ctx, nil}
		c.Out <- req
	} else if nil != orderRsv {
		if orderRsv.Cancelable() {
			if orderRsv.Cancel() {
				ctx.SetCtxValue(config.CtxKeyOrder, nil)
				req.Res = context.Response{"当前订单取消成功", ctx, nil}
				c.Out <- req
			} else {
				req.Res = context.Response{"很抱歉，订单取消失败！请联系客服处理", ctx, nil}
				c.Out <- req
			}
		} else {
			req.Res = context.Response{"没有可以取消的订单", ctx, nil}
			c.Out <- req
		}
	}
}
