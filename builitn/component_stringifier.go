package builitn

import (
	// "encoding/json"
	flow "github.com/wanliu/goflow"
	"log"
)

func NewTextReader() interface{} {
	return new(TextReader)
}

type TextReader struct {
	flow.Component
	In  <-chan string
	Out chan<- string
}

func (s *TextReader) OnIn(input string) {
	s.Out <- input
}

func NewStringifier() interface{} {
	return new(Stringifier)
}

type Stringifier struct {
	flow.Component
	In  <-chan string
	Out chan<- string
}

func (s *Stringifier) OnIn(input string) {
	log.Printf("received from A")
	s.Out <- input + "from A"
	// str, _ := json.Marshal(input)

	// s.Out <- string(str)
}

func NewStringifierB() interface{} {
	return new(StringifierB)
}

type StringifierB struct {
	flow.Component
	In  <-chan string
	Out chan<- string
}

func (s *StringifierB) OnIn(input string) {
	log.Printf("received from B")
	s.Out <- input + " from B"
	// str, _ := json.Marshal(input)

	// s.Out <- string(str)
}
