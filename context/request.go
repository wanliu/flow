package context

import (
	"github.com/hysios/apiai-go"
)

// type Response struct {
// 	Type   string
// 	On     string
// 	Action string
// 	Data   interface{}
// }
type Response struct {
	Reply string
	Ctx   Context
	Data  interface{}
}

type ResReply struct {
	Data interface{}
	Req  *Request
}

// {
// 	On: "orderItem",
// 	Action: "delete/quantity",
// 	Data: map[string]interface{}{"itemName":"伊利纯牛奶"},
// }
type RequestCommand struct {
	On     string
	Action string
	Data   map[string]interface{}
}

type Request struct {
	Ctx         Context
	Id          string
	Text        string
	ApiAiResult apiai.Result
	Res         Response
	Command     *RequestCommand
}
