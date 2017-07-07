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

func NewPriceQueryResolve(ctx Context) *PriceQueryResolve {
	resolve := new(PriceQueryResolve)

	luis := ctx.Value("Result").(ResultParams)

	log.Printf("BEFORE: %v", luis)

	luis.Entities = DistinctEntites(luis.Entities)
	luis.Entities = DeduplicateEntities(luis.Entities)

	log.Printf("AFTER: %v", luis)

	resolve.LuisParams = luis
	resolve.ExtractFromLuis()

	return resolve
}

type PriceQueryResolve struct {
	LuisParams ResultParams
	Products   []*PriceProductResolve
	Current    *PriceProductResolve
}

func (r *PriceQueryResolve) ExtractFromLuis() {
	for _, item := range r.LuisParams.Entities {
		if item.Type == "products" {
			product := PriceProductResolve{
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

func (r *PriceQueryResolve) Next() Resolve {
	for _, pr := range r.Products {
		if !pr.Resolved {
			r.Current = pr
			return pr
		}
	}

	return new(PriceProductResolve)
}

func (r *PriceQueryResolve) Solve(luis ResultParams) (bool, string, string) {

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

func (r PriceQueryResolve) Fullfilled() bool {
	for _, p := range r.Products {
		if !p.Resolved {
			return false
		}
	}

	return true
}

func (r PriceQueryResolve) Answer() string {
	selected := make([]string, 0, 10)

	// TODO 查询后台商品价格
	rand.Seed(time.Now().UTC().UnixNano())

	for _, p := range r.Products {
		p.Price = rand.Float64()
		notition := p.Product + "的价格为" + strconv.FormatFloat(p.Price, 'f', 2, 64) + "元"
		selected = append(selected, notition)
	}

	return strings.Join(selected, ", ")
}

func (r PriceQueryResolve) EmptyProducts() bool {
	return len(r.Products) == 0
}
