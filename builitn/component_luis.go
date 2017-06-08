package builitn

import (
	"fmt"
	"github.com/franela/goreq"
	"net/http/cookiejar"
	"time"

	flow "github.com/wanliu/goflow"
)

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

type IntentScore struct {
	Intent string  `json:"intent"`
	Score  float64 `json:"score"`
}

type Resolution struct {
	Date string
}

type EntityScore struct {
	Entity     string     `json:"entity"`
	Type       string     `json:"type"`
	StartIndex int        `json:"startIndex"`
	EndIndex   int        `json:"endIndex"`
	Score      float64    `json:"score"`
	Resolution Resolution `json:"resolution"`
}

type QueryParams struct {
	Key      string `url:"subscription-key"`
	TimeZone string `url:"timezoneOffset"`
	Query    string `url:"q"`
	Verbose  bool   `url:"verbose"`
}

type ResultParams struct {
	Query            string        `json:"query"`
	TopScoringIntent IntentScore   `json:"topScoringIntent"`
	Intents          []IntentScore `json:"intents"`
	Entities         []EntityScore `json:"entities"`
}

type LuisAnalyze struct {
	flow.Component
	In  <-chan string
	Out chan<- ResultParams
	// Out chan<- string
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
		result ResultParams
	)

	res, err := goreq.Request{
		Uri:         url,
		QueryString: params,
		Timeout:     10 * time.Second,
		CookieJar:   luis.cookieJar,
		Proxy:       luis.Proxy,
	}.Do()

	if err != nil {
		// l.Out <- err.Error()
		l.Out <- *new(ResultParams)
	} else {
		// result, _ = res.Body.ToString()
		res.Body.FromJsonTo(&result)
		l.Out <- result
	}
}
