package builtin

import (
	// "fmt"
	// "strings"
	"log"

	. "github.com/wanliu/flow/context"
)

type Order struct {
	TryGetEntities
	Ctx  <-chan Context
	Type <-chan string
	Out  chan<- ReplyData
}

func NewOrder() interface{} {
	return new(Order)
}

func (tr *Order) OnCtx(ctx Context) {
	// if _, ok := ctx.Value("Result").(ResultParams); ok {
	// orderResolve := NewOpenOrderResolve(ctx)
	orderResolve := NewOpenOrderResolve(ctx)
	orderResolve.AddressFullfilled()
	output := orderResolve.Next().Hint()

	log.Printf("OUTPUT: %v", output)

	replyData := ReplyData{output, ctx}
	tr.Out <- replyData
	// 	tr.Out <- replyData
	// } else {
	// 	// tr.No <- ctx
	// 	replyData := ReplyData{"出现错误，请稍后重试", ctx}
	// 	tr.Out <- replyData
	// }
}
