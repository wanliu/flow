package resolves

import (
	// "log"
	"strconv"
	"strings"
	"time"

	"github.com/hysios/apiai-go"
	"github.com/wanliu/flow/builtin/ai"
	. "github.com/wanliu/flow/builtin/config"
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

	orderResolve.Products.Patch(r.Products)
	orderResolve.Gifts.Patch(r.Gifts)

	if r.Address != "" && r.Origin.Address == "" {
		r.Origin.Address = r.Address
	}

	if r.Customer != "" && r.Origin.Customer == "" {
		r.Origin.Customer = r.Customer
	}

	r.OriginUpdatedAt = r.Origin.UpdatedAt
	// r.Origin.Touch()
}

// 新增 2 种产品, 《伊利畅轻450原味》 已 10 件,
// 《伊利燕麦有机 》已 19 件
func (r PatchOrderResolve) Answer() string {
	if nil == r.Origin {
		return "ERROR"
	}

	shtMns := PatchShortMinutes

	if r.Origin.Fulfiled() {
		return r.LongAnswer()
	} else if r.WithinShortMinute(shtMns) {
		return r.ShortAnswer()
	} else {
		return r.LongAnswer()
	}
}

func (r PatchOrderResolve) WithinShortMinute(mins int) bool {
	return r.OriginUpdatedAt.Add(time.Duration(mins)*time.Minute).UnixNano() > time.Now().UnixNano()
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
		for _, g := range r.Gifts.Products {
			for _, gIn := range r.Origin.Gifts.Products {
				if g.Product == gIn.Product {
					gs = append(gs, "赠品 "+gIn.Product+" 已"+strconv.Itoa(gIn.Quantity)+"件")
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
	desc = desc + "\n-----------订单详情-------------\n"

	desc = desc + r.Origin.AnswerBody()

	return desc
}
