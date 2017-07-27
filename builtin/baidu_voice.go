package builtin

import (
	"encoding/base64"
	"io/ioutil"
	"strings"

	voice "github.com/chenqinghe/baidu-ai-go-sdk/voice"
	flow "github.com/wanliu/goflow"
)

type BaiduVoice struct {
	flow.Component

	Appid     string
	Apikey    string
	Secretkey string

	AppId  <-chan string
	ApiKey <-chan string
	SecKey <-chan string
	Path   <-chan string

	Next chan<- string
	Out  chan<- ReplyData
}

func NewBaiduVoice() interface{} {
	return new(BaiduVoice)
}

func (c *BaiduVoice) OnAppId(id string) {
	c.Appid = id
}

func (c *BaiduVoice) OnApiKey(key string) {
	c.Apikey = key
}

func (c *BaiduVoice) OnSecKey(key string) {
	c.Secretkey = key
}

func (c *BaiduVoice) OnPath(path string) {
	client := voice.NewVoiceClient(c.Apikey, c.Secretkey)

	data, err := ioutil.ReadFile(path)
	bData := base64.StdEncoding.EncodeToString(data)
	leng := len(data)

	var ap voice.ASRParams = voice.ASRParams{
		Format:  "amr",
		Rate:    8000,
		Channel: 1,
		Token:   client.AccessToken,
		Cuid:    "565985655244",
		Lan:     "zh",
		Speech:  bData,
		Len:     leng,
	}

	strs, err := client.SpeechToText(ap)
	if err != nil {
		replyData := ReplyData{err.Error(), nil}
		c.Out <- replyData
		return
	}

	replyData := ReplyData{strings.Join(strs, ", "), nil}
	c.Out <- replyData
}

type BaiduRes struct {
	Err_no    int
	Corpus_no string
	Err_msg   string
	Result    []string
}
