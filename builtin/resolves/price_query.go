package resolves

// import (
// 	"log"
// 	// "math/rand"
// 	// "strconv"
// 	"strings"
// 	// "time"
// 	"encoding/json"
// 	"io/ioutil"
// 	"net/http"

// 	. "github.com/wanliu/flow/builtin/luis"
// 	. "github.com/wanliu/flow/context"
// )

// func NewPriceQueryResolve(ctx Context) *PriceQueryResolve {
// 	resolve := new(PriceQueryResolve)

// 	luis := ctx.Value("Result").(ResultParams)

// 	luis.Entities = DistinctEntites(luis.Entities)
// 	luis.Entities = DeduplicateEntities(luis.Entities)
// 	SortEntities(luis.Entities)

// 	resolve.LuisParams = luis
// 	resolve.ExtractFromLuis()

// 	return resolve
// }

// type PriceQueryResolve struct {
// 	LuisParams ResultParams
// 	Products   []*PriceProductResolve
// 	Current    *PriceProductResolve
// }

// func (r *PriceQueryResolve) ExtractFromLuis() {
// 	for _, item := range r.LuisParams.Entities {
// 		if item.Type == "products" {
// 			product := PriceProductResolve{
// 				Resolved:   false,
// 				Name:       item.Entity,
// 				Stock:      0,
// 				Resolution: item.Resolution,
// 			}

// 			product.CheckResolved()

// 			r.Products = append(r.Products, &product)
// 		} else {
// 			log.Printf("type: %v", item.Type)
// 		}
// 	}
// }

// func (r *PriceQueryResolve) Next() Resolve {
// 	for _, pr := range r.Products {
// 		if !pr.Resolved {
// 			r.Current = pr
// 			return pr
// 		}
// 	}

// 	return new(PriceProductResolve)
// }

// func (r *PriceQueryResolve) Solve(luis ResultParams) (bool, string, string) {

// 	solved, finishedNotition, nextNotition := r.Current.Solve(luis)

// 	if solved {
// 		if r.Fullfilled() {
// 			finishedNotition = finishedNotition + "\n" + r.Answer()
// 			return true, finishedNotition, nextNotition
// 		} else {
// 			next := finishedNotition + "\n" + r.Next().Hint()

// 			return false, finishedNotition, next
// 		}
// 	} else {
// 		return solved, finishedNotition, nextNotition
// 	}
// }

// func (r PriceQueryResolve) Fullfilled() bool {
// 	for _, p := range r.Products {
// 		if !p.Resolved {
// 			return false
// 		}
// 	}

// 	return true
// }

// func (r PriceQueryResolve) Answer() string {
// 	selected := make([]string, 0, 10)

// 	for _, p := range r.Products {
// 		selected = append(selected, "name[]="+p.Product)
// 	}

// 	query := "?auth_token=5f567b5efc3e4d0aa0d9c40922ae07aa&" + strings.Join(selected, "&")

// 	res, err := http.Get("http://192.168.0.155:3000/api/v1/query_items/price" + query)

// 	if err != nil {
// 		// return err.Error()
// 		return "服务暂时不可用，请稍后再试"
// 	} else {
// 		defer res.Body.Close()
// 		body, e := ioutil.ReadAll(res.Body)

// 		if e != nil {
// 			return e.Error()
// 		} else {
// 			var result PriceRes
// 			json.Unmarshal(body, &result)

// 			if res.StatusCode == 422 {
// 				return result.Error
// 			} else {
// 				return result.String()
// 			}
// 		}
// 	}
// }

// func (r PriceQueryResolve) EmptyProducts() bool {
// 	return len(r.Products) == 0
// }

// type PriceRes struct {
// 	Items []ItemRes
// 	Error string
// }

// func (p PriceRes) String() string {
// 	result := make([]string, 0, 10)

// 	for _, i := range p.Items {
// 		result = append(result, i.String())
// 	}

// 	return strings.Join(result, ",")
// }

// type ItemRes struct {
// 	Name  string
// 	Price string
// }

// func (i ItemRes) String() string {
// 	return i.Name + "的价格为" + i.Price + "元"
// }
