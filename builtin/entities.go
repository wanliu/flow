package builtin

import (
	"log"
	"strings"

	. "github.com/wanliu/flow/context"
	flow "github.com/wanliu/goflow"
)

type Product struct {
	Name   string
	Price  float64
	Number int
}

type TryGetEntities struct {
	flow.Component
	_type string
	Ctx   <-chan Context
	Next  chan<- Context
	Type  <-chan string
	No    chan<- Context
}

type TryGetProducts struct {
	TryGetEntities
	Ctx  <-chan Context
	Type <-chan string
	Out  chan<- Context
}

func NewTryGetProducts() interface{} {
	return new(TryGetProducts)
}

type ProInfo []Product

func (tr *TryGetEntities) OnType(typ string) {
	tr._type = typ
}

func (tr *TryGetEntities) OnCtx(ctx Context) {
	if res, ok := ctx.GlobalValue("Result").(ResultParams); ok {

		for _, entity := range res.Entities {
			if entity.Type == tr._type {

			}
		}
	} else {
		tr.No <- ctx
	}
}

func (tr *TryGetProducts) OnCtx(ctx Context) {
	if res, ok := ctx.Value("Result").(ResultParams); ok {
		var products = make([]Product, 0)
		for _, entity := range res.Entities {
			if entity.Type == tr._type {
				products = append(products, Product{
					Name: entity.Resolution.Values[0],
				})
			}
		}

		if len(products) > 0 {
			ctx.SetGlobalValue("products", &products)
			log.Printf("找到 %d 产品 (%s)", len(products), ProInfo(products))
			tr.Out <- ctx

		} else {
			tr.No <- ctx
		}
	} else {
		tr.No <- ctx
	}
}

func (info ProInfo) String() string {
	var out = make([]string, 0, len(info))
	for _, p := range info {
		out = append(out, p.Name)
	}
	return strings.Join(out, ",")
}
