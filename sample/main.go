package main

import (
	"log"
	"strings"

	"github.com/kr/pretty"
	. "github.com/wanliu/flow/builtin"
	goflow "github.com/wanliu/goflow"
)

type testNet struct {
	goflow.Graph
}

type MyInput struct {
	ReadInput
	_ctx Context
	Ctx  <-chan Context
	Next chan<- Context
}

type QuerySave struct {
	goflow.Component
	MultiField
	Result <-chan ResultParams
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
	q.Fields = []string{"Ctx", "Result"}
	q.Process = func() error {
		res, rok := q.Value("Result").(ResultParams)
		ctx, cok := q.Value("Ctx").(Context)
		if rok && cok {
			ctx.SetValue("Result", res)

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

func (q *QuerySave) OnResult(res ResultParams) {
	top := res.TopScoringIntent
	log.Printf("意图解析 -> %s 准确度: %2.2f%%", top.Intent, top.Score*100)
	q.SetValue("Result", res)
}

func newLuisGraph() *testNet {
	net := new(testNet)
	net.InitGraphState()

	var (
		cm    = new(ContextManager)
		input = new(MyInput)
		luis  = new(LuisAnalyze)
		// lo    = new(Log)
		qs       = new(QuerySave)
		ic       = new(IntentCheck)
		products = new(TryGetProducts)
		reset    = new(CtxReset)
	)

	cm.SendHandle = func(ctx, parent Context) error {
		if msg, ok := ctx.GlobalValue("Msg").(string); ok {
			ctx.Send(msg)
		}
		return nil
	}
	net.Add(cm, "CM")
	net.Add(input, "Line")
	net.Add(luis, "Luis")
	// net.Add(lo, "Log")
	net.Add(ic, "OpenOrder")
	net.Add(qs, "Merge")
	net.Add(products, "TryProducts")
	net.Add(reset, "Reset")

	net.AddIIP("8b65b31f-05b0-4da0-ab98-afa62c0e80ae", "Luis", "AppId")
	net.AddIIP("9c6711ad95c846a792248515cb6d1a23", "Luis", "Key")
	net.AddIIP("请输入你的话:", "Line", "Prompt")
	net.AddIIP("开单", "OpenOrder", "Intent")
	net.AddIIP(0.70, "OpenOrder", "Score")
	net.AddIIP("products", "TryProducts", "Type")

	net.Connect("Line", "Out", "Luis", "In")
	net.Connect("Line", "Next", "Merge", "Ctx")
	net.Connect("Luis", "Out", "Merge", "Result")
	net.Connect("Merge", "Out", "CM", "Ctx")
	net.Connect("CM", "Process", "OpenOrder", "Ctx")
	net.Connect("OpenOrder", "Out", "TryProducts", "Ctx")
	net.Connect("TryProducts", "Out", "Reset", "In")
	// net.Connect("OpenOrder", "Out", "Log", "In")

	net.MapInPort("In", "Line", "Ctx")
	// net.MapOutPort("Out", "TryProducts", "Out")

	// net.MapOutPort("Out", "luis", "Out")
	return net
}

func main() {
	net := newLuisGraph()

	in := make(chan Context)
	// out := make(chan Context)

	net.SetInPort("In", in)
	// net.SetOutPort("Out", out)

	goflow.RunNet(net)
	<-net.Ready()
	log.Printf("net: %# v", pretty.Formatter(net))
	ctx := NewContext()
	in <- ctx
	// <-out

	<-net.Wait()
}
