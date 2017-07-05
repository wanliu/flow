package builtin

import (
	. "github.com/wanliu/flow/builtin/luis"
	. "github.com/wanliu/flow/context"
	flow "github.com/wanliu/goflow"
)

type TryGetEntities struct {
	flow.Component
	_type string
	Ctx   <-chan Context
	Next  chan<- Context
	Type  <-chan string
	No    chan<- Context
}

func (tr *TryGetEntities) OnType(typ string) {
	tr._type = typ
}

func (tr *TryGetEntities) OnCtx(ctx Context) {
	if res, ok := ctx.GlobalValue("Result").(ResultParams); ok {

		for _, entity := range res.Entities {
			if entity.Type == tr._type {

			}
		}
	} else {
		tr.No <- ctx
	}
}
