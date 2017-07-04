package builtin

import (
	. "github.com/wanliu/flow/context"
	"log"
)

type Praise struct {
	TryGetEntities
	Ctx  <-chan Context
	Type <-chan string
	Out  chan<- ReplyData
}

func NewPraise() interface{} {
	return new(Praise)
}

func (order *Praise) OnCtx(ctx Context) {
	entities := ctx.Value("Result").(ResultParams).Entities
	output := ""

	replies := map[string]string{
		"女汉子": "你说的我都不好意思了",
		"漂亮":  "你说的我的脸都红了",
		"打满分": "谢谢你的肯定，我一定不会让你失望",
		"聪明":  "都是受你的影响",
		"崇拜":  "啊，我可承受不起呢",
		"牛逼":  "呵呵,多谢夸奖",
		"周到":  "能让你满意，我很高兴",
	}

	for _, e := range entities {
		if e.Type == "称赞" {
			r, hasKey := replies[e.Entity]

			if hasKey {
				output = r
				break
			}
		}
	}
	// "default":

	if output == "" {
		output = "谢谢夸奖, 真是不敢当"
	}

	log.Printf("OUTPUT: %v", output)

	replyData := ReplyData{output, ctx}
	order.Out <- replyData
}
