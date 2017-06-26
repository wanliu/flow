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

func (order *Order) OnCtx(ctx Context) {
	// if _, ok := ctx.Value("Result").(ResultParams); ok {
	// orderResolve := NewOpenOrderResolve(ctx)
	orderResolve := NewOpenOrderResolve(ctx)
	childCtx := ctx.NewContext()
	childCtx.SetValue("orderResolve", orderResolve)
	ctx.Push(childCtx)
	// orderResolve.AddressFullfilled()
	output := orderResolve.Next().Hint()

	log.Printf("OUTPUT: %v", output)

	replyData := ReplyData{output, ctx}
	order.Out <- replyData

	go func(task Context) {
		task.Wait(order.TaskHandle)
		// ctx.Pop()
	}(childCtx)
}

func (order *Order) TaskHandle(ctx Context, raw interface{}) error {
	params := raw.(Context).Value("Result").(ResultParams)

	orderResolve := ctx.Value("orderResolve").(*OpenOrderResolve)
	orderResolve.Solve(params)
	// ctx.Send(raw)
	return nil
}
