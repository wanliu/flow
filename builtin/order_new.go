package builtin

import (
	"log"

	"github.com/wanliu/flow/builtin/config"
	"github.com/wanliu/flow/builtin/resolves"

	"github.com/wanliu/flow/context"
)

type NewOrder struct {
	TryGetEntities
	DefTime    string
	retryCount int
	TablePrint bool

	Ctx        <-chan context.Context
	Deftime    <-chan string
	PrintTable <-chan string
	Out        chan<- ReplyData
	Notice     chan<- context.Context
	Timeout    chan<- context.Context

	RetryOut chan<- context.Context
	RetryIn  <-chan context.Context

	RetryCount <-chan float64
}

func NewNewOrder() interface{} {
	return new(NewOrder)
}

// 默认送货时间
func (c *NewOrder) OnDeftime(t string) {
	c.DefTime = t
}

func (c *NewOrder) OnPrintTable(prt string) {
	if prt == "true" {
		c.TablePrint = true
	}
}

func (c *NewOrder) OnRetryCount(count float64) {
	c.retryCount = int(count)
}

func (c *NewOrder) OnCtx(ctx context.Context) {
	orderResolve := resolves.NewOrderResolve(ctx, c.TablePrint)

	if c.DefTime != "" {
		orderResolve.SetDefTime(c.DefTime)
	}

	output := ""

	if orderResolve.EmptyProducts() {
		if c.retryCount > 0 {
			log.Printf("重新获取开单产品，第1次，共%v次", c.retryCount)
			c.RetryOut <- ctx
		} else {
			if context.GroupChat(ctx) {
				c.GroupAnswer(ctx, orderResolve)
				return
			}

			output = "没有相关的产品"
			replyData := ReplyData{output, ctx, nil}
			c.Out <- replyData
		}
	} else {
		output, table := orderResolve.AnswerWithTable(ctx)

		if orderResolve.Resolved() {
			ctx.SetCtxValue(config.CtxKeyLastOrder, *orderResolve)
			ctx.SetCtxValue(config.CtxKeyOrder, table)
		} else if orderResolve.Failed() {
			ctx.SetCtxValue(config.CtxKeyOrder, table)
		} else if orderResolve.MismatchQuantity() {
			ctx.SetCtxValue(config.CtxKeyOrder, table)
		} else {
			ctx.SetCtxValue(config.CtxKeyOrder, *orderResolve)
		}

		c.Timeout <- ctx

		replyData := ReplyData{output, ctx, table}
		c.Out <- replyData
	}
}

func (c *NewOrder) OnRetryIn(ctx context.Context) {
	orderResolve := resolves.NewOrderResolve(ctx, c.TablePrint)

	if c.DefTime != "" {
		orderResolve.SetDefTime(c.DefTime)
	}

	output := ""

	if orderResolve.EmptyProducts() {
		retriedCount := 1
		retriedCountInt := ctx.Value(config.CtxKeyRetriedCount)

		if retriedCountInt != nil {
			retriedCount = retriedCountInt.(int)
		}

		if retriedCount >= c.retryCount {
			if context.GroupChat(ctx) {
				c.GroupAnswer(ctx, orderResolve)
				return
			}

			output = "没有相关的产品"

			replyData := ReplyData{output, ctx, nil}
			c.Out <- replyData
		} else {
			retriedCount++
			log.Printf("重新获取开单产品，第%v次，共%v次", retriedCount, c.retryCount)

			ctx.SetValue(config.CtxKeyRetriedCount, retriedCount)
			c.RetryOut <- ctx
		}
	} else {
		output, table := orderResolve.AnswerWithTable(ctx)

		if orderResolve.Resolved() {
			ctx.SetCtxValue(config.CtxKeyLastOrder, *orderResolve)
			ctx.SetCtxValue(config.CtxKeyOrder, table)
		} else if orderResolve.Failed() {
			ctx.SetCtxValue(config.CtxKeyOrder, table)
		} else if orderResolve.MismatchQuantity() {
			ctx.SetCtxValue(config.CtxKeyOrder, table)
		} else {
			ctx.SetCtxValue(config.CtxKeyOrder, *orderResolve)
		}

		// c.Notice <- ctx
		c.Timeout <- ctx

		replyData := ReplyData{output, ctx, table}
		c.Out <- replyData
	}

}

func (c *NewOrder) GroupAnswer(ctx context.Context, orderResolve *resolves.OrderResolve) {
	str, tbl := orderResolve.AnswerWithTable(ctx)

	if orderResolve.Fulfiled() {
		replyData := ReplyData{str, ctx, tbl}
		c.Out <- replyData
	} else {
		log.Printf("群聊开单失败, 取消回复。失败原因：%v", str)
	}
}
