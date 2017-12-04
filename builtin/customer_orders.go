package builtin

import (
	"fmt"
	"time"

	"github.com/hysios/apiai-go"
	"github.com/wanliu/brain_data/resolves"
	// "github.com/wanliu/flow/builtin/ai"
	"github.com/wanliu/flow/context"
)

const (
	Per = 5
)

type CustomerOrders struct {
	TryGetEntities
	Ctx  <-chan context.Context
	Type <-chan string
	Out  chan<- ReplyData
}

func NewCustomerOrders() interface{} {
	return new(CustomerOrders)
}

func (c *CustomerOrders) OnCtx(ctx context.Context) {
	rsv := resolves.NewCusOrdersResolve(ctx)

	rsv.Setup(ctx)

	reply := rsv.Answer()

	c.Out <- ReplyData{reply, ctx}
}
