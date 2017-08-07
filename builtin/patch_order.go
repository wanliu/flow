package builtin

import (
	// "log"

	. "github.com/wanliu/flow/builtin/config"
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
		curResolve := ctx.Value(CtxKeyOrder).(OrderResolve)
		patchResolve.Patch(&curResolve)

		if curResolve.Fulfiled() {
			ctx.SetValue(CtxKeyOrder, nil)
			ctx.SetValue(CtxKeyLastOrder, curResolve)
		} else {
			ctx.SetValue(CtxKeyOrder, curResolve)
		}

		output = patchResolve.Answer()
	}

	replyData := ReplyData{output, ctx}
	order.Out <- replyData
}
