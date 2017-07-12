package resolves

import (
	"log"
	"strconv"
	"strings"

	. "github.com/wanliu/flow/builtin/luis"
)

type ItemResolve struct {
	// Original_string string
	Resolved   bool
	Name       string
	Price      float64
	Quantity   int
	Product    string
	Resolution Resolution
}

func (r ItemResolve) Hint() string {
	// choses := "\n" + strings.Join(r.Resolution.Values, "\n")

	var result string

	if r.Product == "" && len(r.Resolution.Values) > 0 {
		index := 1
		choses := "\n"

		for _, value := range r.Resolution.Values {
			choses = choses + strconv.Itoa(index) + ": " + value + "\n"
			index = index + 1
		}

		choses = choses + "\n"

		result = "我们有下列的 " + r.Name + " 产品:" + choses + "请输入序号选择"
	} else if r.Quantity == 0 {
		result = "请告诉我您要购买" + r.Product + "的数量\n"
	}

	return result

}

func (r *ItemResolve) Solve(luis ResultParams) (bool, string, string) {
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
		} else if r.Quantity == 0 {
			if 0 < inNum {
				r.Quantity = inNum
				r.CheckResolved()

				return true, "购买的数量为：" + strconv.Itoa(inNum), ""
			} else {
				return false, "", "购买的数量必须大于零, 请重新输入\n"
			}
		}

		return false, "", "错误的操作，没有可供选择的商品"
	} else {
		log.Printf("luis: %v", luis)
		return false, "", "无效的输入: \"" + luis.Query + "\"。\n" + r.Hint()
	}
}

func (r *ItemResolve) CheckResolved() {
	if len(r.Resolution.Values) == 0 {
		r.Product = r.Name
	} else if len(r.Resolution.Values) == 1 {
		r.Product = r.Resolution.Values[0]
	}

	if r.Product != "" && r.Quantity > 0 {
		r.Resolved = true
	}
}
