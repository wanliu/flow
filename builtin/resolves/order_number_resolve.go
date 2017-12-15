package resolves

import (
	"github.com/wanliu/flow/context"
)

type OrderNumberResolve interface {
	Resolve(orderNo string, ctx context.Context) string
}
