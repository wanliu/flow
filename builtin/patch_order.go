package builtin

import (
	// "log"

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
		curResolve := ctx.Value("Order").(OrderResolve)
		patchResolve.Patch(&curResolve)
		ctx.SetValue("Order", curResolve)

		output = patchResolve.Answer()
	}

	replyData := ReplyData{output, ctx}
	order.Out <- replyData
}
