package builtin

import (
	"log"
	"testing"

	flow "github.com/wanliu/goflow"
)

type testComponentsNet struct {
	flow.Graph
}

func TestComponents(t *testing.T) {
	net := new(testComponentsNet)
	net.InitGraphState()
	// cm := NewContextManager()
	rl := new(ReadLine)
	ot := new(Output)
	// net.Add(cm, "ContextManager")
	net.Add(rl, "ReadLine")
	net.Add(ot, "Display")

	net.Connect("ReadLine", "Out", "Display", "In")

	net.MapInPort("FileName", "ReadLine", "In")
	// net.MapInPort("Ctx", "ContextManager", "Ctx")
	net.MapOutPort("Out", "Display", "Out")
	// net.MapOutPort("Error", "ReadLine", "Error")

	file := make(chan string)
	// ctx := make(chan context.IContext)
	out := make(chan string)
	// err := make(chan error)

	net.SetInPort("FileName", file)
	// net.SetInPort("Ctx", ctx)
	net.SetOutPort("Out", out)
	// net.SetOutPort("Error", err)

	flow.RunNet(net)
	<-net.Ready()
	log.Printf("running net")
	file <- "../test/test.txt"
	close(file)
	for {
		select {
		case <-out:
		case <-net.Wait():
			goto Exit
		}
	}
Exit:
}

func TestComponents2(t *testing.T) {
	net := new(testComponentsNet)
	net.InitGraphState()
	// cm := NewContextManager()
	rl := new(ReadLine)
	o1 := new(Output)
	o2 := new(Output)
	o3 := new(Output)
	o4 := new(Output)
	// net.Add(cm, "ContextManager")
	net.Add(rl, "ReadLine")
	net.Add(o1, "Display1")
	net.Add(o2, "Display2")
	net.Add(o3, "Display3")
	net.Add(o4, "Display4")
	net.Connect("ReadLine", "Out", "Display1", "In")
	net.Connect("Display1", "Out", "Display2", "In")
	net.Connect("Display2", "Out", "Display3", "In")
	net.Connect("Display3", "Out", "Display4", "In")

	net.MapInPort("FileName", "ReadLine", "In")
	// net.MapInPort("Ctx", "ContextManager", "Ctx")
	net.MapOutPort("Out", "Display4", "Out")
	// net.MapOutPort("Error", "ReadLine", "Error")

	file := make(chan string)
	// ctx := make(chan context.IContext)
	out := make(chan string)
	// err := make(chan error)

	net.SetInPort("FileName", file)
	// net.SetInPort("Ctx", ctx)
	net.SetOutPort("Out", out)
	// net.SetOutPort("Error", err)

	flow.RunNet(net)
	<-net.Ready()
	log.Printf("running net")
	file <- "../test/test.txt"
	close(file)
	for {
		select {
		case <-out:
		case <-net.Wait():
			goto Exit
		}
	}
Exit:
}
