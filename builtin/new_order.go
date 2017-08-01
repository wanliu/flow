package builtin

import (
	"log"

	. "github.com/wanliu/flow/builtin/resolves"
	. "github.com/wanliu/flow/context"
)

type NewOrder struct {
	TryGetEntities
	DefTime string
	Ctx     <-chan Context
	Deftime <-chan string
	Out     chan<- ReplyData
}

func NewNewOrder() interface{} {
	return new(NewOrder)
}

// 默认送货时间
func (c *NewOrder) OnDeftime(t string) {
	c.DefTime = t
}

func (c *NewOrder) OnCtx(ctx Context) {
	orderResolve := NewNewOrderResolve(ctx)

	if c.DefTime != "" {
		orderResolve.SetDefTime(c.DefTime)
	}

	output := ""

	if orderResolve.EmptyProducts() {
		output = "没有相关的产品"
	} else {
		output = orderResolve.Answer()
	}

	log.Printf("OUTPUT: %v", output)

	replyData := ReplyData{output, ctx}
	c.Out <- replyData
}
