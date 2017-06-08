package builitn

import (
	"log"
	"testing"

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

type Component2 struct {
	ContextComponent
	In  <-chan string
	Out chan<- string
}

// func (c *Component1) OnIn(_ string) {
// 	c.Term <- struct{}{}
// }

func (c *Component2) OnIn(msg string) {

	log.Printf("component OnIn Msg: %s", msg)
	log.Printf("component: %#p", c)
	c.Out <- msg

}

func TestContextManager(t *testing.T) {
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

	in := make(chan string)
	// ctx := make(chan context.IContext)
	out := make(chan string)

	net.SetInPort("In", in)
	// net.SetInPort("Ctx", ctx)
	net.SetOutPort("Out", out)

	log.Printf("net: %# v", pretty.Formatter(net))
	flow.RunNet(net)
	<-net.Ready()
	log.Printf("running net")

	for i := 0; i < 3; i++ {

		// c, err := context.NewContext()
		// if err != nil {
		// 	log.Fatalf("create context failed: %s", err)
		// }

		in <- "hello"
		// ctx <- c

		//

		<-out
	}

	close(in)
	<-net.Wait()

}
