package resolves

import (
	"fmt"

	"github.com/wanliu/brain_data/database"
	"github.com/wanliu/flow/builtin/config"
	"github.com/wanliu/flow/context"
)

type OrderDeleteConfirm struct {
	OrderNo string
	// OrderId uint
}

func (od OrderDeleteConfirm) SetUp(ctx context.Context) {
	ctx.SetCtxValue(config.CtxKeyConfirm, od)
}

func (od OrderDeleteConfirm) ClearUp(ctx context.Context) {
	ctx.SetCtxValue(config.CtxKeyConfirm, nil)
}

func (od OrderDeleteConfirm) Notice(ctx context.Context) string {
	return fmt.Sprintf("确认删除订单号为 %v 的订单么？", od.OrderNo)
}

func (od OrderDeleteConfirm) Cancel(ctx context.Context) string {
	od.ClearUp(ctx)

	return "已经取消操作"
}

func (od OrderDeleteConfirm) Confirm(ctx context.Context) string {
	order, err := database.GetOrderByNo(od.OrderNo)

	if err != nil {
		return fmt.Sprintf("确认操作已经过期, %v 号订单已经不存在", od.OrderNo)
	}

	if order.Deletable() {
		err = order.Delete()
		if err == nil {
			return fmt.Sprintf("%v 号订单删除成功", od.OrderNo)
		} else {
			return fmt.Sprintf("%v 号订单删除失败，请访问 https://jiejie.io/orders/%v 进行操作", od.OrderNo, order.GlobelId())
		}
	} else {
		return fmt.Sprintf("%v 号订单已经在处理过程中，因此无法删除", od.OrderNo)
	}

	return ""
}
