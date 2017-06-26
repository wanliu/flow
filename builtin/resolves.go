package builtin

import (
	// . "github.com/wanliu/flow/context"
	// goflow "github.com/wanliu/goflow"
	// "fmt"
	"log"
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
	Current  ProductResolve
}

func (psr ProductsResolve) Hint() string {
	return psr.Current.Hint()
}

func (psr ProductsResolve) Solve(luis ResultParams) (bool, error) {
	solved, err := psr.Current.Solve(luis)
	if solved {
		psr.NextProduct()
	}

	return solved, err
}

func (psr *ProductsResolve) add(pr ProductResolve) {
	psr.Products = append(psr.Products, pr)
}

func (psr ProductsResolve) NextProduct() Resolve {
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
	log.Printf("......................SOLVE.......................... %v", pr.Name)
	pr.Resolved = true
	return true, errors.New("err")
}

type AddressResolve struct {
	Address string
	parent  *OpenOrderResolve
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
	Time   time.Time
	parent *OpenOrderResolve
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
