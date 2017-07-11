package builtin

import (
	"fmt"
	"log"
	"time"

	"github.com/franela/goreq"

	. "github.com/wanliu/flow/builtin/luis"
	flow "github.com/wanliu/goflow"
)

func NewLuisAnalyze() interface{} {
	return new(LuisAnalyze)
}

type LuisAnalyze struct {
	flow.Component
	appid string
	key   string
	In    <-chan string
	AppId <-chan string
	Key   <-chan string
	Out   chan<- ResultParams
	// Out chan<- string
}

func (l *LuisAnalyze) OnIn(input string) {
	var (
		// luis   = NewLuis("052297dc-12b9-4044-8220-a21a20d72581", "6b916f7c107643069c242cf881609a82", "")
		luis   = NewLuis(l.appid, l.key, "")
		url    = fmt.Sprintf(luis.Url, luis.AppID)
		params = QueryParams{
			Key:      luis.Key,
			TimeZone: "0",
			Query:    input,
			Verbose:  true,
		}
		result ResultParams
	)

	ch := l.promptBegin()
	res, err := goreq.Request{
		Uri:         url,
		QueryString: params,
		Timeout:     10 * time.Second,
		CookieJar:   luis.CookieJar,
		Proxy:       luis.Proxy,
	}.Do()
	ch <- true
	if err != nil {
		// l.Out <- err.Error()
		log.Printf("luis query error %s", err)
		l.Out <- *new(ResultParams)
	} else {
		// result, _ = res.Body.ToString()
		res.Body.FromJsonTo(&result)
		l.Out <- result
	}
}

func (l *LuisAnalyze) OnAppId(appid string) {
	l.appid = appid
}

func (l *LuisAnalyze) OnKey(key string) {
	l.key = key
}

func (l *LuisAnalyze) promptBegin() chan<- bool {
	var tick bool
	end := make(chan bool)
	fmt.Printf("正在查询。。。")
	go func() {
		for {
			select {
			case <-end:
				fmt.Printf("\r")
				return
			case <-time.After(500 * time.Millisecond):

				if tick {

					fmt.Printf("\r.")
				} else {
					fmt.Printf("\r")
				}

				tick = !tick
			}
		}
	}()
	return end
}
