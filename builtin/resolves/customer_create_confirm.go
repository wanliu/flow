package resolves

import (
	"fmt"

	"github.com/wanliu/brain_data/database"
	"github.com/wanliu/flow/builtin/config"
	"github.com/wanliu/flow/context"
)

type CustomerCreation struct {
	Customer string
	Address  string
	Phone    string
	// OrderId uint
}

func (cc CustomerCreation) SetUp(ctx context.Context) {
	ctx.SetCtxValue(config.CtxKeyConfirm, cc)
}

func (cc CustomerCreation) ClearUp(ctx context.Context) {
	ctx.SetCtxValue(config.CtxKeyConfirm, nil)
}

func (cc CustomerCreation) Notice(ctx context.Context) string {
	return fmt.Sprintf("是否添加\"%v\"为新的客户？？", cc.Customer)
}

func (cc CustomerCreation) Cancel(ctx context.Context) string {
	cc.ClearUp(ctx)

	oInt := ctx.CtxValue(config.CtxKeyOrder)
	// confirm := ctx.CtxValue(config.CtxKeyConfirm)

	if oInt != nil {
		order := oInt.(OrderResolve)
		order.ExtractedCustomer = ""

		if order.Expired(config.SesssionExpiredMinutes) {
			return fmt.Sprintf("已经取消添加\"%v\"为新客户的操作", cc.Customer)
		}

		ctx.SetCtxValue(config.CtxKeyOrder, order)
		return fmt.Sprintf("已经取消添加\"%v\"为新客户的操作", cc.Customer)
	} else {
		return fmt.Sprintf("已经取消添加\"%v\"为新客户的操作, 当前没有正在进行中的订单", cc.Customer)
	}

	return ""
}

func (cc CustomerCreation) Confirm(ctx context.Context) (string, interface{}) {
	person := database.People{
		Name: cc.Customer,
	}

	err := database.CreatePerson(&person)

	if err == nil {
		oInt := ctx.CtxValue(config.CtxKeyOrder)
		// confirm := ctx.CtxValue(config.CtxKeyConfirm)

		if oInt != nil {
			order := oInt.(OrderResolve)

			order.Customer = person.Name

			if order.Expired(config.SesssionExpiredMinutes) {
				return fmt.Sprintf("添加了新的客户\"%v\", 当前没有正在进行中的订单", cc.Customer), nil
			}

			// dataReply := DataReply{
			// 	Type:   "info",
			// 	On:     "order",
			// 	Action: "update",
			// 	Data:   data,
			// }
			reply, data := order.Answer(ctx)
			reply = fmt.Sprintf("添加了新的客户\"%v\"\n%v", cc.Customer, reply)

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
			return fmt.Sprintf("添加了新的客户\"%v\", 当前没有正在进行中的订单", cc.Customer), nil
		}
	} else {
		return fmt.Sprintf("添加新的客户\"%v\"失败，%v", cc.Customer, err.Error()), nil
	}
}
