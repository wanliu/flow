package context

import (
	"sync"

	"fmt"
	"time"
)

type ErrorHandle func(error) error
type CtxOptFunc func(opt *ContextOption) error
type TaskHandle func(ctx Context, raw interface{}) error

type ctxt struct {
	sync.RWMutex
	Stack *Stack
	Reply Replyer

	values map[interface{}]interface{}
	retry  chan bool
	cancel chan bool
	done   chan bool
	rece   chan interface{}

	// 0 for text
	// 1 for text and data
	postMode int

	send      chan string
	sendData  chan interface{}
	sendTable chan *Table
	quit      chan bool

	isRunning bool
	counter   int

	errHandle      ErrorHandle
	replyErrHandle ErrorHandle
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
	CtxValue(interface{}) interface{}
	SetCtxValue(interface{}, interface{})
	GlobalValue(interface{}) interface{}
	SetGlobalValue(interface{}, interface{})

	SetPostMode(int)
	PostMode() int

	Post(string, ...interface{}) error
	PostTable(Table) error

	Run()
	RunCallback(handler ContextReplyHander)
	Close()

	Reset()
	IsRunning() bool
}

var (
	ErrRetry  = fmt.Errorf("Task Retry")
	ErrCancel = fmt.Errorf("Task Cancel")
)

type ContextOption struct {
	Reply Replyer
	Error ErrorHandle
}

func NewContext(args ...CtxOptFunc) (*ctxt, error) {
	opt, err := ctxOption(args)
	if err != nil {
		return nil, err
	}

	root := &ctxt{
		values:    make(map[interface{}]interface{}),
		retry:     make(chan bool),
		cancel:    make(chan bool),
		rece:      make(chan interface{}),
		done:      make(chan bool),
		send:      make(chan string),
		sendData:  make(chan interface{}),
		sendTable: make(chan *Table),
		quit:      make(chan bool),
		errHandle: opt.Error,
		Reply:     opt.Reply,
	}

	root.Stack = NewStack(root)
	return root, nil
}

func ctxOption(args []CtxOptFunc) (*ContextOption, error) {
	var opt = ContextOption{
		Reply: StdoutReply,
	}

	for _, f := range args {
		if err := f(&opt); err != nil {
			return nil, err
		}
	}

	return &opt, nil
}

func UseReply(reply Replyer) CtxOptFunc {
	return func(opt *ContextOption) error {
		opt.Reply = reply
		return nil
	}
}

func OnError(errHandle ErrorHandle) CtxOptFunc {
	return func(opt *ContextOption) error {
		opt.Error = errHandle
		return nil
	}
}

