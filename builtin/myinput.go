package builtin

import (
	"strings"

	. "github.com/wanliu/flow/context"
)

type MyInput struct {
	ReadInput
	_ctx Context
	Ctx  <-chan Context
	Next chan<- Context
}

func NewMyInput() interface{} {
	return new(MyInput)
}

func (in *MyInput) Loop() {
	for {
		select {
		// Handle immediate terminate signal from network
		case <-in.Component.Term:
			return
		case input, ok := <-in.GetLine():
			if ok {
				if strings.Trim(input, " \n") != "" {
					in.Out <- input
					in.Next <- in._ctx
				}
			}
		case ctx := <-in.Ctx:
			in._ctx = ctx
		case prompt, _ := <-in.Prompt:
			in.SetPrompt(prompt)
		}
	}
}
