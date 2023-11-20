package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"os"
	"regexp"
	"slack-reacted-rl/rl"
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

var (
	notionApiKey     string
	notionDatabaseId string
	slackBotToken    string
)

func handler(request events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	// events.LambdaFunctionURLRequest{}
	fmt.Println(request.Body)

	var reqType RequestType
	err := json.Unmarshal([]byte(request.Body), &reqType)
	if err != nil {
		fmt.Println(err)
		return events.LambdaFunctionURLResponse{}, err
	}

	err = loadSSMParameters()
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
		var event RequestEvent
		err := json.Unmarshal([]byte(request.Body), &event)
		if err != nil {
			fmt.Println(err)
			return events.LambdaFunctionURLResponse{}, err
		}

		// リアクションがrlスタンプじゃなければ終了
		if event.Event.Reaction != "rl" {
			return events.LambdaFunctionURLResponse{
				StatusCode: 200,
				Body:       "{\"message\": \"reacted another reaction\"}",
			}, nil
		}

		// Itemからメッセージの内容を取得する
		message, err := getMessage(event.Event.Item.Channel, event.Event.Item.TS)
		if err != nil {
			fmt.Println(err)
			return events.LambdaFunctionURLResponse{}, err
		}
		fmt.Println(message)
		// メッセージからリンクをパースする
		url := parseUrl(message)
		// rlを使ってNotionに書き込む
		err = rl.RL(url, notionApiKey, notionDatabaseId)
		if err != nil {
			fmt.Println(err)
			return events.LambdaFunctionURLResponse{}, err
		}

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

func loadSSMParameters() error {
	sess, err := session.NewSession()
	if err != nil {
		return err
	}

	svc := ssm.New(sess)

	values, err := svc.GetParameters(&ssm.GetParametersInput{
		Names: []*string{
			aws.String(os.Getenv("NOTION_API_KEY_PATH")),
			aws.String(os.Getenv("NOTION_DATABASE_ID_PATH")),
			aws.String(os.Getenv("SLACK_BOT_TOKEN_PATH")),
		},
		WithDecryption: aws.Bool(true),
	})

	if err != nil {
		return err
	}

	notionApiKey = *values.Parameters[0].Value
	notionDatabaseId = *values.Parameters[1].Value
	slackBotToken = *values.Parameters[2].Value

	return nil
}

func parseUrl(message string) string {
	r := regexp.MustCompile(`https?://[\w/:%#\$&\?\(\)~\.=\+\-]+`)
	return r.FindString(message)
}
