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

func (isr *ItemsResolve) Add(pr ItemResolve) {
	isr.Products = append(isr.Products, &pr)
}
