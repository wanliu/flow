package builtin

import (
	"time"

	"github.com/wanliu/brain_data/database"
	"github.com/wanliu/flow/builtin/config"
	"github.com/wanliu/flow/builtin/resolves"

	"github.com/wanliu/flow/context"

	flow "github.com/wanliu/goflow"
)

type ChangeItemUnit struct {
	Ctx <-chan context.Request
	Out chan<- context.Request
}

func NewOChangeItemUnit() interface{} {
	return new(ChangeItemUnit)
}

func OnCtx(req context.Request) {
	ctx := req.Ctx
	orderRsv := resolves.GetCtxOrder(ctx)

	if nil == orderRsv {
		req.Res = context.Response{"操作失败，当前没有正在进行中的订单", ctx, nil}
		c.Out <- req
		return
	} else {
		if orderRsv.Expired(config.SesssionExpiredMinutes) {
			req.Res = context.Response{"操作失败，当前没有正在进行中的订单", ctx, nil}
			c.Out <- req
			return
		}

		cmd := req.Command
		if cmd != nil {
			// delete by command
			data := cmd.Data
			if itemName, ok := data["itemName"].(string); ok {
				if unit, ok := data["unit"].(string); ok {
					itemsResolve := orderRsv.Products
					err := itemsResolve.ChangeUnite(itemName, unit)

					if err != nil {
						req.Res = context.Response{fmt.Sprintf("无效的操作, %v", err.Error()), ctx, nil}
						c.Out <- req
						return
					} else {
						answer, d := orderRsv.Answer(ctx)

						data := map[string]interface{}{
							"type":   "info",
							"on":     "order",
							"action": "update",
							"data":   d,
						}

						reply := fmt.Sprintf("已经更新%v单位为%v", itemName, unit)
						req.Res = context.Response{reply, ctx, data}
						c.Out <- req
					}
				} else {
					req.Res = context.Response{"无效的操作, 请提供单位名称", ctx, nil}
					c.Out <- req
				}
			} else {
				req.Res = context.Response{"无效的删除操作, 请提供商品名称", ctx, nil}
				c.Out <- req
			}
		} else {
			// TODO apiai intent to change item unit
		}
	}
}
