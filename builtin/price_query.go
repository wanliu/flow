package builtin

import (
	// "log"

	// . "github.com/wanliu/flow/builtin/luis"
	// . "github.com/wanliu/flow/builtin/resolves"
	. "github.com/wanliu/flow/context"
)

type PriceQuery struct {
	TryGetEntities
	Ctx  <-chan Context
	Type <-chan string
	Out  chan<- ReplyData
}

func NewPriceQuery() interface{} {
	return new(PriceQuery)
}

func (query *PriceQuery) OnCtx(ctx Context) {
	// priceQuery := NewPriceQueryResolve(ctx)
	// childCtx := ctx.NewContext()
	// childCtx.SetValue("priceQuery", priceQuery)

	// output := ""

	// if priceQuery.EmptyProducts() {
	// 	output = "没有相关的产品"
	// } else if priceQuery.Fullfilled() {
	// 	output = priceQuery.Answer()
	// } else {
	// 	ctx.Push(childCtx)
	// 	output = priceQuery.Next().Hint()
	// }

	// replyData := ReplyData{output, ctx}
	// query.Out <- replyData

	// go func(task Context) {
	// 	task.Wait(query.TaskHandle)
	// }(childCtx)
}

func (query *PriceQuery) TaskHandle(ctx Context, raw interface{}) error {

	// params := raw.(Context).Value("Result").(ResultParams)

	// priceQuery := ctx.Value("priceQuery").(*PriceQueryResolve)

	// solved, finishNotition, nextNotition := priceQuery.Solve(params)

	// if solved {
	// 	log.Printf("测试输出打印: \n%v", finishNotition)

	// 	reply := ReplyData{finishNotition, ctx}
	// 	query.Out <- reply

	// 	ctx.Pop() // 将当前任务踢出队列
	// } else {
	// 	log.Printf("测试输出打印: \n%v\n", nextNotition)

	// 	reply := ReplyData{nextNotition, ctx}
	// 	query.Out <- reply
	// }
	// // ctx.Send(raw)
	return nil
}
