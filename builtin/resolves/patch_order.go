package resolves

import (
	// "log"
	"strconv"
	"strings"
	"time"

	"github.com/hysios/apiai-go"
	"github.com/wanliu/flow/builtin/ai"
	. "github.com/wanliu/flow/context"
)

type PatchOrderResolve struct {
	OrderResolve
	Origin          *OrderResolve
	OriginUpdatedAt time.Time
}

func NewPatchOrderResolve(ctx Context) *PatchOrderResolve {
	resolve := new(PatchOrderResolve)

	aiResult := ctx.Value("Result").(apiai.Result)

	resolve.AiParams = ai.ApiAiOrder{AiResult: aiResult}
	resolve.ExtractFromParams()

	return resolve
}

func (r *PatchOrderResolve) Patch(orderResolve *OrderResolve) {
	r.Origin = orderResolve

	for _, p := range r.Products.Products {
		match := false

		for _, pIn := range orderResolve.Products.Products {
			if p.Product == pIn.Product {
				pIn.Quantity = pIn.Quantity + p.Quantity
				match = true
				break
			}
		}

		if !match {
			orderResolve.Products.Products = append(orderResolve.Products.Products, p)
		}
	}

	if len(r.Gifts.Products) > 0 {
		for _, g := range r.Gifts.Products {
			match := false

			for _, gIn := range orderResolve.Gifts.Products {
				if g.Product == gIn.Product {
					gIn.Quantity = gIn.Quantity + g.Quantity
					break
				}
			}

			if !match {
				orderResolve.Gifts.Products = append(orderResolve.Gifts.Products, g)
			}
		}
	}

	if r.Address != "" && r.Origin.Address == "" {
		r.Origin.Address = r.Address
	}

	r.OriginUpdatedAt = r.Origin.UpdatedAt
	r.Origin.UpdatedAt = time.Now()
}

// 新增 2 种产品, 《伊利畅轻450原味》 已 10 件,
// 《伊利燕麦有机 》已 19 件
func (r PatchOrderResolve) Answer() string {
	if nil == r.Origin {
		return "ERROR"
	}

	if r.PatchInShortMinute() {
		return r.ShortAnswer()
	} else {
		return r.LongAnswer()
	}
}

func (r PatchOrderResolve) PatchInShortMinute() bool {
	return r.OriginUpdatedAt.Add(time.Duration(2)*time.Minute).UnixNano() > time.Now().UnixNano()
}

func (r PatchOrderResolve) ShortAnswer() string {
	desc := "新增" + CnNum(len(r.Products.Products)) + "种产品"

	if len(r.Gifts.Products) > 0 {
		desc = desc + ", " + CnNum(len(r.Gifts.Products)) + "种赠品" + "\n"
	} else {
		desc = desc + "\n"
	}

	ps := []string{}
	for _, p := range r.Products.Products {
		for _, pIn := range r.Origin.Products.Products {
			if p.Product == pIn.Product {
				ps = append(ps, pIn.Product+"已"+strconv.Itoa(pIn.Quantity)+"件")
				break
			}
		}
	}
	desc = desc + strings.Join(ps, ", ") + "\n"

	gs := []string{}
	if len(r.Gifts.Products) > 0 {
		desc = desc + "-------赠品-------\n"

		for _, g := range r.Gifts.Products {
			for _, gIn := range r.Origin.Gifts.Products {
				if g.Product == gIn.Product {
					gs = append(gs, g.Product+"已"+strconv.Itoa(g.Quantity)+"件")
					break
				}
			}
		}
	}
	desc = desc + strings.Join(gs, ", ") + "\n"

	return desc
}

func (r PatchOrderResolve) LongAnswer() string {
	desc := r.ShortAnswer()

	desc = desc + r.Origin.AnswerBody()

	return desc
}
