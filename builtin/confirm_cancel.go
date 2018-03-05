package builtin

import (
	// "log"

	"github.com/wanliu/flow/builtin/resolves"

	"github.com/wanliu/flow/context"

	config "github.com/wanliu/flow/builtin/config"
)

type Confirm struct {
	TryGetEntities

	expMins float64

	ExpMins <-chan float64

	Confirm <-chan context.Request
	Cancel  <-chan context.Request
	Expire  <-chan context.Request

	Out chan<- ReplyData
}

func NewConfirm() interface{} {
	return new(Confirm)
}

func (c *Confirm) OnExpMins(mins float64) {
	c.expMins = mins
}

func (c Confirm) OnExpire(req context.Request) {
	ctx := req.Ctx
	ctx.SetCtxValue(config.CtxKeyConfirm, nil)
}

func (c *Confirm) OnConfirm(req context.Request) {
	ctx := req.Ctx
	cIn := ctx.CtxValue(config.CtxKeyConfirm)

	if cIn != nil {
		cfm := cIn.(resolves.Data)
		reply, data := cfm.Confirm(ctx)
		c.Out <- ReplyData{reply, ctx, data}
	} else {
		// 群聊无待确认项目时，不回应
		// if context.GroupChat(ctx) {
		// 	log.Printf("不回应非开单相关的普通群聊")
		// 	return
		// }

		reply := ReplyData{"确认操作已经过期", ctx, nil}
		c.Out <- reply
	}
}

func (c *Confirm) OnCancel(req context.Request) {
	ctx := req.Ctx
	cIn := ctx.CtxValue(config.CtxKeyConfirm)

	if cIn != nil {
		cfm := cIn.(resolves.Data)
		reply := cfm.Cancel(ctx)
		c.Out <- ReplyData{reply, ctx, nil}
	} else {
		// 群聊无待确认项目时，不回应
		// if context.GroupChat(ctx) {
		// 	log.Printf("不回应非开单相关的普通群聊")
		// 	return
		// }

		reply := ReplyData{"确认操作已经过期", ctx, nil}
		c.Out <- reply
	}
}
