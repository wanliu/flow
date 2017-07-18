package ai

import (
	"fmt"

	"github.com/hysios/apiai-go"
)

func ApiAiQuery(queryString, token, sessionId string) (apiai.Result, error) {
	client, err := apiai.NewClient(
		&apiai.ClientConfig{
			Token:      token,
			QueryLang:  "zh-CN", //Default en
			SpeechLang: "zh-CN", //Default en-US
		},
	)
	if err != nil {
		fmt.Printf("AI CONFIG ERROR: %v\n", err)
		return apiai.Result{}, err
	}

	fmt.Printf("Query: %v, token: %v, sessionid: %v\n", queryString, token, sessionId)
	//Set the query string and your current user identifier.
	qr, err := client.Query(apiai.Query{Query: []string{queryString}, SessionId: sessionId})
	if err != nil {
		fmt.Printf("AI REQUEST ERROR: %v\n", err)
		return apiai.Result{}, err
	}

	return qr.Result, nil
}
