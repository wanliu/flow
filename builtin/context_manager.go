package builtin

import (
	"sync"
	"time"

	flow "github.com/wanliu/goflow"
)

type SendHandle func(ctx, parent Context) error

type ContextManager struct {
	flow.Component
	Ctx        <-chan Context
	Process    chan<- Context
	SendHandle SendHandle
	// Signal chan<- Context
}

func (cm *ContextManager) Init() {
	if cm.SendHandle == nil {
		cm.SendHandle = func(ctx, _ Context) error {
			ctx.Send("Hello")
			return nil
		}
	}
}

func (cm *ContextManager) OnCtx(ctx Context) {
	if child := ctx.Peek(); child != nil {
		cm.SendHandle(child, ctx)
	} else {
		cm.Process <- ctx
	}
}

type ContextComponent struct {
	flow.Component
	// 任务恢复信号
	// Resume <-chan context.IContext
	// 任务的首次触发端口
	Enter <-chan Context
	// 任务触发执行端口
	Process chan<- Context
	// 任务完成端口
	Next chan<- Context

	TaskHandle TaskHandle
}

type ContextString struct {
	flow.Component
	fieldName string
	Field     <-chan string
	Ctx       <-chan Context
	Next      chan<- Context
	Out       chan<- string
}

type ContextInt struct {
	flow.Component
	fieldName string

	Field <-chan string
	Ctx   <-chan Context
	Next  chan<- Context
	Out   chan<- int
}

type ContextBool struct {
	flow.Component
	fieldName string
	Field     <-chan string
	Ctx       <-chan Context
	Next      chan<- Context
	Out       chan<- bool
}

type CtxControl struct {
	flow.Component
	sync.Once
	cond    *sync.Cond
	Ctx     <-chan Context
	Cancel  <-chan bool
	Done    <-chan bool
	Retry   <-chan bool
	_done   bool
	_retry  bool
	_cancel bool
	ctx     Context
}

func (cc *ContextComponent) Init() {
	if cc.TaskHandle == nil {
		cc.TaskHandle = func(ctx Context, raw interface{}) error {
			if cc.IsRunning {
				cc.Process <- ctx
			}

			return nil
		}
	}
	// cc.Mode = flow.ComponentModeSync
}

func (cc *ContextComponent) OnEnter(ctx Context) {
	childCtx := ctx.NewContext()
	ctx.Push(childCtx)

	go func(task Context) {
		task.Wait(cc.TaskHandle)
		ctx.Pop()
		cc.Next <- ctx

	}(childCtx)
}

func (cc *CtxControl) run() {
	go cc.Do(func() {
		var (
			l sync.Mutex
		)
		cc.cond = sync.NewCond(&l)

		cc.cond.L.Lock()
		for !cc.condition() {
			cc.cond.Wait()
		}

		switch {
		case cc._retry:
			cc.ctx.Retry()
		case cc._cancel:
			cc.ctx.Cancel()
		case cc._done:
			cc.ctx.Done()
		}

		cc.Reset()

		cc.cond.L.Unlock()
	})
	time.Sleep(1 * time.Millisecond)
}

func (cc *CtxControl) Reset() {
	cc.ctx = nil
	cc._done = false
	cc._cancel = false
	cc._retry = false
}

func (cc *CtxControl) condition() bool {
	return cc.ctx != nil && (cc._retry || cc._done || cc._cancel)
}

func (cc *CtxControl) OnCtx(ctx Context) {
	cc.run()
	cc.ctx = ctx
	cc.cond.Signal()
}

func (cc *CtxControl) OnCancel(do bool) {
	cc.run()
	cc._cancel = true
	cc.cond.Signal()
}

func (cc *CtxControl) OnRetry(do bool) {
	cc.run()
	cc._retry = true
	cc.cond.Signal()
}

func (cc *CtxControl) OnDone(do bool) {
	cc.run()
	cc._done = true
	cc.cond.Signal()
}

func (cv *ContextString) OnCtx(ctx Context) {
	if str, ok := ctx.GlobalValue(cv.fieldName).(string); ok {
		cv.Out <- str
		cv.Next <- ctx
	} else {
		// error
	}
}

func (cv *ContextString) OnField(name string) {
	cv.fieldName = name
}

func (cv *ContextInt) OnCtx(ctx Context) {
	if i, ok := ctx.GlobalValue(cv.fieldName).(int); ok {
		cv.Next <- ctx
		cv.Out <- i
	} else {
		// error
	}
}

func (cv *ContextInt) OnField(name string) {
	cv.fieldName = name
}

func (cv *ContextBool) OnCtx(ctx Context) {
	if b, ok := ctx.GlobalValue(cv.fieldName).(bool); ok {
		cv.Next <- ctx
		cv.Out <- b
	} else {
		// error
	}
}

func (cv *ContextBool) OnField(name string) {
	cv.fieldName = name
}

func NewContextManager() *ContextManager {
	return &ContextManager{}
}
