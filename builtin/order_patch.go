package builtin

import (
	"github.com/wanliu/flow/builtin/config"
	. "github.com/wanliu/flow/builtin/resolves"
	. "github.com/wanliu/flow/context"
)

type PatchOrder struct {
	TryGetEntities
	DefTime string
	Ctx     <-chan Context
	Out     chan<- ReplyData
}

func NewPatchOrder() interface{} {
	return new(PatchOrder)
}

func (order *PatchOrder) OnCtx(ctx Context) {
	patchResolve := NewPatchOrderResolve(ctx)

	var data interface{}
	output := ""

	if patchResolve.EmptyProducts() {
		output = "没有相关的产品"
	} else {
		curResolve := ctx.CtxValue(config.CtxKeyOrder).(OrderResolve)
		patchResolve.Patch(&curResolve)

		var d interface{}

		output, d = curResolve.Answer(ctx)
		data = map[string]interface{}{
			"type":   "info",
			"on":     "order",
			"action": "update",
			"data":   d,
		}

		if curResolve.Resolved() {
			ctx.SetCtxValue(config.CtxKeyOrder, nil)
			ctx.SetCtxValue(config.CtxKeyLastOrder, curResolve)
		} else if curResolve.Failed() {
			ctx.SetCtxValue(config.CtxKeyOrder, nil)
		} else {
			ctx.SetCtxValue(config.CtxKeyOrder, curResolve)
		}

	}

	replyData := ReplyData{
		Reply: output,
		Ctx:   ctx,
		Data:  data,
	}
	order.Out <- replyData
}
