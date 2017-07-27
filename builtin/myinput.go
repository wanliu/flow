package builtin

import (
	"log"
	"strings"

	"github.com/hysios/apiai-go"
	. "github.com/wanliu/flow/context"
	goflow "github.com/wanliu/goflow"
)

type MyInput struct {
	ReadInput
	_ctx Context
	Ctx  <-chan Context
	Next chan<- Context
}

func NewMyInput() interface{} {
	return new(MyInput)
}

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

func (in *MyInput) Loop() {
	for {
		select {
		// Handle immediate terminate signal from network
		case <-in.Component.Term:
			return
		case input, ok := <-in.GetLine():
			if ok {
				if strings.Trim(input, " \n") != "" {
					in.Out <- input
					in.Next <- in._ctx
				}
			}
		case ctx := <-in.Ctx:
			in._ctx = ctx
		case prompt, _ := <-in.Prompt:
			in.SetPrompt(prompt)
		}
	}
}

func (q *QuerySave) Init() {
	q.Fields = []string{"Ctx", "Result", "Text"}
	q.Process = func() error {
		res, rok := q.Value("Result").(apiai.Result)
		ctx, cok := q.Value("Ctx").(Context)
		txt, tok := q.Value("Text").(string)
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

		if cok && tok {
			ctx.SetValue("Text", txt)

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

func (q *QuerySave) OnText(text string) {
	q.SetValue("Text", text)
}
