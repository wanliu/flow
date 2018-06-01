package resolves

import (
	"github.com/wanliu/brain_data/database"
)

type ItemResolve struct {
	Resolved bool
	Name     string
	Price    float64
	Quantity int
	Product  string
	Unit     string
	// Resolution Resolution
}

func (ir *ItemResolve) CheckUnit() bool {
	item, err := database.NewOrderItem("", ir.Product, uint(ir.Quantity), ir.Unit, ir.Price)

	if err != nil {
		return false
	}

	ir.Unit = item.Unit
	ir.Price = item.Price

	return true
}
