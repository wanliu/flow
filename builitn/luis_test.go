package builitn

import (
	"log"
	"strings"
	"testing"

	flow "github.com/wanliu/goflow"
)

type testLuisNet struct {
	flow.Graph
}

func TestLuis(t *testing.T) {
	net := new(testLuisNet)
	net.InitGraphState()
	// cm := NewContextManager()
	luis := new(LuisAnalyze)
	stringifier := new(Stringifier)
	// net.Add(cm, "ContextManager")
	net.Add(luis, "LuisAnalyze")
	net.Add(stringifier, "Stringifier")

	// if failed {
	// 	t.Fatalf("asdfasdf'", ...)
	// } else {
	// 	t.Skip()
	// }

	net.Connect("LuisAnalyze", "Out", "Stringifier", "In")

	net.MapInPort("Input", "LuisAnalyze", "In")
	// net.MapInPort("Ctx", "ContextManager", "Ctx")
	net.MapOutPort("Output", "Stringifier", "Out")
	// net.MapOutPort("Error", "ReadLine", "Error")

	input := make(chan string)
	// ctx := make(chan context.IContext)
	output := make(chan string)
	// err := make(chan error)

	net.SetInPort("Input", input)
	// net.SetInPort("Ctx", ctx)
	net.SetOutPort("Output", output)
	// net.SetOutPort("Error", err)

	flow.RunNet(net)

	<-net.Ready()
	log.Printf("running net")

	input <- "纯牛奶１２件"

	close(input)

	for {
		select {
		case msg := <-output:
			if !strings.Contains(msg, "intents") {
				t.Fatalf("failed")
			} else {
				goto Exit
			}
		case <-net.Wait():
			goto Exit
		}
	}
Exit:
}
