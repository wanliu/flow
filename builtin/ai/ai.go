package ai

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/hysios/apiai-go"
)

func ApiAiQuery(queryString, token, sessionId, proxyUrl string) (apiai.Result, error) {
	var (
		client *apiai.ApiClient
		err    error
	)

	if proxyUrl != "" {
		client, err = apiai.NewClient(
			&apiai.ClientConfig{
				Token:      token,
				QueryLang:  "zh-CN", //Default en
				SpeechLang: "zh-CN", //Default en-US
				ProxyURL:   proxyUrl,
			},
		)
	} else {
		client, err = apiai.NewClient(
			&apiai.ClientConfig{
				Token:      token,
				QueryLang:  "zh-CN", //Default en
				SpeechLang: "zh-CN", //Default en-US
			},
		)
	}

	if err != nil {
		fmt.Printf("AI CONFIG ERROR: %v\n", err)
		return apiai.Result{}, err
	}

	rand.Seed(time.Now().UnixNano())
	randId := strconv.Itoa(rand.Intn(10000000))

	fmt.Printf("Query: %v, token: %v, sessionid: %v, proxy: %v\n", queryString, token, randId, proxyUrl)
	//Set the query string and your current user identifier.
	qr, err := client.Query(apiai.Query{Query: []string{queryString}, SessionId: randId})
	if err != nil {
		fmt.Printf("AI REQUEST ERROR: %v\n", err)
		return apiai.Result{}, err
	}

	return qr.Result, nil
}
