package context

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/olekukonko/tablewriter"
)

var StdoutReply = NewReply(os.Stdout)

type Reply struct {
	Writer io.Writer
}

type TestReply struct {
	Reply
	rece chan string
}

type Replyer interface {
	Text(string) error
	Table(*Table) error
}

func NewReply(w io.Writer) *Reply {
	return &Reply{
		Writer: w,
	}
}

func (r *Reply) Text(text string) error {
	fmt.Fprintf(r.Writer, text)
	return nil
}

func (r *Reply) Table(table *Table) error {
	t := tablewriter.NewWriter(r.Writer)
	t.SetHeader(table.Headers)
	t.SetFooter(table.Footers)
	t.AppendBulk(table.Rows)

	t.Render()

	return nil
}

func NewTestReply() *TestReply {
	return &TestReply{
		rece: make(chan string),
	}
}

func (tr *TestReply) Text(text string) error {
	// tr.Reply.Text(text)
	log.Printf("Reply: %s", text)
	go func() {
		tr.rece <- text
	}()
	return nil
}

func (tr *TestReply) Table(table *Table) error {
	var buf bytes.Buffer

	tr.Writer = &buf
	tr.Reply.Table(table)
	log.Printf("Reply Table:\n%s", buf.String())
	go func() {
		tr.rece <- buf.String()
	}()
	return nil
}

func (tr *TestReply) MatchText(pattern string) (bool, error) {
	select {
	case txt := <-tr.rece:
		return regexp.MatchString(pattern, txt)
	}
	// return false, fmt.Errorf("Empty Match")
}

func (tr *TestReply) MatchTable(text string) bool {
	select {
	case txt := <-tr.rece:
		return strings.Compare(text, txt) == 0
	}
}
