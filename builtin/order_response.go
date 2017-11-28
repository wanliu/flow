package builtin

import (
	"fmt"
	"strconv"

	"github.com/wanliu/brain_data/confirm"
	"github.com/wanliu/brain_data/database"
	"github.com/wanliu/brain_data/wrapper"
	"github.com/wanliu/flow/builtin/resolves"
	"github.com/wanliu/flow/context"
)

func OrderResponse(r *resolves.OrderResolve, ctx context.Context) {
	if r.Fulfiled() {
		return postOrderAndAnswer(r, ctx)
	} else {
		return answerHead(r) + answerFooter(r, ctx, "", "")
	}
}

func postOrderAndAnswer(r *OrderResolve, ctx context.Context) string {
	items := make([]database.OrderItem, 0, 0)
	gifts := make([]database.GiftItem, 0, 0)

	for _, pr := range r.Products.Products {
		item, err := database.NewOrderItem("", pr.Product, uint(pr.Quantity), pr.Unit, pr.Price)
		if err != nil {
			r.IsFailed = true
			return fmt.Sprintf("%v, 订单创建失败", err.Error())
		}
		items = append(items, *item)
	}

	for _, pr := range r.Gifts.Products {
		gift, err := database.NewGiftItem("", pr.Product, uint(pr.Quantity), pr.Unit)
		if err != nil {
			r.IsFailed = true
			return fmt.Sprintf("%v, 订单创建失败", err.Error())
		}

		gifts = append(gifts, *gift)
	}

	if r.User == nil {
		// return "无法创建订单，请与工作人员联系！"
		order, err := wrapper.CreateFlowOrder(r.Address, r.Note, r.Time, r.Customer, 0, r.Storehouse, items, gifts)

		if err != nil {
			r.IsFailed = true
			return fmt.Sprintf("%v, 订单创建失败", err.Error())
		} else {
			r.IsResolved = true
			r.BrainOrder = &order
			r.Id = order.ID
			return answerHead(r) + answerBody(r) + answerFooter(r, ctx, order.No, order.GlobelId())
		}
	} else {
		// order, err := r.User.CreateSaledOrder(r.Address, r.Note, r.Time, 0, 0, items, gifts)
		order, err := wrapper.CreateFlowOrder(r.Address, r.Note, r.Time, r.Customer, r.User.ID, r.Storehouse, items, gifts)

		if err != nil {
			r.IsFailed = true
			return fmt.Sprintf("%v, 订单创建失败", err.Error())
		} else {
			r.IsResolved = true
			r.BrainOrder = &order
			r.Id = order.ID
			return answerHead(r) + answerBody(r) + answerFooter(r, ctx, order.No, order.GlobelId())
		}
	}
}

func answerHead(r *OrderResolve) string {
	desc := "订单正在处理, 已经添加" + resolves.CnNum(len(r.Products.Products)) + "种产品"

	if r.Fulfiled() {
		desc = "订单已经生成, 共" + resolves.CnNum(len(r.Products.Products)) + "种产品"
	}

	if len(r.Gifts.Products) > 0 {
		desc = desc + ", " + resolves.CnNum(len(r.Gifts.Products)) + "种赠品" + "\n"
	} else {
		desc = desc + "\n"
	}

	return desc
}

func answerBody(r *OrderResolve) string {
	desc := ""

	if r.IsResolved && r.BrainOrder != nil {
		for _, i := range r.BrainOrder.OrderItems {
			desc = desc + fmt.Sprintf("%v %v %v\n", i.ProductName, i.Quantity, i.Unit)
		}

		if len(r.BrainOrder.GiftItems) > 0 {
			desc = desc + "申请的赠品:\n"

			for _, g := range r.BrainOrder.GiftItems {
				desc = desc + fmt.Sprintf("%v %v %v\n", g.ProductName, g.Quantity, g.Unit)
			}
		}

		desc = desc + fmt.Sprintf("时间:%v\n", r.BrainOrder.DeliveryTime.Format("2006年01月02日"))

		if r.BrainOrder.Note != "" {
			desc = desc + "备注：" + r.BrainOrder.Note + "\n"
		}
	} else {
		for _, p := range r.Products.Products {
			// desc = desc + p.Product + " " + strconv.Itoa(p.Quantity) + p.Unit + "\n"
			desc = desc + fmt.Sprintf("%v %v %v\n", p.Product, p.Quantity, p.Unit)

		}

		if len(r.Gifts.Products) > 0 {
			desc = desc + "申请的赠品:\n"

			for _, g := range r.Gifts.Products {
				// desc = desc + g.Product + " " + strconv.Itoa(g.Quantity) + g.Unit + "\n"
				desc = desc + fmt.Sprintf("%v %v %v\n", g.Product, g.Quantity, g.Unit)
			}
		}

		// desc = desc + "时间:" + r.Time.Format("2006年01月02日") + "\n"
		desc = desc + fmt.Sprintf("时间:%v\n", r.Time.Format("2006年01月02日"))

		if r.Note != "" {
			desc = desc + "备注：" + r.Note + "\n"
		}
	}

	return desc
}

func answerFooter(r *OrderResolve, ctx context.Context, no, id interface{}) string {
	desc := ""

	if r.Fulfiled() {
		desc = desc + r.AddressInfo()
		desc = desc + "订单已经生成，订单号为：" + fmt.Sprint(no) + "\n"
		desc = desc + "订单入口: http://jiejie.wanliu.biz/order/QueryDetail/" + fmt.Sprint(id)
	} else {
		if r.ExtractedCustomer != "" && r.Customer == "" {
			customerConfirm := confirm.CustomerCreation{
				Name: r.ExtractedCustomer,
			}

			customerConfirm.SetUp(ctx)
			// desc = desc + fmt.Sprintf("\"%v\"为无效的客户，还缺少客户信息\n", r.ExtractedCustomer)
			desc = desc + fmt.Sprintf("\"%v\"为无效的客户，还缺少客户信息\n", r.ExtractedCustomer)
		} else {
			desc = desc + "还缺少客户信息\n"
		}
	}

	return desc
}
