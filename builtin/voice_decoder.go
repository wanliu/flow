package builtin

// 微信中的语音amr文件其实是silk格式，为了使百度语音能够识别，先转化成为mp3格式，再转回amr格式。
// 直接转成amr格式肯定也是可以的，这个作为下一步的工作去摸索
// silk to mp3: github.com/kn007/silk-v3-decoder
// mp3 to amr: github.com/seka17/mp3-amr-converter
// 需要安装依赖 ffmpeg

import (
	"encoding/base64"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

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

func (c VoiceDecoder) OnIn(src string) {
	dst, _ := base64.URLEncoding.DecodeString(src)

	rand.Seed(time.Now().UnixNano())
	randId := strconv.Itoa(rand.Intn(10000000))

	tempDir, _ := filepath.Abs("./.tmp/")
	if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		log.Printf("[AUDIO]创建临时文件夹 %v", tempDir)
		os.Mkdir(tempDir, 0777)
		tempDir = tempDir + "/"
	}

	filename := tempDir + strconv.Itoa(int(time.Now().Unix())) + randId + ".mp3"
	filename, _ = filepath.Abs(filename)

	log.Printf("[AUDIO]输出mp3文件到 %v", filename)
	ioutil.WriteFile(filename, dst, 0644)

	pathPre := strings.Replace(filename, ".mp3", "", -1)
	conPath := pathPre + "_copy.amr"

	log.Printf("[AUDIO]输出amr文件到 %v", conPath)
	comm := exec.Command("ffmpeg", "-i", filename, "-ab", AudioBitRate, "-ac", NumberOfAudioChannels, "-ar", AudioSamplingRateAMR, conPath)
	if err := comm.Run(); err != nil {
		if _, err := os.Stat(conPath); os.IsNotExist(err) {
			log.Printf("[AUDIO ERROR]输出amr文件失败")

			replyData := ReplyData{"解码音频文件失败", nil}
			c.Out <- replyData
			return
		}
	}

	c.Next <- conPath
}
