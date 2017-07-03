package builtin

import (
	. "github.com/wanliu/flow/context"
	"log"
)

type Unimplemented struct {
	TryGetEntities
	Ctx  <-chan Context
	Type <-chan string
	Out  chan<- ReplyData
}

func NewUnimplemented() interface{} {
	return new(Unimplemented)
}

func (order *Unimplemented) OnCtx(ctx Context) {
	output := "对不起，请问有什么可以帮您？"

	log.Printf("OUTPUT: %v", output)

	replyData := ReplyData{output, ctx}
	order.Out <- replyData
}
