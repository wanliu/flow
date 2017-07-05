package resolves

import (
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
	Parent     *ItemsResolve
	// Current    string
}

func (pr ItemResolve) Hint() string {
	// choses := "\n" + strings.Join(pr.Resolution.Values, "\n")

	var result string

	if pr.Product == "" && len(pr.Resolution.Values) > 0 {
		index := 1
		choses := "\n"

		for _, value := range pr.Resolution.Values {
			choses = choses + strconv.Itoa(index) + ": " + value + "\n"
			index = index + 1
		}

		choses = choses + "\n"

		result = "我们有下列的 " + pr.Name + " 产品:" + choses + "请输入序号选择"
	} else if pr.Quantity == 0 {
		result = "请告诉我您要购买的数量\n"
	}

	return result

}

func (pr ItemResolve) Solve(luis ResultParams) (bool, string, string) {
	if luis.TopScoringIntent.Intent == "选择" {
		// TODO 无法识别全角数字
		number := strings.Trim(luis.Entities[0].Resolution.Value, " ")
		chose, _ := strconv.ParseInt(number, 10, 64)
		inNum := int(chose)

		for _, product := range pr.Parent.Products {
			if product.Name == pr.Name {
				if product.Product == "" {
					if len(product.Resolution.Values) >= inNum {
						prdName := product.Resolution.Values[chose-1]

						product.Product = prdName
						product.CheckResolved()

						return true, "已选择" + prdName, "err"
					} else {
						return false, "", "超出选择范围\n" + product.Hint()
					}
				} else if product.Quantity == 0 {
					if 0 < inNum {
						product.Quantity = inNum
						product.CheckResolved()

						return true, "购买的数量为：" + strconv.Itoa(inNum), ""
					} else {
						return false, "", "购买的数量必须大于零, 请重新输入\n"
					}

				}
			}
		}

		return false, "", "错误的操作，没有可供选择的商品"
	} else {
		return false, "", "无效的输入\n" + pr.Hint()
	}
}

func (pr *ItemResolve) CheckResolved() {
	if len(pr.Resolution.Values) == 0 {
		pr.Product = pr.Name
	}

	if pr.Product != "" && pr.Quantity > 0 {
		pr.Resolved = true
	}
}
