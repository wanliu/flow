package ai

import (
	"time"
)

type Item struct {
	Product     string
	Quantity    int
	Price       float64
	Unit        string
	ResolveList []string
	Spec        string
	Taste       string
}

type AiOrder interface {
	Score() float64

	Query() string

	Items() []Item
	Products() []Item
	Quantities() []Item

	GiftItems() []Item
	GiftProducts() []Item
	GiftQuantities() []Item

	Address() string
	Customer() string
	Duration() string
	Count() int

	Time() time.Time

	Fulfiled() bool

	Note() string
}