func (ctx *ctxt) NewContext() Context {
	childCtx, _ := NewContext()
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

// TODO 混合私聊和群聊情况，私聊可以去读群聊，而群聊不可以读取私聊信息？
func (ctx *ctxt) CtxValue(name interface{}) interface{} {
	if GroupChat(ctx) {
		name = fmt.Sprintf("Group:%v", name)
	}

	ctx.RLock()
	defer ctx.RUnlock()
	return ctx.values[name]
}

func (ctx *ctxt) SetCtxValue(name, value interface{}) {
	if GroupChat(ctx) {
		name = fmt.Sprintf("Group:%v", name)
	}

	ctx.Lock()
	defer ctx.Unlock()
	ctx.values[name] = value
}

func (ctx *ctxt) Value(name interface{}) interface{} {
	ctx.RLock()
	defer ctx.RUnlock()

	return ctx.values[name]
}

func (ctx *ctxt) SetValue(name, value interface{}) {
	ctx.Lock()
	defer ctx.Unlock()

	ctx.values[name] = value
}

func (ctx *ctxt) GlobalValue(name interface{}) interface{} {
	return ctx.Stack.Root.Value(name)
	// return ctx.Stack.Root.values[name]
}

func (ctx *ctxt) SetGlobalValue(name, value interface{}) {
	ctx.Stack.Root.SetValue(name, value)
	// ctx.Stack.Root.values[name] = value
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

// func (ctx *ctxt) Post(msg string, args ...interface{}) error {
// 	// ctx.Stack.
// 	log.Printf(msg, args...)
// 	return nil
// }

func (ctx *ctxt) SetPostMode(mode int) {
	ctx.postMode = mode
}

func (ctx ctxt) PostMode() int {
	return ctx.postMode
}

func (ctx *ctxt) Post(text string, args ...interface{}) error {
	ctx.counter++

	if 0 == ctx.PostMode() {
		// ctx.Stack.Root.send <- fmt.Sprintf(text, args...)
		ctx.Stack.Root.send <- text
	} else {
		r := map[string]interface{}{
			"reply": text,
		}

		if len(args) > 0 {
			r["data"] = args[0]
		}

		ctx.Stack.Root.sendData <- r
	}

	return nil
}

func (ctx *ctxt) PostTable(table Table) error {
	ctx.counter++
	ctx.Stack.Root.sendTable <- &table
	return nil
}

func (ctx *ctxt) Run() {
	var err error
LOOP:
	for {
		select {
		case txt := <-ctx.send:
			ctx.counter--
			err = ctx.Reply.Text(txt, ctx)
		case _ = <-ctx.sendData:
			ctx.counter--
			err = ctx.Reply.Text("unimplimented data channel", ctx)
		case table := <-ctx.sendTable:
			ctx.counter--
			err = ctx.Reply.Table(table, ctx)
		case <-ctx.quit:
			// ctx.waitingEnd()
			break LOOP
		}
		ctx.doReplyError(err)
	}
}

type ContextReplyHander func(txt *string, table *Table, data interface{})

func (ctx *ctxt) RunCallback(handler ContextReplyHander) {
	var err error
LOOP:
	for {
		select {
		case txt := <-ctx.send:
			ctx.counter--
			handler(&txt, nil, nil)
		case data := <-ctx.sendData:
			ctx.counter--
			handler(nil, nil, data)
		case table := <-ctx.sendTable:
			ctx.counter--
			handler(nil, table, nil)
		case <-ctx.quit:
			// ctx.waitingEnd()
			break LOOP
		}
		ctx.doReplyError(err)
	}
}

func (ctx *ctxt) IsRunning() bool {
	return ctx.isRunning
}

func (ctx *ctxt) Reset() {
	for _, child := range ctx.Stack.Children {
		if child.IsRunning() {
			go child.Cancel()
		}
	}

	if ctx.Stack.Root.IsRunning() {
		go ctx.Stack.Root.Cancel()
	}
	ctx.Lock()
	defer ctx.Unlock()
	root := ctx.Stack.Root
	root.values = make(map[interface{}]interface{})
	root.retry = make(chan bool)
	root.cancel = make(chan bool)
	root.rece = make(chan interface{})
	root.done = make(chan bool)
	root.isRunning = false
	root.Stack.Children = make([]Context, 0)
}

func (ctx *ctxt) Close() {
	var (
		tick = 1
		tt   = time.NewTimer(15 * time.Second)
		exit bool
	)

	for !exit {
		select {
		case <-time.After(time.Duration(tick) * time.Millisecond):
			if ctx.counter > 0 {
				tick *= 2
			} else {
				exit = true
				break
			}
		case <-tt.C:
			exit = true
			break
		}
	}
	ctx.quit <- true
}

func (ctx *ctxt) OnError(handle ErrorHandle) {
	ctx.errHandle = handle
}

func (ctx *ctxt) OnReplyError(handle ErrorHandle) {
	ctx.replyErrHandle = handle
}

func (ctx *ctxt) doError(err error) error {
	if ctx.errHandle != nil {
		return ctx.errHandle(err)
	}
	return err
}

func (ctx *ctxt) doReplyError(err error) error {
	if ctx.errHandle != nil {
		return ctx.errHandle(err)
	}
	return err
}
