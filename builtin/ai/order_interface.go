package ai

import (
	"time"
)

type Item struct {
	Product     string
	Quantity    int
	Price       float64
	ResolveList []string
}

type AiOrder interface {
	Items() []Item
	Products() []Item
	Quantities() []Item

	GiftItems() []Item
	GiftProducts() []Item
	GiftQuantities() []Item

	Address() string
	Customer() string

	Time() time.Time

	Fulfiled() bool

	Note() string
}
