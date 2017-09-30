package resolves

import (
	"strconv"
	"strings"
	"time"

	"github.com/hysios/apiai-go"
	"github.com/wanliu/brain_data/database"
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
	Customer  string
	Time      time.Time
	DefTime   string
	Current   Resolve
	Note      string
	UpdatedAt time.Time
	Editing   bool
	Canceled  bool

	User *database.User
}

func NewOrderResolve(ctx Context) *OrderResolve {
	resolve := new(OrderResolve)
	resolve.Touch()

	aiResult := ctx.Value("Result").(apiai.Result)

	resolve.AiParams = ai.ApiAiOrder{AiResult: aiResult}
	resolve.ExtractFromParams()

	if viewer := ctx.Value("Viewer"); viewer != nil {
		user := viewer.(*database.User)
		resolve.User = user
	}

	return resolve
}

func (r *OrderResolve) Solve(aiResult apiai.Result) string {
	return r.Answer()
}

func (r *OrderResolve) Touch() {
	r.UpdatedAt = time.Now()
}

func (r OrderResolve) Modifable(expireMin int) bool {
	return !r.Expired(expireMin) || r.Submited()
}

// TODO
func (r OrderResolve) Cancelable() bool {
	return true
}

// TODO
func (r *OrderResolve) Cancel() bool {
	r.Canceled = true
	return true
}

func (r OrderResolve) Fulfiled() bool {
	return len(r.Products.Products) > 0 && (r.Address != "" || r.Customer != "")
}

func (r OrderResolve) Expired(expireMin int) bool {
	return r.UpdatedAt.Add(time.Duration(expireMin)*time.Minute).UnixNano() < time.Now().UnixNano()
}

// TODO
func (r OrderResolve) Submited() bool {
	return false
}

// 从ｌｕｉｓ数据构造结构数据
func (r *OrderResolve) ExtractFromParams() {
	r.ExtractItems()
	r.ExtractGiftItems()
	r.ExtractAddress()
	r.ExtractCustomer()
	r.ExtractTime()
	r.ExtractNote()
}

func (r *OrderResolve) ExtractItems() {
	for _, i := range r.AiParams.Items() {
		name := strings.Replace(i.Product, "%", "%%", -1)
		item := &ItemResolve{
			Resolved: true,
			Name:     name,
			Price:    i.Price,
			Quantity: i.Quantity,
			Product:  name,
		}

		r.Products.Products = append(r.Products.Products, item)
	}
}

func (r *OrderResolve) ExtractGiftItems() {
	for _, i := range r.AiParams.GiftItems() {
		name := strings.Replace(i.Product, "%", "%%", -1)
		item := &ItemResolve{
			Resolved: true,
			Name:     name,
			Price:    i.Price,
			Quantity: i.Quantity,
			Product:  name,
		}

		r.Gifts.Products = append(r.Gifts.Products, item)
	}
}

func (r *OrderResolve) ExtractAddress() {
	r.Address = r.AiParams.Address()
}

func (r *OrderResolve) ExtractCustomer() {
	r.Customer = r.AiParams.Customer()
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
	if r.Fulfiled() {
		return r.PostOrderAndAnswer()
	} else {
		return r.AnswerHead() + r.AnswerFooter("", "")
	}
}

func (r *OrderResolve) PostOrderAndAnswer() string {
	items := make([]database.OrderItem, 0, 0)
	gifts := make([]database.GiftItem, 0, 0)

	for _, pr := range r.Products.Products {
		item, err := database.NewOrderItem("", pr.Product, pr.Product, pr.Price)
		if err != nil {
			return err.Error()
		}
		items = append(items, item)
	}

	for _, pr := range r.Gifts.Products {
		gift, err := database.NewGiftItem("", pr.Product, pr.Quantity)
		if err != nil {
			return err.Error()
		}

		gifts = append(gifts, gift)
	}

	order, err := r.User.CreateSaledOrder(r.Address, r.Note, r.Time, 0, items, gifts)

	if err != nil {
		return err.Error()
	} else {
		return r.AnswerHead() + r.AnswerBody() + r.AnswerFooter(order.No, order.ID)
	}
}

func (r OrderResolve) AddressInfo() string {
	if r.Address != "" && r.Customer != "" {
		return "地址:" + r.Address + r.Customer + "\n"
	} else if r.Address != "" {
		return "地址:" + r.Address + "\n"
	} else if r.Customer != "" {
		return "客户:" + r.Customer + "\n"
	} else {
		return ""
	}
}

func (r OrderResolve) AnswerHead() string {
	desc := "订单正在处理, 已经添加" + CnNum(len(r.Products.Products)) + "种产品"

	if r.Fulfiled() {
		desc = "订单已经生成, 共" + CnNum(len(r.Products.Products)) + "种产品"
	}

	if len(r.Gifts.Products) > 0 {
		desc = desc + ", " + CnNum(len(r.Gifts.Products)) + "种赠品" + "\n"
	} else {
		desc = desc + "\n"
	}

	return desc
}

func (r OrderResolve) AnswerBody() string {
	desc := ""

	for _, p := range r.Products.Products {
		desc = desc + p.Product + " " + strconv.Itoa(p.Quantity) + "件\n"
	}

	if len(r.Gifts.Products) > 0 {
		desc = desc + "申请的赠品:\n"

		for _, g := range r.Gifts.Products {
			desc = desc + g.Product + " " + strconv.Itoa(g.Quantity) + "件\n"
		}
	}

	desc = desc + "时间:" + r.Time.Format("2006年01月02日") + "\n"

	if r.Note != "" {
		desc = desc + "备注：" + r.Note + "\n"
	}

	return desc
}

func (r OrderResolve) AnswerFooter(no, id interface{}) string {
	desc := ""

	if r.Fulfiled() {
		desc = desc + r.AddressInfo()
		desc = desc + "订单已经生成，订单号为：" + fmt.Sprint(no) + "\n"
		desc = desc + "订单入口: http://wanliu.biz/orders/" + fmt.Sprint(id)
	} else {
		desc = desc + "还缺少收货地址或客户信息\n"
	}

	return desc
}

func CnNum(num int) string {
	switch num {
	case 1:
		return "一"
	case 2:
		return "两"
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
