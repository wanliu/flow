package resolves

import (
// "strconv"
// "strings"

// "github.com/hysios/apiai-go"
)

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
