package builtin

// 微信中的语音amr文件其实是silk格式，为了使百度语音能够识别，先转化成为mp3格式，再转回amr格式。
// 直接转成amr格式肯定也是可以的，这个作为下一步的工作去摸索
// silk to mp3: github.com/kn007/silk-v3-decoder
// mp3 to amr: github.com/seka17/mp3-amr-converter
// 需要安装依赖 ffmpeg

import (
	"encoding/base64"
	"io/ioutil"
	// "math/rand"
	// "os"
	// "os/exec"
	// "path/filepath"
	// "strconv"
	// "strings"
	// "time"

	flow "github.com/wanliu/goflow"
)

type Mp3Reader struct {
	flow.Component

	In  <-chan string
	Out chan<- string
}

func NewMp3Reader() interface{} {
	return new(Mp3Reader)
}

func (c Mp3Reader) OnIn(src string) {
	dst, _ := ioutil.ReadFile(src)
	data := base64.URLEncoding.EncodeToString(dst)

	c.Out <- data
}
