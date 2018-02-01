package resolves

import (
	"fmt"

	"github.com/wanliu/flow/builtin/config"
	"github.com/wanliu/flow/context"
)

type AddressConfirm struct {
	Values []string
	// order  *OrderResolve
}

func (ac AddressConfirm) SetUp(ctx context.Context) {
	ctx.SetCtxValue(config.CtxKeyConfirm, ac)
}

func (ac AddressConfirm) ClearUp(ctx context.Context) {
	ctx.SetCtxValue(config.CtxKeyConfirm, nil)
}

func (ac AddressConfirm) Notice(ctx context.Context) string {
	oIn := ctx.CtxValue(config.CtxKeyOrder)
	// confirm := ctx.CtxValue(config.CtxKeyConfirm)

	if oIn != nil {
		order := oIn.(OrderResolve)

		if order.Expired(config.SesssionExpiredMinutes) {
			return "当前没有正在进行中的订单"
		}

		if len(ac.Values) > 0 {
			return fmt.Sprintf("是否将 “%v” 做为收货客户?", ac.Values[0])
		} else {
			// ctx.SetCtxValue(config.CtxKeyConfirm, nil)
			ac.ClearUp(ctx)

			if order.Customer == "" {
				return "取消操作完成，当前订单收货客户尚未确认，请输入收货客户"
			}
		}
	} else {
		return "当前没有正在进行中的订单"
	}

	return ""
}

func (ac AddressConfirm) Cancel(ctx context.Context) string {
	oIn := ctx.CtxValue(config.CtxKeyOrder)
	// confirm := ctx.CtxValue(config.CtxKeyConfirm)

	if oIn != nil {
		order := oIn.(OrderResolve)

		if order.Expired(config.SesssionExpiredMinutes) {
			return "当前没有正在进行中的订单"
		}

		if len(ac.Values) > 1 {
			newValues := ac.Values[1:]
			ac.Values = newValues
			// ctx.SetCtxValue(config.CtxKeyConfirm, ac)
			ac.SetUp(ctx)

			return fmt.Sprintf("是否将 “%v” 做为收货客户?", newValues[0])
		} else {
			// ctx.SetCtxValue(config.CtxKeyConfirm, nil)
			ac.ClearUp(ctx)

			if order.Customer == "" {
				return "取消操作完成，当前订单收货客户尚未确认，请输入收货客户"
			}
		}
	} else {
		return "当前没有正在进行中的订单"
	}

	return ""
}

func (ac AddressConfirm) Confirm(ctx context.Context) (string, interface{}) {
	oIn := ctx.CtxValue(config.CtxKeyOrder)
	// confirm := ctx.CtxValue(config.CtxKeyConfirm)

	if oIn != nil {
		order := oIn.(OrderResolve)

		if order.Expired(config.SesssionExpiredMinutes) {
			return "当前没有正在进行中的订单", nil
		}

		// cConfirm := confirm.(AddressConfirm)

		if order.Customer == "" {
			if len(ac.Values) > 0 {
				order.ExtractedCustomer = ac.Values[0]
				order.CheckExtractedCustomer()

				reply, data := order.Answer(ctx)
				reply = fmt.Sprintf("已经确认\"%v\"为收货客户\n%v", ac.Values[0], reply)

				if order.Resolved() {
					ctx.SetCtxValue(config.CtxKeyOrder, nil)
					ctx.SetCtxValue(config.CtxKeyLastOrder, order)
				} else if order.Failed() {
					ctx.SetCtxValue(config.CtxKeyOrder, nil)
				} else {
					ctx.SetCtxValue(config.CtxKeyOrder, order)
				}

				return reply, data

			} else {
				// ctx.SetCtxValue(config.CtxKeyConfirm, nil)
				ac.ClearUp(ctx)
			}
		}
	} else {
		return "当前没有正在进行中的订单", nil
	}

	return "确认操作已经过期", nil
}
