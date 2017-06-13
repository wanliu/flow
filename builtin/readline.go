package builtin

import (
	"fmt"
	"log"
	"strings"

	flow "github.com/wanliu/goflow"
)

type ReadInput struct {
	flow.Component
	_prompt string
	Prompt  <-chan string
	Out     chan<- string
}

// func (r *ReadInput) OnPrompt(pro string) {

// 	r._prompt = pro
// }

func (r *ReadInput) Loop() {
	for {
		select {
		// Handle immediate terminate signal from network
		case <-r.Component.Term:
			return
		case input, ok := <-r.GetLine():
			if ok {
				if strings.Trim(input, " \n") != "" {
					r.Out <- input
				}
			}
		case prompt, _ := <-r.Prompt:
			r._prompt = prompt
		}
	}
}

func (r *ReadInput) GetLine() <-chan string {
	ch := make(chan string)

	go func() {
		for {
			var in string
			fmt.Printf(r._prompt)
			_, err := fmt.Scanf("%q", &in)
			if err != nil {
				log.Printf("err: %s", err)
			}
			log.Printf("%s", in)
			ch <- in
		}
	}()
	return ch
}
