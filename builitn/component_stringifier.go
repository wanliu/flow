package builitn

import (
	"encoding/json"
	flow "github.com/wanliu/goflow"
)

func NewStringifier() interface{} {
	return new(Stringifier)
}

type Stringifier struct {
	flow.Component
	In  <-chan ResultParams
	Out chan<- string
}

func (s *Stringifier) OnIn(input ResultParams) {
	str, _ := json.Marshal(input)

	s.Out <- string(str)
}
