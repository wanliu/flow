package builitn

import (
	"log"
	// "strings"
	// "github.com/kr/pretty"
	"testing"
	// "time"

	flow "github.com/wanliu/goflow"
)

type testChannelNet struct {
	flow.Graph
}

func TestChannel(t *testing.T) {
	net := new(testChannelNet)
	net.InitGraphState()

	luis := new(TextReader)
	stringifier := new(Stringifier)
	stringifierB := new(StringifierB)
	ot := new(Output)

	net.Add(ot, "Display")
	net.Add(luis, "TextReader")
	net.Add(stringifier, "Stringifier")
	net.Add(stringifierB, "StringifierB")

	net.Connect("TextReader", "Out", "Stringifier", "In")
	net.Connect("Stringifier", "Out", "Display", "In")
	net.Connect("TextReader", "Out", "StringifierB", "In")
	net.Connect("StringifierB", "Out", "Display", "In")

	net.MapInPort("Input", "TextReader", "In")
	net.MapOutPort("Output", "Display", "Out")

	input := make(chan string)
	output := make(chan string)

	net.SetInPort("Input", input)
	net.SetOutPort("Output", output)

	// log.Printf("net: %# v", pretty.Formatter(net))
	flow.RunNet(net)

	<-net.Ready()
	log.Printf("running net")

	input <- "纯牛奶12件"
	// input <- "优酸乳20件"

	close(input)

	for {
		select {
		case <-output:
			// log.Print("================")
			// default:
			// 	time.Sleep(time.Second * 2)
		}
	}
}
