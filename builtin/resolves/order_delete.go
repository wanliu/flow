// 开单产品选择
package resolves

import (
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
	ctx.SetValue(config.CtxKeyOrderDel, odr)
}

func (odr OrderDeleteResolve) ClearUp(ctx context.Context) {
	ctx.SetValue(config.CtxKeyOrderDel, nil)
}
