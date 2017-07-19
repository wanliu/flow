package resolves

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	. "github.com/wanliu/flow/builtin/luis"
	. "github.com/wanliu/flow/context"
)

func NewStockQueryResolve(ctx Context) *StockQueryResolve {
	resolve := new(StockQueryResolve)

	luis := ctx.Value("Result").(ResultParams)

	luis.Entities = DistinctEntites(luis.Entities)
	luis.Entities = DeduplicateEntities(luis.Entities)
	SortEntities(luis.Entities)

	resolve.LuisParams = luis
	resolve.ExtractFromLuis()

	return resolve
}

type StockQueryResolve struct {
	LuisParams ResultParams
	Products   []*StockProductResolve
	Current    *StockProductResolve
}

// TODO 无法识别全角数字
func (r *StockQueryResolve) ExtractFromLuis() {
	r.ExtractProducts()
	quantities := r.ExtractQuantity()

	for i, q := range quantities {
		if len(r.Products) >= i+1 {
			pr := r.Products[i]
			pr.Quantity = q
		}
	}

	for _, p := range r.Products {
		p.CheckResolved()
	}
}

func (r *StockQueryResolve) ExtractQuantity() []int {
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

func (r *StockQueryResolve) ExtractProducts() {
	for _, item := range r.LuisParams.Entities {
		if item.Type == "products" {
			product := &StockProductResolve{
				Resolved:   false,
				Name:       item.Entity,
				Quantity:   0, // 默认值
				Stock:      0,
				Resolution: item.Resolution,
			}

			r.Products = append(r.Products, product)
		} else {
			log.Printf("type: %v", item.Type)
		}
	}
}

func (r *StockQueryResolve) Next() Resolve {
	for _, pr := range r.Products {
		if !pr.Resolved {
			r.Current = pr
			return pr
		}
	}

	return new(StockProductResolve)
}

func (r *StockQueryResolve) Solve(luis ResultParams) (bool, string, string) {

	solved, finishedNotition, nextNotition := r.Current.Solve(luis)

	if solved {
		if r.Fullfilled() {
			finishedNotition = finishedNotition + "\n" + r.Answer()
			return true, finishedNotition, nextNotition
		} else {
			next := finishedNotition + "\n" + r.Next().Hint()

			return false, finishedNotition, next
		}
	} else {
		return solved, finishedNotition, nextNotition
	}
}

func (r StockQueryResolve) Fullfilled() bool {
	for _, p := range r.Products {
		if !p.Resolved {
			return false
		}
	}

	return true
}

func (r StockQueryResolve) Answer() string {
	selected := make([]string, 0, 10)

	for _, p := range r.Products {
		selected = append(selected, "name[]="+p.Product)
	}

	query := "?auth_token=5f567b5efc3e4d0aa0d9c40922ae07aa&" + strings.Join(selected, "&")

	res, err := http.Get("http://192.168.0.155:3000/api/v1/query_items/stock" + query)

	if err != nil {
		// return err.Error()
		return "服务暂时不可用，请稍后再试"
	} else {
		defer res.Body.Close()
		body, e := ioutil.ReadAll(res.Body)

		if e != nil {
			return e.Error()
		} else {
			var result StockRes
			json.Unmarshal(body, &result)

			if res.StatusCode == 422 {
				return result.Error
			} else {
				result.Compare(r)
				return result.String()
			}
		}
	}
}

func (r StockQueryResolve) EmptyProducts() bool {
	return len(r.Products) == 0
}

type StockRes struct {
	Items []*ItemStockRes
	Error string
}

func (s *StockRes) Compare(r StockQueryResolve) {
	for _, i := range s.Items {
		for _, p := range r.Products {
			if p.Product == i.Name && p.Quantity != 0 {
				i.Quantity = p.Quantity
				break
			}
		}
	}
}

func (s StockRes) String() string {
	result := make([]string, 0, 10)

	for _, i := range s.Items {
		result = append(result, i.String())
	}

	return strings.Join(result, ",")
}

type ItemStockRes struct {
	Name          string
	Current_stock int
	Quantity      int
}

func (i ItemStockRes) String() string {
	if i.Quantity == 0 {
		if i.Current_stock <= 0 {
			return i.Name + "没有货"
		} else if i.Current_stock <= 20 {
			return i.Name + "的库存不多"
		} else {
			return i.Name + "有货"
		}
	} else {
		if i.Current_stock >= i.Quantity {
			return i.Name + "有超过" + strconv.Itoa(i.Quantity) + "的货可以出售货"
		} else {
			return i.Name + "的库存不足" + strconv.Itoa(i.Quantity)
		}
	}
}
