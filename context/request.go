package context

import (
	"github.com/hysios/apiai-go"
)

type Request struct {
	Ctx         Context
	RequestId   string
	Text        string
	ApiAiResult apiai.Result
}
