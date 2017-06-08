package builitn

import (
	"log"

	"github.com/wanliu/context"
	flow "github.com/wanliu/goflow"
)

type ContextManager struct {
	flow.Component
	Ctx     <-chan context.IContext
	Query   <-chan interface{}
	manager *__contextManager
}

type __contextManager struct {
	childs []Contexter
}

type Contexter interface {
	Send(interface{})
}

func (cm *ContextManager) OnCtx(_ context.IContext) {
	// cm.current = ctx
	cm.Term <- struct{}{}
}

// func (cm *ContextManager) OnQuery(query interface{}) {
// 	cm.current
// }

type ContextComponent struct {
	flow.Component
	Ctx <-chan context.IContext
}

func (cc *ContextComponent) Init() {
	log.Printf("ContextComponent Init %#p", cc)
	// cc.Ctx = make(chan context.IContext)
}

func NewContextManager() *ContextManager {
	return &ContextManager{
		manager: new(__contextManager),
	}
}

// func (cc *ContextComponent) RegisterChild(cc *ContextComponent) {

// }

// func (cc *ContextComponent) GetManager() *ContextManager {
// 	ctxm := cc.Net.Get("__context_manager")
// 	if cm, ok := ctxm.(*ContextManager); ok {
// 		return cm
// 	}

// 	return nil
// }
