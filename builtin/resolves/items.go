package resolves

import (
	"fmt"
)

type ItemsResolve struct {
	Products []*ItemResolve
	Current  *ItemResolve
}

func (r *ItemsResolve) Add(pr ItemResolve) {
	if pr.CheckUnit() {
		for _, p := range r.Products {
			if pr.Product == p.Product {
				p.Quantity = p.Quantity + pr.Quantity
				return
			}
		}

		r.Products = append(r.Products, &pr)
	}
}

func (r *ItemsResolve) Remove(itemName string) bool {
	newProducts := []*ItemResolve{}
	included := false

	for _, item := range r.Products {
		if item.Product == itemName {
			included = true
		} else {
			newProducts = append(newProducts, item)
		}
	}

	if included {
		r.Products = newProducts
	}

	return included
}

func (r *ItemsResolve) Patch(isr ItemsResolve) {
	for _, p := range isr.Products {
		match := false

		for _, pIn := range r.Products {
			if p.Product == pIn.Product {
				pIn.Quantity = pIn.Quantity + p.Quantity
				match = true
				break
			}
		}

		if !match {
			r.Products = append(r.Products, p)
		}
	}
}

func (r ItemsResolve) MismatchQuantity() bool {
	for _, p := range r.Products {
		if p.Product == "" {
			return true
		}

		if p.Quantity == 0 {
			return true
		}
	}

	return false
}

func (r ItemsResolve) Empty() bool {
	if len(r.Products) == 0 {
		return true
	}

	for _, p := range r.Products {
		if p.Product != "" {
			return false
		}
	}

	return true
}

func (r *ItemsResolve) ChangeUint(itemName, unit string) error {
	var item ItemResolve

	for _, p := range r.Products {
		if p.Product == itemName {
			item = p
			break
		}
	}

	newItem := ItemResolve{
		Product:  item.Product,
		Quantity: item.Quantity,
		Unit:     unit,
	}

	if newItem.ValidUnit() {
		item.Unit = unit
		item.CheckUnit()
		return nil
	}

	return fmt.Errorf("%v不能以为%v单位出售", itemName, unit)
}
