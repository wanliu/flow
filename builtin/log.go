package builtin

import (
	"log"

	flow "github.com/wanliu/goflow"
)

type Log struct {
	flow.Component
	In <-chan string
}

func (l *Log) OnIn(msg string) {
	log.Printf("%s", msg)
}
