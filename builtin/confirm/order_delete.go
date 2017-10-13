package confirm

import (
	"fmt"

	"github.com/wanliu/brain_data/database"
	"github.com/wanliu/flow/builtin/config"
	"github.com/wanliu/flow/context"
)

type OrderDelete struct {
	OrderNo string
	// OrderId uint
}

func (od OrderDelete) SetUp(ctx context.Context) {
	ctx.SetValue(config.CtxKeyConfirm, od)
}

func (od OrderDelete) ClearUp(ctx context.Context) {
	ctx.SetValue(config.CtxKeyConfirm, nil)
}

func (od OrderDelete) Notice(ctx context.Context) string {
	return fmt.Sprintf("确认删除订单号为 %v 的订单么？", od.OrderNo)
}

func (od OrderDelete) Cancel(ctx context.Context) string {
	od.ClearUp(ctx)

	return "已经取消操作"
}

func (od OrderDelete) Confirm(ctx context.Context) string {
	order, err := database.GetOrderByNo(od.OrderNo)

	if err != nil {
		return fmt.Sprintf("确认操作已经过期, %v 号订单已经不存在", od.OrderNo)
	}

	if order.Deletable() {
		return fmt.Sprintf("%v 号订单已经在处理过程中，因此无法删除", od.OrderNo)
	} else {
		err = order.Delete()
		if err == nil {
			return fmt.Sprintf("%v 号订单删除成功", od.OrderNo)
		} else {
			return fmt.Sprintf("%v 号订单删除失败，请访问 http://jiejie.wanliu.biz/orders/%v 进行操作", od.OrderNo, order.GlobelId())
		}
	}

	return ""
}
