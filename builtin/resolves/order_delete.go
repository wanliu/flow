// 开单产品选择
package resolves

import (
	"fmt"

	"github.com/wanliu/brain_data/database"
	"github.com/wanliu/flow/builtin/config"
	"github.com/wanliu/flow/context"
)

type OrderDeleteResolve struct {
	OrderNo string
}

func (odr OrderDeleteResolve) Hint() string {
	return "请告诉我你要取消/删除订单的订单号"
}

func (odr OrderDeleteResolve) SetUp(ctx context.Context) {
	// need order number to resolve
	ctx.SetValue(config.CtxKeyOrderNum, odr)
}

func (odr OrderDeleteResolve) ClearUp(ctx context.Context) {
	// need order number to resolve, clear mark
	ctx.SetValue(config.CtxKeyOrderNum, nil)
}

func (odr OrderDeleteResolve) Resolve(orderNo string, ctx context.Context) string {
	order, err := database.GetOrderByNo(orderNo)

	if err != nil {
		return fmt.Sprintf("找不到订单号为 %v 的订单", orderNo)
	} else {
		err = order.Delete()
		if err == nil {
			odr.ClearUp(ctx)
			return fmt.Sprintf("%v 号订单删除成功", orderNo)
		} else {
			return fmt.Sprintf("%v 号订单删除失败，请访问 http://jiejie.wanliu.biz/orders/%v 进行操作", orderNo, order.GlobelId())
		}
	}

	return ""
}
