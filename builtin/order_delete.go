package builtin

import (
	"fmt"
	// "log"
	// "time"

	// "github.com/wanliu/flow/builtin/config"
	// "github.com/hysios/apiai-go"
	"github.com/wanliu/brain_data/database"
	"github.com/wanliu/flow/builtin/resolves"

	"github.com/wanliu/flow/context"

	flow "github.com/wanliu/goflow"
)

func NewOrderDelete() interface{} {
	return new(OrderDelete)
}

type OrderDelete struct {
	flow.Component

	Ctx <-chan context.Request
	Out chan<- ReplyData
}

func (c *OrderDelete) OnCtx(req context.Request) {
	// if context.GroupChat(ctx) {
	// 	log.Printf("不回应非开单相关的普通群聊")
	// 	return
	// }

	ctx := req.Ctx
	aiResult := req.ApiAiResult

	if numInt, exist := aiResult.Params["order-numder"]; exist {
		orderNo := numInt.(string)
		if orderNo == "" {
			c.setupResolve(ctx)
		} else {

			order, err := database.GetOrderByNo(orderNo)

			if err != nil {
				reply := fmt.Sprintf("找不到订单号为 %v 的订单", orderNo)
				c.Out <- ReplyData{reply, ctx, nil}
			} else {
				reply := ""

				err = order.Delete()
				if err == nil {
					reply = fmt.Sprintf("%v 号订单删除成功", orderNo)
				} else {
					reply = fmt.Sprintf("%v 号订单删除失败，请访问 https://jiejie.io/orders/%v 进行操作", orderNo, order.GlobelId())
				}

				c.Out <- ReplyData{reply, ctx, nil}
			}
		}
	} else {
		c.setupResolve(ctx)
	}
}

func (c *OrderDelete) setupResolve(ctx context.Context) {
	deleteResolve := resolves.OrderDeleteResolve{}
	deleteResolve.SetUp(ctx)

	c.Out <- ReplyData{deleteResolve.Hint(), ctx, nil}
}
