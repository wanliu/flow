package builtin

import (
	"log"

	"github.com/hysios/apiai-go"
	. "github.com/wanliu/flow/context"
	goflow "github.com/wanliu/goflow"
)

func NewQuerySave() interface{} {
	return new(QuerySave)
}

type QuerySave struct {
	goflow.Component
	MultiField
	Result <-chan apiai.Result
	Ctx    <-chan Context
	Out    chan<- Context
}

func (q *QuerySave) Init() {
	q.Fields = []string{"Ctx", "Result"}
	q.Process = func() error {
		res, rok := q.Value("Result").(apiai.Result)
		ctx, cok := q.Value("Ctx").(Context)
		if rok && cok {
			ctx.SetValue("Result", res)
			intent := res.Metadata.IntentName
			score := res.Score
			query := res.ResolvedQuery

			log.Printf("意图解析\"%s\" -> %s 准确度: %2.2f%%", query, intent, score*100)
			go func() {
				q.Out <- ctx
			}()
		}

		return nil
	}
}

func (q *QuerySave) OnCtx(ctx Context) {
	q.SetValue("Ctx", ctx)
}

func (q *QuerySave) OnResult(res apiai.Result) {
	q.SetValue("Result", res)
}
