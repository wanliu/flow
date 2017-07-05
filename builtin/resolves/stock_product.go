// 问货/问价产品选择
package resolves

import (
	"strconv"
	"strings"

	. "github.com/wanliu/flow/builtin/luis"
)

type StockProductResolve struct {
	Resolved   bool
	Name       string
	Price      float64
	Stock      int
	Product    string
	Resolution Resolution
	Parent     *StockQueryResolve
}

func (r StockProductResolve) Solve(luis ResultParams) (bool, string, string) {
	if luis.TopScoringIntent.Intent == "选择" {
		// TODO 无法识别全角数字
		number := strings.Trim(luis.Entities[0].Resolution.Value, " ")
		chose, _ := strconv.ParseInt(number, 10, 64)
		inNum := int(chose)

		for _, product := range r.Parent.Products {
			if product.Name == r.Name {
				if product.Product == "" {
					if len(product.Resolution.Values) >= inNum {
						prdName := product.Resolution.Values[chose-1]
						product.Product = prdName
						product.CheckResolved()

						return true, "已选择" + prdName, "err"
					} else {
						return false, "", "超出选择范围\n" + product.Hint()
					}
				}
			}
		}

		return false, "", "错误的操作，没有可供选择的商品"
	} else {
		return false, "", "无效的输入\n" + r.Hint()
	}
}

func (r StockProductResolve) Hint() string {
	result := ""

	if r.Product == "" && len(r.Resolution.Values) > 0 {
		index := 1
		choses := "\n"

		for _, value := range r.Resolution.Values {
			choses = choses + strconv.Itoa(index) + ": " + value + "\n"
			index = index + 1
		}

		choses = choses + "\n"

		result = "我们有下列的 " + r.Name + " 产品:" + choses + "请输入序号选择你要查询的商品"
	}

	return result
}

func (r *StockProductResolve) CheckResolved() {
	if len(r.Resolution.Values) == 0 {
		r.Product = r.Name
	}

	if r.Product != "" {
		r.Resolved = true
	}
}
