package builtin

import (
	"fmt"
	"strings"

	"github.com/wanliu/flow/builtin/ai"
	"github.com/wanliu/flow/builtin/config"
	"github.com/wanliu/flow/builtin/resolves"
	"github.com/wanliu/flow/context"

	flow "github.com/wanliu/goflow"
)

type OrderItemDelete struct {
	flow.Component

	Ctx <-chan context.Request
	Out chan<- context.Request
}

func NewOrderItemDelete() interface{} {
	return new(OrderItemDelete)
}

func (c *OrderItemDelete) OnCtx(req context.Request) {
	ctx := req.Ctx
	orderRsv := resolves.GetCtxOrder(ctx)

	if nil != orderRsv {
		if orderRsv.Expired(config.SesssionExpiredMinutes) {
			req.Res = context.Response{"会话已经过时，当前没有正在进行中的订单", ctx, nil}
			c.Out <- req
			return
		}

		cmd := req.Command
		if cmd != nil {
			// delete by command
			data := cmd.Data
			if itemName, ok := data["itemName"].(string); ok {
				itemsResolve := orderRsv.Products

				removed := itemsResolve.Remove(itemName)
				if removed {
					orderRsv.Products = itemsResolve

					answer, d := orderRsv.Answer(ctx)

					data := map[string]interface{}{
						"type":   "info",
						"on":     "order",
						"action": "update",
						"data":   d,
					}

					reply := fmt.Sprintf("已经删除%v, %v", itemName, answer)
					req.Res = context.Response{reply, ctx, data}
					c.Out <- req
				} else {
					reply := fmt.Sprintf("无效的操作，%v不存在", itemName)
					req.Res = context.Response{reply, ctx, nil}
					c.Out <- req
				}
			} else {
				req.Res = context.Response{"无效的删除操作", ctx, nil}
				c.Out <- req
			}
		} else {
			// delete by intent
			aiResult := req.ApiAiResult

			aiExtract := ai.ApiAiOrder{AiResult: aiResult}
			deletedItems := []string{}
			products := []string{}
			itemsResolve := orderRsv.Products

			for _, product := range aiExtract.Products() {
				name := product.Product
				products = append(products, name)

				removed := itemsResolve.Remove(name)
				if removed {
					deletedItems = append(deletedItems, name)
				}
			}

			if len(deletedItems) > 0 {
				orderRsv.Products = itemsResolve
				answer, d := orderRsv.Answer(ctx)

				data := map[string]interface{}{
					"type":   "info",
					"on":     "order",
					"action": "update",
					"data":   d,
				}

				reply := fmt.Sprintf("已经删除%v, %v", strings.Join(deletedItems, ","), answer)
				req.Res = context.Response{reply, ctx, data}
				c.Out <- req
			} else {
				_, d := orderRsv.Answer(ctx)

				data := map[string]interface{}{
					"type":   "info",
					"on":     "order",
					"action": "update",
					"data":   d,
				}

				reply := fmt.Sprintf("无效的操作，%v不存在", strings.Join(products, ","))
				req.Res = context.Response{reply, ctx, data}
				c.Out <- req
			}
		}
	} else {
		req.Res = context.Response{"无效的操作，当前没有正在进行中的订单", ctx, nil}
		c.Out <- req
	}
}
