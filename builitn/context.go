package builitn

import (
	"fmt"
)

type ctxt struct {
	Stack     *Stack
	values    map[interface{}]interface{}
	retry     chan bool
	cancel    chan bool
	done      chan bool
	rece      chan interface{}
	isRunning bool
}

type Context interface {
	NewContext() Context
	Peek() Context
	Push(Context)
	Pop() Context
	Wait(TaskHandle)
	Retry()
	Cancel()
	Done()
	Send(interface{})
	Value(interface{}) interface{}
	SetValue(interface{}, interface{})
}

var (
	ErrRetry  = fmt.Errorf("Task Retry")
	ErrCancel = fmt.Errorf("Task Cancel")
)

type TaskHandle func(ctx Context, raw interface{}) error

type Stack struct {
	Root     *ctxt
	Children []Context
}

func NewStack(root *ctxt) *Stack {
	return &Stack{
		Root:     root,
		Children: make([]Context, 0),
	}
}

func (s *Stack) Peek() Context {
	if len(s.Children) == 0 {
		return nil
	} else {
		return s.Children[len(s.Children)-1]
	}
}

func (s *Stack) Push(ctx Context) {
	s.Children = append(s.Children, ctx)
}

func (s *Stack) Pop() Context {
	if len(s.Children) == 0 {
		return nil
	} else {
		ctx := s.Peek()
		s.Children = s.Children[:len(s.Children)-1]
		return ctx
	}
}

func NewContext() *ctxt {
	root := &ctxt{
		values: make(map[interface{}]interface{}),
		retry:  make(chan bool),
		cancel: make(chan bool),
		rece:   make(chan interface{}),
		done:   make(chan bool),
	}

	root.Stack = NewStack(root)
	return root
}

func (ctx *ctxt) NewContext() Context {
	childCtx := NewContext()
	childCtx.Stack = ctx.Stack

	return childCtx
}

func (ctx *ctxt) Peek() Context {

	return ctx.Stack.Peek()
}

func (ctx *ctxt) Push(cc Context) {
	ctx.Stack.Push(cc)
}

func (ctx *ctxt) Pop() Context {
	return ctx.Stack.Pop()
}

func (ctx *ctxt) Value(name interface{}) interface{} {
	return ctx.values[name]
}

func (ctx *ctxt) SetValue(name, value interface{}) {
	ctx.values[name] = value
}

func (ctx *ctxt) Wait(task TaskHandle) {
	if ctx.isRunning {
		return
	}

	ctx.isRunning = true
	for ctx.isRunning {
		select {
		case info := <-ctx.rece: // do process
			task(ctx, info)
		case <-ctx.retry:
		case <-ctx.done:
			ctx.isRunning = false
			break
		case <-ctx.cancel:
			ctx.isRunning = false
			break
		}
	}
}

func (ctx *ctxt) Cancel() {
	if ctx.isRunning {
		ctx.cancel <- true
	}
}

func (ctx *ctxt) Retry() {
	if ctx.isRunning {
		ctx.retry <- true
	}
}

func (ctx *ctxt) Done() {
	if ctx.isRunning {
		ctx.done <- true
	}
}

func (ctx *ctxt) Send(raw interface{}) {
	if ctx.isRunning {
		ctx.rece <- raw
	}
}
