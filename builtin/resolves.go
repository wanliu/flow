package builtin

import (
	// . "github.com/wanliu/flow/context"
	// goflow "github.com/wanliu/goflow"
	// "fmt"
	_ "errors"
	_ "log"
	"strings"
	"time"
)

type Resolve interface {
	Hint() string
	Solve(ResultParams) (bool, string, string) // 是否全部完成，完成提示，下一步动作提醒
}

type ProductsResolve struct {
	Products []*ProductResolve
	Current  *ProductResolve
}

func (psr ProductsResolve) Hint() string {
	return psr.Current.Hint()
}

func (psr ProductsResolve) Solve(luis ResultParams) (bool, string, string) {
	solved, finishNotition, nextNotition := psr.Current.Solve(luis)
	if solved {
		if psr.Fullfilled() {
			// selected = "您已经选择了:"
			selected := make([]string, 10)

			for _, resolve := range psr.Products {
				selected = append(selected, resolve.Product)
			}

			notition := "您已经选择了 " + strings.Join(selected, ", ") + "等" + string(len(selected)) + "件商品"

			return solved, notition, ""
		} else {
			solve := psr.NextProduct()

			hint := solve.Hint()
			return false, finishNotition, hint
		}
	} else {
		return solved, finishNotition, nextNotition
	}

}

func (psr *ProductsResolve) add(pr ProductResolve) {
	pr.Parent = psr
	psr.Products = append(psr.Products, &pr)
}

func (psr *ProductsResolve) NextProduct() Resolve {
	for _, pr := range psr.Products {
		if !pr.Resolved {
			psr.Current = pr
			return pr
		}
	}

	return ProductResolve{}
}

func (psr ProductsResolve) Fullfilled() bool {
	if len(psr.Products) == 0 {
		return false
	}

	for _, product := range psr.Products {
		if !product.Resolved {
			return false
		}
	}

	return true
}

type ProductResolve struct {
	// Original_string string
	Resolved   bool
	Name       string
	Price      float64
	Number     int
	Product    string
	Resolution Resolution
	Current    Resolve
	Parent     *ProductsResolve
}

func (pr ProductResolve) Hint() string {
	// choses := "\n" + strings.Join(pr.Resolution.Values, "\n")

	var result string

	if pr.Product == "" && len(pr.Resolution.Values) > 0 {
		index := 1
		choses := "\n"

		for _, value := range pr.Resolution.Values {
			choses = choses + string(index) + ": " + value + "\n"
			index = index + 1
		}

		choses = choses + "\n"

		result = "我们有下列的 " + pr.Name + " 产品:" + choses + "请输入序号选择"
	} else if pr.Number == 0 {
		result = "请告诉我您要购买的数量\n"
	}

	return result

}

func (pr ProductResolve) Solve(luis ResultParams) (bool, string, string) {
	// log.Printf("......................SOLVE.......................... %v, %V", pr.Name, len(pr.Parent.Products))
	if luis.TopScoringIntent.Intent == "选择" {
		for _, product := range pr.Parent.Products {
			if product.Name == pr.Name {
				// log.Printf("------------------------------- FOUND ------------------------%v", product)
				product.Resolved = true
				// log.Printf("------------------------------- MODIFIED ------------------------%v", product)
			}
		}
		// pr.Resolved = true
		return true, "已选择ｘｘｘｘ商品", "err"
	} else {
		return false, "", "无效的输入\n" + pr.Hint()
	}

}

type AddressResolve struct {
	// Address string
	Parent *OpenOrderResolve
}

func (ar AddressResolve) Hint() string {
	return "请告诉我送货地址"
}

func (pr AddressResolve) Solve(luis ResultParams) (bool, string, string) {
	// pr.Address = "some where"
	if luis.TopScoringIntent.Intent == "地址" {
		pr.Parent.Address = "some where"
		return true, "已经定好了送货地址", "err"
	} else {
		return false, "", "无效的输入\n" + pr.Hint()
	}
}

// func (ar AddressResolve) Fullfilled() bool {
// 	return ar.Address != ""
// }

type OrderTimeResolve struct {
	// Time   time.Time
	Parent *OpenOrderResolve
}

func (ar OrderTimeResolve) Hint() string {
	return "请告诉我送货时间"
}

func (pr OrderTimeResolve) Solve(luis ResultParams) (bool, string, string) {
	if luis.TopScoringIntent.Intent == "时间" {
		pr.Parent.Time = time.Now()
		return true, "已经定好了送货时间", "err"
	} else {
		return false, "", "无效的输入\n" + pr.Hint()
	}
}

// func (ar OrderTimeResolve) Fullfilled() bool {
// 	return !ar.Time.IsZero()
// }

func (pr *ProductResolve) CheckResolved() {
	if len(pr.Resolution.Values) == 0 {
		pr.Product = pr.Name
	}
}
