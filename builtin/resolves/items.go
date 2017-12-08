package resolves

import ()

type ItemsResolve struct {
	Products []*ItemResolve
	Current  *ItemResolve
}

func (r *ItemsResolve) Add(pr ItemResolve) {
	r.Products = append(r.Products, &pr)
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
