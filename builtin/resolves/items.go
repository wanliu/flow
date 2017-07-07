package resolves

import (
	"strconv"
	"strings"

	. "github.com/wanliu/flow/builtin/luis"
)

type ItemsResolve struct {
	Products []*ItemResolve
	Current  *ItemResolve
}

func (isr ItemsResolve) Hint() string {
	return isr.Current.Hint()
}

func (isr ItemsResolve) Solve(luis ResultParams) (bool, string, string) {
	solved, finishNotition, nextNotition := isr.Current.Solve(luis)
	if solved {
		if isr.Fullfilled() {
			// selected = "您已经选择了:"
			selected := make([]string, 10)

			for _, resolve := range isr.Products {
				selected = append(selected, resolve.Product)
			}

			notition := "您已经选择了 " + strings.Join(selected, ", ") + "等" + strconv.Itoa(len(selected)) + "件商品"

			return solved, notition, ""
		} else {
			solve := isr.NextProduct()

			hint := solve.Hint()
			return false, finishNotition, hint
		}
	} else {
		return solved, finishNotition, nextNotition
	}

}

func (isr *ItemsResolve) Add(pr ItemResolve) {
	isr.Products = append(isr.Products, &pr)
}

func (isr *ItemsResolve) NextProduct() Resolve {
	for _, pr := range isr.Products {
		if !pr.Resolved {
			isr.Current = pr
			return pr
		}
	}

	return new(ItemResolve)
}

func (isr ItemsResolve) Fullfilled() bool {
	if len(isr.Products) == 0 {
		return false
	}

	for _, product := range isr.Products {
		if !product.Resolved {
			return false
		}
	}

	return true
}
