package builtin

import (
	"log"

	. "github.com/wanliu/flow/builtin/luis"
	. "github.com/wanliu/flow/builtin/resolves"
	. "github.com/wanliu/flow/context"
)

type StockQuery struct {
	TryGetEntities
	Ctx  <-chan Context
	Type <-chan string
	Out  chan<- ReplyData
}

func NewStockQuery() interface{} {
	return new(StockQuery)
}

func (query *StockQuery) OnCtx(ctx Context) {
	stockQuery := NewStockQueryResolve(ctx)
	childCtx := ctx.NewContext()
	childCtx.SetValue("stockQuery", stockQuery)

	output := ""

	if stockQuery.EmptyProducts() {
		output = "没有相关的产品"
	} else if stockQuery.Fullfilled() {
		output = stockQuery.Answer()
	} else {
		ctx.Push(childCtx)
		output = stockQuery.Next().Hint()
	}

	replyData := ReplyData{output, ctx}
	query.Out <- replyData

	go func(task Context) {
		task.Wait(query.TaskHandle)
	}(childCtx)
}

func (query *StockQuery) TaskHandle(ctx Context, raw interface{}) error {

	params := raw.(Context).Value("Result").(ResultParams)

	stockQuery := ctx.Value("stockQuery").(*StockQueryResolve)

	solved, finishNotition, nextNotition := stockQuery.Solve(params)

	if solved {
		log.Printf("测试输出打印: \n%v", finishNotition)

		reply := ReplyData{finishNotition, ctx}
		query.Out <- reply

		ctx.Pop() // 将当前任务踢出队列
	} else {
		log.Printf("测试输出打印: \n%v\n", nextNotition)

		reply := ReplyData{nextNotition, ctx}
		query.Out <- reply
	}
	// ctx.Send(raw)
	return nil
}
