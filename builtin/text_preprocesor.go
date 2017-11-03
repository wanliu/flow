package builtin

import (
	"regexp"
	"strings"
	"time"

	flow "github.com/wanliu/goflow"
)

type TextPreprocesor struct {
	flow.Component

	MultiField

	Out chan<- string
	In  <-chan string
}

func NewTextPreprocesor() interface{} {
	return new(TextPreprocesor)
}

func (c *TextPreprocesor) OnIn(input string) {
	output := numberAfterLetter(input)
	output = dateTransfer(output)
	c.Out <- output
}

func numberAfterLetter(s string) string {
	r := regexp.MustCompile("[a-zA-Z][0-9]")

	is := r.FindStringIndex(s)

	for len(is) == 2 {
		i := (is[0] + is[1]) / 2
		s = s[:i] + " " + s[i:]

		is = r.FindStringIndex(s)
	}

	return s
}

// 今早 => 今天早上
// 今下 => 今天下午
// 明早 => 明天早上
// 明下 => 明天下午
// 今早上 => 今天早上
// 今下午 => 今天下午
// 明早上 => 明天早上
// 明下午 => 明天下午
// 后天
// [星期日,星期天,周日]
func dateTransfer(s string) string {
	s = strings.Replace(s, "今早", "今天早上", -1)
	s = strings.Replace(s, "今下", "今天下午", -1)
	s = strings.Replace(s, "明早", "明天早上", -1)
	s = strings.Replace(s, "明下", "明天下午", -1)
	s = strings.Replace(s, "今早上", "今天早上", -1)
	s = strings.Replace(s, "今下午", "今天下午", -1)
	s = strings.Replace(s, "明早上", "明天早上", -1)
	s = strings.Replace(s, "明下午", "明天下午", -1)

	dat := time.Now().AddDate(0, 0, 2)
	datStr := dat.Format("1月2日")
	s = strings.Replace(s, "后天", datStr, -1)

	step := int(7 - time.Now().Weekday())

	if step == 7 {
		step = 0
	}

	dat = time.Now().AddDate(0, 0, step)
	datStr = dat.Format("1月2日")
	s = strings.Replace(s, "星期日", datStr, -1)
	s = strings.Replace(s, "星期天", datStr, -1)
	s = strings.Replace(s, "周日", datStr, -1)

	return s
}
