package builtin

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	flow "github.com/wanliu/goflow"
)

var REPLACE_DICT map[string]string = map[string]string{
	"今早":    "今天早上",
	"今下":    "今天下午",
	"明早":    "明天早上",
	"明下":    "明天下午",
	"今早上":   "今天早上",
	"今下午":   "今天下午",
	"明早上":   "明天早上",
	"明下午":   "明天下午",
	"1.1红":  "1100红",
	"1.1原":  "1100原",
	"上次":    "上一次",
	",":     "，",
	"16Q":   "十六Q",
	"16QQ":  "十六QQ",
	"190Q":  "一九零Q",
	"190QQ": "一九零QQ",
	"200Q":  "二百Q",
	"1l":    "一L",
	"毫升":    "ml",
}

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
	output := atFilter(input)
	output = replaceDeliver(output)
	output = numberAfterLetter(output)
	output = dateTransfer(output)
	output = dictTransfer(output)
	output = replaceUnit(output)
	c.Out <- output
}

func replaceDeliver(s string) string {
	r := regexp.MustCompile(`急送[^，,.。]`)

	is := r.FindStringIndex(s)

	for len(is) == 2 {
		b := is[0] + 6
		s = s[:b] + "," + s[b:]

		is = r.FindStringIndex(s)
	}

	return s
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

// [星期日,星期天,周日]
func dateTransfer(s string) string {
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

func dictTransfer(s string) string {
	for k, v := range REPLACE_DICT {
		s = strings.Replace(s, k, v, -1)
	}

	return s
}

func atFilter(s string) string {
	r := regexp.MustCompile("^@[\u4e00-\u9fa5\\w]+[\\s\u2005]")
	is := r.FindStringIndex(s)

	if len(is) == 2 {
		i := is[1]
		s = s[i:]
	}

	return s
}

// 件 条 个 支 => 龘件 龘条 龘个 龘支
func replaceUnit(s string) string {
	palceholder := "龘"
	units := []string{
		"件", "提", "条", "瓶", "排", "箱", "桶", "支", "袋", "个",
	}

	for _, unit := range units {
		r := regexp.MustCompile("[\\d一二两三四五六七八九十零]\\s*" + unit)
		is := r.FindStringIndex(s)
		for len(is) == 2 {
			total := is[1] - is[0]
			unitlen := len(unit)
			s = fmt.Sprintf("%v%v%v", s[:is[0]+total-unitlen], palceholder, s[is[1]-unitlen:])
			is = r.FindStringIndex(s)
		}
	}

	return s
}
