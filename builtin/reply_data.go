package builtin

import (
	"github.com/wanliu/flow/context"
)

type ReplyData struct {
	Reply string
	Ctx   context.Context
	Data  interface{}
}
