package resolves

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	. "github.com/wanliu/flow/builtin/luis"
	. "github.com/wanliu/flow/context"
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
	Products   ItemsResolve
	// Products   ProductsResolve
	// Address    AddressResolve
	// Time       OrderTimeResolve
	Address string
	Time    time.Time
	Current Resolve
}

func NewOpenOrderResolve(ctx Context) *OpenOrderResolve {
	resolve := new(OpenOrderResolve)

	luis := ctx.Value("Result").(ResultParams)

	luis.Entities = DistinctEntites(luis.Entities)
	luis.Entities = DeduplicateEntities(luis.Entities)

	resolve.LuisParams = luis
	resolve.ExtractFromLuis()

	return resolve
}

func (t *OpenOrderResolve) Solve(luis ResultParams) (bool, string, string) {
	solved, finishNotition, nextNotition := t.Current.Solve(luis)

	if solved {
		if t.Fullfilled() {
			return true, finishNotition + "\n" + t.Answer(), ""
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
			pr := t.Products.Products[i]
			pr.Quantity = q
		}
	}

	for _, p := range t.Products.Products {
		p.CheckResolved()
	}
}

func (t *OpenOrderResolve) ExtractQuantity() []int {
	result := make([]int, 0, 10)

	for _, item := range t.LuisParams.Entities {
		if item.Type == "builtin.number" {
			number := strings.Trim(item.Resolution.Value, " ")
			q, _ := strconv.ParseInt(number, 10, 64)
			result = append(result, int(q))
		}
	}

	return result
}

func (t *OpenOrderResolve) ExtractProducts() {
	for _, item := range t.LuisParams.Entities {
		if item.Type == "products" {
			product := ItemResolve{
				Resolved:   false,
				Name:       item.Entity,
				Price:      0,
				Quantity:   0, // 默认值
				Product:    "",
				Resolution: item.Resolution,
			}

			t.Products.Add(product)
		} else {
			log.Printf("type: %v", item.Type)
		}
	}
}

func (t *OpenOrderResolve) ExtractAddress() {
	for _, item := range t.LuisParams.Entities {

		if item.Type == "address" || item.Type == "地点" {
			log.Printf("Address ... %v", item.Entity)
			t.Address = item.Entity
		}
	}
}

func (t *OpenOrderResolve) ExtractTime() {
	for _, item := range t.LuisParams.Entities {
		if item.Type == "builtin.datetime.date" {
			luisTime, err := time.Parse("2006-01-02", item.Resolution.Date)

			if err != nil {
				log.Printf("::::::::ERROR: %v", err)
			} else {
				// dTime = luisTime
				t.Time = luisTime
			}
		}
	}
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

func (t OpenOrderResolve) Answer() string {
	desc := ""

	desc = desc + "=== 订单输入完成 ===\n"
	desc = desc + "本订单包含如下商品：" + "\n"

	params := url.Values{
		"auth_token":   {"5f567b5efc3e4d0aa0d9c40922ae07aa"},
		"street":       {t.Address},
		"deliver_time": {t.Time.String()},
	}

	for i, p := range t.Products.Products {
		desc = desc + p.Product + " " + strconv.Itoa(p.Quantity) + "件\n"

		nk := "items[" + strconv.Itoa(i) + "][name]"
		nv := p.Product
		params.Add(nk, nv)

		qk := "items[" + strconv.Itoa(i) + "][quantity]"
		qv := strconv.Itoa(p.Quantity)
		params.Add(qk, qv)
	}

	desc = desc + "地址:" + t.Address + "\n"
	desc = desc + "送货时间" + t.Time.String() + "\n"
	desc = desc + "=== 结束 ===\n"

	res, err := http.PostForm("http://192.168.0.155:3000/api/v1/temp_orders", params)

	if err != nil {
		// return err.Error()
		return "服务暂时不可用，请稍后再试"
	} else {
		defer res.Body.Close()
		body, e := ioutil.ReadAll(res.Body)

		if e != nil {
			return e.Error()
		} else {
			var result Res
			json.Unmarshal(body, &result)

			if res.StatusCode == 422 {
				return result.Error
			} else {
				return desc + "请通过以下地址完成订单操作：" + result.Confirm_path
			}

		}
	}
}

type Res struct {
	Id           int
	Confirm_path string
	Error        string
}
