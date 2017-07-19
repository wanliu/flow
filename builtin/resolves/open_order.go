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
type OpenOrderResolve struct {
	AiParams ai.AiOrder
	Products ItemsResolve
	Gifts    ItemsResolve
	Address  string
	Time     time.Time
	DefTime  string
	Current  Resolve
	Note     string
}

func NewOpenOrderResolve(ctx Context) *OpenOrderResolve {
	resolve := new(OpenOrderResolve)

	aiResult := ctx.Value("Result").(apiai.Result)

	resolve.AiParams = ai.ApiAiOrder{AiResult: aiResult}
	resolve.ExtractFromLuis()

	return resolve
}

func (r *OpenOrderResolve) Solve(aiResult apiai.Result) string {
	return r.Answer()
}

// 从ｌｕｉｓ数据构造结构数据
func (r *OpenOrderResolve) ExtractFromLuis() {
	r.ExtractItems()
	r.ExtractGiftItems()
	r.ExtractAddress()
	r.ExtractTime()
	r.ExtractNote()
}

func (r *OpenOrderResolve) ExtractItems() {
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

func (r *OpenOrderResolve) ExtractGiftItems() {
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

func (r *OpenOrderResolve) ExtractAddress() {
	r.Address = r.AiParams.Address()
}

func (r *OpenOrderResolve) ExtractTime() {
	r.Time = r.AiParams.Time()
}

func (r *OpenOrderResolve) ExtractNote() {
	r.Note = r.AiParams.Note()
}

func (r *OpenOrderResolve) SetDefTime(t string) {
	r.DefTime = t

	if r.Time.IsZero() && r.DefTime != "" {
		r.SetTimeByDef()
	}
}

func (r *OpenOrderResolve) SetTimeByDef() {
	if r.DefTime == "今天" {
		r.Time = time.Now()
	} else if r.DefTime == "明天" {
		r.Time = time.Now().Add(24 * time.Hour)
	}
}

func (r OpenOrderResolve) EmptyProducts() bool {
	return len(r.Products.Products) == 0
}

func (r OpenOrderResolve) Answer() string {
	desc := ""

	desc = desc + "=== 订单输入完成 ===\n"
	desc = desc + "本订单包含如下商品：" + "\n"

	for _, p := range r.Products.Products {
		desc = desc + p.Product + " " + strconv.Itoa(p.Quantity) + "件\n"
	}

	if len(r.Gifts.Products) > 0 {
		desc = desc + "-------赠品-------\n"

		for _, g := range r.Gifts.Products {
			desc = desc + g.Product + " " + strconv.Itoa(g.Quantity) + "件\n"
		}
	}

	desc = desc + "地址:" + r.Address + "\n"
	desc = desc + "送货时间" + r.Time.Format("2006年01月02日") + "\n"

	if r.Note != "" {
		desc = desc + "备注：" + r.Note + "\n"
	}

	desc = desc + "=== 结束 ===\n"

	return desc
}

type Res struct {
	Id           int
	Confirm_path string
	Error        string
}
