package builtin

import (
	"encoding/json"
	"log"
	"os/exec"
	"path/filepath"

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
	pyPath, _ := filepath.Abs("./lib/baidu_sdk.py")
	comm := exec.Command("python", pyPath, c.Appid, c.Apikey, c.Secretkey, path)
	output, err := comm.CombinedOutput()

	if err != nil {
		replyData := ReplyData{err.Error(), nil}
		c.Out <- replyData
		return
	}

	var res BaiduRes
	json.Unmarshal([]byte(output), &res)
	log.Println(res)

	replyData := ReplyData{string(res.Result[0]), nil}
	c.Out <- replyData
}

type BaiduRes struct {
	Err_no    int
	Corpus_no string
	Err_msg   string
	Result    []string
}
