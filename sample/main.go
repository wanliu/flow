package main

import (
	. "github.com/wanliu/flow/builtin"
	goflow "github.com/wanliu/goflow"
)

type testNet struct {
	goflow.Graph
}

func newLuisGraph() *testNet {
	net := new(testNet)
	net.InitGraphState()

	cm := new(ContextManager)
	ri := new(ReadInput)
	luis := new(LuisAnalyze)
	lo := new(Log)
	stringer := NewStringifier()

	cm.SendHandle = func(ctx, parent Context) error {
		if msg, ok := ctx.GlobalValue("Msg").(string); ok {
			ctx.Send(msg)
		}
		return nil
	}
	net.Add(cm, "cm")
	net.Add(ri, "ri")
	net.Add(luis, "luis")
	net.Add(lo, "log")
	net.Add(stringer, "stringer")

	net.AddIIP("052297dc-12b9-4044-8220-a21a20d72581", "luis", "AppId")
	net.AddIIP("6b916f7c107643069c242cf881609a82", "luis", "Key")
	net.AddIIP("请输入你的话:", "ri", "Prompt")

	net.Connect("ri", "Out", "luis", "In")
	net.Connect("luis", "Out", "stringer", "In")
	net.Connect("stringer", "Out", "log", "In")
	net.MapInPort("Ctx", "cm", "Ctx")

	// net.MapOutPort("Out", "luis", "Out")
	return net
}

func main() {
	net := newLuisGraph()

	in := make(chan Context)
	// out1 := make(chan string)
	// out := make(chan ResultParams)
	// out2 := make(chan string)

	net.SetInPort("Ctx", in)
	// net.SetOutPort("Out", out)

	goflow.RunNet(net)
	<-net.Ready()
	// log.Printf("net: %# v", pretty.Formatter(net))

	ctx := NewContext()
	in <- ctx

	<-net.Wait()
}
