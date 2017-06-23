package builtin

import (
	. "github.com/wanliu/flow/context"
	// goflow "github.com/wanliu/goflow"
	// "fmt"
	"log"
	// "strings"
)

// 处理开单的逻辑结构, 不需要是组件
// 作为context的一个部分，或者存在一个Value中
type OpenOrderResolve struct {
	// goflow.Component
	// Ctx        <-chan Context
	// Out        chan<- ReplyData
	// Address    string
	// Time       time.Time
	// Products   []ProductResolve
	LuisParams ResultParams
	Products   ProductsResolve
	Address    AddressResolve
	Time       OrderTimeResolve
	current    Resolve
}

func NewOpenOrderResolve(ctx Context) *OpenOrderResolve {
	resolve := new(OpenOrderResolve)

	resolve.LuisParams = ctx.Value("Result").(ResultParams)
	resolve.ExtractFromLuis()

	return resolve
}

// 从ｌｕｉｓ数据构造结构数据
func (t *OpenOrderResolve) ExtractFromLuis() {
	// log.Printf("====:: %v", t.LuisParams.Entities)

	t.ExtractProducts()
	t.ExtractAddress()
	t.ExtractTime()
	// t.ExtractQuantity()

	// log.Printf("----> %v", t.Products)
}

func (t *OpenOrderResolve) ExtractProducts() {
	for _, item := range t.LuisParams.Entities {
		if item.Type == "products" {
			product := ProductResolve{
				Resolved:   false,
				Name:       item.Entity,
				Price:      0,
				Number:     1, // 默认值
				Product:    "",
				Resolution: item.Resolution,
			}

			product.CheckResolved()

			t.Products.add(product)
		} else {
			log.Printf("type: %v", item.Type)
		}
	}
}

func (t *OpenOrderResolve) ExtractAddress() {

}

func (t *OpenOrderResolve) ExtractTime() {

}

// func (t *OpenOrderResolve) ExtractQuantity() {
// 	// type: builtin.number
// }

// 去除重复产品
func (t *OpenOrderResolve) UniqProducts() {

}

func (t OpenOrderResolve) ProductsFullfilled() bool {
	return t.Products.Fullfilled()
}

func (t OpenOrderResolve) TimeFullfilled() bool {
	return t.Time.Fullfilled()
}

func (t OpenOrderResolve) AddressFullfilled() bool {
	return t.Address.Fullfilled()
}

// 是否条件全部满足
func (t OpenOrderResolve) Fullfilled() bool {
	return t.ProductsFullfilled() &&
		t.TimeFullfilled() &&
		t.AddressFullfilled()
}

// 下一个为满足项目
func (t *OpenOrderResolve) Next() Resolve {
	if !t.ProductsFullfilled() {
		unsolved := t.NextProduct()
		t.current = unsolved
		return unsolved
	} else if !t.AddressFullfilled() {
		unsolved := new(AddressResolve)
		t.current = unsolved
		return unsolved
	} else if !t.TimeFullfilled() {
		unsolved := new(OrderTimeResolve)
		t.current = unsolved
		return unsolved
	} else {
		return nil
	}
}

func (t OpenOrderResolve) NextNotify() string {
	unsolved := t.Next()
	return unsolved.Hint()
}

//
func (t OpenOrderResolve) PostService() string {
	return ""
}

func (t OpenOrderResolve) NextProduct() Resolve {
	return t.Products.NextProduct()
}
