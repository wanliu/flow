package builtin

import (
	"context"
	"io/ioutil"
	"log"

	speech "cloud.google.com/go/speech/apiv1"
	"google.golang.org/api/option"
	speechpb "google.golang.org/genproto/googleapis/cloud/speech/v1"

	flow "github.com/wanliu/goflow"
)

type GoogleVoice struct {
	flow.Component

	Apikeyfile string

	ApiKeyFile <-chan string
	Path       <-chan string

	Next chan<- string
	Out  chan<- string
}

func NewGoogleVoice() interface{} {
	return new(GoogleVoice)
}

func (c *GoogleVoice) OnApiKeyFile(filepath string) {
	c.Apikeyfile = filepath
}

func (c *GoogleVoice) OnPath(path string) {

	ctx := context.Background()
	cl, err := speech.NewClient(ctx, option.WithServiceAccountFile(c.Apikeyfile))
	if err != nil {
		// TODO: Handle error.
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		// replyData := ReplyData{err.Error(), nil}
		c.Out <- err.Error()
		return
	}

	content := &speechpb.RecognitionAudio_Content{
		Content: data,
	}

	req := &speechpb.RecognizeRequest{
		Config: &speechpb.RecognitionConfig{
			Encoding:        speechpb.RecognitionConfig_AMR,
			SampleRateHertz: 8000,
			LanguageCode:    "cmn-Hans-CN",
		},
		Audio: &speechpb.RecognitionAudio{
			AudioSource: content,
		},
	}

	resp, err := cl.Recognize(ctx, req)
	if err != nil {
		c.Out <- err.Error()
		return
	}
	// _ = resp
	log.Printf("resp %#v", resp)

	// data, err := ioutil.ReadFile(path)
	// bData := base64.StdEncoding.EncodeToString(data)
	// leng := len(data)

	// var ap voice.ASRParams = voice.ASRParams{
	// 	Format:  "amr",
	// 	Rate:    8000,
	// 	Channel: 1,
	// 	Token:   client.AccessToken,
	// 	Cuid:    "565985655244",
	// 	Lan:     "zh",
	// 	Speech:  bData,
	// 	Len:     leng,
	// }

	// strs, err := client.SpeechToText(ap)
	// if err != nil {
	// 	// replyData := ReplyData{err.Error(), nil}
	// 	c.Out <- err.Error()
	// 	return
	// }

	// replyData := ReplyData{strings.Join(strs, ", "), nil}
	c.Out <- "nothing"
}

type GoogleRes struct {
	Err_no    int
	Corpus_no string
	Err_msg   string
	Result    []string
}
