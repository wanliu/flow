package confirm

import (
	"fmt"

	"github.com/wanliu/flow/builtin/config"
	"github.com/wanliu/flow/builtin/resolves"
	"github.com/wanliu/flow/context"
)

type AddressConfirm struct {
	Values []string
	// order  *resolves.OrderResolve
}

func (ac *AddressConfirm) SetUp(ctx context.Context) {
	ctx.SetValue(config.CtxKeyConfirm, *ac)
}

func (ac *AddressConfirm) ClearUp(ctx context.Context) {
	ctx.SetValue(config.CtxKeyConfirm, nil)
}

func (ac *AddressConfirm) Notice(ctx context.Context) string {
	oIn := ctx.Value(config.CtxKeyOrder)
	// confirm := ctx.Value(config.CtxKeyConfirm)

	if oIn != nil {
		order := oIn.(resolves.OrderResolve)

		if order.Expired(config.SesssionExpiredMinutes) {
			return "当前没有正在进行中的订单"
		}

		if len(ac.Values) > 0 {
			return fmt.Sprintf("是否将 “%v” 做为收货地址?", ac.Values[0])
		} else {
			// ctx.SetValue(config.CtxKeyConfirm, nil)
			ac.ClearUp(ctx)

			if order.Address == "" {
				return "取消操作完成，当前订单收货地址尚未确认，请输入收货地址"
			}
		}
	} else {
		return "当前没有正在进行中的订单"
	}

	return ""
}

func (ac *AddressConfirm) Cancel(ctx context.Context) string {
	oIn := ctx.Value(config.CtxKeyOrder)
	// confirm := ctx.Value(config.CtxKeyConfirm)

	if oIn != nil {
		order := oIn.(resolves.OrderResolve)

		if order.Expired(config.SesssionExpiredMinutes) {
			return "当前没有正在进行中的订单"
		}

		if len(ac.Values) > 1 {
			newValues := ac.Values[1:]
			ac.Values = newValues
			// ctx.SetValue(config.CtxKeyConfirm, ac)
			ac.SetUp(ctx)

			return fmt.Sprintf("是否将 “%v” 做为收货地址?", newValues[0])
		} else {
			// ctx.SetValue(config.CtxKeyConfirm, nil)
			ac.ClearUp(ctx)

			if order.Address == "" {
				return "取消操作完成，当前订单收货地址尚未确认，请输入收货地址"
			}
		}
	} else {
		return "当前没有正在进行中的订单"
	}

	return ""
}

func (ac *AddressConfirm) Comfirm(ctx context.Context) string {
	oIn := ctx.Value(config.CtxKeyOrder)
	// confirm := ctx.Value(config.CtxKeyConfirm)

	if oIn != nil {
		order := oIn.(resolves.OrderResolve)

		if order.Expired(config.SesssionExpiredMinutes) {
			return "当前没有正在进行中的订单"
		}

		// cConfirm := confirm.(AddressConfirm)

		if order.Address == "" {
			if len(ac.Values) > 0 {
				order.Address = ac.Values[0]

				if order.Fulfiled() {
					ctx.SetValue(config.CtxKeyOrder, nil)
					ctx.SetValue(config.CtxKeyLastOrder, order)
				} else {
					ctx.SetValue(config.CtxKeyOrder, order)
				}

				return fmt.Sprintf("已经确认\"%v\"为收货地址\n%v", ac.Values[0], order.Answer())

			} else {
				// ctx.SetValue(config.CtxKeyConfirm, nil)
				ac.ClearUp(ctx)
			}
		}
	} else {
		return "当前没有正在进行中的订单"
	}

	return "确认操作已经过期"
}
