package builtin

import (
	"regexp"

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
