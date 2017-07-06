package builtin

import (
	// "fmt"
	// "strings"
	"log"

	. "github.com/wanliu/flow/builtin/luis"
	. "github.com/wanliu/flow/builtin/resolves"
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
	orderResolve := NewOpenOrderResolve(ctx)
	childCtx := ctx.NewContext()
	childCtx.SetValue("orderResolve", orderResolve)

	output := ""

	if orderResolve.EmptyProducts() {
		output = "没有相关的产品"
	} else {
		ctx.Push(childCtx)
		output = orderResolve.Next().Hint()
	}

	log.Printf("OUTPUT: %v", output)

	replyData := ReplyData{output, ctx}
	order.Out <- replyData

	go func(task Context) {
		task.Wait(order.TaskHandle)
	}(childCtx)
}

func (order *Order) TaskHandle(ctx Context, raw interface{}) error {

	params := raw.(Context).Value("Result").(ResultParams)

	orderResolve := ctx.Value("orderResolve").(*OpenOrderResolve)

	solved, finishNotition, nextNotition := orderResolve.Solve(params)

	if solved {
		log.Printf("测试输出打印: \n%v", finishNotition)

		reply := ReplyData{finishNotition, ctx}
		order.Out <- reply

		ctx.Pop() // 将当前任务踢出队列
	} else {
		log.Printf("测试输出打印: \n%v\n", nextNotition)

		reply := ReplyData{nextNotition, ctx}
		order.Out <- reply
	}
	// ctx.Send(raw)
	return nil
}
