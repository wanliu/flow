package builtin

import (
	"log"

	. "github.com/wanliu/flow/context"
	goflow "github.com/wanliu/goflow"
)

func NewVoiceSave() interface{} {
	return new(VoiceSave)
}

type VoiceSave struct {
	goflow.Component
	MultiField
	Text <-chan string
	Ctx  <-chan Context
	Out  chan<- Context
}

func (q *VoiceSave) Init() {
	q.Fields = []string{"Ctx", "Text"}
	q.Process = func() error {
		txt, rok := q.Value("Text").(string)
		ctx, cok := q.Value("Ctx").(Context)
		if rok && cok {
			ctx.SetValue("Text", txt)

			go func() {
				q.Out <- ctx
			}()
		}

		return nil
	}
}

func (q *VoiceSave) OnCtx(ctx Context) {
	log.Printf(".........2........")
	q.SetValue("Ctx", ctx)
}

func (q *VoiceSave) OnText(txt string) {
	log.Printf(".........1........%v", txt)
	q.SetValue("Text", txt)
}
