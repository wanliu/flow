package main

import (
	"fmt"

	"github.com/wanliu/flow/builtin"
	"github.com/wanliu/flow/context"
)

func main() {
	list := builtin.ComponentList()
	fmt.Printf("Component length: %v", len(list))
	fmt.Print(context.GroupChat)
}
