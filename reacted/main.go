package main

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type RequestType struct {
	Type string `json:"type"`
}

type URLVerification struct {
	Token     string `json:"token"`
	Challenge string `json:"challenge"`
}

type RequestEvent struct {
	Event ReactionEvent `json:"event"`
}

type ReactionEvent struct {
	Type     string `json:"type"`
	User     string `json:"user"`
	Reaction string `json:"reaction"`
	ItemUser string `json:"item_user"`
	Item     Item   `json:"item"`
	EventTS  string `json:"event_ts"`
}

type Item struct {
	Type    string `json:"type"`
	Channel string `json:"channel"`
	TS      string `json:"ts"`
}

func handler(request events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	// events.LambdaFunctionURLRequest{}
	fmt.Println(request.Body)
	var reqType RequestType
	err := json.Unmarshal([]byte(request.Body), &reqType)
	if err != nil {
		fmt.Println(err)
		return events.LambdaFunctionURLResponse{}, err
	}

	if reqType.Type == "url_verification" {
		var urlVerification URLVerification
		err := json.Unmarshal([]byte(request.Body), &urlVerification)
		if err != nil {
			fmt.Println(err)
			return events.LambdaFunctionURLResponse{}, err
		}
		return events.LambdaFunctionURLResponse{
			StatusCode: 200,
			Body:       urlVerification.Challenge,
		}, nil
	}

	if reqType.Type == "event_callback" {
		var event ReactionEvent
		err := json.Unmarshal([]byte(request.Body), &event)
		if err != nil {
			fmt.Println(err)
			return events.LambdaFunctionURLResponse{}, err
		}
		// Itemからメッセージの内容を取得する

		// メッセージからリンクをパースする

		// rlを使ってNotionに書き込む

		return events.LambdaFunctionURLResponse{
			StatusCode: 200,
			Body:       "{\"message\": \"reaction_added\"}",
		}, nil
	}

	greeting := "Hello World"

	return events.LambdaFunctionURLResponse{
		StatusCode: 200,
		Body:       "{\"message\": \"" + greeting + "\"}",
	}, nil
}

func main() {
	lambda.Start(handler)
}
