package builtin

import (
	// "log"

	// . "github.com/wanliu/flow/builtin/resolves"
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

// 默认送货时间
func (order *PatchOrder) OnDeftime(t string) {
	order.DefTime = t
}

func (order *PatchOrder) OnCtx(ctx Context) {
	// orderResolve := NewPatchOrderResolve(ctx)

	// output := ""

	// if orderResolve.EmptyProducts() {
	// 	output = "没有相关的产品"
	// } else {
	// 	output = orderResolve.Answer()
	// }

	// log.Printf("OUTPUT: %v", output)

	// replyData := ReplyData{output, ctx}
	// order.Out <- replyData
}
