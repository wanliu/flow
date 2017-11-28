package resolves

import (
	"fmt"
	// "strconv"
	"strings"
	"time"

	"github.com/hysios/apiai-go"
	"github.com/wanliu/brain_data/database"
	// "github.com/wanliu/brain_data/wrapper"
	"github.com/wanliu/flow/builtin/ai"

	. "github.com/wanliu/flow/context"
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
	// Current   Resolve

	Id uint

	User       *database.User
	BrainOrder *database.Order
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

// 是否达到可以生成订单的条件
func (r OrderResolve) Fulfiled() bool {
	return len(r.Products.Products) > 0 && (r.Address != "" || r.Customer != "")
}

// 是否已经成功生成订单
func (c OrderResolve) Resolved() bool {
	return c.IsResolved
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
	return len(r.Products.Products) == 0
}

func (r OrderResolve) AddressInfo() string {
	// if r.Address != "" && r.Customer != "" {
	// 	return "地址:" + r.Address + r.Customer + "\n"
	// } else if r.Address != "" {
	// 	return "地址:" + r.Address + "\n"
	// } else if r.Customer != "" {

	if r.Customer != "" {
		return fmt.Sprintf("客户:%v\n", r.Customer)
	} else {
		return ""
	}
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
