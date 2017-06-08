package builitn

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
