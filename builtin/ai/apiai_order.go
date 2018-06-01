package ai

import (
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/hysios/apiai-go"
)

const (
	MODE_PRODUCTS      = 1
	MODE_PRODUCT_ITEMS = 2
	MODE_PROD_TASTES   = 3
)

type ApiAiOrder struct {
	AiResult apiai.Result
}

func (aa ApiAiOrder) Mode() int {
	mode := MODE_PRODUCTS

	if _, ok := aa.AiResult.Params["productItems"]; ok {
		mode = MODE_PRODUCT_ITEMS
	}

	if _, ok := aa.AiResult.Params["prodTastes"]; ok {
		mode = MODE_PROD_TASTES
	}

	return mode
}

func (aa ApiAiOrder) Score() float64 {
	return aa.AiResult.Score
}

func (aa ApiAiOrder) Query() string {
	return aa.AiResult.ResolvedQuery
}

func (aa ApiAiOrder) Items() []Item {
	if aa.Mode() == MODE_PRODUCT_ITEMS {
		return aa.ExtractProductItems("productItems")
	} else if aa.Mode() == MODE_PROD_TASTES {
		products := aa.prodTastItems()
		quantities := aa.Quantities()

		return composeItems(products, quantities)
	} else {
		products := aa.Products()
		quantities := aa.Quantities()

		return composeItems(products, quantities)
	}

	return make([]Item, 0)

	// for i, q := range quantities {
	// 	if len(products) >= i+1 {
	// 		products[i].Quantity = q.Quantity
	// 		products[i].Unit = q.Unit
	// 	}
	// }

	// return products
}

func (aa ApiAiOrder) Products() []Item {
	return aa.ExtractProducts("products")
}

func (aa ApiAiOrder) prodTastItems() []Item {
	return aa.ExtractProducts("prodTastes")
}

func (aa ApiAiOrder) Quantities() []Item {
	return aa.ExtractQuantities("quantity")
}

func (aa ApiAiOrder) GiftItems() []Item {
	if aa.Mode() == MODE_PRODUCT_ITEMS {
		return aa.ExtractProductItems("giftItems")
	}

	gifts := aa.GiftProducts()
	quantities := aa.GiftQuantities()

	return composeItems(gifts, quantities)

	// for i, q := range quantities {
	// 	if len(gifts) >= i+1 {
	// 		gifts[i].Quantity = q.Quantity
	// 		gifts[i].Unit = q.Unit
	// 	}
	// }

	// return gifts
}

func composeItems(products []Item, quantities []Item) []Item {
	result := make([]Item, 0, 0)
	l := len(products)
	qlen := len(quantities)

	if l < qlen {
		l = qlen
	}

	for i := 0; i < l; i++ {
		item := Item{}

		if len(products) >= i+1 {
			p := products[i]
			item.Product = p.Product

			if p.Taste != "" {
				item.Taste = p.Taste
			}
		}

		if len(quantities) >= i+1 {
			q := quantities[i]
			item.Quantity = q.Quantity
			item.Unit = q.Unit
		}

		result = append(result, item)
	}

	return result
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

func (aa ApiAiOrder) Count() int {
	if c, exist := aa.AiResult.Params["number"]; exist {
		switch c.(type) {
		case float64:
			fval := c.(float64)
			return int(fval)
		case float32:
			fval := c.(float32)
			return int(fval)
		case int:
			return c.(int)
		case string:
			sval := c.(string)
			ival, err := strconv.Atoi(sval)
			if err != nil {
				return 0
			} else {
				return ival
			}
		}
	}

	return 0
}

func (aa ApiAiOrder) Duration() string {
	if c, exist := aa.AiResult.Params["duration"]; exist {
		// return c.(string)
		switch v := c.(type) {
		case string:
			return v
			// case []interface{}:
			// 	if len(v) > 0 {
			// 		i := v[0]
			// 		return i.(string)
			// 	}
		}
	}

	return ""
}

func (aa ApiAiOrder) Time() time.Time {
	if t, exist := aa.AiResult.Params["date"]; exist {
		switch v := t.(type) {
		case string:
			if aiTime, err := time.Parse("2006-01-02", v); err == nil {
				return aiTime
			}
		}
	}

	return time.Time{}
}

func (aa ApiAiOrder) Fulfiled() bool {
	return true
}

func (aa ApiAiOrder) Note() string {
	if imp, exist := aa.AiResult.Params["important"]; exist {
		// return imp.(string)
		switch v := imp.(type) {
		case string:
			return v
			// case []interface{}:
			// 	if len(v) > 0 {
			// 		i := v[0]
			// 		return i.(string)
			// 	}
		}
	}

	return ""
}

// MODE 1 extracting, by key products
func (aa ApiAiOrder) ExtractProducts(t string) []Item {
	result := make([]Item, 0, 50)

	if products, exist := aa.AiResult.Params[t]; exist {
		ps := reflect.ValueOf(products)

		for i := 0; i < ps.Len(); i++ {
			p := ps.Index(i).Interface()

			switch v := p.(type) {
			case string:
				// name := p.(string)
				item := Item{Product: v}
				result = append(result, item)
			case map[string]interface{}:
				// itemMap := p.(map[string]interface{})
				name, _ := v["product"].(string)
				taste, _ := v["taste"].(string)

				if name != "" {
					item := Item{
						Product: name,
					}

					if taste != "" {
						item.Taste = taste
					}

					result = append(result, item)
				}
			}
		}
	}

	return result
}

// MODE 2 extracting, by key productItems
func (aa ApiAiOrder) ExtractProductItems(s string) []Item {
	result := make([]Item, 0)

	if prodItems, ok := aa.AiResult.Params[s]; ok {
		ps := reflect.ValueOf(prodItems)

		for i := 0; i < ps.Len(); i++ {
			var name, unit, spec, taste string
			var quantity int

			p := ps.Index(i)
			prodItem := p.Interface().(map[string]interface{})
			name, _ = prodItem["product"].(string)
			spec, _ = prodItem["spec"].(string)
			taste, _ = prodItem["taste"].(string)

			quanMap, _ := prodItem["quantity"].(map[string]interface{})
			numberFloat, ok := quanMap["number"].(float64)
			if ok {
				quantity = int(numberFloat)
			}

			unit, _ = quanMap["unit"].(string)

			item := Item{
				Product:  name,
				Quantity: quantity,
				Unit:     unit,
				Spec:     spec,
				Taste:    taste,
			}
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

			switch v := q.(type) {
			case string:
				// qs := q.(string)
				qi := extractQuantity(v)
				item := Item{Quantity: qi}
				result = append(result, item)
			case float64:
				// qf := q.(float64)
				item := Item{Quantity: int(v)}
				result = append(result, item)
			case map[string]interface{}:
				// quanMap := q.(map[string]interface{})
				log.Printf("quantity: %v\n", v)
				qf, ok := v["number"].(float64)
				if !ok {
					qf, _ = v["quantity"].(float64)
				}

				quantity := int(qf)
				item := Item{Quantity: quantity}

				if unit, ok := v["unit"].(string); ok {
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
