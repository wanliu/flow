package builtin

import (
	. "github.com/wanliu/flow/context"
	"log"
)

type Abuse struct {
	TryGetEntities
	Ctx  <-chan Context
	Type <-chan string
	Out  chan<- ReplyData
}

func NewAbuse() interface{} {
	return new(Abuse)
}

// entity: 贬低
func (order *Abuse) OnCtx(ctx Context) {
	// entities := ctx.Value("Result").(ResultParams).Entities
	output := "请不要脏话哦"

	log.Printf("OUTPUT: %v", output)

	replyData := ReplyData{output, ctx}
	order.Out <- replyData
}
