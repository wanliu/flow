package builtin

import (
	flow "github.com/wanliu/goflow"
)

func NewVoice() interface{} {
	return new(Voice)
}

type Voice struct {
	flow.Component

	token string

	Token <-chan string
	In    <-chan string
	Out   chan<- ReplyData
}

// NOOP
func (s *Voice) OnIn(input string) {
}

func (s *Voice) OnToken(t string) {
	s.token = t
}
