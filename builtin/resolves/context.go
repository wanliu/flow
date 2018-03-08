package resolves

import (
	"github.com/wanliu/flow/builtin/config"
	"github.com/wanliu/flow/context"
)

func SetCtxRsv(ctx context.Context, key string, rsv *OrderResolve) {
	ctx.SetValue(key, rsv)
}

func GetCtxRsv(ctx context.Context, key string) *OrderResolve {
	rsvInt := ctx.Value(key)
	if rsvInt == nil {
		return nil
	}

	switch v := rsvInt.(type) {
	case *OrderResolve:
		return v
	}

	return nil
}

func SetCtxOrder(ctx context.Context, rsv *OrderResolve) {
	SetCtxRsv(ctx, config.CtxKeyOrder, rsv)
}

func GetCtxOrder(ctx context.Context) *OrderResolve {
	return GetCtxRsv(ctx, config.CtxKeyOrder)
}

func SetCtxLastOrder(ctx context.Context, rsv *OrderResolve) {
	SetCtxRsv(ctx, config.CtxKeyLastOrder, rsv)
}

func GetCtxLastOrder(ctx context.Context) *OrderResolve {
	return GetCtxRsv(ctx, config.CtxKeyLastOrder)
}

func ClearCtxOrder(ctx context.Context) {
	SetCtxOrder(ctx, nil)
}

func ClearCtxLastOrder(ctx context.Context) {
	SetCtxLastOrder(ctx, nil)
}
