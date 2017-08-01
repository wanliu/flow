package resolves

import (
	"strconv"
	"time"

	"github.com/hysios/apiai-go"
	"github.com/wanliu/flow/builtin/ai"
	. "github.com/wanliu/flow/context"
)

// 处理开单的逻辑结构, 不需要是组件
// 作为context的一个部分，或者存在一个Value中
type OrderResolve struct {
	AiParams  ai.AiOrder
	Products  ItemsResolve
	Gifts     ItemsResolve
	Address   string
	Time      time.Time
	DefTime   string
	Current   Resolve
	Note      string
	UpdatedAt time.Time
}

func NewOrderResolve(ctx Context) *OrderResolve {
	resolve := new(OrderResolve)

	aiResult := ctx.Value("Result").(apiai.Result)

	resolve.AiParams = ai.ApiAiOrder{AiResult: aiResult}
	resolve.ExtractFromParams()

	return resolve
}

func (r *OrderResolve) Solve(aiResult apiai.Result) string {
	return r.Answer()
}

func (r OrderResolve) Modifable() bool {
	return true
	// return UpdatedAt.Add(time.Duration(30)*time.Minute) >= time.Now() || r.Unsubmited()
}

// 从ｌｕｉｓ数据构造结构数据
func (r *OrderResolve) ExtractFromParams() {
	r.ExtractItems()
	r.ExtractGiftItems()
	r.ExtractAddress()
	r.ExtractTime()
	r.ExtractNote()
}

func (r *OrderResolve) ExtractItems() {
	for _, i := range r.AiParams.Items() {
		item := &ItemResolve{
			Resolved: true,
			Name:     i.Product,
			Price:    i.Price,
			Quantity: i.Quantity,
			Product:  i.Product,
		}

		r.Products.Products = append(r.Products.Products, item)
	}
}

func (r *OrderResolve) ExtractGiftItems() {
	for _, i := range r.AiParams.GiftItems() {
		item := &ItemResolve{
			Resolved: true,
			Name:     i.Product,
			Price:    i.Price,
			Quantity: i.Quantity,
			Product:  i.Product,
		}

		r.Gifts.Products = append(r.Gifts.Products, item)
	}
}

func (r *OrderResolve) ExtractAddress() {
	r.Address = r.AiParams.Address()
}

func (r *OrderResolve) ExtractTime() {
	r.Time = r.AiParams.Time()
}

func (r *OrderResolve) ExtractNote() {
	r.Note = r.AiParams.Note()
}

func (r *OrderResolve) SetDefTime(t string) {
	r.DefTime = t

	if r.Time.IsZero() && r.DefTime != "" {
		r.SetTimeByDef()
	}
}

func (r *OrderResolve) SetTimeByDef() {
	if r.DefTime == "今天" {
		r.Time = time.Now()
	} else if r.DefTime == "明天" {
		r.Time = time.Now().Add(24 * time.Hour)
	}
}

func (r OrderResolve) EmptyProducts() bool {
	return len(r.Products.Products) == 0
}

func (r OrderResolve) Answer() string {
	desc := "订单已经生成, 共" + CnNum(len(r.Products.Products)) + "种产品"

	if len(r.Gifts.Products) > 0 {
		desc = desc + ", " + CnNum(len(r.Gifts.Products)) + "种赠品" + "\n"
	} else {
		desc = desc + "\n"
	}

	for _, p := range r.Products.Products {
		desc = desc + p.Product + " " + strconv.Itoa(p.Quantity) + "件\n"
	}

	if len(r.Gifts.Products) > 0 {
		desc = desc + "-------赠品-------\n"

		for _, g := range r.Gifts.Products {
			desc = desc + g.Product + " " + strconv.Itoa(g.Quantity) + "件\n"
		}
	}

	desc = desc + "时间:" + r.Time.Format("2006年01月02日") + "\n"

	if r.Note != "" {
		desc = desc + "备注：" + r.Note + "\n"
	}

	if r.Address != "" {
		desc = desc + "地址:" + r.Address + "\n"
	} else {
		desc = desc + "告诉我地址或客户是谁，我就安排送货了\n"
	}

	desc = desc + "订单入口: http://wanliu.biz/orders/"

	return desc
}

func CnNum(num int) string {
	switch num {
	case 1:
		return "一"
	case 2:
		return "二"
	case 3:
		return "三"
	case 4:
		return "四"
	case 5:
		return "五"
	case 6:
		return "六"
	case 7:
		return "七"
	case 8:
		return "八"
	case 9:
		return "九"
	case 10:
		return "十"
	default:
		return strconv.Itoa(num)
	}
}
