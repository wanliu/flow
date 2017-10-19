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

	output := ""

	if patchResolve.EmptyProducts() {
		output = "没有相关的产品"
	} else {
		curResolve := ctx.Value(config.CtxKeyOrder).(OrderResolve)
		patchResolve.Patch(&curResolve)

		output = curResolve.Answer()

		if curResolve.Resolved() {
			ctx.SetValue(config.CtxKeyOrder, nil)
			ctx.SetValue(config.CtxKeyLastOrder, curResolve)
		} else if curResolve.Failed() {
			ctx.SetValue(config.CtxKeyOrder, nil)
		} else {
			ctx.SetValue(config.CtxKeyOrder, curResolve)
		}

	}

	replyData := ReplyData{output, ctx}
	order.Out <- replyData
}
