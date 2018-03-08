package builtin

import (
	// "github.com/wanliu/flow/builtin/config"
	"github.com/wanliu/flow/builtin/resolves"
	"github.com/wanliu/flow/context"
)

type PatchOrder struct {
	TryGetEntities
	DefTime string
	Ctx     <-chan context.Request
	Out     chan<- context.Request
}

func NewPatchOrder() interface{} {
	return new(PatchOrder)
}

func (order *PatchOrder) OnCtx(req context.Request) {
	ctx := req.Ctx
	patchResolve := resolves.NewPatchOrderResolve(req)

	var data interface{}
	output := ""

	if patchResolve.EmptyProducts() && patchResolve.EmptyGifts() {
		output = "没有相关的产品"
	} else {
		orderRsv := resolves.GetCtxOrder(ctx)
		patchResolve.Patch(orderRsv)

		var d interface{}

		output, d = orderRsv.Answer(ctx)
		data = map[string]interface{}{
			"type":   "info",
			"on":     "order",
			"action": "update",
			"data":   d,
		}

		if orderRsv.Resolved() {
			resolves.ClearCtxOrder(ctx)
			resolves.SetCtxLastOrder(ctx, orderRsv)
		} else if orderRsv.Failed() {
			resolves.ClearCtxOrder(ctx)
		}

	}

	req.Res = context.Response{
		Reply: output,
		Ctx:   ctx,
		Data:  data,
	}
	order.Out <- req
}
