package resolves

import (
	// "encoding/json"
	// "io/ioutil"
	"log"
	// "net/http"
	// "net/url"
	"reflect"
	"strconv"
	"strings"
	// "regexp"
	"time"

	"github.com/hysios/apiai-go"
	. "github.com/wanliu/flow/context"
)

// 处理开单的逻辑结构, 不需要是组件
// 作为context的一个部分，或者存在一个Value中
type OpenOrderResolve struct {
	AiParams apiai.Result
	Products ItemsResolve
	Address  string
	Time     time.Time
	DefTime  string
	Current  Resolve
}

func NewOpenOrderResolve(ctx Context) *OpenOrderResolve {
	resolve := new(OpenOrderResolve)

	aiResult := ctx.Value("Result").(apiai.Result)

	// aiResult.Entities = DistinctEntites(aiResult.Entities)
	// aiResult.Entities = DeduplicateEntities(aiResult.Entities)
	// SortEntities(aiResult.Entities)

	resolve.AiParams = aiResult
	resolve.ExtractFromLuis()

	return resolve
}

func (r *OpenOrderResolve) Solve(aiResult apiai.Result) string {
	return r.Answer()
}

// 从ｌｕｉｓ数据构造结构数据
func (r *OpenOrderResolve) ExtractFromLuis() {
	r.ExtractItems()
	r.ExtractAddress()
	r.ExtractTime()
}

// TODO 无法识别全角数字
func (r *OpenOrderResolve) ExtractItems() {
	r.ExtractProducts()
	quantities := r.ExtractQuantity()

	for i, q := range quantities {
		if len(r.Products.Products) >= i+1 {
			pr := r.Products.Products[i]
			pr.Quantity = q
		}
	}

	// for _, p := range r.Products.Products {
	// 	p.CheckResolved()
	// }
}

func (r *OpenOrderResolve) ExtractQuantity() []int {
	result := make([]int, 0, 50)

	if quantities, exist := r.AiParams.Params["quantity"]; exist {
		qs := reflect.ValueOf(quantities)

		for i := 0; i < qs.Len(); i++ {
			q := qs.Index(i).Interface()

			switch t := q.(type) {
			case string:
				qs := q.(string)
				qi := extractQuantity(qs)
				result = append(result, qi)
			case float64:
				qf := q.(float64)
				result = append(result, int(qf))
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

	log.Printf("NUMS: %v", numsCgd)

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

func (r *OpenOrderResolve) ExtractProducts() {
	if products, exist := r.AiParams.Params["products"]; exist {
		ps := reflect.ValueOf(products)

		for i := 0; i < ps.Len(); i++ {
			p := ps.Index(i)
			product := ItemResolve{
				Resolved: false,
				Name:     p.Interface().(string),
				Price:    0,
				Quantity: 0, // 默认值
				Product:  p.Interface().(string),
			}

			r.Products.Add(product)
		}

		// ps := products.([]string)

		// for i, _ := range ps {
		// 	product := ItemResolve{
		// 		Resolved: false,
		// 		Name:     p,
		// 		Price:    0,
		// 		Quantity: 0, // 默认值
		// 		Product:  p,
		// 	}

		// 	r.Products.Add(product)
		// }
	}
}

func (r *OpenOrderResolve) ExtractAddress() {
	if address, exist := r.AiParams.Params["street-address"]; exist {
		r.Address = address.(string)
	}
}

func (r *OpenOrderResolve) ExtractTime() {
	if t, exist := r.AiParams.Params["date"]; exist {
		if aiTime, err := time.Parse("2006-01-02", t.(string)); err == nil {
			r.Time = aiTime
		}
	}
}

func (r *OpenOrderResolve) SetDefTime(t string) {
	r.DefTime = t

	if r.Time.IsZero() && r.DefTime != "" {
		r.SetTimeByDef()
	}
}

func (r *OpenOrderResolve) SetTimeByDef() {
	if r.DefTime == "今天" {
		r.Time = time.Now()
	} else if r.DefTime == "明天" {
		r.Time = time.Now().Add(24 * time.Hour)
	}
}

func (r OpenOrderResolve) EmptyProducts() bool {
	return len(r.Products.Products) == 0
}

func (r OpenOrderResolve) Answer() string {
	desc := ""

	desc = desc + "=== 订单输入完成 ===\n"
	desc = desc + "本订单包含如下商品：" + "\n"

	// params := url.Values{
	// 	"auth_token":   {"5f567b5efc3e4d0aa0d9c40922ae07aa"},
	// 	"street":       {r.Address},
	// 	"deliver_time": {r.Time.Format("2006年01月02日")},
	// }

	for _, p := range r.Products.Products {
		desc = desc + p.Product + " " + strconv.Itoa(p.Quantity) + "件\n"

		// nk := "items[" + strconv.Itoa(i) + "][name]"
		// nv := p.Product
		// params.Add(nk, nv)

		// qk := "items[" + strconv.Itoa(i) + "][quantity]"
		// qv := strconv.Itoa(p.Quantity)
		// params.Add(qk, qv)
	}

	desc = desc + "地址:" + r.Address + "\n"
	desc = desc + "送货时间" + r.Time.Format("2006年01月02日") + "\n"
	desc = desc + "=== 结束 ===\n"

	return desc

	// res, err := http.PostForm("http://192.168.0.155:3000/api/v1/temp_orders", params)

	// if err != nil {
	// 	// return err.Error()
	// 	return "服务暂时不可用，请稍后再试"
	// } else {
	// 	defer res.Body.Close()
	// 	body, e := ioutil.ReadAll(res.Body)

	// 	if e != nil {
	// 		return e.Error()
	// 	} else {
	// 		var result Res
	// 		json.Unmarshal(body, &result)

	// 		if res.StatusCode == 422 {
	// 			return result.Error
	// 		} else {
	// 			return desc + "请通过以下地址完成订单操作：" + result.Confirm_path
	// 		}

	// 	}
	// }
}

type Res struct {
	Id           int
	Confirm_path string
	Error        string
}
