package resolves

import (
	"fmt"
	"time"

	"github.com/hysios/apiai-go"
	"github.com/wanliu/brain_data/database"
	"github.com/wanliu/flow/builtin/ai"
	"github.com/wanliu/flow/builtin/config"
	"github.com/wanliu/flow/builtin/resolves/templates"
	"github.com/wanliu/flow/context"
)

type CustomerOrdersResolve struct {
	Total      int
	CursorId   *uint
	Done       bool
	Prefetched int
	Fetched    int
	Per        int

	CustomerName string
	Customer     *database.People
	ProductName  string
	Product      *database.Product

	Count      int
	QuertyTime time.Time
	Duration   string
	BeginT     *time.Time
	EndT       *time.Time

	LastActiveTime *time.Time
}

func NewCusOrdersResolve(ctx context.Context, perPage int) *CustomerOrdersResolve {
	var product string

	if perPage <= 0 {
		perPage = 5
	}

	aiResult := ctx.Value("Result").(apiai.Result)
	aiParams := ai.ApiAiOrder{AiResult: aiResult}

	customerName := aiParams.Customer()
	queryTime := aiParams.Time()
	count := aiParams.Count()
	duration := aiParams.Duration()

	if productsInt, exist := aiResult.Params["products"]; exist {
		product = productsInt.(string)
	}

	beginT, endT := getBeginAndEndTime(duration, queryTime)

	rsv := CustomerOrdersResolve{
		CustomerName: customerName,
		ProductName:  product,
		Duration:     duration,
		QuertyTime:   queryTime,
		BeginT:       beginT,
		EndT:         endT,
		Count:        count,
		Per:          perPage,
		Total:        -1,
	}

	rsv.FetchCustomer()
	rsv.FetchProduct()

	return &rsv
}

func (r *CustomerOrdersResolve) Answer() string {
	var orders *[]database.Order

	if r.CustomerName == "" {
		r.Done = true
		return "未提供客户，无法查询客户订单"
	}

	if r.Done {
		return "查询的订单已经显示完毕"
	}

	if r.Customer == nil {
		r.Done = true
		return fmt.Sprintf("客户\"%v\"不存在，无法查询客户订单", r.CustomerName)
	}

	if r.ProductName != "" && r.Product == nil {
		r.Done = true
		return fmt.Sprintf("商品\"%v\"不存在，无法查询包含该商品的订单", r.ProductName)
	}

	if r.Total == -1 {
		r.FetchTotal()
	}

	if r.Total > 0 {
		orders = r.FetchOrders()
	}

	r.Done = r.IsDone()

	return templates.RenderCustomerOrders(r, orders)

	// result := r.AnswerHeader()

	// if r.Total == -1 {
	// 	r.FetchTotal()
	// }

	// if r.Total == 0 {
	// 	result = result + "没有订单可以显示"
	// 	return result
	// }

	// orders := r.FetchOrders()
	// if r.Count > 0 {
	// 	count := r.Count

	// 	if r.Count > r.Total {
	// 		count = r.Total
	// 	}

	// 	result = result + fmt.Sprintf("共%v个订单，显示最近的%v个订单，以下为第%v到第%v个：\n", r.Total, count, r.Prefetched+1, r.Fetched)
	// } else {
	// 	result = result + fmt.Sprintf("共%v个订单，以下为第%v到第%v个：\n", r.Total, r.Prefetched+1, r.Fetched)
	// }
	// r.Done = r.IsDone()

	// if orders != nil && len(*orders) > 0 {
	// 	for _, order := range *orders {
	// 		result = result + "------------------------\n"
	// 		result = result + fmt.Sprintf("订单号：%v\n总金额：%v\n送货时间：%v\n", order.No, fmt.Sprintf("%.2f", order.Amount), order.DeliveryTime.Format("2006年01月02日"))
	// 		if order.Note != "" {
	// 			result = result + fmt.Sprint("备注：%v\n", order.Note)
	// 		}

	// 		result = result + "商品:\n"
	// 		for _, item := range order.OrderItems {
	// 			result = result + fmt.Sprintf("  %v %v%v\n", item.ProductName, item.Quantity, item.Unit)
	// 		}

	// 		if len(order.GiftItems) > 0 {
	// 			result = result + "赠品:\n"
	// 			for _, gift := range order.GiftItems {
	// 				result = result + fmt.Sprintf("  %v %v%v\n", gift.ProductName, gift.Quantity, gift.Unit)
	// 			}
	// 		}
	// 	}

	// 	result = result + "------------------------\n"

	// 	if !r.Done {
	// 		result = result + "输入\"继续\"，或者\"下一页\"，查看剩下的订单\n"
	// 	} else {
	// 		total := 0
	// 		if r.Count > 0 && r.Count <= r.Total {
	// 			total = r.Count
	// 		} else {
	// 			total = r.Total
	// 		}

	// 		result = result + fmt.Sprintf("%v个订单已经全部显示完毕！\n", total)
	// 	}
	// } else {
	// 	result = result + "没有订单可以显示"
	// }

	// return result
}

