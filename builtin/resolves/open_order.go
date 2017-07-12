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
	DefTime string
	Current Resolve
}

func NewOpenOrderResolve(ctx Context) *OpenOrderResolve {
	resolve := new(OpenOrderResolve)

	luis := ctx.Value("Result").(ResultParams)

	luis.Entities = DistinctEntites(luis.Entities)
	luis.Entities = DeduplicateEntities(luis.Entities)
	SortEntities(luis.Entities)

	resolve.LuisParams = luis
	resolve.ExtractFromLuis()

	return resolve
}

func (r *OpenOrderResolve) Solve(luis ResultParams) (bool, string, string) {
	solved, finishNotition, nextNotition := r.Current.Solve(luis)

	if solved {
		if r.Fullfilled() {
			return true, finishNotition + "\n" + r.Answer(), ""
		} else {
			r.Current = r.Next()
			hint := r.Current.Hint()

			return false, finishNotition, finishNotition + "\n" + hint
		}
	} else {
		return solved, finishNotition, nextNotition
	}

}

func (r OpenOrderResolve) Hint() string {
	return r.Current.Hint()
}

// 从ｌｕｉｓ数据构造结构数据
func (r *OpenOrderResolve) ExtractFromLuis() {
	// log.Printf("====:: %v", r.LuisParams.Entities)

	// r.ExtractProducts()
	r.ExtractItems()
	r.ExtractAddress()
	r.ExtractTime()
	// r.ExtractQuantity()

	// log.Printf("----> %v", r.Products)
}

// TODO 无法识别全角数字
func (r *OpenOrderResolve) ExtractItems() {
	r.ExtractProducts()
	quantities := r.ExtractQuantity()

	for i, q := range quantities {
		if len(r.Products.Products) >= i+1 {
			pr := r.Products.Products[i]
			pr.Quantity = q
		}
	}

	for _, p := range r.Products.Products {
		p.CheckResolved()
	}
}

func (r *OpenOrderResolve) ExtractQuantity() []int {
	result := make([]int, 0, 10)

	for _, item := range r.LuisParams.Entities {
		if item.Type == "builtin.number" {
			number := strings.Trim(item.Resolution.Value, " ")
			q, _ := strconv.ParseInt(number, 10, 64)
			result = append(result, int(q))
		}
	}

	return result
}

func (r *OpenOrderResolve) ExtractProducts() {
	for _, item := range r.LuisParams.Entities {
		if item.Type == "products" {
			product := ItemResolve{
				Resolved:   false,
				Name:       item.Entity,
				Price:      0,
				Quantity:   0, // 默认值
				Product:    "",
				Resolution: item.Resolution,
			}

			r.Products.Add(product)
		} else {
			log.Printf("type: %v", item.Type)
		}
	}
}

func (r *OpenOrderResolve) ExtractAddress() {
	for _, item := range r.LuisParams.Entities {

		if item.Type == "地址" {
			r.Address = item.Entity
		}
	}
}

func (r *OpenOrderResolve) ExtractTime() {
	for _, item := range r.LuisParams.Entities {
		if item.Type == "builtin.datetime.date" {
			luisTime, err := time.Parse("2006-01-02", item.Resolution.Date)

			if err != nil {
				log.Printf("::::::::ERROR: %v", err)
			} else {
				// dTime = luisTime
				r.Time = luisTime
			}
		}
	}
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

func (r OpenOrderResolve) ProductsFullfilled() bool {
	return r.Products.Fullfilled()
}

func (r OpenOrderResolve) TimeFullfilled() bool {
	// return r.Time.Fullfilled()
	return !r.Time.IsZero()
}

func (r OpenOrderResolve) AddressFullfilled() bool {
	// return r.Address.Fullfilled()
	return r.Address != ""
}

// 是否条件全部满足
func (r OpenOrderResolve) Fullfilled() bool {
	return r.ProductsFullfilled() &&
		r.TimeFullfilled() &&
		r.AddressFullfilled()
}

// 下一个为满足项目
func (r *OpenOrderResolve) Next() Resolve {
	if !r.ProductsFullfilled() {
		unsolved := r.NextProduct()
		r.Current = unsolved
		return unsolved
	} else if !r.AddressFullfilled() {
		unsolved := &AddressResolve{Parent: r}
		r.Current = unsolved
		return unsolved
	} else if !r.TimeFullfilled() {
		unsolved := &OrderTimeResolve{Parent: r}
		r.Current = unsolved
		return unsolved
	} else {
		return nil
	}
}

func (r OpenOrderResolve) EmptyProducts() bool {
	return len(r.Products.Products) == 0
}

func (r OpenOrderResolve) NextNotify() string {
	unsolved := r.Next()
	return unsolved.Hint()
}

//
func (r OpenOrderResolve) PostService() string {
	return ""
}

func (r OpenOrderResolve) NextProduct() Resolve {
	return r.Products.NextProduct()
}

func (r OpenOrderResolve) Answer() string {
	desc := ""

	desc = desc + "=== 订单输入完成 ===\n"
	desc = desc + "本订单包含如下商品：" + "\n"

	params := url.Values{
		"auth_token":   {"5f567b5efc3e4d0aa0d9c40922ae07aa"},
		"street":       {r.Address},
		"deliver_time": {r.Time.Format("2006年01月02日")},
	}

	for i, p := range r.Products.Products {
		desc = desc + p.Product + " " + strconv.Itoa(p.Quantity) + "件\n"

		nk := "items[" + strconv.Itoa(i) + "][name]"
		nv := p.Product
		params.Add(nk, nv)

		qk := "items[" + strconv.Itoa(i) + "][quantity]"
		qv := strconv.Itoa(p.Quantity)
		params.Add(qk, qv)
	}

	desc = desc + "地址:" + r.Address + "\n"
	desc = desc + "送货时间" + r.Time.Format("2006年01月02日") + "\n"
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
