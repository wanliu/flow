package builtin

import (
	"fmt"
	"strings"

	. "github.com/wanliu/flow/context"
)

type Product struct {
	Name   string
	Price  float64
	Number int
}

type ReplyData struct {
	Reply string
	Ctx   Context
}

type TryGetProducts struct {
	TryGetEntities
	Ctx  <-chan Context
	Type <-chan string
	Out  chan<- ReplyData
}

func NewTryGetProducts() interface{} {
	return new(TryGetProducts)
}

type ProInfo []Product

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
			// log.Printf("找到 %d 产品 (%s)", len(products), ProInfo(products))
			// tr.No <- ctx
			reply := fmt.Sprintf("找到 %d 产品 (%s)", len(products), ProInfo(products))
			replyData := ReplyData{reply, ctx}
			tr.Out <- replyData
		} else {
			// tr.No <- ctx
			replyData := ReplyData{"没有相关的产品", ctx}
			tr.Out <- replyData
		}
	} else {
		// tr.No <- ctx
		replyData := ReplyData{"出现错误，请稍后重试", ctx}
		tr.Out <- replyData
	}
}

func (info ProInfo) String() string {
	var out = make([]string, 0, len(info))
	for _, p := range info {
		out = append(out, p.Name)
	}
	return strings.Join(out, ",")
}
