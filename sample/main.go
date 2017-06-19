package main

import (
	"log"

	"github.com/kr/pretty"
	. "github.com/wanliu/flow/builtin"
	. "github.com/wanliu/flow/context"
	goflow "github.com/wanliu/goflow"
)

type testNet struct {
	goflow.Graph
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
	ctx, _ := NewContext()
	in <- ctx
	// <-out

	<-net.Wait()
}
