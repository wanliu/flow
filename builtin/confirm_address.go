package builtin

import (
	// "log"

	. "github.com/wanliu/flow/builtin/resolves"
	. "github.com/wanliu/flow/context"

	config "github.com/wanliu/flow/builtin/config"
)

// type ConfirmData {
// 	Resolve *Resolve
// 	ResolveType string

// 	Action string
// 	Value interface{}
// }

type AddressConfirm struct {
	Order *OrderResolve
	Value string
}

type Confirm struct {
	TryGetEntities

	expMins float64

	ExpMins <-chan float64

	Confirm <-chan Context
	Cancel  <-chan Context
	Expire  <-chan Context

	Out chan<- ReplyData
}

func NewConfirm() interface{} {
	return new(Confirm)
}

func (c *Confirm) OnExpMins(mins float64) {
	c.expMins = mins
}

func (c Confirm) OnExpire(ctx Context) {
	ctx.SetValue(config.CtxKeyConfirm, nil)
}

func (c *Confirm) OnConfirm(ctx Context) {
	currentOrder := ctx.Value(config.CtxKeyOrder)
	confirm := ctx.Value(config.CtxKeyConfirm)

	if currentOrder != nil && confirm != nil {
		cOrder := currentOrder.(OrderResolve)

		if cOrder.Expired(config.SesssionExpiredMinutes) {
			c.Out <- ReplyData{"当前没有正在进行中的订单", ctx}
			return
		}

		cConfirm := confirm.(AddressConfirm)

		if cOrder.Address == "" {
			cOrder.Address = cConfirm.Value

			if cOrder.Fulfiled() {
				ctx.SetValue(config.CtxKeyOrder, nil)
				ctx.SetValue(config.CtxKeyLastOrder, cOrder)
			} else {
				ctx.SetValue(config.CtxKeyOrder, cOrder)
			}

			reply := "已经确认\"" + cConfirm.Value + "\"为收货地址" + "\n" + cOrder.Answer()
			c.Out <- ReplyData{reply, ctx}
		} else {
			reply := ReplyData{"确认操作已经过期", ctx}
			c.Out <- reply
		}
	} else {
		reply := ReplyData{"确认操作已经过期", ctx}
		c.Out <- reply
	}
}

func (c *Confirm) OnCancel(ctx Context) {
	currentOrder := ctx.Value(config.CtxKeyOrder)
	confirm := ctx.Value(config.CtxKeyConfirm)

	if currentOrder != nil && confirm != nil {
		cOrder := currentOrder.(OrderResolve)

		ctx.SetValue(config.CtxKeyConfirm, nil)

		reply := "取消操作完成"

		if cOrder.Address == "" {
			reply = reply + "，当前订单收货地址尚未确认，请输入收货地址"
		}

		c.Out <- ReplyData{reply, ctx}
	} else {
		reply := ReplyData{"取消操作已经过期", ctx}
		c.Out <- reply
	}
}
