package builtin

import (
	"log"
	"testing"
	"time"

	"github.com/kr/pretty"
	flow "github.com/wanliu/goflow"
)

type testNet struct {
	flow.Graph
}

// type Component1 struct {
// 	ContextComponent
// 	In <-chan string
// }

type ContextManager2 struct {
	ContextManager
	Ctx <-chan Context
}

type Component2 struct {
	ContextComponent
	In  <-chan string
	Out chan<- string
}

type Component3 struct {
	ContextComponent
	Enter <-chan Context
	Out   chan<- string
}

type Component4 struct {
	ContextComponent
	Enter <-chan Context
	Out   chan<- string
}

type IsContain struct {
	flow.Component
	contain string
	Contain <-chan string
	In      <-chan Context
	Yes     chan<- bool
	Next    chan<- Context
}

func (c *Component3) Init() {
	c.TaskHandle = func(ctx Context, raw interface{}) error {
		if msg, ok := raw.(string); ok {
			c.Out <- msg
		}
		return nil
	}
}

func (c *Component4) Init() {
	c.TaskHandle = func(ctx Context, raw interface{}) error {
		if msg, ok := raw.(string); ok {
			c.Process <- ctx
			c.Out <- msg
		}
		return nil
	}
}

func (cm *ContextManager2) Init() {
	cm.SendHandle = func(ctx, parent Context) error {
		log.Printf("cm ctx: %#v, parent: %#v", ctx, parent)
		if msg, ok := ctx.GlobalValue("Msg").(string); ok {
			ctx.Send(msg)
		}
		return nil
	}
}

func (cm *ContextManager2) OnCtx(ctx Context) {
	cm.ContextManager.OnCtx(ctx)
}

func (c *IsContain) OnIn(ctx Context) {
	if msg, ok := ctx.GlobalValue("Msg").(string); ok {
		if msg == c.contain {
			c.Yes <- true
			c.Next <- ctx
		}
	}
}

func (c *IsContain) OnContain(con string) {
	c.contain = con
}

func newGraph() *testNet {
	net := new(testNet)
	net.InitGraphState()
	// cm := NewContextManager()
	c1 := new(Component2)
	c2 := new(Component2)

	// net.Add(cm, "ContextManager")
	net.Add(c1, "c1")
	net.Add(c2, "c2")

	net.Connect("c1", "Out", "c2", "In")

	net.MapInPort("In", "c1", "In")
	// net.MapInPort("Ctx", "ContextManager", "Ctx")
	net.MapOutPort("Out", "c2", "Out")
	return net
}

func newContextGraph() *testNet {
	net := new(testNet)
	net.InitGraphState()

	cm := new(ContextManager)
	c1 := new(ContextComponent)
	// c2 := new(ContextComponent)

	net.Add(cm, "cm")
	net.Add(c1, "c1")
	// net.Add(c2, "c2")

	net.Connect("cm", "Process", "c1", "Enter")
	net.MapInPort("In", "cm", "Ctx")
	net.MapOutPort("Out", "c1", "Process")
	return net
}

func newContextGraph2() *testNet {
	net := new(testNet)
	net.InitGraphState()

	cm := new(ContextManager)
	c1 := new(Component3)
	// c2 := new(ContextComponent)

	net.Add(cm, "cm")
	net.Add(c1, "c1")
	// net.Add(c2, "c2")

	net.Connect("cm", "Process", "c1", "Enter")
	net.MapInPort("In", "cm", "Ctx")
	net.MapOutPort("Out", "c1", "Out")
	return net
}

func newContextGraph2Task() *testNet {
	net := new(testNet)
	net.InitGraphState()

	cm := new(ContextManager)

	cm.SendHandle = func(ctx, parent Context) error {
		if msg, ok := ctx.GlobalValue("Msg").(string); ok {
			ctx.Send(msg)
		}
		return nil
	}
	c1 := new(ContextComponent)
	// c2 := new(Component4)
	is := new(IsContain)
	ctrl := new(CtxControl)
	str := new(ContextString)

	net.Add(cm, "cm")
	net.Add(c1, "c1")
	// net.Add(c2, "c2")
	net.Add(is, "is")
	net.Add(ctrl, "ctrl")
	net.Add(str, "str")
	net.AddIIP("Msg", "str", "Field")
	net.AddIIP("Hello", "is", "Contain")

	net.Connect("cm", "Process", "c1", "Enter")
	// net.Connect("c1", "Process", "c1", "Done")
	net.Connect("c1", "Process", "is", "In")
	net.Connect("is", "Yes", "ctrl", "Done")
	net.Connect("is", "Next", "ctrl", "Ctx")
	net.Connect("c1", "Next", "str", "Ctx")

	net.MapInPort("In", "cm", "Ctx")
	net.MapOutPort("Out", "str", "Out")
	return net
}

