package resolves

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hysios/apiai-go"
	"github.com/wanliu/brain_data/database"
	"github.com/wanliu/brain_data/wrapper"
	"github.com/wanliu/flow/builtin/ai"
	"github.com/wanliu/flow/builtin/config"

	"github.com/wanliu/flow/context"
)

// 处理开单的逻辑结构, 不需要是组件
// 作为context的一个部分，或者存在一个Value中
type OrderResolve struct {
	AiParams          ai.AiOrder
	Products          ItemsResolve
	Gifts             ItemsResolve
	Address           string
	ExtractedCustomer string
	Customer          string
	Time              time.Time
	DefTime           string
	Note              string
	Storehouse        string
	UpdatedAt         time.Time
	Editing           bool
	Canceled          bool
	IsResolved        bool
	IsFailed          bool
	OrderSyncQueue    string
	// Current   Resolve

	Id uint

	User       *database.User
	BrainOrder *database.Order
}

func NewOrderResolve(ctx context.Context) *OrderResolve {
	resolve := new(OrderResolve)
	resolve.Touch()

	aiResult := ctx.Value(config.ValueKeyResult).(apiai.Result)

	resolve.AiParams = ai.ApiAiOrder{AiResult: aiResult}
	resolve.ExtractFromParams()

	if syncQueue := ctx.Value(config.CtxKeySyncQueue); syncQueue != nil {
		resolve.OrderSyncQueue = syncQueue.(string)
	}

	if viewer := ctx.Value("Viewer"); viewer != nil {
		user := viewer.(*database.User)
		resolve.User = user
	}

	return resolve
}

// func (r *OrderResolve) Solve() string {
// 	return r.Answer()
// }

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

// 是否达到可以生成订单的条件
func (r OrderResolve) Fulfiled() bool {
	return len(r.Products.Products) > 0 && (r.Address != "" || r.Customer != "")
}

// 是否已经成功生成订单
func (c OrderResolve) Resolved() bool {
	return c.IsResolved
}

func (c OrderResolve) MismatchQuantity() bool {
	return c.Products.MismatchQuantity() || c.Gifts.MismatchQuantity()
}

func (c OrderResolve) Failed() bool {
	return c.IsFailed
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
		unit := strings.Replace(i.Unit, " ", "", -1)

		item := &ItemResolve{
			Resolved: true,
			Name:     name,
			Price:    i.Price,
			Quantity: i.Quantity,
			Product:  name,
			Unit:     unit,
		}

		r.Products.Products = append(r.Products.Products, item)
	}
}

func (r *OrderResolve) ExtractGiftItems() {
	for _, i := range r.AiParams.GiftItems() {
		name := strings.Replace(i.Product, "%", "%%", -1)
		unit := strings.Replace(i.Unit, " ", "", -1)

		item := &ItemResolve{
			Resolved: true,
			Name:     name,
			Price:    i.Price,
			Quantity: i.Quantity,
			Product:  name,
			Unit:     unit,
		}

		r.Gifts.Products = append(r.Gifts.Products, item)
	}
}

func (r *OrderResolve) ExtractAddress() {
	r.Address = r.AiParams.Address()
}

func (r *OrderResolve) ExtractCustomer() {
	r.ExtractedCustomer = r.AiParams.Customer()
	r.CheckExtractedCustomer()
}

func (r *OrderResolve) CheckExtractedCustomer() {
	if r.ExtractedCustomer != "" {
		var count int

		database.DB.Model(&database.People{}).Where("name = ?", r.ExtractedCustomer).Count(&count)

		if count > 0 {
			r.Customer = r.ExtractedCustomer
		}
	}
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
	return r.Products.Empty()
}

// 商品与数量不匹配, 识别为:
// 商品:
// [ 伊利金典纯牛奶 , 红谷粒多1*12瓶, 伊利安幕希酸奶原味 , 190QQ星儿童成长牛奶健固型, xxxxx]
// 数量:
// [ 16 提,  3 盒 , 3 提 ,  6 提 ]
func (r *OrderResolve) Answer(ctx context.Context) string {
	if r.MismatchQuantity() {
		return r.MismatchAnswer()
	} else if r.Fulfiled() {
		return r.PostOrderAndAnswer()
	} else {
		return r.AnswerUnfulfiled(ctx)
	}
}

