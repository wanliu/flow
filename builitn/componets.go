package builitn

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http/cookiejar"
	"os"
	"time"

	"github.com/franela/goreq"

	flow "github.com/wanliu/goflow"
)

type GetElement struct {
	flow.Component
}

type Split struct {
	flow.Component
	In  <-chan string
	Out chan<- string
}

type Output struct {
	flow.Component
	In      <-chan string
	Options <-chan map[string]interface{}
	Out     chan<- string
}

type ReadFile struct {
	flow.Component
	Read <-chan string
	// Encoding <-chan string
	Out   chan<- string
	Error chan<- error
}

type ReadLine struct {
	flow.Component
	ReadLine <-chan string
	Out      chan<- string
	Error    chan<- error
}

func NewGetElement() interface{} {
	return new(GetElement)
}

func NewSplit() interface{} {
	return new(Split)
}

func NewOutput() interface{} {
	return new(Output)
}

func NewReadFile() interface{} {
	return new(ReadFile)
}

func NewReadLine() interface{} {
	return new(ReadLine)
}

func NewLuisAnalyze() interface{} {
	return new(LuisAnalyze)
}

type Luis struct {
	Url       string
	AppID     string
	Key       string
	Secret    string
	cookieJar *cookiejar.Jar
	Proxy     string
}

func NewLuis(appid, key, secret string) *Luis {
	jar, _ := cookiejar.New(nil)
	// https://westus.api.cognitive.microsoft.com/luis/v2.0/apps/052297dc-12b9-4044-8220-a21a20d72581?subscription-key=6b916f7c107643069c242cf881609a82&timezoneOffset=0.0&verbose=true&q=
	return &Luis{
		Url:       "https://westus.api.cognitive.microsoft.com/luis/v2.0/apps/%s",
		AppID:     appid,
		Key:       key,
		Secret:    secret,
		cookieJar: jar,
		Proxy:     "http://192.168.0.151:1087",
	}
}

type LuisInput struct {
	Query string
}

type QueryParams struct {
	Key      string `url:"subscription-key"`
	TimeZone string `url:"timezoneOffset"`
	Query    string `url:"q"`
	Verbose  bool   `url:"verbose"`
}

type LuisAnalyze struct {
	flow.Component
	In  <-chan string
	Out chan<- string
}

func (l *LuisAnalyze) OnIn(input string) {
	var (
		luis   = NewLuis("052297dc-12b9-4044-8220-a21a20d72581", "6b916f7c107643069c242cf881609a82", "")
		url    = fmt.Sprintf(luis.Url, luis.AppID)
		params = QueryParams{
			Key:      luis.Key,
			TimeZone: "0",
			Query:    input,
			Verbose:  true,
		}
	)

	res, err := goreq.Request{
		Uri:         url,
		QueryString: params,
		Timeout:     10 * time.Second,
		CookieJar:   luis.cookieJar,
		Proxy:       luis.Proxy,
	}.Do()

	if err != nil {
		l.Out <- err.Error()
	} else {
		result, _ := res.Body.ToString()
		l.Out <- result
	}
}

func (o *Output) OnIn(msg string) {
	log.Printf("output: %s", msg)
	o.Out <- msg
}

func (rf *ReadFile) Init() {
	// rf.Error = make(chan error)
}

func (rf *ReadFile) OnRead(filename string) {
	if buf, err := ioutil.ReadFile(filename); err != nil {
		rf.Error <- err
	} else {
		rf.Out <- string(buf)
	}
}

func (rl *ReadLine) OnIn(filename string) {
	if f, err := os.Open(filename); err != nil {
		rl.Error <- err
	} else {
		defer f.Close()
		var reader = bufio.NewReader(f)

		if line, _, err := reader.ReadLine(); err != nil {
			rl.Error <- err
		} else {
			rl.Out <- string(line)
		}
	}
}

func (sp *Split) OnIn(msg string) {
	sp.Out <- msg
}
