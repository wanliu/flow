package confirm

import (
	"github.com/wanliu/flow/context"
)

// type ConfirmData {
// 	Resolve *Resolve
// 	ResolveType string

// 	Action string
// 	Value interface{}
// }

type Data interface {
	Notice(ctx context.Context) string
	Cancel(ctx context.Context) string
	Confirm(ctx context.Context) string
	SetUp(ctx context.Context)
	ClearUp(ctx context.Context)
}
