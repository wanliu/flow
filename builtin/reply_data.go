package builtin

import (
	"github.com/wanliu/flow/context"
)

type ReplyData struct {
	Reply string
	Table *context.Table
	Ctx   context.Context
}
