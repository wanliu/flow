package builtin

import (
	"log"
	"strconv"

	"github.com/wanliu/flow/builtin/config"
	"github.com/wanliu/flow/builtin/resolves"

	"github.com/wanliu/flow/context"

	flow "github.com/wanliu/goflow"
)

type OperationNotice struct {
	flow.Component
	Ctx <-chan context.Request
	Out chan<- context.Request
}

func NewOperationNotice() interface{} {
	return new(OperationNotice)
}

func (s OperationNotice) OnCtx(req context.Request) {
	ctx := req.Ctx

	if context.GroupChat(ctx) {
		log.Printf("不回应非开单相关的普通群聊")
		return
	}

	// currentOrder := ctx.CtxValue(CtxKeyOrder)
	orderRsv := resolves.GetCtxOrder(ctx)

	if nil != orderRsv {
		if !orderRsv.Fulfiled() {
			expMins := config.SesssionExpiredMinutes

			if nil != ctx.CtxValue(config.CtxKeyExpiredMinutes) {
				expMins = ctx.CtxValue(config.CtxKeyExpiredMinutes).(int)
			}

			req.Res = context.Response{"你可以继续提交产品到订单，也可以立刻取消当前任务（" + strconv.Itoa(expMins) + "分钟以内）", ctx, nil}
			s.Out <- req
		}
	}
}
