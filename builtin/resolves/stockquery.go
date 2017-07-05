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

	resolve.LuisParams = ctx.Value("Result").(ResultParams)
	resolve.ExtractFromLuis()

	return resolve
}

type StockQueryResolve struct {
	LuisParams ResultParams
	Products   []*ProductResolve
	Current    *ProductResolve
}

func (sqr *StockQueryResolve) ExtractFromLuis() {
	for _, item := range sqr.LuisParams.Entities {
		if item.Type == "products" {
			product := ProductResolve{
				Resolved:   false,
				Name:       item.Entity,
				Stock:      0,
				Resolution: item.Resolution,
			}

			product.CheckResolved()

			product.Parent = sqr
			sqr.Products = append(sqr.Products, &product)
		} else {
			log.Printf("type: %v", item.Type)
		}
	}
}

func (sqr *StockQueryResolve) Next() Resolve {
	for _, pr := range sqr.Products {
		if !pr.Resolved {
			sqr.Current = pr
			return pr
		}
	}

	return ProductResolve{}
}

func (sqr *StockQueryResolve) Solve(luis ResultParams) (bool, string, string) {

	solved, finishedNotition, nextNotition := sqr.Current.Solve(luis)

	if solved {
		if sqr.Fullfilled() {
			finishedNotition = finishedNotition + "\n" + sqr.Answer()
			return true, finishedNotition, nextNotition
		} else {
			next := finishedNotition + "\n" + sqr.Next().Hint()

			return false, finishedNotition, next
		}
	} else {
		return solved, finishedNotition, nextNotition
	}
}

func (sqr StockQueryResolve) Fullfilled() bool {
	for _, p := range sqr.Products {
		if !p.Resolved {
			return false
		}
	}

	return true
}

func (sqr StockQueryResolve) Answer() string {
	selected := make([]string, 0, 10)

	// TODO 查询后台商品价格
	rand.Seed(time.Now().UTC().UnixNano())

	for _, p := range sqr.Products {
		p.Stock = rand.Intn(100)

		if p.Stock <= 50 {
			selected = append(selected, p.Product+"已经没货")
		} else {
			selected = append(selected, p.Product+"还有库存："+strconv.Itoa(p.Stock))
		}
	}

	return strings.Join(selected, ", ")
}

func (sqr StockQueryResolve) EmptyProducts() bool {
	return len(sqr.Products) == 0
}