func newLuisGraph() *testNet {
	net := new(testNet)
	net.InitGraphState()

	cm := new(ContextManager)
	ri := new(ReadInput)
	luis := new(LuisAnalyze)

	cm.SendHandle = func(ctx, parent Context) error {
		if msg, ok := ctx.GlobalValue("Msg").(string); ok {
			ctx.Send(msg)
		}
		return nil
	}
	net.Add(cm, "cm")
	net.Add(ri, "ri")
	net.Add(luis, "luis")

	net.AddIIP("052297dc-12b9-4044-8220-a21a20d72581", "luis", "AppId")
	net.AddIIP("6b916f7c107643069c242cf881609a82", "luis", "Key")
	net.AddIIP("请输入你的话:", "ri", "Prompt")

	net.Connect("ri", "Out", "luis", "In")
	net.MapInPort("Ctx", "cm", "Ctx")
	net.MapOutPort("Out", "luis", "Out")
	return net
}

func (c *Component2) OnIn(msg string) {
	c.Out <- msg
}

func TestContextManager(t *testing.T) {
	net := newGraph()
	in := make(chan string)
	// ctx := make(chan context.IContext)
	out := make(chan string)

	net.SetInPort("In", in)
	// net.SetInPort("Ctx", ctx)
	net.SetOutPort("Out", out)
	flow.RunNet(net)
	<-net.Ready()
	log.Printf("running net")

	for i := 0; i < 3; i++ {
		in <- "hello"
		<-out
	}

	close(in)
	<-net.Wait()

}

func BenchmarkGraph(b *testing.B) {
	for i := 0; i < b.N; i++ {
		net := newGraph()
		in := make(chan string)
		// ctx := make(chan context.IContext)
		out := make(chan string)

		net.SetInPort("In", in)
		// net.SetInPort("Ctx", ctx)
		net.SetOutPort("Out", out)
		flow.RunNet(net)

		<-net.Ready()
		in <- "hello"
		<-out

		close(in)
		<-net.Wait()
	}
}

func BenchmarkGraph2(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			net := newGraph()
			in := make(chan string)
			// ctx := make(chan context.IContext)
			out := make(chan string)

			net.SetInPort("In", in)
			// net.SetInPort("Ctx", ctx)
			net.SetOutPort("Out", out)
			flow.RunNet(net)

			<-net.Ready()
			in <- "hello"
			<-out

			close(in)
			<-net.Wait()
		}
	})
}

func BenchmarkGraph3(b *testing.B) {
	net := newGraph()
	in := make(chan string)
	// ctx := make(chan context.IContext)
	out := make(chan string)

	net.SetInPort("In", in)
	// net.SetInPort("Ctx", ctx)
	net.SetOutPort("Out", out)
	flow.RunNet(net)

	<-net.Ready()
	for i := 0; i < b.N; i++ {

		in <- "hello"
		<-out
	}

	close(in)
	<-net.Wait()

}

func BenchmarkGraph4(b *testing.B) {
	net := newGraph()
	in := make(chan string)
	// ctx := make(chan context.IContext)
	out := make(chan string)

	net.SetInPort("In", in)
	// net.SetInPort("Ctx", ctx)
	net.SetOutPort("Out", out)
	flow.RunNet(net)

	<-net.Ready()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {

			in <- "hello"
			<-out

		}
	})
	close(in)
	<-net.Wait()
}

func TestContextManager2(t *testing.T) {
	net := newContextGraph()

	in := make(chan Context)
	out := make(chan Context)

	net.SetInPort("In", in)
	net.SetOutPort("Out", out)

	flow.RunNet(net)
	<-net.Ready()
	// log.Printf("%# v", pretty.Formatter(net))

	ctx := NewContext()
	in <- ctx
	time.Sleep(1 * time.Millisecond)
	in <- ctx
	go func() {
		time.Sleep(100 * time.Millisecond)
		close(in)
	}()
	// net.Stop()
	// log.Printf("%# v", pretty.Formatter(net))
	for {
		select {
		case ctx := <-out:
			log.Printf("out: %#v", ctx)
		case <-net.Wait():
			goto Exit
		}
	}
Exit:
}

