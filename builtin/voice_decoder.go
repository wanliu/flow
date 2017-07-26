package builtin

// 微信中的语音amr文件其实是silk格式，为了使百度语音能够识别，先转化成为mp3格式，再转回amr格式。
// 直接转成amr格式肯定也是可以的，这个作为下一步的工作去摸索
// silk to mp3: github.com/kn007/silk-v3-decoder
// mp3 to amr: github.com/seka17/mp3-amr-converter
// 需要安装依赖 ffmpeg

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	flow "github.com/wanliu/goflow"
)

const (
	MP3FORMAT             = "mp3"
	AudioBitRate          = "12.2k" // in Hz
	NumberOfAudioChannels = "1"
	AudioSamplingRateAMR  = "8000"
)

type VoiceDecoder struct {
	flow.Component

	In   <-chan string
	Next chan<- string
	Out  chan<- ReplyData

	// WechatPath string
	// Mp3Path    string
	// AmrPath    string
}

func NewVoiceDecoder() interface{} {
	return new(VoiceDecoder)
}

func (c VoiceDecoder) OnIn(input string) {
	if !strings.HasSuffix(input, "amr") {
		replyData := ReplyData{"错误，不支持的音频格式", nil}
		c.Out <- replyData
		return
	}

	if _, err := os.Stat(input); os.IsNotExist(err) {
		replyData := ReplyData{"错误，制定的音频文件不存在", nil}
		c.Out <- replyData
		return
	}

	shPath, _ := filepath.Abs("./cmd/silk-v3-decoder/converter.sh")
	comm := exec.Command("/bin/sh", shPath, input, MP3FORMAT)
	if err := comm.Run(); err != nil {
		log.Printf("===1 %v", err)
		replyData := ReplyData{err.Error(), nil}
		c.Out <- replyData
		return
	}

	pathPre := strings.Replace(input, ".amr", "", -1)
	conPath := pathPre + "_copy.amr"
	comm = exec.Command("ffmpeg", "-i", pathPre+".mp3", "-ab", AudioBitRate, "-ac", NumberOfAudioChannels, "-ar", AudioSamplingRateAMR, conPath)
	if err := comm.Run(); err != nil {
		if _, err := os.Stat(conPath); os.IsNotExist(err) {
			replyData := ReplyData{"解码音频文件失败", nil}
			c.Out <- replyData
			return
		}

	}

	c.Next <- conPath
}
