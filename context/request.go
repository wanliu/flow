package context

import (
	"github.com/hysios/apiai-go"
)

type Request struct {
	Ctx         Context
	Id          string
	Text        string
	ApiAiResult apiai.Result
}
