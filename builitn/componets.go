package builitn

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"

	flow "github.com/wanliu/goflow"
)

type GetElement struct {
	flow.Component
}

type Split struct {
	flow.Component
	In  <-chan string
	Out chan<- string
}

type Output struct {
	flow.Component
	In      <-chan string
	Options <-chan map[string]interface{}
	Out     chan<- string
}

type ReadFile struct {
	flow.Component
	Read <-chan string
	// Encoding <-chan string
	Out   chan<- string
	Error chan<- error
}

type ReadLine struct {
	flow.Component
	ReadLine <-chan string
	Out      chan<- string
	Error    chan<- error
}

func NewGetElement() interface{} {
	return new(GetElement)
}

func NewSplit() interface{} {
	return new(Split)
}

func NewOutput() interface{} {
	return new(Output)
}

func NewReadFile() interface{} {
	return new(ReadFile)
}

func (o *Output) OnIn(msg string) {
	log.Printf("output: %s", msg)
	o.Out <- msg
}

func (rf *ReadFile) Init() {
	// rf.Error = make(chan error)
}

func (rf *ReadFile) OnRead(filename string) {
	if buf, err := ioutil.ReadFile(filename); err != nil {
		rf.Error <- err
	} else {
		rf.Out <- string(buf)
	}
}

func (rl *ReadLine) OnReadLine(filename string) {
	if f, err := os.Open(filename); err != nil {
		rl.Error <- err
	}

	defer f.Close()
	reader = bufio.NewReader(f)
	rl.Out <- string(reader.ReadLine())
}

func (sp *Split) OnIn(msg string) {
	sp.Out <- msg
}