func (r *OrderResolve) MismatchAnswer() string {
	if r.Products.MismatchQuantity() {
		return r.MismatchProductsAnswer()
	} else if r.Gifts.MismatchQuantity() {
		return r.MismatchGiftsAnswer()
	} else {
		return "数量匹配"
	}
}

func (r *OrderResolve) MismatchProductsAnswer() string {
	result := "商品与数量不匹配, 识别为:\n"

	ps := make([]string, 0)
	pq := make([]string, 0)

	for _, p := range r.Products.Products {
		if p.Product != "" {
			ps = append(ps, p.Product)
		}

		if p.Quantity != 0 {
			pq = append(pq, fmt.Sprint(p.Quantity))
		}
	}

	result = result + "商品:\n"
	result = result + "[" + strings.Join(ps, ", ") + "]\n"

	result = result + "数量:\n"
	result = result + "[" + strings.Join(pq, ", ") + "]\n"

	return result
}

func (r *OrderResolve) MismatchGiftsAnswer() string {
	result := "赠品与数量不匹配, 识别为:\n"

	ps := make([]string, 0)
	pq := make([]string, 0)

	for _, p := range r.Gifts.Products {
		if p.Product != "" {
			ps = append(ps, p.Product)
		}

		if p.Quantity != 0 {
			pq = append(pq, fmt.Sprint(p.Quantity))
		}
	}

	result = result + "赠品:\n"
	result = result + "[" + strings.Join(ps, ", ") + "]\n"

	result = result + "数量:\n"
	result = result + "[" + strings.Join(pq, ", ") + "]\n"

	return result
}

func (r *OrderResolve) PostOrderAndAnswer() string {
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
		order, err := wrapper.CreateFlowOrder(r.OrderSyncQueue, r.Address, r.Note, r.Time, r.Customer, 0, r.Storehouse, items, gifts)

		if err != nil {
			r.IsFailed = true
			return fmt.Sprintf("%v, 订单创建失败", err.Error())
		} else {
			r.IsResolved = true
			return RenderSolvedOrder(order)

			// r.BrainOrder = &order
			// r.Id = order.ID
			// return r.AnswerHead() + r.AnswerBody() + r.AnswerFooter(ctx, order.No, order.GlobelId())
		}
	} else {
		// order, err := r.User.CreateSaledOrder(r.Address, r.Note, r.Time, 0, 0, items, gifts)
		order, err := wrapper.CreateFlowOrder(r.OrderSyncQueue, r.Address, r.Note, r.Time, r.Customer, r.User.ID, r.Storehouse, items, gifts)

		if err != nil {
			r.IsFailed = true
			return fmt.Sprintf("%v, 订单创建失败", err.Error())
		} else {
			r.IsResolved = true
			return RenderSolvedOrder(order)

			// r.BrainOrder = &order
			// r.Id = order.ID
			// return r.AnswerHead() + r.AnswerBody() + r.AnswerFooter(ctx, order.No, order.GlobelId())
		}
	}
}

func (r OrderResolve) AnswerUnfulfiled(ctx context.Context) string {
	return r.AnswerHead() + r.AnswerBody() + r.AnswerFooter(ctx)
}

func (r OrderResolve) AnswerHead() string {
	// desc := "订单正在处理, 已经添加" + CnNum(len(r.Products.Products)) + "种产品"

	// if len(r.Gifts.Products) > 0 {
	// 	desc = desc + ", " + CnNum(len(r.Gifts.Products)) + "种赠品" + "\n"
	// } else {
	// 	desc = desc + "\n"
	// }

	// return desc
	return RenderOrderHeader(r)
}

func (r OrderResolve) AnswerBody() string {
	return RenderOrderBody(r)
}

func (r OrderResolve) AnswerFooter(ctx context.Context) string {
	desc := ""

	if r.ExtractedCustomer != "" && r.Customer == "" {
		confirm := CustomerCreation{
			Customer: r.ExtractedCustomer,
		}

		confirm.SetUp(ctx)

		desc = desc + fmt.Sprintf("\"%v\"为无效的客户，%v\n", r.ExtractedCustomer, confirm.Notice(ctx))
		// desc = desc + fmt.Sprintf("\"%v\"为无效的客户，还缺少客户信息\n", r.ExtractedCustomer)
	} else {
		desc = desc + "还缺少客户信息\n"
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
