package builtin

import (
	"log"
	"strconv"
	"time"

	"github.com/wanliu/flow/builtin/config"
	"github.com/wanliu/flow/builtin/resolves"

	"github.com/wanliu/flow/context"
)

const (
	Per = 5
)

type CustomerOrders struct {
	TryGetEntities
	Type <-chan string

	expireMins int

	Ctx       <-chan context.Context
	Page      <-chan context.Context
	ExpireMin <-chan interface{}

	Out chan<- ReplyData
}

func NewCustomerOrders() interface{} {
	return new(CustomerOrders)
}

func (c *CustomerOrders) OnCtx(ctx context.Context) {
	// if context.GroupChat(ctx) {
	// 	log.Printf("不回应非开单相关的普通群聊")
	// 	return
	// }

	rsv := resolves.NewCusOrdersResolve(ctx, Per)

	reply := rsv.Answer()

	rsv.Setup(ctx)
	c.ResetTick(rsv, ctx)

	c.Out <- ReplyData{reply, ctx, nil}
}

func (c *CustomerOrders) OnPage(ctx context.Context) {
	in := ctx.CtxValue(config.CtxKeyCusOrders)
	if in == nil {
		if context.GroupChat(ctx) {
			log.Printf("不回应无context的非@翻页")
			return
		}

		c.Out <- ReplyData{"当前没有正在进行的查询", ctx, nil}
	} else {
		rsv := in.(*resolves.CustomerOrdersResolve)

		c.ResetTick(rsv, ctx)
		reply := rsv.Answer()
		c.Out <- ReplyData{reply, ctx, nil}

		rsv.ClearIfDone(ctx)
	}
}

func (c *CustomerOrders) OnExpireMin(mins interface{}) {
	switch mins.(type) {
	case float64:
		c.expireMins = int(mins.(float64))
	case string:
		s := mins.(string)
		f, err := strconv.ParseFloat(s, 64)

		if err == nil {
			c.expireMins = int(f)
		}
	}
}

func (c *CustomerOrders) ResetTick(rsv *resolves.CustomerOrdersResolve, ctx context.Context) {
	now := time.Now()
	rsv.LastActiveTime = &now

	go func() {
		expiredMins := config.SesssionExpiredMinutes
		settedMins := c.expireMins

		if settedMins != 0 {
			expiredMins = settedMins
		}

		time.Sleep(time.Duration(expiredMins) * time.Minute)

		now := time.Now()

		if rsv != nil {
			t := rsv.LastActiveTime.Add(time.Duration(expiredMins) * time.Minute)

			if now.After(t) {
				rsv.Clear(ctx)
				log.Printf("[Notify]超过过期时间，清除客户订单查询任务")
			} else {
				log.Printf("[Notify]客户订单查询任务延期")
			}
		}
	}()
}