func BenchmarkContextGraph(b *testing.B) {

	net := newContextGraph()

	in := make(chan Context)
	out := make(chan Context)

	net.SetInPort("In", in)
	net.SetOutPort("Out", out)

	flow.RunNet(net)
	<-net.Ready()

	for i := 0; i < b.N; i++ {
		ctx := NewContext()
		in <- ctx
		time.Sleep(1 * time.Millisecond)
		in <- ctx

		b.Logf("out: %v", <-out)
	}

	close(in)
}

func BenchmarkContextGraph2(b *testing.B) {

	net := newContextGraph()

	in := make(chan Context)
	out := make(chan Context)

	net.SetInPort("In", in)
	net.SetOutPort("Out", out)

	flow.RunNet(net)
	<-net.Ready()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			ctx := NewContext()
			in <- ctx
			time.Sleep(1 * time.Millisecond)
			in <- ctx

			b.Logf("out: %v", <-out)
		}
	})

	close(in)
}

func TestContextManager3(t *testing.T) {
	net := newContextGraph2()

	in := make(chan Context)
	out := make(chan string)

	net.SetInPort("In", in)
	net.SetOutPort("Out", out)

	flow.RunNet(net)
	<-net.Ready()
	// log.Printf("%# v", pretty.Formatter(net))

	ctx := NewContext()
	in <- ctx
	time.Sleep(1 * time.Millisecond)
	in <- ctx
	time.Sleep(1 * time.Millisecond)
	ctx = NewContext()
	in <- ctx
	time.Sleep(1 * time.Millisecond)
	in <- ctx
	time.Sleep(1 * time.Millisecond)
	ctx = NewContext()
	in <- ctx
	time.Sleep(1 * time.Millisecond)
	in <- ctx

	go func() {
		time.Sleep(100 * time.Millisecond)
		close(in)
	}()
	// net.Stop()
	// log.Printf("%# v", pretty.Formatter(net))
	for {
		select {
		case ctx, ok := <-out:
			if ok {
				log.Printf("out: %#v", ctx)
			} else {
				log.Printf("out close.")
			}
		case <-net.Wait():
			goto Exit
		}
	}
Exit:
}

func TestContextManager4(t *testing.T) {
	net := newContextGraph2Task()

	in := make(chan Context)
	// out1 := make(chan string)
	out := make(chan string)
	// out2 := make(chan string)

	net.SetInPort("In", in)
	net.SetOutPort("Out", out)

	flow.RunNet(net)
	<-net.Ready()

	ctx1 := NewContext()
	in <- ctx1
	time.Sleep(1 * time.Millisecond)
	ctx1.SetValue("Msg", "Hello")

	in <- ctx1

	time.Sleep(2 * time.Millisecond)
	ctx1.SetValue("Msg", "World")

	in <- ctx1

	go func() {
		time.Sleep(100 * time.Millisecond)
		close(in)
	}()

	var exit bool
	for !exit {
		select {
		case msg, ok := <-out:
			if ok {
				log.Printf("out msg: %s", msg)
				exit = true
			} else {
				goto Exit
			}

		// case msg := <-out2:
		// 	log.Printf("out2 msg: %s", msg)
		case <-net.Wait():
			goto Exit
		}
	}

Exit:
}

func ListWaits(net *testNet) {
	log.Printf("Wait List ---")
	for _, wait := range net.ListWaits() {
		log.Printf("name: %s, counter: %d", wait.Name, wait.Counter)
	}
}

func TestLuisTalk(t *testing.T) {
	net := newLuisGraph()

	in := make(chan Context)
	// out1 := make(chan string)
	out := make(chan ResultParams)
	// out2 := make(chan string)

	net.SetInPort("Ctx", in)
	net.SetOutPort("Out", out)

	flow.RunNet(net)
	<-net.Ready()
	log.Printf("net: %# v", pretty.Formatter(net))

	ctx := NewContext()
	in <- ctx

	<-net.Wait()
}
