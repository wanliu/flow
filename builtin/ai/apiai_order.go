package ai

import (
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/hysios/apiai-go"
)

type ApiAiOrder struct {
	AiResult apiai.Result
}

func (aa ApiAiOrder) Score() float64 {
	return aa.AiResult.Score
}

func (aa ApiAiOrder) Query() string {
	return aa.AiResult.ResolvedQuery
}

func (aa ApiAiOrder) Items() []Item {
	products := aa.Products()
	quantities := aa.Quantities()

	for i, q := range quantities {
		if len(products) >= i+1 {
			products[i].Quantity = q.Quantity
			products[i].Unit = q.Unit
		}
	}

	return products
}

func (aa ApiAiOrder) Products() []Item {
	return aa.ExtractProducts("products")
}

func (aa ApiAiOrder) Quantities() []Item {
	return aa.ExtractQuantities("quantity")
}

func (aa ApiAiOrder) GiftItems() []Item {
	gifts := aa.GiftProducts()
	quantities := aa.GiftQuantities()

	for i, q := range quantities {
		if len(gifts) >= i+1 {
			gifts[i].Quantity = q.Quantity
			gifts[i].Unit = q.Unit
		}
	}

	return gifts
}

func (aa ApiAiOrder) GiftProducts() []Item {
	return aa.ExtractProducts("gifts")
}

func (aa ApiAiOrder) GiftQuantities() []Item {
	return aa.ExtractQuantities("giftNumber")
}

func (aa ApiAiOrder) Address() string {
	if a, exist := aa.AiResult.Params["street-address"]; exist {

		rt := reflect.TypeOf(a)
		vals := reflect.ValueOf(a)

		switch rt.Kind() {
		case reflect.Slice:
			if vals.Len() > 0 {
				return vals.Index(0).Interface().(string)
			}
		case reflect.Array:
			if vals.Len() > 0 {
				return vals.Index(0).Interface().(string)
			}
		case reflect.String:
			return vals.Interface().(string)
		}
	}

	return ""
}

func (aa ApiAiOrder) Customer() string {
	if c, exist := aa.AiResult.Params["customer"]; exist {

		rt := reflect.TypeOf(c)
		vals := reflect.ValueOf(c)

		switch rt.Kind() {
		case reflect.Slice:
			if vals.Len() > 0 {
				return vals.Index(0).Interface().(string)
			}
		case reflect.Array:
			if vals.Len() > 0 {
				return vals.Index(0).Interface().(string)
			}
		case reflect.String:
			return vals.Interface().(string)
		}
	}

	return ""
}

func (aa ApiAiOrder) Time() time.Time {
	if t, exist := aa.AiResult.Params["date"]; exist {
		if aiTime, err := time.Parse("2006-01-02", t.(string)); err == nil {
			return aiTime
		}
	}

	return time.Time{}
}

func (aa ApiAiOrder) Fulfiled() bool {
	return true
}

func (aa ApiAiOrder) Note() string {
	if imp, exist := aa.AiResult.Params["important"]; exist {
		return imp.(string)
	}

	return ""
}

func (aa ApiAiOrder) ExtractProducts(t string) []Item {
	result := make([]Item, 0, 50)

	if products, exist := aa.AiResult.Params[t]; exist {
		ps := reflect.ValueOf(products)

		for i := 0; i < ps.Len(); i++ {
			p := ps.Index(i)
			name := p.Interface().(string)
			item := Item{Product: name}
			result = append(result, item)
		}
	}

	return result
}

func (aa ApiAiOrder) ExtractQuantities(t string) []Item {
	result := make([]Item, 0, 50)

	if quantities, exist := aa.AiResult.Params[t]; exist {
		qs := reflect.ValueOf(quantities)

		for i := 0; i < qs.Len(); i++ {
			q := qs.Index(i).Interface()

			switch t := q.(type) {
			case string:
				qs := q.(string)
				qi := extractQuantity(qs)
				item := Item{Quantity: qi}
				result = append(result, item)
			case float64:
				qf := q.(float64)
				item := Item{Quantity: int(qf)}
				result = append(result, item)
			case map[string]interface{}:
				log.Printf("quantity: %v\n", t)
				qf := t["number"].(float64)
				quantity := int(qf)
				item := Item{Quantity: quantity}

				if unit, ok := t["unit"].(string); ok {
					unit = strings.Replace(unit, "龘", "", -1)
					unit = strings.Replace(unit, " ", "", -1)
					item.Unit = unit
				}

				result = append(result, item)
			default:
				log.Println("Unknown Quantity type: %v", t)
			}
		}
	}

	return result
}

func extractQuantity(s string) int {
	nums := strings.TrimFunc(s, TrimToNum)
	numsCgd := DBCtoSBC(nums)

	if len(numsCgd) > 0 {
		q, _ := strconv.Atoi(numsCgd)
		return q
	} else {
		return 0
	}
	// re := regexp.MustCompile("[0-9]+")
}

func TrimToNum(r rune) bool {
	if n := r - '0'; n >= 0 && n <= 9 {
		return false
	} else if m := r - '０'; m >= 0 && m <= 9 {
		return false
	}

	return true
}

func DBCtoSBC(s string) string {
	retstr := ""
	for _, i := range s {
		inside_code := i
		if inside_code == 12288 {
			inside_code = 32
		} else {
			inside_code -= 65248
		}
		if inside_code < 32 || inside_code > 126 {
			retstr += string(i)
		} else {
			retstr += string(inside_code)
		}
	}
	return retstr
}
