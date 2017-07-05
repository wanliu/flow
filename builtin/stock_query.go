package builtin

import (
	// "fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	. "github.com/wanliu/flow/builtin/luis"
	. "github.com/wanliu/flow/builtin/resolves"
	. "github.com/wanliu/flow/context"
)

type StockQuery struct {
	TryGetEntities
	Ctx  <-chan Context
	Type <-chan string
	Out  chan<- ReplyData
}

func NewStockQuery() interface{} {
	return new(StockQuery)
}

func (query *StockQuery) OnCtx(ctx Context) {
	stockQuery := NewStockQueryResolve(ctx)
	childCtx := ctx.NewContext()
	childCtx.SetValue("stockQuery", stockQuery)

	output := ""

	if stockQuery.EmptyProducts() {
		output = "没有相关的产品"
	} else {
		ctx.Push(childCtx)
		output = stockQuery.Next().Hint()
	}

	log.Printf("OUTPUT: %v", output)

	replyData := ReplyData{output, ctx}
	query.Out <- replyData

	go func(task Context) {
		task.Wait(query.TaskHandle)
		// if ctx.
		// ctx.Pop()
	}(childCtx)
}

func (query *StockQuery) TaskHandle(ctx Context, raw interface{}) error {

	params := raw.(Context).Value("Result").(ResultParams)

	stockQuery := ctx.Value("stockQuery").(*StockQueryResolve)

	solved, finishNotition, nextNotition := stockQuery.Solve(params)

	if solved {
		log.Printf("测试输出打印: \n%v", finishNotition)

		reply := ReplyData{finishNotition, ctx}
		query.Out <- reply

		ctx.Pop() // 将当前任务踢出队列
	} else {
		log.Printf("测试输出打印: \n%v\n", nextNotition)

		reply := ReplyData{nextNotition, ctx}
		query.Out <- reply
	}
	// ctx.Send(raw)
	return nil
}

func NewStockQueryResolve(ctx Context) *StockQueryResolve {
	resolve := new(StockQueryResolve)

	resolve.LuisParams = ctx.Value("Result").(ResultParams)
	resolve.ExtractFromLuis()

	return resolve
}

type StockQueryResolve struct {
	LuisParams ResultParams
	Products   []*ProductResolve
	Current    *ProductResolve
}

func (sqr *StockQueryResolve) ExtractFromLuis() {
	for _, item := range sqr.LuisParams.Entities {
		if item.Type == "products" {
			product := ProductResolve{
				Resolved:   false,
				Name:       item.Entity,
				Stock:      0,
				Resolution: item.Resolution,
			}

			product.CheckResolved()

			product.Parent = sqr
			sqr.Products = append(sqr.Products, &product)
		} else {
			log.Printf("type: %v", item.Type)
		}
	}
}

func (sqr *StockQueryResolve) Next() Resolve {
	for _, pr := range sqr.Products {
		if !pr.Resolved {
			sqr.Current = pr
			return pr
		}
	}

	return ProductResolve{}
}

func (sqr *StockQueryResolve) Solve(luis ResultParams) (bool, string, string) {

	solved, finishedNotition, nextNotition := sqr.Current.Solve(luis)

	if solved {
		if sqr.Fullfilled() {
			finishedNotition = finishedNotition + "\n" + sqr.Answer()
			return true, finishedNotition, nextNotition
		} else {
			next := finishedNotition + "\n" + sqr.Next().Hint()

			return false, finishedNotition, next
		}
	} else {
		return solved, finishedNotition, nextNotition
	}
}

func (sqr StockQueryResolve) Fullfilled() bool {
	for _, p := range sqr.Products {
		if !p.Resolved {
			return false
		}
	}

	return true
}

func (sqr StockQueryResolve) Answer() string {
	selected := make([]string, 0, 10)

	// TODO 查询后台商品价格
	rand.Seed(time.Now().UTC().UnixNano())

	for _, p := range sqr.Products {
		p.Stock = rand.Intn(100)

		if p.Stock <= 50 {
			selected = append(selected, p.Product+"已经没货")
		} else {
			selected = append(selected, p.Product+"还有库存："+strconv.Itoa(p.Stock))
		}
	}

	return strings.Join(selected, ", ")
}

func (sqr StockQueryResolve) EmptyProducts() bool {
	return len(sqr.Products) == 0
}

type ProductResolve struct {
	Resolved   bool
	Name       string
	Price      float64
	Stock      int
	Product    string
	Resolution Resolution
	Parent     *StockQueryResolve
}

func (pr ProductResolve) Solve(luis ResultParams) (bool, string, string) {
	if luis.TopScoringIntent.Intent == "选择" {
		// TODO 无法识别全角数字
		number := strings.Trim(luis.Entities[0].Resolution.Value, " ")
		chose, _ := strconv.ParseInt(number, 10, 64)
		inNum := int(chose)

		for _, product := range pr.Parent.Products {
			if product.Name == pr.Name {
				if product.Product == "" {
					if len(product.Resolution.Values) >= inNum {
						prdName := product.Resolution.Values[chose-1]
						product.Product = prdName
						product.CheckResolved()

						return true, "已选择" + prdName, "err"
					} else {
						return false, "", "超出选择范围\n" + product.Hint()
					}
				}
			}
		}

		return false, "", "错误的操作，没有可供选择的商品"
	} else {
		return false, "", "无效的输入\n" + pr.Hint()
	}
}

func (pr ProductResolve) Hint() string {
	result := ""

	if pr.Product == "" && len(pr.Resolution.Values) > 0 {
		index := 1
		choses := "\n"

		for _, value := range pr.Resolution.Values {
			choses = choses + strconv.Itoa(index) + ": " + value + "\n"
			index = index + 1
		}

		choses = choses + "\n"

		result = "我们有下列的 " + pr.Name + " 产品:" + choses + "请输入序号选择你要查询的商品"
	}

	return result
}

func (pr *ProductResolve) CheckResolved() {
	if len(pr.Resolution.Values) == 0 {
		pr.Product = pr.Name
	}

	if pr.Product != "" {
		pr.Resolved = true
	}
}
