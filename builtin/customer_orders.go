package builtin

import (
	"fmt"
	"time"

	"github.com/hysios/apiai-go"
	"github.com/wanliu/brain_data/database"
	"github.com/wanliu/flow/builtin/ai"
	"github.com/wanliu/flow/context"
)

const (
	Per = 5
)

type CustomerOrders struct {
	TryGetEntities
	Ctx  <-chan context.Context
	Type <-chan string
	Out  chan<- ReplyData

	Total    int
	CursorId *uint
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
	duration := aiParams.Duration()

	beginT, endT := getBeginAndEndTime(duration, queryTime)
	// fmt.Printf("[NUMBER]---->%v", count)

	// if count == 0 {
	// 	count = 2
	// }

	// if count > 5 {
	// 	count = 5
	// }

	// if customer == "" {
	// 	c.Out <- ReplyData{"请提供要查询的客户", ctx}
	// 	return
	// }

	// var person database.People
	// err := database.DB.Where("name = ?", customer).First(&person).Error
	// if nil != err || database.DB.NewRecord(person) {
	// 	c.Out <- ReplyData{fmt.Sprintf("客户\"%v\"不存在", customer), ctx}
	// 	return
	// }

	// var orders []database.Order
	// result := ""

	// if queryTime.IsZero() {
	// 	person.GetRecentOrders(&orders, nil, nil, count)
	// 	if len(orders) == 0 {
	// 		reply := fmt.Sprintf("客户\"%v\"最近没有订单", customer)
	// 		c.Out <- ReplyData{reply, ctx}
	// 		return
	// 	}

	// 	if count > len(orders) {
	// 		result = fmt.Sprintf("客户\"%v\"只有%v个订单：\n", customer, len(orders))
	// 	} else {
	// 		result = fmt.Sprintf("客户\"%v\"最近的%v个订单：\n", customer, len(orders))
	// 	}
	// } else {
	// 	person.GetRecentOrders(&orders, &queryTime, nil, count)
	// 	date := queryTime.Format("2006年01月02日")

	// 	if len(orders) == 0 {
	// 		reply := fmt.Sprintf("客户\"%v\"在%v没有订单", customer, date)
	// 		c.Out <- ReplyData{reply, ctx}
	// 		return
	// 	}

	// 	if count > len(orders) {
	// 		result = fmt.Sprintf("客户\"%v\"%v只有%v个订单：\n", customer, date, len(orders))
	// 	} else {
	// 		result = fmt.Sprintf("客户\"%v\"%v最近的%v个订单：\n", customer, date, len(orders))
	// 	}
	// }

	// for _, order := range orders {
	// 	result = result + "------------------------\n"
	// 	result = result + fmt.Sprintf("订单号：%v\n总金额：%v\n送货时间：%v\n", order.No, order.Amount, order.DeliveryTime.Format("2006年01月02日"))
	// 	if order.Note != "" {
	// 		result = result + fmt.Sprint("备注：%v\n", order.Note)
	// 	}

	// 	result = result + "商品:\n"
	// 	for _, item := range order.OrderItems {
	// 		result = result + fmt.Sprintf("  %v %v%v\n", item.ProductName, item.Quantity, item.Unit)
	// 	}

	// 	if len(order.GiftItems) > 0 {
	// 		result = result + "赠品:\n"
	// 		for _, gift := range order.GiftItems {
	// 			result = result + fmt.Sprintf("  %v %v%v\n", gift.ProductName, gift.Quantity, gift.Unit)
	// 		}
	// 	}
	// }

	c.Out <- ReplyData{result, ctx}
}

func getBeginAndEndTime(duration string, queryTime time.Time) (*time.Time, *time.Time) {
	if duration == "" {
		if queryTime.IsZero() {
			return nil, nil
		} else {
			b := BeginOfDay(queryTime)
			e := EndOfDay(queryTime)
			return &b, &e
		}
	} else {
		return getTimeFromDuration(duration, queryTime)
	}
}

func getTimeFromDuration(duration string) (*time.Time, *time.Time) {
	if duration == "" {
		return nil, nil
	}

	now := time.Now()

	switch duration {
	case "本月":
		bt := BeginOfMonth(now)

		return &bt, &now
	case "上个月":
		lt := BeginOfMonth(now)
		tt := lt.Add(-48 * time.Hour)
		bt := BeginOfMonth(tt)
		et := EndOfMonth(tt)

		return &bt, &et
	case "1月":
		return getBETimeOfMonth(1)
	case "2月":
		return getBETimeOfMonth(2)
	case "3月":
		return getBETimeOfMonth(3)
	case "4月":
		return getBETimeOfMonth(4)
	case "5月":
		return getBETimeOfMonth(5)
	case "6月":
		return getBETimeOfMonth(6)
	case "7月":
		return getBETimeOfMonth(7)
	case "8月":
		return getBETimeOfMonth(8)
	case "9月":
		return getBETimeOfMonth(9)
	case "10月":
		return getBETimeOfMonth(10)
	case "11月":
		return getBETimeOfMonth(11)
	case "12月":
		return getBETimeOfMonth(12)
	case "本周":
		bt := BeginOfWeek(now)
		return &bt, &now
	case "上周":
		wt := BeginOfWeek(now)
		tt := tt.Add(-7 * 24 * time.Hour)
		bt := BeginOfWeek(tt)

		return &bt, &now
	case "近一个月":
		bt := now.Add(-30 * 24 * time.Hour)
		return &bt, &now
	case "近一个星期":
		bt := now.Add(-7 * 24 * time.Hour)
		return &bt, &now
	}
}

func getBETimeOfMonth(m int) (*time.Time, *time.Time) {
	year, month, _ := now.Date()

	if month >= m {
		tt := time.Date(year, m, 1, 0, 0, 0, 0, t.Location())
		bt := BeginOfMonth(tt)
		et := EndOfMonth(tt)
	} else {
		year = year - 1
		tt := time.Date(year, m, 1, 0, 0, 0, 0, t.Location())
		bt := BeginOfMonth(tt)
		et := EndOfMonth(tt)
	}

	return &bt, &et
}

func BeginOfMonth(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, 0, 0, 0, 0, 0, t.Location())
}

func BeginOfWeek(t time.Time) time.Time {
	step := int(time.Now().Weekday() - 1)
	tt := now.Add(-step * 24 * time.Hour)

	return BeginOfDay(tt)
}

func EndOfMonth(t time.Time) time.Time {
	year, month, day := t.Date()

	if month == 12 {
		year = year + 1
		month = 1
	} else {
		month = month + 1
	}

	return time.Date(year, month, 0, 0, 0, 0, 0, t.Location())
}

func BeginOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

func EndOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 23, 59, 59, 0, t.Location())
}