// func (r *CustomerOrdersResolve) AnswerHeader() string {
// 	if r.Duration != "" {
// 		if r.ProductName != "" {
// 			return fmt.Sprintf("查询客户\"%v\"%v包含\"%v\"的订单\n", r.CustomerName, r.Duration, r.ProductName)
// 		} else {
// 			return fmt.Sprintf("查询客户\"%v\"%v的订单\n", r.CustomerName, r.Duration)
// 		}
// 	} else if !r.QuertyTime.IsZero() {
// 		if r.ProductName != "" {
// 			return fmt.Sprintf("查询客户\"%v\"%v包含\"%v\"的订单\n", r.CustomerName, r.QuertyTime.Format("2006年1月2日"), r.ProductName)
// 		} else {
// 			return fmt.Sprintf("查询客户\"%v\"%v的订单\n", r.CustomerName, r.QuertyTime.Format("2006年1月2日"))
// 		}
// 	} else {
// 		if r.ProductName != "" {
// 			return fmt.Sprintf("查询客户\"%v\"最近包含\"%v\"的订单\n", r.CustomerName, r.ProductName)
// 		} else {
// 			return fmt.Sprintf("查询客户\"%v\"最近的订单\n", r.CustomerName)
// 		}
// 	}
// }

func (r *CustomerOrdersResolve) FetchCustomer() {
	if r.Customer == nil {
		r.Customer, _ = database.GetPersonByName(r.CustomerName)
	}
}

func (r *CustomerOrdersResolve) FetchProduct() {
	if r.Product == nil && r.ProductName != "" {
		r.Product, _ = database.GetProductByName(r.ProductName)
	}
}

func (r *CustomerOrdersResolve) FetchTotal() {
	if r.Customer != nil {
		pids := []uint{}

		if r.Product != nil {
			pids = []uint{r.Product.ID}
		}

		r.Total = r.Customer.GetRecentOrdersTotal(pids, r.BeginT, r.EndT)
	} else {
		r.Total = 0
	}
}

func (r *CustomerOrdersResolve) FetchOrders() *[]database.Order {
	var orders []database.Order

	if r.Customer == nil {
		return nil
	}

	per := r.Per

	if r.Count > 0 && r.Count-r.Fetched < per {
		per = r.Count - r.Fetched
	}

	pids := make([]uint, 0, 0)

	if r.Product != nil {
		pids = []uint{r.Product.ID}
	}

	r.Customer.GetRecentOrders(&orders, pids, r.BeginT, r.EndT, r.CursorId, 0, per)

	l := len(orders)
	if l > 0 {
		r.Prefetched = r.Fetched

		r.Fetched = r.Fetched + l
		lastOrder := orders[l-1]
		lid := lastOrder.ID
		r.CursorId = &lid
	}

	return &orders
}

// A.有指定数目时：
//   1.订单数大于要求的数目，要求的数目已经达到
//   2.订单数不足要求的数目，但是所有的订单显示完毕
// B.无指定数目时: 但是所有的订单显示完毕
func (r *CustomerOrdersResolve) IsDone() bool {
	if r.Count != 0 {
		if r.Total >= r.Count {
			return r.Fetched >= r.Count
		} else {
			return r.Fetched >= r.Total
		}
	} else {
		return r.Fetched >= r.Total
	}
}

func (r *CustomerOrdersResolve) Setup(ctx context.Context) {
	if !r.Done {
		ctx.SetValue(config.CtxKeyCusOrders, r)
	}
}

func (r *CustomerOrdersResolve) ClearIfDone(ctx context.Context) {
	if r.Done {
		r.Clear(ctx)
	}
}

func (r *CustomerOrdersResolve) Clear(ctx context.Context) {
	in := ctx.Value(config.CtxKeyCusOrders)
	if in != nil {
		rsv := in.(*CustomerOrdersResolve)

		if r == rsv {
			ctx.SetValue(config.CtxKeyCusOrders, nil)
		}
	}
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
		return getTimeFromDuration(duration)
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
		// tt := lt.Add(-48 * time.Hour)
		tt := lt.AddDate(0, 0, -2)
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
		// tt := tt.Add(-7 * 24 * time.Hour)
		tt := wt.AddDate(0, 0, -7)
		bt := BeginOfWeek(tt)
		et := EndOfWeek(tt)

		return &bt, &et
	case "近一个月":
		// bt := now.Add(-30 * 24 * time.Hour)
		bt := now.AddDate(0, 0, -30)
		return &bt, &now
	case "近一个星期":
		// bt := now.Add(-7 * 24 * time.Hour)
		bt := now.AddDate(0, 0, -7)
		return &bt, &now
	}

	return nil, nil
}

func getBETimeOfMonth(m int) (*time.Time, *time.Time) {
	var (
		bt time.Time
		et time.Time
	)

	t := time.Now()
	year, month, _ := t.Date()

	if int(month) >= m {
		tt := time.Date(year, time.Month(m), 1, 0, 0, 0, 0, t.Location())
		bt = BeginOfMonth(tt)
		et = EndOfMonth(tt)
	} else {
		year = year - 1
		tt := time.Date(year, time.Month(m), 1, 0, 0, 0, 0, t.Location())
		bt = BeginOfMonth(tt)
		et = EndOfMonth(tt)
	}

	return &bt, &et
}

func BeginOfMonth(t time.Time) time.Time {
	year, month, _ := t.Date()
	return time.Date(year, month, 1, 0, 0, 0, 0, t.Location())
}

func BeginOfWeek(t time.Time) time.Time {
	step := int(t.Weekday() - 1)
	tt := t.AddDate(0, 0, -step)

	return BeginOfDay(tt)
}

func EndOfWeek(t time.Time) time.Time {
	step := int(7 - t.Weekday())
	tt := t.AddDate(0, 0, step)

	return EndOfDay(tt)
}

func EndOfMonth(t time.Time) time.Time {
	year, month, _ := t.Date()

	if month == 12 {
		year = year + 1
		month = 1
	} else {
		month = month + 1
	}

	return time.Date(year, month, 0, 23, 59, 59, 999999999, t.Location())
}

func BeginOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

func EndOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 23, 59, 59, 999999999, t.Location())
}
