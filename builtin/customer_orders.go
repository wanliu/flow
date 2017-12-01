package builtin

import (
	"fmt"

	"github.com/hysios/apiai-go"
	"github.com/wanliu/brain_data/database"
	"github.com/wanliu/flow/builtin/ai"
	"github.com/wanliu/flow/context"
)

type CustomerOrders struct {
	TryGetEntities
	Ctx  <-chan context.Context
	Type <-chan string
	Out  chan<- ReplyData
}

func NewCustomerOrders() interface{} {
	return new(CustomerOrders)
}

func (c *CustomerOrders) OnCtx(ctx context.Context) {
	aiResult := ctx.Value("Result").(apiai.Result)
	aiParams := ai.ApiAiOrder{AiResult: aiResult}

	customer := aiParams.Customer()
	queryTime := aiParams.Time()
	count := aiParams.Count()

	fmt.Printf("[NUMBER]---->%v", count)

	if count == 0 {
		count = 2
	}

	if customer == "" {
		c.Out <- ReplyData{"请提供要查询的客户", ctx}
		return
	}

	var person database.People
	err := database.DB.Where("name = ?", customer).First(&person).Error
	if nil != err || database.DB.NewRecord(person) {
		c.Out <- ReplyData{fmt.Sprintf("客户\"%v\"不存在", customer), ctx}
		return
	}

	var orders []database.Order
	result := ""

	if queryTime.IsZero() {
		person.GetRecentOrders(&orders, nil, count)
		if len(orders) == 0 {
			reply := fmt.Sprintf("客户\"%v\"最近没有订单", customer)
			c.Out <- ReplyData{reply, ctx}
			return
		}

		if count > len(orders) {
			result = fmt.Sprintf("客户\"%v\"只有%v个订单：\n", customer, len(orders))
		} else {
			result = fmt.Sprintf("客户\"%v\"最近的%v个订单：\n", customer, len(orders))
		}
	} else {
		person.GetRecentOrders(&orders, &queryTime, count)
		date := queryTime.Format("2006年01月02日")

		if len(orders) == 0 {
			reply := fmt.Sprintf("客户\"%v\"在%v没有订单", customer, date)
			c.Out <- ReplyData{reply, ctx}
			return
		}

		if count > len(orders) {
			result = fmt.Sprintf("客户\"%v\"%v只有%v个订单：\n", customer, date, len(orders))
		} else {
			result = fmt.Sprintf("客户\"%v\"%v最近的%v个订单：\n", customer, date, len(orders))
		}
	}

	for _, order := range orders {
		result = result + "------------------------\n"
		result = result + fmt.Sprintf("订单号：%v\n总金额：%v\n送货时间：%v\n", order.No, order.Amount, order.DeliveryTime.Format("2006年01月02日"))
		if order.Note != "" {
			result = result + fmt.Sprint("备注：%v\n", order.Note)
		}

		result = result + "商品:\n"
		for _, item := range order.OrderItems {
			result = result + fmt.Sprintf("  %v %v%v\n", item.ProductName, item.Quantity, item.Unit)
		}

		if len(order.GiftItems) > 0 {
			result = result + "赠品:\n"
			for _, gift := range order.GiftItems {
				result = result + fmt.Sprintf("  %v %v%v\n", gift.ProductName, gift.Quantity, gift.Unit)
			}
		}
	}

	c.Out <- ReplyData{result, ctx}
}
