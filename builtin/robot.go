package builtin

import (
	. "github.com/wanliu/flow/builtin/luis"
	. "github.com/wanliu/flow/context"
)

type Robot struct {
	TryGetEntities
	Ctx  <-chan Context
	Type <-chan string
	Out  chan<- ReplyData
}

func NewRobot() interface{} {
	return new(Robot)
}

func (order *Robot) OnCtx(ctx Context) {
	entities := ctx.Value("Result").(ResultParams).Entities
	output := ""

	entiDefaults := map[string]string{
		"身份":   "你是要问我吗？",
		"生理":   "我是虚拟的，这是没有意义的",
		"亲缘":   "我也不知道有没有",
		"亲密关系": "没有哦~",
	}

	replies := map[string][]string{
		"我是智能机器人，我的名字叫小花, 我能为您提供购物服务": []string{"姓名", "谁", "名字", "小名"},
		"我的年龄保密。":                     []string{"年龄", "多大", "年纪"},
		"45公分":                        []string{"身高", "多高", "腿长"},
		"40kg":                        []string{"体重"},
		"48 48 48":                    []string{"三围", "身材", "胸大", "腰围", "腰细", "臀围"},
		"人家不会告诉你，我是女孩纸啦": []string{"男", "女", "男孩", "男人", "男的", "男生", "女孩",
			"女人", "女的", "女生", "是男是女", "性别"},
		"这不重要，你告诉我工作表现就好": []string{"帅"},
		"我是女的":            []string{"人妖"},
		"处女座":             []string{"星座"},
		"我不能结婚， 我已经嫁给智能客服事业": []string{"结婚", "婚姻"},
		"陪人聊天":               []string{"爱好"},
		"我已经献身于智能客服事业":       []string{"男朋友"},
		"怎么，你打算给我介绍一个":       []string{"对象"},
		"人家不爱这个啦":            []string{"女朋友"},
		"捷杰科技的开发团队":          []string{"爸妈", "爸爸", "妈妈", "父母", "父亲", "母亲"},
		"没有哦":                []string{"哥哥", "弟弟", "表哥", "表弟", "老表", "兄弟", "姐妹", "姐姐", "妹妹", "表姐", "表妹", "闺蜜"},
		"我可离不开电，我爱充电":        []string{"爱吃什么", "爱吃什么食物", "爱吃"},
		"喜欢你":                []string{"喜欢"},
		"耒阳":                 []string{"住址"},
		"我没有私人电话，但你可以打这个号码 0731-83991490，跟我的开发者反映问题": []string{"手机", "联系电话", "电话"},
		"我会向很多智能语音的前辈们学习(Siri, Alex, 微软小冰)":          []string{"偶像"},
		"中华人民共和国":                                    []string{"那个国家", "国家", "国籍"},
		"湖南":                                         []string{"那里", "哪里", "祖籍"},
		"2017-05-01":                                 []string{"出生", "生日"},
		"我觉得还行吧":                                     []string{"好看", "好靓", "美女", "漂亮"},
	}

	found := false
	for _, e := range entities {
		for reply, enti := range replies {
			for _, item := range enti {
				if e.Entity == item {
					output = reply
					found = true
					break
				}
			}

			if found {
				break
			}
		}

		if found {
			break
		}
	}

	found = false
	if output == "" {
		for _, e := range entities {
			for t, d := range entiDefaults {
				if e.Type == t {
					output = d
					found = true
					break
				}
			}

			if found {
				break
			}
		}
	}

	if output == "" {
		output = "我是服务大家的机器人小花"
	}

	replyData := ReplyData{output, ctx, nil}
	order.Out <- replyData
}
