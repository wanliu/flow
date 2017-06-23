package builtin

import (
	// . "github.com/wanliu/flow/context"
	// goflow "github.com/wanliu/goflow"
	// "fmt"
	// "log"
	// "strings"
	"errors"
	"time"
)

type Resolve interface {
	Hint() string
	Solve(ResultParams) (bool, error)
}

type ProductsResolve struct {
	Products []ProductResolve
}

func (prs *ProductsResolve) add(pr ProductResolve) {
	prs.Products = append(prs.Products, pr)
}

func (prs ProductsResolve) NextProduct() Resolve {
	for _, product := range prs.Products {
		if !product.Resolved {
			return product
		}
	}

	return ProductResolve{}
}

func (prs ProductsResolve) Fullfilled() bool {
	if len(prs.Products) == 0 {
		return false
	}

	for _, product := range prs.Products {
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
	current    *Resolve
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

func (pr ProductResolve) Solve(luis ResultParams) (bool, error) {
	return true, errors.New("err")
}

type AddressResolve struct {
	Address string
}

func (ar AddressResolve) Hint() string {
	return "请告诉我送货地址"
}

func (pr AddressResolve) Solve(luis ResultParams) (bool, error) {
	return true, errors.New("err")
}

func (ar AddressResolve) Fullfilled() bool {
	return ar.Address != ""
}

type OrderTimeResolve struct {
	Time time.Time
}

func (ar OrderTimeResolve) Hint() string {
	return "请告诉我送货时间"
}

func (pr OrderTimeResolve) Solve(luis ResultParams) (bool, error) {
	return true, errors.New("err")
}

func (ar OrderTimeResolve) Fullfilled() bool {
	return !ar.Time.IsZero()
}

func (pr *ProductResolve) CheckResolved() {
	if len(pr.Resolution.Values) == 0 {
		pr.Product = pr.Name
	}
}
