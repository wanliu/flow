package resolves

import (
	"github.com/wanliu/flow/builtin/config"
	"github.com/wanliu/flow/context"
)

func setCtxRsv(ctx context.Context, key string, rsv *OrderResolve) {
	ctx.SetValue(key, rsv)
}

func getCtxRsv(ctx context.Context, key string) *OrderResolve {
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
	setCtxRsv(ctx, config.CtxKeyOrder, rsv)
}

func GetCtxOrder(ctx context.Context) *OrderResolve {
	return getCtxRsv(ctx, config.CtxKeyOrder)
}

func SetCtxLastOrder(ctx context.Context, rsv *OrderResolve) {
	setCtxRsv(ctx, config.CtxKeyLastOrder, rsv)
}

func GetCtxLastOrder(ctx context.Context) *OrderResolve {
	return getCtxRsv(ctx, config.CtxKeyLastOrder)
}

func ClearCtxOrder(ctx context.Context) {
	SetCtxOrder(ctx, nil)
}

func ClearCtxLastOrder(ctx context.Context) {
	SetCtxLastOrder(ctx, nil)
}
