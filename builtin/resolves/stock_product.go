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
}

func (r *StockProductResolve) Solve(luis ResultParams) (bool, string, string) {
	if luis.TopScoringIntent.Intent == "选择" {
		input, exist := FetchEntity("builtin.number", luis.Entities)

		if !exist {
			return false, "", "无效的输入: \"" + luis.Query + "\"。\n" + r.Hint()
		}

		number := strings.Trim(input.Resolution.Value, " ")
		chose, _ := strconv.ParseInt(number, 10, 64)
		inNum := int(chose)

		if r.Product == "" {
			if len(r.Resolution.Values) >= inNum && inNum > 0 {
				prdName := r.Resolution.Values[chose-1]
				r.Product = prdName
				r.CheckResolved()

				return true, "已选择" + prdName, "err"
			} else {
				return false, "", "超出选择范围\n" + r.Hint()
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
	} else if len(r.Resolution.Values) == 1 {
		r.Product = r.Resolution.Values[0]
	}

	if r.Product != "" {
		r.Resolved = true
	}
}
