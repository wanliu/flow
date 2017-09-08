package graphqlClient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

const (
	URL = "http://localhost:8000/graphql"
	// Cookie = "csrf_token=88a8a8a8a8a8a8a8a8a8a8a8=; token=9s9s9s9s9s9s9s9s9s9s9s9s9s9s9s9s9s9s9s9s9s9s9s9s9s9s9s9s9s9s9s9s9s9s9s9s9s9s9s9s9s9s9s9s9s9s9s9s9s9s9s9s9s9s9s9s9s9s9s9s9s9s9s9s9s9s9s9s9s9s9s9s=="
)

func QueryToRequest(queryString, variablesString string) string {
	if variablesString == "" {
		return `{"query":` + strconv.QuoteToASCII(queryString) + `}`
	} else {
		return `{"query":` + strconv.QuoteToASCII(queryString) + `,"variables":` + variablesString + `}`
	}
}

func MakeGraphqlRequest(requestString string) (*CreateOrderResponse, error) {
	// http://stackoverflow.com/questions/24455147/how-do-i-send-a-json-string-in-a-post-request-in-go
	fmt.Println("URL:", URL)

	var str = []byte(requestString)
	req, err := http.NewRequest("POST", URL, bytes.NewBuffer(str))
	req.Header.Set("Content-Type", "application/json")
	// req.Header.Set("Cookie", Cookie)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var cr *CreateOrderResponse
	switch resp.StatusCode {
	case http.StatusOK:
		body, _ := ioutil.ReadAll(resp.Body)
		err = json.Unmarshal(body, &cr)
		if err != nil {
			return nil, err
		}

		return cr, nil
	default:
		return nil, fmt.Errorf("apiai: wops something happens because status code is %v", resp.StatusCode)
	}
}

// func main() {
// 	queries := []string{
// 		`mutation Update {
// 			user_update(
// 				user: {
// 					fname: "John",
// 					lname: "Smith",
// 				},
// 			) {
// 				fname
// 				lname
// 			}
// 		}`,
// 		`query User {
// 			user {
// 				fname
// 				lname
// 			}
// 		}`,
// 	}

// 	for _, query := range queries {
// 		makeRequest(queryToRequest(strconv.QuoteToASCII(query)))
// 	}
// }
