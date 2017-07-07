package resolves

import (
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	. "github.com/wanliu/flow/builtin/luis"
	. "github.com/wanliu/flow/context"
)

func NewStockQueryResolve(ctx Context) *StockQueryResolve {
	resolve := new(StockQueryResolve)

	luis := ctx.Value("Result").(ResultParams)

	luis.Entities = DistinctEntites(luis.Entities)
	luis.Entities = DeduplicateEntities(luis.Entities)

	resolve.LuisParams = luis
	resolve.ExtractFromLuis()

	return resolve
}

type StockQueryResolve struct {
	LuisParams ResultParams
	Products   []*StockProductResolve
	Current    *StockProductResolve
}

func (r *StockQueryResolve) ExtractFromLuis() {
	for _, item := range r.LuisParams.Entities {
		if item.Type == "products" {
			product := StockProductResolve{
				Resolved:   false,
				Name:       item.Entity,
				Stock:      0,
				Resolution: item.Resolution,
			}

			product.CheckResolved()

			r.Products = append(r.Products, &product)
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

	// TODO 查询后台商品价格
	rand.Seed(time.Now().UTC().UnixNano())

	for _, p := range r.Products {
		p.Stock = rand.Intn(100)

		if p.Stock <= 50 {
			selected = append(selected, p.Product+"已经没货")
		} else {
			selected = append(selected, p.Product+"还有库存："+strconv.Itoa(p.Stock))
		}
	}

	return strings.Join(selected, ", ")
}

func (r StockQueryResolve) EmptyProducts() bool {
	return len(r.Products) == 0
}
