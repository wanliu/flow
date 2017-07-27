package builtin

import (
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
		ctx, cok := q.Value("Ctx").(Context)
		txt, tok := q.Value("Text").(string)

		if cok && tok {
			ctx.SetValue("Text", txt)

			go func() {
				q.Out <- ctx
			}()
		}

		return nil
	}
}

func (q *VoiceSave) OnCtx(ctx Context) {
	q.SetValue("Ctx", ctx)
}

func (q *VoiceSave) OnText(text string) {
	q.SetValue("Text", text)
}
