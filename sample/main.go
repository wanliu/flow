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
			q.Out <- ctx
		}

		return nil
	}
}

func (q *QuerySave) OnCtx(ctx Context) {
	log.Printf("QuerySave Ctx: %#v", ctx)
	q.SetValue("Ctx", ctx)
}

func (q *QuerySave) OnResult(res ResultParams) {
	log.Printf("QuerySave Result: %#v", res)
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
		qs = new(QuerySave)
		ic = new(IntentCheck)
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

	net.AddIIP("052297dc-12b9-4044-8220-a21a20d72581", "Luis", "AppId")
	net.AddIIP("6b916f7c107643069c242cf881609a82", "Luis", "Key")
	net.AddIIP("请输入你的话:", "Line", "Prompt")
	net.AddIIP("开单", "OpenOrder", "Intent")
	net.AddIIP(0.75, "OpenOrder", "Score")

	net.Connect("Line", "Out", "Luis", "In")
	net.Connect("Line", "Next", "Merge", "Ctx")
	net.Connect("Luis", "Out", "Merge", "Result")
	net.Connect("Merge", "Out", "CM", "Ctx")
	net.Connect("CM", "Process", "OpenOrder", "Ctx")

	// net.Connect("OpenOrder", "Out", "Log", "In")

	net.MapInPort("In", "Line", "Ctx")
	net.MapOutPort("Out", "OpenOrder", "Out")

	// net.MapOutPort("Out", "luis", "Out")
	return net
}

func main() {
	net := newLuisGraph()

	in := make(chan Context)
	out := make(chan Context)

	net.SetInPort("In", in)
	net.SetOutPort("Out", out)

	goflow.RunNet(net)
	<-net.Ready()
	log.Printf("net: %# v", pretty.Formatter(net))
	ctx := NewContext()
	in <- ctx
	<-out

	<-net.Wait()
}
