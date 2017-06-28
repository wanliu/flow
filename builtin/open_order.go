package builtin

import (
	. "github.com/wanliu/flow/context"
	// goflow "github.com/wanliu/goflow"
	// "fmt"
	// "errors"
	"log"
	"strconv"
	"strings"
	"time"
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
	// Address    AddressResolve
	// Time       OrderTimeResolve
	Address string
	Time    time.Time
	Current Resolve
}

func NewOpenOrderResolve(ctx Context) *OpenOrderResolve {
	resolve := new(OpenOrderResolve)

	resolve.LuisParams = ctx.Value("Result").(ResultParams)
	resolve.ExtractFromLuis()

	return resolve
}

func (t *OpenOrderResolve) Solve(luis ResultParams) (bool, string, string) {
	solved, finishNotition, nextNotition := t.Current.Solve(luis)

	if solved {
		if t.Fullfilled() {
			return true, finishNotition + "\n" + t.String(), ""
		} else {
			t.Current = t.Next()
			hint := t.Current.Hint()

			return false, finishNotition, finishNotition + "\n" + hint
		}
	} else {
		return solved, finishNotition, nextNotition
	}

}

func (t OpenOrderResolve) Hint() string {
	return t.Current.Hint()
}

// 从ｌｕｉｓ数据构造结构数据
func (t *OpenOrderResolve) ExtractFromLuis() {
	// log.Printf("====:: %v", t.LuisParams.Entities)

	// t.ExtractProducts()
	t.ExtractItems()
	t.ExtractAddress()
	t.ExtractTime()
	// t.ExtractQuantity()

	// log.Printf("----> %v", t.Products)
}

// TODO 无法识别全角数字
func (t *OpenOrderResolve) ExtractItems() {
	t.ExtractProducts()
	quantities := t.ExtractQuantity()

	for i, q := range quantities {
		if len(t.Products.Products) >= i+1 {
			t.Products.Products[i].Quantity = q
		}
	}
}

func (t *OpenOrderResolve) ExtractQuantity() []int {
	result := make([]int, 0, 10)

	for _, item := range t.LuisParams.Entities {
		if item.Type == "builtin.number" {
			number := strings.Trim(item.Entity, " ")
			q, _ := strconv.ParseInt(number, 10, 64)
			result = append(result, int(q))
		}
	}

	return result
}

func (t *OpenOrderResolve) ExtractProducts() {
	for _, item := range t.LuisParams.Entities {
		if item.Type == "products" {
			product := ProductResolve{
				Resolved:   false,
				Name:       item.Entity,
				Price:      0,
				Quantity:   1, // 默认值
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
	// return t.Time.Fullfilled()
	return !t.Time.IsZero()
}

func (t OpenOrderResolve) AddressFullfilled() bool {
	// return t.Address.Fullfilled()
	return t.Address != ""
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
		t.Current = unsolved
		return unsolved
	} else if !t.AddressFullfilled() {
		unsolved := AddressResolve{Parent: t}
		t.Current = unsolved
		return unsolved
	} else if !t.TimeFullfilled() {
		unsolved := OrderTimeResolve{Parent: t}
		t.Current = unsolved
		return unsolved
	} else {
		return nil
	}
}

func (t OpenOrderResolve) EmptyProducts() bool {
	return len(t.Products.Products) == 0
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

func (t OpenOrderResolve) String() string {
	result := ""

	result = result + "=== 订单输入完成 ===\n"
	result = result + "本订单包含如下商品：" + "\n"

	for _, p := range t.Products.Products {
		result = result + p.Product + " " + strconv.Itoa(p.Quantity) + "件\n"
	}

	result = result + "地址:" + t.Address + "\n"
	result = result + "送货时间" + t.Time.String() + "\n"
	result = result + "=== 结束 ===\n"

	return result
}
